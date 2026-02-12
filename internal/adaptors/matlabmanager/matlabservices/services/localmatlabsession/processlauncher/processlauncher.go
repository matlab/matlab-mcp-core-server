// Copyright 2025-2026 The MathWorks, Inc.

package processlauncher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const gracefulShutdownTimeout = 2 * time.Minute

type MATLABProcessLauncher struct{}

func New() *MATLABProcessLauncher {
	return &MATLABProcessLauncher{}
}

func (l *MATLABProcessLauncher) Launch(
	ctx context.Context,
	logger entities.Logger,
	sessionRoot string,
	matlabRoot string,
	workingDir string,
	args []string,
	env []string,
) (int, func(), <-chan struct{}, error) {
	stdIO, stdIOCleanup, err := createLocalStdioForNewProcess(logger, sessionRoot)
	if err != nil {
		return 0, nil, nil, err
	}

	// Use WithoutCancel to preserve existing behaviour: startup is not cancellable.
	// The context is threaded through for future use but does not affect startup.
	process, err := startMatlab(context.WithoutCancel(ctx), logger, matlabRoot, workingDir, args, env, stdIO)
	if err != nil {
		stdIOCleanup()
		return 0, nil, nil, fmt.Errorf("failed to start MATLAB process: %w", err)
	}

	processExited := make(chan struct{})
	waitResult := make(chan error, 1)

	go func() {
		_, err := process.Wait()
		waitResult <- err
		close(processExited)
	}()

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			logger.Debug("Waiting for MATLAB process to exit gracefully")

			select {
			case err := <-waitResult:
				logger.Debug("Done waiting for MATLAB process to exit")
				if err != nil && !errors.Is(err, os.ErrProcessDone) {
					logger.Warn("MATLAB process did not exit gracefully, forcefully kill it")
					killMATLABProcess(logger, process)
				}
			case <-time.After(gracefulShutdownTimeout):
				logger.Warn("Timed out waiting for MATLAB process to exit gracefully, forcefully kill it")
				killMATLABProcess(logger, process)
			}

			stdIOCleanup()
		})
	}

	return process.Pid, cleanup, processExited, nil
}

func killMATLABProcess(logger entities.Logger, process *os.Process) {
	err := process.Kill()
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		logger.WithError(err).Warn("Failed to kill MATLAB process")
	}
}

type stdIO struct {
	stdIn        *os.File
	stdOut       *os.File
	stdErr       *os.File
	writeToStdIn *os.File
}

func (s *stdIO) cleanup(logger entities.Logger) {
	for _, file := range []*os.File{s.stdIn, s.stdOut, s.stdErr, s.writeToStdIn} {
		if file != nil {
			err := file.Close()
			if err != nil {
				logger.WithError(err).Warn(fmt.Sprintf("Failed to close %v", file))
			}
		}
	}
}

func createLocalStdioForNewProcess(logger entities.Logger, sessionRoot string) (*stdIO, func(), error) {
	stdIO := &stdIO{}

	stdOut, err := os.Create(filepath.Join(sessionRoot, "matlab_stdout.log")) //nolint:gosec // We construct this path, and file
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdOut log file: %w", err)
	}

	stdIO.stdOut = stdOut

	stdErr, err := os.Create(filepath.Join(sessionRoot, "matlab_stderr.log")) //nolint:gosec // We construct this path, and file
	if err != nil {
		stdIO.cleanup(logger)
		return nil, nil, fmt.Errorf("failed to create stdErr log file: %w", err)
	}

	stdIO.stdErr = stdErr

	stdIn, writeToStdIn, err := os.Pipe()
	if err != nil {
		stdIO.cleanup(logger)
		return nil, nil, fmt.Errorf("failed to create stdIn pipe: %w", err)
	}

	stdIO.stdIn = stdIn
	stdIO.writeToStdIn = writeToStdIn

	return stdIO, func() {
		stdIO.cleanup(logger)
	}, nil
}
