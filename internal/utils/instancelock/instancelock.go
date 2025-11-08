// Copyright 2025 The MathWorks, Inc.

package instancelock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const lockFileName = "matlab-mcp-core-server.lock"

// InstanceLock manages a lock file to prevent multiple instances from running
type InstanceLock struct {
	lockFilePath string
	pid          int
}

// New creates a new instance lock. The lock file will be created in the user's temp directory.
func New() (*InstanceLock, error) {
	tempDir := os.TempDir()
	lockFilePath := filepath.Join(tempDir, lockFileName)

	return &InstanceLock{
		lockFilePath: lockFilePath,
		pid:          os.Getpid(),
	}, nil
}

// TryLock attempts to acquire the lock. Returns true if lock was acquired, false if another instance is running.
func (l *InstanceLock) TryLock() (bool, error) {
	// Check if lock file exists
	if _, err := os.Stat(l.lockFilePath); err == nil {
		// Lock file exists, read the PID
		pidBytes, err := os.ReadFile(l.lockFilePath)
		if err != nil {
			// If we can't read it, assume it's stale and try to remove it
			os.Remove(l.lockFilePath)
			return l.createLock()
		}

		existingPID, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
		if err != nil {
			// Invalid PID in lock file, remove it
			os.Remove(l.lockFilePath)
			return l.createLock()
		}

		// Check if the process is still running
		if l.isProcessRunning(existingPID) {
			// Another instance is running
			return false, nil
		}

		// Process is dead, remove stale lock file
		os.Remove(l.lockFilePath)
	}

	// No lock file or stale lock, create new one
	return l.createLock()
}

// createLock creates the lock file with current PID
func (l *InstanceLock) createLock() (bool, error) {
	pidStr := strconv.Itoa(l.pid)
	err := os.WriteFile(l.lockFilePath, []byte(pidStr), 0o644)
	if err != nil {
		return false, fmt.Errorf("failed to create lock file: %w", err)
	}
	return true, nil
}

// Unlock removes the lock file
func (l *InstanceLock) Unlock() error {
	return os.Remove(l.lockFilePath)
}

// isProcessRunning checks if a process with the given PID is still running
func (l *InstanceLock) isProcessRunning(pid int) bool {
	return checkProcessRunningPlatformSpecific(pid)
}

