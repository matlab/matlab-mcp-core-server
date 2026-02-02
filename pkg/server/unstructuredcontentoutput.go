// Copyright 2026 The MathWorks, Inc.

package server

import (
	"context"

	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type ToolWithUnstructuredContentOutput[ToolInput any] struct {
	definition tools.Definition
	handler    HandlerForToolWithUnstructuredContentOutput[ToolInput]
}

type HandlerForToolWithUnstructuredContentOutput[ToolInput any] func(ctx context.Context, request *tools.CallRequest, inputs ToolInput) (tools.RichContent, i18n.Error)

func NewToolWithUnstructuredContentOutput[ToolInput any](definition tools.Definition, handler HandlerForToolWithUnstructuredContentOutput[ToolInput]) *ToolWithUnstructuredContentOutput[ToolInput] {
	return &ToolWithUnstructuredContentOutput[ToolInput]{
		definition: definition,
		handler:    handler,
	}
}

func (t *ToolWithUnstructuredContentOutput[ToolInput]) toInternal(loggerFactoryInstance loggerFactory) internaltools.Tool {
	return basetool.NewToolWithUnstructuredContent(
		t.definition.Name,
		t.definition.Title,
		t.definition.Description,
		t.definition.Annotations,
		loggerFactoryInstance,
		adaptorForHandlerForToolWithUnstructuredContentOutput(t.handler),
	)
}

func adaptorForHandlerForToolWithUnstructuredContentOutput[ToolInput any](handler HandlerForToolWithUnstructuredContentOutput[ToolInput]) basetool.HandlerWithUnstructuredContentOutput[ToolInput] {
	return func(ctx context.Context, logger entities.Logger, inputs ToolInput) (internaltools.RichContent, error) {
		richContent, err := handler(ctx, newToolCallRequest(), inputs)
		if err != nil {
			return internaltools.RichContent{}, err
		}

		return internaltools.RichContent{
			TextContent:  richContent.TextContent,
			ImageContent: nil,
		}, nil
	}
}
