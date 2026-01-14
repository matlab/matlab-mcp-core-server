// Copyright 2025-2026 The MathWorks, Inc.

package process

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

func NewProcess(
	osLayer OSLayer,
	logger entities.Logger,
	directory Directory,
	config Config,
) messages.Error {
	return newProcess(
		osLayer,
		logger,
		directory,
		config,
	)
}
