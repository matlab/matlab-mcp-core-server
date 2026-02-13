// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type ToolThatLogsInput struct {
	Name string `json:"name"`
}

func NewToolThatLogs() server.Tool {
	return server.NewToolWithUnstructuredContentOutput(
		tools.NewDefinition(
			"tool-that-logs",
			"Tool That Logs",
			"A tool that logs a message",
			tools.NewReadOnlyAnnotations(),
		),
		func(ctx context.Context, request tools.CallRequest, inputs ToolThatLogsInput) (tools.RichContent, i18n.Error) {
			logger := request.Logger()

			logger.Info("Logging from unstructured tool: " + inputs.Name)

			return tools.RichContent{
				TextContent: []string{"Hello " + inputs.Name},
			}, nil
		},
	)
}
