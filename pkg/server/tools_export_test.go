// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
)

type ToolArray []Tool

func (t ToolArray) ToInternal(loggerFactoryInstance basetool.LoggerFactory, config internalconfig.GenericConfig, messageCatalog definition.MessageCatalog) []internaltools.Tool {
	return toolArray(t).toInternal(loggerFactoryInstance, config, messageCatalog)
}
