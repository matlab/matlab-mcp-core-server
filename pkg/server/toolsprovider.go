// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type ToolsProviderResources[Dependencies any] struct {
	logger       logger.Logger
	dependencies Dependencies
}

func newToolsProviderResources[Dependencies any](resources definition.ToolsProviderResources) ToolsProviderResources[Dependencies] {
	var dependencies Dependencies
	castDependencies, ok := resources.Dependencies.(Dependencies)
	if ok {
		dependencies = castDependencies
	}

	return ToolsProviderResources[Dependencies]{
		logger:       newLoggerAdaptor(resources.Logger),
		dependencies: dependencies,
	}
}

func (r ToolsProviderResources[Dependencies]) Logger() logger.Logger {
	return r.logger
}

func (r ToolsProviderResources[Dependencies]) Dependencies() Dependencies {
	return r.dependencies
}

type ToolsProvider[Dependencies any] func(toolsProviderResources ToolsProviderResources[Dependencies]) []Tool

func (p ToolsProvider[Dependencies]) toInternal() definition.ToolsProvider {
	return func(resources definition.ToolsProviderResources) []internaltools.Tool {
		if p == nil {
			return nil
		}

		tools := p(newToolsProviderResources[Dependencies](resources))

		return toolArray(tools).toInternal(
			resources.LoggerFactory,
			resources.Config,
			resources.MessageCatalog,
		)
	}
}
