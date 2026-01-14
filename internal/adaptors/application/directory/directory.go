// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

const (
	defaultLogDirPattern = "matlab-mcp-core-server-"
	markerFileName       = ".matlab-mcp-core-server"
)

type Config interface {
	BaseDir() string
	ServerInstanceID() string
}

type directory struct {
	baseDir string
	id      string

	osFacade OSLayer
}

func newDirectory(
	config Config,
	filenameFactory FilenameFactory,
	osFacade OSLayer,
) (*directory, messages.Error) {
	baseDir := config.BaseDir()

	if baseDir == "" {
		var err error
		if baseDir, err = osFacade.MkdirTemp("", defaultLogDirPattern); err != nil {
			return nil, messages.New_StartupErrors_FailedToCreateSubdirectory_Error(os.TempDir())
		}
	} else {
		if err := osFacade.MkdirAll(baseDir, 0o700); err != nil {
			return nil, messages.New_StartupErrors_FailedToCreateDirectory_Error(baseDir)
		}
	}

	serverInstanceID := config.ServerInstanceID()

	if serverInstanceID == "" {
		markerFilePath := filepath.Join(baseDir, markerFileName)
		_, id, err := filenameFactory.CreateFileWithUniqueSuffix(markerFilePath, "")
		if err != nil {
			return nil, messages.New_StartupErrors_FailedToCreateFile_Error(markerFilePath)
		}

		serverInstanceID = id
	}

	return &directory{
		baseDir: baseDir,
		id:      serverInstanceID,

		osFacade: osFacade,
	}, nil
}

func (d *directory) BaseDir() string {
	return d.baseDir
}

func (d *directory) ID() string {
	return d.id
}

func (d *directory) CreateSubDir(pattern string) (string, messages.Error) {
	if !strings.HasSuffix(pattern, "-") {
		pattern = fmt.Sprintf("%s-", pattern)
	}

	pattern = fmt.Sprintf("%s%s-", pattern, d.id)

	tempDir, err := d.osFacade.MkdirTemp(d.baseDir, pattern)
	if err != nil {
		return "", messages.New_StartupErrors_FailedToCreateSubdirectory_Error(d.baseDir)
	}

	return tempDir, nil
}

func (d *directory) RecordToLogger(logger entities.Logger) {
	logger.
		With("log-dir", d.baseDir).
		With("id", d.id).
		Info("Application directory state")
}
