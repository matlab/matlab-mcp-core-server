// Copyright 2025 The MathWorks, Inc.

package matlablocator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/osfacade"
)

const mcpMATLABPathEnvKey = "MCP_MATLAB_PATH"

// Environment defines the interface for environment variable operations
type Environment interface {
	Getenv(key string) string
}

// FileSystem defines the interface for file system operations
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	EvalSymlinks(path string) (string, error)
}

type Locator struct {
	Env Environment
	FS  FileSystem
}

func GetPath() (string, error) {
	return (&Locator{
		Env: osfacade.RealEnvironment{},
		FS:  filefacade.RealFileSystem{},
	}).GetPath()
}

func (l *Locator) GetPath() (string, error) {
	matlabPath := l.Env.Getenv(mcpMATLABPathEnvKey)
	if matlabPath == "" {
		return "", fmt.Errorf("%v environment variable is empty", mcpMATLABPathEnvKey)
	}

	if !filepath.IsAbs(matlabPath) {
		return "", fmt.Errorf("%v must be an absolute path, got %q", mcpMATLABPathEnvKey, matlabPath)
	}

	if _, err := l.FS.Stat(matlabPath); err != nil {
		return "", fmt.Errorf("MATLAB path does not exist: %s", matlabPath)
	}

	// Resolve symlinks and clean the path for a canonical absolute path.
	if resolved, err := l.FS.EvalSymlinks(matlabPath); err == nil {
		matlabPath = resolved
	}

	matlabPath = filepath.Clean(matlabPath)

	return matlabPath, nil
}
