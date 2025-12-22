// Copyright 2025 The MathWorks, Inc.

package config

import (
	"runtime/debug"

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

type Config struct {
	osLayer            OSLayer
	specifiedArguments parser.SpecifiedArguments
}

func New(osLayer OSLayer, parser Parser) (*Config, messages.Error) {
	specifiedArguments, err := parser.Parse(osLayer.Args()[1:])

	if err != nil {
		return nil, err
	}
	return &Config{osLayer: osLayer, specifiedArguments: specifiedArguments}, nil
}

// Version returns the application version string from Go's build info.
func (c *Config) Version() string {
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

func (c *Config) VersionMode() bool {
	return c.specifiedArguments.VersionMode
}

func (c *Config) HelpMode() bool {
	return c.specifiedArguments.HelpMode
}

func (c *Config) DisableTelemetry() bool {
	return c.specifiedArguments.DisableTelemetry
}

func (c *Config) UseSingleMATLABSession() bool {
	return c.specifiedArguments.UseSingleMATLABSession
}

func (c *Config) LogLevel() entities.LogLevel {
	return c.specifiedArguments.LogLevel
}

func (c *Config) PreferredLocalMATLABRoot() string {
	return c.specifiedArguments.PreferredLocalMATLABRoot
}

func (c *Config) PreferredMATLABStartingDirectory() string {
	return c.specifiedArguments.PreferredMATLABStartingDirectory
}

func (c *Config) BaseDir() string {
	return c.specifiedArguments.BaseDirectory
}

func (c *Config) WatchdogMode() bool {
	return c.specifiedArguments.WatchdogMode
}

func (c *Config) ServerInstanceID() string {
	return c.specifiedArguments.ServerInstanceID
}

func (c *Config) InitializeMATLABOnStartup() bool {
	return c.specifiedArguments.InitializeMATLABOnStartup
}

func (c *Config) RecordToLogger(logger entities.Logger) {
	logger.
		With(flags.DisableTelemetry, c.specifiedArguments.DisableTelemetry).
		With(flags.UseSingleMATLABSession, c.specifiedArguments.UseSingleMATLABSession).
		With(flags.LogLevel, c.specifiedArguments.LogLevel).
		With(flags.PreferredLocalMATLABRoot, c.specifiedArguments.PreferredLocalMATLABRoot).
		With(flags.PreferredMATLABStartingDirectory, c.specifiedArguments.PreferredMATLABStartingDirectory).
		With(flags.InitializeMATLABOnStartup, c.specifiedArguments.InitializeMATLABOnStartup).
		Info("Configuration state")
}
