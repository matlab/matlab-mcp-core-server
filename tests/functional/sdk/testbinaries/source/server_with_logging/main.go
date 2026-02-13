// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

func main() {
	serverDefinition := server.Definition[Dependencies]{
		Name:         "server-with-logging",
		Title:        "Server With Logging",
		Instructions: "This is a test server for validating logging behaviour.",

		DependenciesProvider: func(dependenciesProviderResources server.DependenciesProviderResources) (Dependencies, i18n.Error) {
			return DependenciesProvider(dependenciesProviderResources)
		},
		ToolsProvider: func(toolsProviderResources server.ToolsProviderResources[Dependencies]) []server.Tool {
			return ToolsProvider(toolsProviderResources)
		},
	}
	serverInstance := server.New(serverDefinition)

	exitCode := serverInstance.StartAndWaitForCompletion(context.Background())

	os.Exit(exitCode)
}
