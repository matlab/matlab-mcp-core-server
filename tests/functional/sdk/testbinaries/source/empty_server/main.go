// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

func main() {
	serverDefinition := server.Definition[any]{
		Name:         "empty-server",
		Title:        "Empty Server",
		Instructions: "This is the Empty Server test binary",
	}
	serverInstance := server.New(serverDefinition)

	exitCode := serverInstance.StartAndWaitForCompletion(context.Background())

	os.Exit(exitCode)
}
