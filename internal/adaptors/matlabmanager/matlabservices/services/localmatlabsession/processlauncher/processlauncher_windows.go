// Copyright 2025 The MathWorks, Inc.
//go:build windows

package processlauncher

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"golang.org/x/sys/windows"
)

const EXPECTED_CHILD_MATLAB_PROCESS_NAME = "MATLAB.exe"

func startMatlab(logger entities.Logger, matlabRoot string, workingDir string, args []string, env []string, stdIO *stdIO) (*os.Process, error) {
	matlabPath := filepath.Join(matlabRoot, "bin", "matlab.exe")

	if _, err := os.Stat(matlabPath); err != nil {
		return nil, err
	}

	// Careful on Windows, we need to quote the args
	cmdLine := `"` + matlabPath + `"`
	for _, arg := range args {
		cmdLine += ` "` + arg + `"`
	}

	envBlock, err := buildEnvironmentBlock(env)
	if err != nil {
		return nil, fmt.Errorf("failed to build environment block: %w", err)
	}

	cmdLinePtr, err := windows.UTF16PtrFromString(cmdLine)
	if err != nil {
		return nil, fmt.Errorf("error converting command line: %w", err)
	}

	workingDirPtr, err := windows.UTF16PtrFromString(workingDir)
	if err != nil {
		return nil, fmt.Errorf("error converting working directory: %w", err)
	}

	var si windows.StartupInfo
	var pi windows.ProcessInformation

	si.Cb = uint32(unsafe.Sizeof(si))

	si.StdInput = windows.Handle(stdIO.stdIn.Fd())
	si.StdOutput = windows.Handle(stdIO.stdOut.Fd())
	si.StdErr = windows.Handle(stdIO.stdErr.Fd())

	// Set window to start minimized to reduce visual impact
	// STARTF_USESHOWWINDOW = 0x00000001, SW_MINIMIZE = 6
	si.Flags = 0x00000001 // STARTF_USESHOWWINDOW
	si.ShowWindow = 6     // SW_MINIMIZE

	creationFlags := uint32(windows.CREATE_NEW_PROCESS_GROUP | windows.DETACHED_PROCESS | windows.CREATE_UNICODE_ENVIRONMENT)

	err = windows.CreateProcess(
		nil,           // appName
		cmdLinePtr,    // commandLine
		nil,           // procSecurity
		nil,           // threadSecurity
		true,          // inheritHandles
		creationFlags, // creationFlags
		envBlock,      // env
		workingDirPtr, // currentDir
		&si,           // startupInfo
		&pi,           // outProcInfo
	)

	if err != nil {
		return nil, fmt.Errorf("error creating MATLAB process: %w", err)
	}

	// Close thread handle as we don't need it
	if err := windows.CloseHandle(pi.Thread); err != nil {
		logger.WithError(err).Warn("failed to close thread handle")
	}

	matlabLauncherProcess, err := os.FindProcess(int(pi.ProcessId))
	if err != nil {
		return nil, fmt.Errorf("error finding MATLAB launcher process: %w", err)
	}

	// On Windows, MATLAB uses a launcher architecture:
	// - The MCP launches matlab.exe (launcher process)
	// - The launcher then spawns MATLAB.exe (actual MATLAB process)
	// - This is expected Windows behavior, not a bug
	// - We need to find the child process of the launcher process to track the actual MATLAB instance
	// - There should be only one child process, which is the actual MATLAB process
	// Note: If you see TWO FULL MATLAB instances (not just matlab.exe + MATLAB.exe), this suggests
	// multiple MCP server processes are running, which would cause duplicate session initialization.
	matlabProcess, err := waitForMATLABProcess(logger, matlabLauncherProcess)
	if err != nil {
		return nil, err
	}

	return matlabProcess, nil
}

