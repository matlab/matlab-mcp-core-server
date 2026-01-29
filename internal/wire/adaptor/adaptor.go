// Copyright 2026 The MathWorks, Inc.

package adaptor

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

type Application interface {
	ModeSelector() ModeSelector
	MessageCatalog() MessageCatalog
}

type ModeSelector interface {
	StartAndWaitForCompletion(ctx context.Context) error
}

type MessageCatalog interface {
	GetFromGeneralError(err error) (string, bool)
	Get(key messages.MessageKey) string
}

type adaptor struct {
	*wire.Application
}

func (w *adaptor) ModeSelector() ModeSelector {
	return w.Application.ModeSelector
}

func (w *adaptor) MessageCatalog() MessageCatalog {
	return w.Application.MessageCatalog
}
