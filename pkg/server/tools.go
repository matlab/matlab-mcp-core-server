// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
)

type Tool interface {
	toInternal(
		loggerFactory basetool.LoggerFactory,
		config internalconfig.GenericConfig,
		messageCatalog definition.MessageCatalog,
	) internaltools.Tool
}

type toolArray []Tool

func (t toolArray) toInternal(
	loggerFactoryInstance basetool.LoggerFactory,
	config internalconfig.GenericConfig,
	messageCatalog definition.MessageCatalog,
) []internaltools.Tool {
	internalTools := make([]internaltools.Tool, len(t))

	for i, tool := range t {
		internalTools[i] = tool.toInternal(
			loggerFactoryInstance,
			config,
			messageCatalog,
		)
	}

	return internalTools
}
