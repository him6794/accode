package models

import (
	"sync"
)

type Session struct {
	ID           string
	OutputBuffer string
	InputChan    chan string
	OutputChan   chan OutputMessage
	Completed    bool
	NeedsInput   bool
	StopChan     chan struct{}
	mu           sync.RWMutex
}

type OutputMessage struct {
	Output     string
	NeedsInput bool
	Completed  bool
}

func NewSession(id string) *Session {
	return &Session{
		ID:         id,
		InputChan:  make(chan string, 10),
		OutputChan: make(chan OutputMessage, 100),
		StopChan:   make(chan struct{}, 1),
	}
}

func (s *Session) AppendOutput(output string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.OutputBuffer += output
}

func (s *Session) GetAndClearOutput() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	output := s.OutputBuffer
	s.OutputBuffer = ""
	return output
}

func (s *Session) SetNeedsInput(needs bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.NeedsInput = needs
}

func (s *Session) GetNeedsInput() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.NeedsInput
}

func (s *Session) SetCompleted(completed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Completed = completed
}

func (s *Session) GetCompleted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Completed
}