// buildEnvironmentBlock builds a UTF-16 environment block suitable for CreateProcess.
// - Each entry is "name=value" NUL-terminated
// - The block ends with an extra NUL
// - Entries are sorted case-insensitively by name (Windows requirement for Unicode blocks)
// - Duplicate names are resolved by keeping the last occurrence.
// - Returning the first pointer of the array of the environment block as required by the CreateProcess function.
func buildEnvironmentBlock(env []string) (*uint16, error) {
	if len(env) == 0 {
		return nil, nil
	}

	type kv struct{ nameUpper, entry string }

	// Deduplicate by case-insensitive name: last wins
	m := make(map[string]string, len(env))
	for _, e := range env {
		if strings.IndexByte(e, 0) >= 0 {
			return nil, fmt.Errorf("environment entry contains NUL: %q", e)
		}
		name, value, ok := strings.Cut(e, "=")
		if !ok || name == "" {
			return nil, fmt.Errorf("invalid environment entry (expected name=value): %q", e)
		}
		if name[0] == '=' {
			return nil, fmt.Errorf("invalid environment variable name (starts with '='): %q", e)
		}
		m[strings.ToUpper(name)] = name + "=" + value
	}

	// Collect and sort by case-insensitive name
	entries := make([]kv, 0, len(m))
	for upperName, entry := range m {
		entries = append(entries, kv{nameUpper: upperName, entry: entry})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].nameUpper < entries[j].nameUpper })

	// Flatten to UTF-16: each entry is NUL-terminated; add an extra NUL at the end
	var block []uint16
	for _, e := range entries {
		w, err := windows.UTF16FromString(e.entry) // includes trailing NUL
		if err != nil {
			return nil, fmt.Errorf("failed to encode env entry %q: %w", e.entry, err)
		}
		block = append(block, w...)
	}
	block = append(block, 0) // final extra NUL
	return &block[0], nil
}

func waitForMATLABProcess(logger entities.Logger, matlabLauncherProcess *os.Process) (*os.Process, error) {
	timeout := time.After(20 * time.Second)
	tick := time.Tick(1 * time.Second)

	if matlabLauncherProcess.Pid < 0 || matlabLauncherProcess.Pid > 0xFFFFFFFF {
		return nil, fmt.Errorf("invalid process ID: %d", matlabLauncherProcess.Pid)
	}

	matlabLauncherProcessID := uint32(matlabLauncherProcess.Pid)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for matlab process to begin")
		case <-tick:
			matlabProcess := getMATLABChildProcess(logger, matlabLauncherProcessID)
			if matlabProcess != nil {
				return matlabProcess, nil
			}
		}
	}
}

func getMATLABChildProcess(logger entities.Logger, matlabLauncherProcessID uint32) *os.Process {
	snapshot, err := windows.CreateToolhelp32Snapshot(
		windows.TH32CS_SNAPPROCESS,
		0, // When CreateToolhelp32Snapshot is called with TH32CS_SNAPPROCESS, processID is ignored, use 0
	)
	if err != nil {
		logger.WithError(err).Warn("error creating process snapshot")
		return nil
	}

	defer func() {
		if err := windows.CloseHandle(snapshot); err != nil {
			logger.WithError(err).Warn("failed to close snapshot handle")
		}
	}()

	for processEntry := range snapshotProcessIterator(snapshot) {
		if processEntry.ParentProcessID != matlabLauncherProcessID {
			continue
		}

		if windows.UTF16ToString(processEntry.ExeFile[:]) != EXPECTED_CHILD_MATLAB_PROCESS_NAME {
			continue
		}

		childProcess, err := os.FindProcess(int(processEntry.ProcessID))
		if err != nil {
			logger.WithError(err).Warn("we found a child process in the snapshot but could os.FindProcess it")
			continue
		}

		return childProcess
	}

	return nil
}

func snapshotProcessIterator(snapshot windows.Handle) iter.Seq[windows.ProcessEntry32] {
	return func(yield func(pe windows.ProcessEntry32) bool) {
		var pe windows.ProcessEntry32
		pe.Size = uint32(unsafe.Sizeof(pe))

		for err := windows.Process32First(snapshot, &pe); err == nil; err = windows.Process32Next(snapshot, &pe) {
			if !yield(pe) {
				return
			}
		}
	}
}
