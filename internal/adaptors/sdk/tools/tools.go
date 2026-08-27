// Copyright 2026 The MathWorks, Inc.

package tools

import (
	internalconfig "github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	publictypes "github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type ToolCallRequestFactory interface {
	New(
		internalLogger entities.Logger,
		internalConfig internalconfig.GenericConfig,
		internalMessageCatalog definition.MessageCatalog,
	) publictypes.ToolCallRequest
}

type ConvertibleTool interface {
	publictypes.Tool
	ToInternal(
		toolCallRequestFactory ToolCallRequestFactory,
		loggerFactory basetool.LoggerFactory,
		telemetryFactory basetool.TelemetryFactory,
		config internalconfig.GenericConfig,
		messageCatalog definition.MessageCatalog,
	) internaltools.Tool
}
