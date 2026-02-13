// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type GreetStructuredToolInput struct {
	Name string `json:"name"`
}

type GreetStructuredToolOutput struct {
	Response       string `json:"response"`
	ParameterValue string `json:"configValue"`
}

func NewGreetStructuredTool() server.Tool {
	return server.NewToolWithStructuredContentOutput(
		tools.NewDefinition(
			"greet-structured",
			"Greet (Structured Content Output)",
			"Greets a user by name (Structured Content Output)",
			tools.NewReadOnlyAnnotations(),
		),
		func(ctx context.Context, request tools.CallRequest, inputs GreetStructuredToolInput) (GreetStructuredToolOutput, i18n.Error) {
			cfg := request.Config()

			customParameter := CustomParameter()
			customParameterValue, err := config.Get(cfg, customParameter)
			if err != nil {
				return GreetStructuredToolOutput{}, err
			}

			return GreetStructuredToolOutput{
				Response:       "Hello " + inputs.Name,
				ParameterValue: customParameterValue,
			}, nil
		},
	)
}
