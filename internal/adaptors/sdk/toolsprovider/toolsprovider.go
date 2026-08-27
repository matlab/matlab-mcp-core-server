// Copyright 2026 The MathWorks, Inc.

package toolsprovider

import (
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
	pkgtools "github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/tools"
)

type ToolsProvider[Dependencies any] = publictypes.ToolsProvider[Dependencies]

type ToolCallRequestFactory = pkgtools.ToolCallRequestFactory

type ResourcesFactory[Dependencies any] interface {
	New(internal definition.ToolsProviderResources) publictypes.ToolsProviderResources[Dependencies]
}

type Factory[Dependencies any] struct {
	resourcesFactory       ResourcesFactory[Dependencies]
	toolCallRequestFactory ToolCallRequestFactory
}

func NewFactory[Dependencies any](
	resourcesFactory ResourcesFactory[Dependencies],
	toolCallRequestFactory ToolCallRequestFactory,
) *Factory[Dependencies] {
	return &Factory[Dependencies]{
		resourcesFactory:       resourcesFactory,
		toolCallRequestFactory: toolCallRequestFactory,
	}
}

func (f *Factory[Dependencies]) New(provider ToolsProvider[Dependencies]) definition.ToolsProvider {
	return func(internalResources definition.ToolsProviderResources) []internaltools.Tool {
		if provider == nil {
			return nil
		}

		resources := f.resourcesFactory.New(internalResources)
		tools := provider(resources)

		internalTools := []internaltools.Tool{}
		for _, tool := range tools {
			convertible, ok := tool.(pkgtools.ConvertibleTool)
			if !ok {
				internalResources.Logger.Error("Tool does not implement ConvertibleTool, skipping")
				continue
			}

			internalTools = append(internalTools, convertible.ToInternal(
				f.toolCallRequestFactory,
				internalResources.LoggerFactory,
				internalResources.TelemetryFactory,
				internalResources.Config,
				internalResources.MessageCatalog,
			))
		}

		return internalTools
	}
}
