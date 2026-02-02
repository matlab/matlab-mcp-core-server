// Copyright 2026 The MathWorks, Inc.

package server

import (
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolProviderResources[Dependencies any] struct{}

type Tool interface {
	toInternal(loggerFactory loggerFactory) internaltools.Tool
}

type toolArray []Tool

func (t toolArray) toInternal(loggerFactoryInstance loggerFactory) []internaltools.Tool {
	internalTools := make([]internaltools.Tool, len(t))

	for i, tool := range t {
		internalTools[i] = tool.toInternal(loggerFactoryInstance)
	}

	return internalTools
}

type loggerFactory interface {
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

func newToolCallRequest() *tools.CallRequest {
	return &tools.CallRequest{}
}
