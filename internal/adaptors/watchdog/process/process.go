// Copyright 2025-2026 The MathWorks, Inc.

package process

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type Config interface {
	LogLevel() entities.LogLevel
}

type Directory interface {
	BaseDir() string
	ID() string
}

func newProcess(
	osLayer OSLayer,
	logger entities.Logger,
	directory Directory,
	config Config,
) messages.Error {
	programPath, err := osLayer.Executable()
	if err != nil {
		logger.WithError(err).Error("Failed to get executable path")
		return messages.New_StartupErrors_FailedToGetExecutablePath_Error()
	}

	cmd := osLayer.Command(programPath,
		"--"+flags.WatchdogMode,
		"--"+flags.BaseDir, directory.BaseDir(),
		"--"+flags.ServerInstanceID, directory.ID(),
		"--"+flags.LogLevel, string(config.LogLevel()),
	)

	cmd.SetSysProcAttr(getSysProcAttrForDetachingAProcess())

	if err := cmd.Start(); err != nil {
		logger.WithError(err).Error("Failed to start watchdog process")
		return messages.New_StartupErrors_FailedToStartWatchdogProcess_Error()
	}

	return nil
}
