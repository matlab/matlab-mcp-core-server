// Copyright 2025-2026 The MathWorks, Inc.

package resources

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server interface {
	AddResource(resource *mcp.Resource, handler mcp.ResourceHandler)
}

type Resource interface {
	AddToServer(server Server) error
}
