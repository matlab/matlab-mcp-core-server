// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) SetServerTransport(serverTransport mcp.Transport) {
	s.serverTransport = serverTransport
}
