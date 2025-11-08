// Copyright 2025 The MathWorks, Inc.

package instancelock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	return l.TryLockWithKill(false)
}

// TryLockWithKill attempts to acquire the lock, optionally killing the existing instance if one is running.
// If killExisting is true and an existing instance is found, it will be terminated and the lock acquired.
// Returns true if lock was acquired, false if another instance is running and killExisting is false.
func (l *InstanceLock) TryLockWithKill(killExisting bool) (bool, error) {
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

		// Don't kill our own process (shouldn't happen, but safety check)
		if existingPID == l.pid {
			// We already have the lock
			return true, nil
		}

		// Check if the process is still running
		if l.isProcessRunning(existingPID) {
			// Another instance is running
			if killExisting {
				// Kill the existing instance
				if err := l.killProcess(existingPID); err != nil {
					return false, fmt.Errorf("failed to kill existing instance (PID %d): %w", existingPID, err)
				}
				// Wait for the process to exit (with retries)
				// We check up to 10 times with 100ms delay between checks (max 1 second wait)
				for i := 0; i < 10; i++ {
					if !l.isProcessRunning(existingPID) {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
				// Remove the lock file (process should be dead now)
				os.Remove(l.lockFilePath)
				return l.createLock()
			}
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

// killProcess terminates the process with the given PID
func (l *InstanceLock) killProcess(pid int) error {
	return killProcessPlatformSpecific(pid)
}

