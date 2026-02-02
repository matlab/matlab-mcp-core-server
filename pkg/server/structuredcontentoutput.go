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

type ToolWithStructuredContentOutput[ToolInput, ToolOutput any] struct {
	definition tools.Definition
	handler    HandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput]
}

type HandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput any] func(ctx context.Context, request *tools.CallRequest, inputs ToolInput) (ToolOutput, i18n.Error)

func NewToolWithStructuredContentOutput[ToolInput, ToolOutput any](definition tools.Definition, handler HandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput]) *ToolWithStructuredContentOutput[ToolInput, ToolOutput] {
	return &ToolWithStructuredContentOutput[ToolInput, ToolOutput]{
		definition: definition,
		handler:    handler,
	}
}

func (t *ToolWithStructuredContentOutput[ToolInput, ToolOutput]) toInternal(loggerFactoryInstance loggerFactory) internaltools.Tool {
	return basetool.NewToolWithStructuredContent(
		t.definition.Name,
		t.definition.Title,
		t.definition.Description,
		t.definition.Annotations,
		loggerFactoryInstance,
		adaptorForHandlerForToolWithStructuredContentOutput(t.handler),
	)
}

func adaptorForHandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput any](handler HandlerForToolWithStructuredContentOutput[ToolInput, ToolOutput]) basetool.HandlerWithStructuredContentOutput[ToolInput, ToolOutput] {
	return func(ctx context.Context, logger entities.Logger, inputs ToolInput) (ToolOutput, error) {
		return handler(ctx, newToolCallRequest(), inputs)
	}
}
