// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
)

func NewServer(
	httpServerFactory HTTPServerFactory,
	logger entities.Logger,
	h handler.Handler,
) (*Server, error) {
	return newServer(
		httpServerFactory,
		logger,
		h,
	)
}
