// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type GreetToolInput struct {
	Name string `json:"name"`
}

func NewGreetTool() server.Tool {
	return server.NewToolWithUnstructuredContentOutput(
		tools.NewDefinition(
			"greet",
			"Greet",
			"Greets a user by name",
			tools.NewReadOnlyAnnotations(),
		),
		func(ctx context.Context, request tools.CallRequest, inputs GreetToolInput) (tools.RichContent, i18n.Error) {
			cfg := request.Config()

			customParameter := CustomParameter()
			customParameterValue, err := config.Get(cfg, customParameter)
			if err != nil {
				return tools.RichContent{}, err
			}

			return tools.RichContent{
				TextContent: []string{"Hello " + inputs.Name + " " + customParameterValue},
			}, nil
		},
	)
}
