// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type toolCallRequestAdaptor struct {
	logger logger.Logger
	config config.Config
}

func newToolCallRequestAdaptor(
	logger entities.Logger,
	internalConfig internalconfig.GenericConfig,
	messageCatalog definition.MessageCatalog,
) *toolCallRequestAdaptor {
	return &toolCallRequestAdaptor{
		logger: newLoggerAdaptor(logger),
		config: newConfigAdaptor(internalConfig, messageCatalog),
	}
}

func (a *toolCallRequestAdaptor) Logger() logger.Logger {
	return a.logger
}

func (a *toolCallRequestAdaptor) Config() config.Config {
	return a.config
}
