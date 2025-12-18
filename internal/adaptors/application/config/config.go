// Copyright 2025 The MathWorks, Inc.

package config

import (
	"runtime/debug"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/spf13/pflag"
)

type OSLayer interface {
	Args() []string
	ReadBuildInfo() (info *debug.BuildInfo, ok bool)
}

type Config struct {
	osLayer OSLayer

	versionMode                      bool
	disableTelemetry                 bool
	useSingleMATLABSession           bool
	logLevel                         entities.LogLevel
	preferredLocalMATLABRoot         string
	preferredMATLABStartingDirectory string
	baseDirectory                    string
	watchdogMode                     bool
	serverInstanceID                 string
	initializeMATLABOnStartup        bool
}

func New(
	osLayer OSLayer,
) (*Config, error) {
	flagSet := pflag.NewFlagSet(pflag.CommandLine.Name(), pflag.ExitOnError)
	err := setupFlags(flagSet)
	if err != nil {
		return nil, err
	}
	return createConfigWithFlagValues(osLayer, flagSet, osLayer.Args()[1:])
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
	return c.versionMode
}

func (c *Config) DisableTelemetry() bool {
	return c.disableTelemetry
}

func (c *Config) UseSingleMATLABSession() bool {
	return c.useSingleMATLABSession
}

func (c *Config) LogLevel() entities.LogLevel {
	return c.logLevel
}

func (c *Config) PreferredLocalMATLABRoot() string {
	return c.preferredLocalMATLABRoot
}

func (c *Config) PreferredMATLABStartingDirectory() string {
	return c.preferredMATLABStartingDirectory
}

func (c *Config) BaseDir() string {
	return c.baseDirectory
}

func (c *Config) WatchdogMode() bool {
	return c.watchdogMode
}

func (c *Config) ServerInstanceID() string {
	return c.serverInstanceID
}

func (c *Config) InitializeMATLABOnStartup() bool {
	return c.initializeMATLABOnStartup
}

func (c *Config) RecordToLogger(logger entities.Logger) {
	logger.
		With(flags.DisableTelemetry, c.disableTelemetry).
		With(flags.UseSingleMATLABSession, c.useSingleMATLABSession).
		With(flags.LogLevel, c.logLevel).
		With(flags.PreferredLocalMATLABRoot, c.preferredLocalMATLABRoot).
		With(flags.PreferredMATLABStartingDirectory, c.preferredMATLABStartingDirectory).
		Info("Configuration state")
}
