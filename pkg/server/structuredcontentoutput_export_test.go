// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
)

func (t *ToolWithStructuredContentOutput[ToolInput, ToolOutput]) ToInternal(loggerFactoryInstance basetool.LoggerFactory, config internalconfig.GenericConfig, messageCatalog definition.MessageCatalog) basetool.ToolWithStructuredContentOutput[ToolInput, ToolOutput] {
	return t.toInternal(loggerFactoryInstance, config, messageCatalog).(basetool.ToolWithStructuredContentOutput[ToolInput, ToolOutput])
}
