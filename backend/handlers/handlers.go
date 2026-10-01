package handlers

import (
	"accode-go/executor"
	"accode-go/models"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	sessions   = make(map[string]*models.Session)
	sessionsMu sync.RWMutex
	sessionID  int
)

func RunCode(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	sessionsMu.Lock()
	id := strconv.Itoa(sessionID)
	sessionID++
	session := models.NewSession(id)
	sessions[id] = session
	sessionsMu.Unlock()

	go executeCode(session, req.Code)

	c.JSON(http.StatusOK, gin.H{"sessionId": id})
}

func executeCode(session *models.Session, code string) {
	session.SetNeedsInput(false)
	session.SetCompleted(false)

	var stopped atomic.Bool
	executionDone := make(chan struct{})
	defer close(executionDone)

	go func() {
		select {
		case <-session.StopChan:
			stopped.Store(true)
		case <-executionDone:
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			session.AppendOutput(fmt.Sprintf("execution panic: %v\n", r))
			session.SetCompleted(true)
		}
	}()

	outputFunc := func(s string) {
		if stopped.Load() {
			return
		}
		session.AppendOutput(s)
	}

	inputFunc := func() string {
		session.SetNeedsInput(true)
		defer session.SetNeedsInput(false)

		for {
			if stopped.Load() {
				return ""
			}
			select {
			case input := <-session.InputChan:
				return input
			case <-session.StopChan:
				stopped.Store(true)
				return ""
			case <-time.After(100 * time.Millisecond):
			}
		}
	}

	execEngine := executor.NewExecutor(outputFunc, inputFunc)
	execEngine.SetStopCheck(func() bool {
		return stopped.Load()
	})

	if err := execEngine.Execute(code); err != nil && !stopped.Load() {
		session.AppendOutput(fmt.Sprintf("%v\n", err))
	}

	if stopped.Load() {
		session.AppendOutput("execution stopped by user\n")
	}

	session.SetNeedsInput(false)
	session.SetCompleted(true)
}

func GetStatus(c *gin.Context) {
	sid := c.Param("session_id")

	sessionsMu.RLock()
	session, exists := sessions[sid]
	sessionsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	output := session.GetAndClearOutput()
	needsInput := session.GetNeedsInput()
	completed := session.GetCompleted()

	c.JSON(http.StatusOK, gin.H{
		"output":     output,
		"needsInput": needsInput,
		"completed":  completed,
	})
}

func SubmitInput(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
		Input     string `json:"input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	sessionsMu.RLock()
	session, exists := sessions[req.SessionID]
	sessionsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	if !session.GetNeedsInput() {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "session is not waiting for input"})
		return
	}

	select {
	case session.InputChan <- req.Input:
		c.JSON(http.StatusOK, gin.H{"message": "input submitted"})
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "input submit timeout"})
	}
}

func StopExecution(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload, missing sessionId"})
		return
	}

	sessionsMu.RLock()
	session, exists := sessions[req.SessionID]
	sessionsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	if session.GetCompleted() {
		c.JSON(http.StatusOK, gin.H{"message": "execution already completed"})
		return
	}

	select {
	case session.StopChan <- struct{}{}:
		c.JSON(http.StatusOK, gin.H{"message": "execution stop requested"})
	default:
		c.JSON(http.StatusOK, gin.H{"message": "execution is already stopping"})
	}
}

func CleanupSessions() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sessionsMu.Lock()
		for id, session := range sessions {
			if session.GetCompleted() {
				delete(sessions, id)
			}
		}
		sessionsMu.Unlock()
	}
}
