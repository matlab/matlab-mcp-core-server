// Copyright 2026 The MathWorks, Inc.

package definition

import (
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type ToolsProviderResources struct {
	Logger         entities.Logger
	Config         config.GenericConfig
	MessageCatalog MessageCatalog

	Dependencies any

	// We're forced to do this, because we can't deps inject in basetool, due to generics in golang MCP SDK
	LoggerFactory    basetool.LoggerFactory
	TelemetryFactory basetool.TelemetryFactory
}

type ToolsProvider func(resources ToolsProviderResources) []tools.Tool

func NewToolsProviderResources(
	logger entities.Logger,
	config config.GenericConfig,
	messageCatalog MessageCatalog,
	dependencies any,
	loggerFactory basetool.LoggerFactory,
	telemetryFactory basetool.TelemetryFactory,
) ToolsProviderResources {
	return ToolsProviderResources{
		Logger:         logger,
		Config:         config,
		MessageCatalog: messageCatalog,

		Dependencies: dependencies,

		LoggerFactory:    loggerFactory,
		TelemetryFactory: telemetryFactory,
	}
}
