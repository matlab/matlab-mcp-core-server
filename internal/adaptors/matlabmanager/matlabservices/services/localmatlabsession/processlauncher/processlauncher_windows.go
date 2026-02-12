// Copyright 2025-2026 The MathWorks, Inc.
//go:build windows

package processlauncher

import (
	"context"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processlauncher/utils/winenvironmentbuilder"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"golang.org/x/sys/windows"
)

const EXPECTED_CHILD_MATLAB_PROCESS_NAME = "MATLAB.exe"

func startMatlab(ctx context.Context, logger entities.Logger, matlabRoot string, workingDir string, args []string, env []string, stdIO *stdIO) (*os.Process, error) {
	matlabPath := filepath.Join(matlabRoot, "bin", "matlab.exe")

	if _, err := os.Stat(matlabPath); err != nil {
		return nil, err
	}

	// Careful on Windows, we need to quote the args
	cmdLine := `"` + matlabPath + `"`
	for _, arg := range args {
		cmdLine += ` "` + arg + `"`
	}

	envBlock, err := winenvironmentbuilder.Build(env)
	if err != nil {
		return nil, fmt.Errorf("failed to build environment block: %w", err)
	}

	envPtr := envBlock.PointerToFirstElement()
	if envPtr == nil {
		return nil, fmt.Errorf("failed to get environment block pointer")
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

	creationFlags := uint32(windows.CREATE_NEW_PROCESS_GROUP | windows.DETACHED_PROCESS | windows.CREATE_UNICODE_ENVIRONMENT)

	err = windows.CreateProcess(
		nil,           // appName
		cmdLinePtr,    // commandLine
		nil,           // procSecurity
		nil,           // threadSecurity
		true,          // inheritHandles
		creationFlags, // creationFlags
		envPtr,        // env
		workingDirPtr, // currentDir
		&si,           // startupInfo
		&pi,           // outProcInfo
	)

	if err != nil {
		return nil, fmt.Errorf("error creating MATLAB process: %w", err)
	}

	// Close handles as we don't need them after FindProcess
	if closeErr := windows.CloseHandle(pi.Thread); closeErr != nil {
		logger.WithError(closeErr).Warn("failed to close thread handle")
	}
	if closeErr := windows.CloseHandle(pi.Process); closeErr != nil {
		logger.WithError(closeErr).Warn("failed to close process handle")
	}

	matlabLauncherProcess, err := os.FindProcess(int(pi.ProcessId))
	if err != nil {
		return nil, fmt.Errorf("error finding MATLAB launcher process: %w", err)
	}

	// On Windows, the process we launch is a launcher process that then launches the actual MATLAB process.
	// Therefore, we need to find the child process of the launcher process.
	// There should be only one process, and that would be the actual MATLAB process.
	matlabProcess, err := waitForMATLABProcess(ctx, logger, matlabLauncherProcess)
	if err != nil {
		return nil, err
	}

	return matlabProcess, nil
}

func waitForMATLABProcess(_ context.Context, logger entities.Logger, matlabLauncherProcess *os.Process) (*os.Process, error) {
	const pollInterval = 1 * time.Second
	const timeout = 20 * time.Second

	if matlabLauncherProcess.Pid < 0 || matlabLauncherProcess.Pid > 0xFFFFFFFF {
		return nil, fmt.Errorf("invalid process ID: %d", matlabLauncherProcess.Pid)
	}

	matlabLauncherProcessID := uint32(matlabLauncherProcess.Pid)

	waitCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	matlabProcess, err := retry.Retry(waitCtx, func() (*os.Process, bool, error) {
		process := getMATLABChildProcess(logger, matlabLauncherProcessID)
		return process, process != nil, nil
	}, retry.NewLinearRetryStrategy(pollInterval))

	if err != nil {
		return nil, fmt.Errorf("timeout waiting for matlab process to begin")
	}

	return matlabProcess, nil
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
