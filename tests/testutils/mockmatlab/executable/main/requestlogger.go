// Copyright 2026 The MathWorks, Inc.

package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type requestLogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"messageType"`
	Content     string    `json:"content"`
}

type requestLogger struct {
	mu   sync.Mutex
	file *os.File
}

func newRequestLogger(path string) *requestLogger {
	f, err := os.Create(path) //nolint:gosec // Test utility writing to session temp dir
	if err != nil {
		log.Printf("failed to create request log: %v", err)
		return &requestLogger{}
	}
	return &requestLogger{file: f}
}

func (rl *requestLogger) log(messageType, content string) {
	if rl.file == nil {
		return
	}

	entry := requestLogEntry{
		Timestamp:   time.Now(),
		MessageType: messageType,
		Content:     content,
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := json.NewEncoder(rl.file).Encode(entry); err != nil {
		log.Printf("failed to write request log entry: %v", err)
	}
}

func (rl *requestLogger) close() {
	if rl.file == nil {
		return
	}
	if err := rl.file.Close(); err != nil {
		log.Printf("failed to close request log: %v", err)
	}
}
