// Copyright 2026 The MathWorks, Inc.

package adaptor

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ModeSelector interface {
	StartAndWaitForCompletion(ctx context.Context) error
}

type MessageCatalog interface {
	GetFromGeneralError(err error) (string, bool)
	Get(key messages.MessageKey) string
}

type LoggerFactory interface {
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type adaptor struct {
	*wire.Application
}

func newAdaptor(application *wire.Application) *adaptor {
	return &adaptor{
		Application: application,
	}
}

func (w *adaptor) ModeSelector() ModeSelector {
	return w.Application.ModeSelector
}

func (w *adaptor) MessageCatalog() MessageCatalog {
	return w.Application.MessageCatalog
}

func (w *adaptor) LoggerFactory() LoggerFactory {
	return w.Application.LoggerFactory
}
