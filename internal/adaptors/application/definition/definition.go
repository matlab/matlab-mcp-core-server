// Copyright 2026 The MathWorks, Inc.

package definition

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type LoggerFactory interface {
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

type toolsProvider func(loggerFactory LoggerFactory) []tools.Tool

type Definition struct {
	name         string
	title        string
	instructions string

	toolsProvider toolsProvider
}

func New(name, title, instructions string, toolsProvider toolsProvider) Definition {
	return Definition{
		name:         name,
		title:        title,
		instructions: instructions,

		toolsProvider: toolsProvider,
	}
}

func (d Definition) Name() string {
	return d.name
}

func (d Definition) Title() string {
	return d.title
}

func (d Definition) Instructions() string {
	return d.instructions
}

func (d Definition) Tools(loggerFactory LoggerFactory) []tools.Tool {
	if d.toolsProvider == nil {
		return nil
	}

	return d.toolsProvider(loggerFactory)
}
