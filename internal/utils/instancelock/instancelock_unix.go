// Copyright 2025 The MathWorks, Inc.
//go:build !windows

package instancelock

import (
	"os"
	"syscall"
)

// checkProcessRunningPlatformSpecific performs Unix-specific process existence check
func checkProcessRunningPlatformSpecific(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix, sending signal 0 checks if process exists without actually sending a signal
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// killProcessPlatformSpecific terminates a process on Unix
func killProcessPlatformSpecific(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// Send SIGTERM first for graceful shutdown
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	// Note: We could wait and then send SIGKILL if needed, but for simplicity
	// we'll just send SIGTERM and let the OS handle cleanup
	return nil
}

