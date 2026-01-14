// Copyright 2025-2026 The MathWorks, Inc.

package config

import (
	"runtime/debug"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type Parser interface {
	Parse(args []string) (parser.SpecifiedArguments, messages.Error)
}

type OSLayer interface {
	Args() []string
	ReadBuildInfo() (info *debug.BuildInfo, ok bool)
}

type Config interface {
	Version() string
	HelpMode() bool
	VersionMode() bool
	WatchdogMode() bool
	BaseDir() string
	ServerInstanceID() string
	UseSingleMATLABSession() bool
	InitializeMATLABOnStartup() bool
	RecordToLogger(logger entities.Logger)
	LogLevel() entities.LogLevel
	PreferredLocalMATLABRoot() string
	PreferredMATLABStartingDirectory() string
}

type Factory struct {
	parser  Parser
	osLayer OSLayer

	initOnce       sync.Once
	initError      messages.Error
	configInstance *config
}

func NewFactory(parser Parser, osLayer OSLayer) *Factory {
	return &Factory{
		parser:  parser,
		osLayer: osLayer,
	}
}

func (f *Factory) Config() (Config, messages.Error) {
	f.initOnce.Do(func() {
		configInstance, err := newConfig(f.osLayer, f.parser)
		if err != nil {
			f.initError = err
			return
		}

		f.configInstance = configInstance
	})

	if f.initError != nil {
		return nil, f.initError
	}

	return f.configInstance, nil
}

type config struct {
	osLayer            OSLayer
	specifiedArguments parser.SpecifiedArguments
}

func newConfig(osLayer OSLayer, parser Parser) (*config, messages.Error) {
	specifiedArguments, err := parser.Parse(osLayer.Args()[1:])
	if err != nil {
		return nil, err
	}

	return &config{
		osLayer:            osLayer,
		specifiedArguments: specifiedArguments,
	}, nil
}

// Version returns the application version string from Go's build info.
func (c *config) Version() string {
	buildInfo, ok := c.osLayer.ReadBuildInfo()
	if !ok {
		return "(unknown)"
	}

	version := buildInfo.Main.Version
	if version == "" {
		version = "(devel)"
	}

	return buildInfo.Main.Path + " " + version
}

func (c *config) VersionMode() bool {
	return c.specifiedArguments.VersionMode
}

func (c *config) HelpMode() bool {
	return c.specifiedArguments.HelpMode
}

func (c *config) DisableTelemetry() bool {
	return c.specifiedArguments.DisableTelemetry
}

func (c *config) UseSingleMATLABSession() bool {
	return c.specifiedArguments.UseSingleMATLABSession
}

func (c *config) LogLevel() entities.LogLevel {
	return c.specifiedArguments.LogLevel
}

func (c *config) PreferredLocalMATLABRoot() string {
	return c.specifiedArguments.PreferredLocalMATLABRoot
}

func (c *config) PreferredMATLABStartingDirectory() string {
	return c.specifiedArguments.PreferredMATLABStartingDirectory
}

func (c *config) BaseDir() string {
	return c.specifiedArguments.BaseDirectory
}

func (c *config) WatchdogMode() bool {
	return c.specifiedArguments.WatchdogMode
}

func (c *config) ServerInstanceID() string {
	return c.specifiedArguments.ServerInstanceID
}

func (c *config) InitializeMATLABOnStartup() bool {
	return c.specifiedArguments.InitializeMATLABOnStartup
}

func (c *config) RecordToLogger(logger entities.Logger) {
	logger.
		With(flags.DisableTelemetry, c.specifiedArguments.DisableTelemetry).
		With(flags.UseSingleMATLABSession, c.specifiedArguments.UseSingleMATLABSession).
		With(flags.LogLevel, c.specifiedArguments.LogLevel).
		With(flags.PreferredLocalMATLABRoot, c.specifiedArguments.PreferredLocalMATLABRoot).
		With(flags.PreferredMATLABStartingDirectory, c.specifiedArguments.PreferredMATLABStartingDirectory).
		With(flags.InitializeMATLABOnStartup, c.specifiedArguments.InitializeMATLABOnStartup).
		Info("Configuration state")
}
