// Copyright 2025-2026 The MathWorks, Inc.

package process

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type OSLayer interface {
	Command(name string, arg ...string) osfacade.Cmd
	Executable() (string, error)
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type Directory interface {
	BaseDir() string
	ID() string
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Process struct {
	osLayer OSLayer
	cmd     osfacade.Cmd
	logger  entities.Logger
}

func New(
	osLayer OSLayer,
	loggerFactory LoggerFactory,
	directory Directory,
	configFactory ConfigFactory,
) (*Process, error) {
	logger, err := loggerFactory.GetGlobalLogger()
	if err != nil {
		return nil, err
	}

	config, err := configFactory.Config()
	if err != nil {
		return nil, err
	}

	programPath, execErr := osLayer.Executable()
	if execErr != nil {
		logger.WithError(execErr).Error("Failed to get executable path")
		return nil, execErr
	}
	cmd := osLayer.Command(programPath,
		"--"+flags.WatchdogMode,
		"--"+flags.BaseDir, directory.BaseDir(),
		"--"+flags.ServerInstanceID, directory.ID(),
		"--"+flags.LogLevel, string(config.LogLevel()),
	)

	cmd.SetSysProcAttr(getSysProcAttrForDetachingAProcess())

	process := &Process{
		osLayer: osLayer,
		cmd:     cmd,
		logger:  logger,
	}

	return process, nil
}

func (p *Process) Start() error {
	if err := p.cmd.Start(); err != nil {
		p.logger.WithError(err).Error("Failed to start watchdog process")
		return err
	}

	return nil
}
