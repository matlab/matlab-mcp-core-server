// Copyright 2026 The MathWorks, Inc.

package main

import (
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

type ToolsProviderResources interface { //nolint:iface // Same interface is happenstance
	Logger() logger.Logger
}

func ToolsProvider(resources ToolsProviderResources) []server.Tool {
	resources.Logger().Info("Creating Tools")

	return []server.Tool{
		NewToolThatLogs(),
		NewStructuredToolThatLogs(),
	}
}
