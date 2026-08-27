// Copyright 2026 The MathWorks, Inc.

package tools

import (
	"context"

	internalconfig "github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	publictypes "github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type UnstructuredHandler[ToolInput any] func(ctx context.Context, request publictypes.ToolCallRequest, inputs ToolInput) (publictypes.RichContent, publictypes.Error)

type ToolWithUnstructuredContentOutput[ToolInput any] struct {
	publictypes.ToolSeal
	definition publictypes.ToolDefinition
	handler    UnstructuredHandler[ToolInput]
}

var _ ConvertibleTool = &ToolWithUnstructuredContentOutput[any]{}

func NewUnstructured[ToolInput any](definition publictypes.ToolDefinition, handler UnstructuredHandler[ToolInput]) *ToolWithUnstructuredContentOutput[ToolInput] {
	return &ToolWithUnstructuredContentOutput[ToolInput]{
		definition: definition,
		handler:    handler,
	}
}

func (t *ToolWithUnstructuredContentOutput[ToolInput]) ToInternal(
	toolCallRequestFactory ToolCallRequestFactory,
	loggerFactoryInstance basetool.LoggerFactory,
	telemetryFactory basetool.TelemetryFactory,
	config internalconfig.GenericConfig,
	messageCatalog definition.MessageCatalog,
) internaltools.Tool {
	annotations, ok := t.definition.Annotations.(ConvertibleAnnotation)
	if !ok {
		annotations = newDefaultAnnotation()
	}

	return basetool.NewToolWithUnstructuredContent(
		t.definition.Name,
		t.definition.Title,
		t.definition.Description,
		annotations,
		loggerFactoryInstance,
		telemetryFactory,
		adaptUnstructuredHandler(toolCallRequestFactory, config, messageCatalog, t.handler),
	)
}

func adaptUnstructuredHandler[ToolInput any](
	toolCallRequestFactory ToolCallRequestFactory,
	config internalconfig.GenericConfig,
	messageCatalog definition.MessageCatalog,
	handler UnstructuredHandler[ToolInput],
) basetool.HandlerWithUnstructuredContentOutput[ToolInput] {
	return func(ctx context.Context, logger entities.Logger, inputs ToolInput) (internaltools.RichContent, error) {
		callRequest := toolCallRequestFactory.New(
			logger,
			config,
			messageCatalog,
		)

		richContent, err := handler(ctx, callRequest, inputs)
		if err != nil {
			return internaltools.RichContent{}, err
		}

		return internaltools.RichContent{
			TextContent:  richContent.TextContent,
			ImageContent: nil,
		}, nil
	}
}
