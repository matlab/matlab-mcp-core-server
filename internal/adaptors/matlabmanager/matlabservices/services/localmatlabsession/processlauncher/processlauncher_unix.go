// Copyright 2025-2026 The MathWorks, Inc.
//go:build !windows

package processlauncher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"golang.org/x/sys/unix"
)

func startMatlab(_ context.Context, _ entities.Logger, matlabRoot string, workingDir string, args []string, env []string, stdIO *stdIO) (*os.Process, error) {
	matlabPath := filepath.Join(matlabRoot, "bin", "matlab")
	if _, err := os.Stat(matlabPath); err != nil {
		return nil, err
	}

	attr := &os.ProcAttr{
		Dir:   workingDir,
		Env:   env,
		Files: []*os.File{stdIO.stdIn, stdIO.stdOut, stdIO.stdErr},
		Sys: &unix.SysProcAttr{
			Setsid: true, // Create a new session
		},
	}

	// Careful here, for start process, we need the path first. From the doc:
	//   > StartProcess starts a new process with the program, arguments and attributes specified by name, argv and attr.
	//   > The argv slice will become os.Args in the new process, so it normally starts with the program name.
	args = append([]string{matlabPath}, args...)

	process, err := os.StartProcess(matlabPath, args, attr)
	if err != nil {
		return nil, fmt.Errorf("error starting MATLAB: %w", err)
	}

	return process, nil
}
