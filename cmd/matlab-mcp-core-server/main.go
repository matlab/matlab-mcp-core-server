// Copyright 2025-2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/server"

	_ "embed"
)

//go:embed assets/instructions.txt
var instructions string

func main() {
	serverDefinition := server.Definition[any]{
		Name:         "matlab-mcp-core-server",
		Title:        "MATLAB MCP Core Server",
		Instructions: instructions,
	}
	serverInstance := server.New(serverDefinition)

	ctx := context.Background()
	exitCode := serverInstance.StartAndWaitForCompletion(ctx)

	os.Exit(exitCode)
}
