// Copyright 2025 The MathWorks, Inc.
//go:build windows

package instancelock

import (
	"golang.org/x/sys/windows"
)

// checkProcessRunningPlatformSpecific performs Windows-specific process existence check
func checkProcessRunningPlatformSpecific(pid int) bool {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(handle)

	var exitCode uint32
	err = windows.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false
	}

	// STILL_ACTIVE = 259 on Windows
	const STILL_ACTIVE = 259
	return exitCode == STILL_ACTIVE
}

// killProcessPlatformSpecific terminates a process on Windows
func killProcessPlatformSpecific(pid int) error {
	handle, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	err = windows.TerminateProcess(handle, 1)
	if err != nil {
		return err
	}

	return nil
}

