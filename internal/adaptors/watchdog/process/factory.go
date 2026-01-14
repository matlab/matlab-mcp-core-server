// Copyright 2025-2026 The MathWorks, Inc.

package process

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
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

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Factory struct {
	osLayer          OSLayer
	loggerFactory    LoggerFactory
	directoryFactory DirectoryFactory
	configFactory    ConfigFactory
}

func New(
	osLayer OSLayer,
	loggerFactory LoggerFactory,
	directoryFactory DirectoryFactory,
	configFactory ConfigFactory,
) *Factory {
	return &Factory{
		osLayer:          osLayer,
		loggerFactory:    loggerFactory,
		directoryFactory: directoryFactory,
		configFactory:    configFactory,
	}
}

func (p *Factory) StartNewProcess() messages.Error {
	logger, err := p.loggerFactory.GetGlobalLogger()
	if err != nil {
		return err
	}

	directory, err := p.directoryFactory.Directory()
	if err != nil {
		return err
	}

	config, err := p.configFactory.Config()
	if err != nil {
		return err
	}

	return newProcess(
		p.osLayer,
		logger,
		directory,
		config,
	)
}
