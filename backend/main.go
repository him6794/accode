package main

import (
	"accode-go/handlers"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	allowOrigins := []string{
		"https://accode.justin0711.com",
		"http://accode.justin0711.com",
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	if raw := strings.TrimSpace(os.Getenv("ACCODE_CORS_ORIGINS")); raw != "" {
		customOrigins := make([]string, 0)
		for _, origin := range strings.Split(raw, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				customOrigins = append(customOrigins, trimmed)
			}
		}
		if len(customOrigins) > 0 {
			allowOrigins = customOrigins
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/run", handlers.RunCode)
		api.GET("/status/:session_id", handlers.GetStatus)
		api.POST("/input", handlers.SubmitInput)
		api.POST("/stop", handlers.StopExecution)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/run", handlers.RunCode)
	r.GET("/status/:session_id", handlers.GetStatus)
	r.POST("/submit_input", handlers.SubmitInput)
	r.POST("/stop", handlers.StopExecution)

	go handlers.CleanupSessions()

	addr := "0.0.0.0:5001"
	if raw := strings.TrimSpace(os.Getenv("ACCODE_API_ADDR")); raw != "" {
		if strings.Contains(raw, ":") {
			addr = raw
		} else {
			addr = "0.0.0.0:" + raw
		}
	}

	log.Printf("ACcode API listening on http://%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start API server: %v", err)
	}
}
