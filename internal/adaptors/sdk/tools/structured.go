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

type StructuredHandler[ToolInput, ToolOutput any] func(ctx context.Context, request publictypes.ToolCallRequest, inputs ToolInput) (ToolOutput, publictypes.Error)

type ToolWithStructuredContentOutput[ToolInput, ToolOutput any] struct {
	publictypes.ToolSeal
	definition publictypes.ToolDefinition
	handler    StructuredHandler[ToolInput, ToolOutput]
}

var _ ConvertibleTool = &ToolWithStructuredContentOutput[any, any]{}

func NewStructured[ToolInput, ToolOutput any](definition publictypes.ToolDefinition, handler StructuredHandler[ToolInput, ToolOutput]) *ToolWithStructuredContentOutput[ToolInput, ToolOutput] {
	return &ToolWithStructuredContentOutput[ToolInput, ToolOutput]{
		definition: definition,
		handler:    handler,
	}
}

func (t *ToolWithStructuredContentOutput[ToolInput, ToolOutput]) ToInternal(
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

	return basetool.NewToolWithStructuredContent(
		t.definition.Name,
		t.definition.Title,
		t.definition.Description,
		annotations,
		loggerFactoryInstance,
		telemetryFactory,
		adaptStructuredHandler(toolCallRequestFactory, config, messageCatalog, t.handler),
	)
}

func adaptStructuredHandler[ToolInput, ToolOutput any](
	toolCallRequestFactory ToolCallRequestFactory,
	config internalconfig.GenericConfig,
	messageCatalog definition.MessageCatalog,
	handler StructuredHandler[ToolInput, ToolOutput],
) basetool.HandlerWithStructuredContentOutput[ToolInput, ToolOutput] {
	return func(ctx context.Context, logger entities.Logger, inputs ToolInput) (ToolOutput, error) {
		callRequest := toolCallRequestFactory.New(
			logger,
			config,
			messageCatalog,
		)

		return handler(ctx, callRequest, inputs)
	}
}
