// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

type ToolInput struct {
	Name string `json:"name"`
}

type ToolOutput struct {
	Response string `json:"response"`
}

func NewToolWithUnstructuredContentOutput() server.Tool {
	return server.NewToolWithUnstructuredContentOutput(
		tools.NewDefinition("greet", "Greet", "Greets a user by name", tools.NewReadOnlyAnnotations()),
		func(ctx context.Context, request *tools.CallRequest, inputs ToolInput) (tools.RichContent, i18n.Error) {
			name := inputs.Name
			return tools.RichContent{
				TextContent: []string{"Hello " + name},
			}, nil
		},
	)
}

func NewToolWithStructureContentOutput() server.Tool {
	return server.NewToolWithStructuredContentOutput(
		tools.NewDefinition("greet-structured", "Greet (Structured Content Output)", "Greets a user by name (Structured Content Output)", tools.NewReadOnlyAnnotations()),
		func(ctx context.Context, request *tools.CallRequest, inputs ToolInput) (ToolOutput, i18n.Error) {
			name := inputs.Name
			return ToolOutput{
				Response: "Hello " + name,
			}, nil
		},
	)
}

func main() {
	serverDefinition := server.Definition[any]{
		Name:         "server-with-custom-tools",
		Title:        "Server With Custom Tools",
		Instructions: "This is a test server with custom tools",

		ToolsProvider: func(_ server.ToolProviderResources[any]) []server.Tool {
			return []server.Tool{
				NewToolWithUnstructuredContentOutput(),
				NewToolWithStructureContentOutput(),
			}
		},
	}
	serverInstance := server.New(serverDefinition)

	exitCode := serverInstance.StartAndWaitForCompletion(context.Background())

	os.Exit(exitCode)
}
