// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
)

func (t *ToolWithUnstructuredContentOutput[ToolInput]) ToInternal(loggerFactoryInstance basetool.LoggerFactory, config internalconfig.GenericConfig, messageCatalog definition.MessageCatalog) basetool.ToolWithUnstructuredContentOutput[ToolInput] {
	return t.toInternal(loggerFactoryInstance, config, messageCatalog).(basetool.ToolWithUnstructuredContentOutput[ToolInput])
}
