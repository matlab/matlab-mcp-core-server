// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type StructuredToolThatLogsInput struct {
	Name string `json:"name"`
}

type StructuredToolThatLogsOutput struct {
	Response string `json:"response"`
}

func NewStructuredToolThatLogs() server.Tool {
	return server.NewToolWithStructuredContentOutput(
		tools.NewDefinition(
			"structured-tool-that-logs",
			"Structured Tool That Logs",
			"A structured tool that logs a message",
			tools.NewReadOnlyAnnotations(),
		),
		func(ctx context.Context, request tools.CallRequest, inputs StructuredToolThatLogsInput) (StructuredToolThatLogsOutput, i18n.Error) {
			logger := request.Logger()

			logger.Info("Logging from structured tool: " + inputs.Name)

			return StructuredToolThatLogsOutput{
				Response: "Hello " + inputs.Name,
			}, nil
		},
	)
}
