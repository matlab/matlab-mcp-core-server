// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
)

func NewToolCallRequestAdaptor(
	logger entities.Logger,
	internalConfig internalconfig.GenericConfig,
	messageCatalog definition.MessageCatalog,
) tools.CallRequest {
	return newToolCallRequestAdaptor(logger, internalConfig, messageCatalog)
}
