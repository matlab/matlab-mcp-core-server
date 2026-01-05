// Copyright 2025-2026 The MathWorks, Inc.

package annotations

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// annotations represents tool safety classification metadata.
// All fields are required and use plain bool types to ensure complete specification.
// This design insulates the codebase from MCP SDK's optional field semantics.
// The type is unexported to enforce construction via factory functions only.
type annotations struct {
	readOnly    bool
	destructive bool
	idempotent  bool
	openWorld   bool
}

// ToToolAnnotations converts to the MCP SDK protocol type.
// Handles the SDK's use of *bool for certain fields.
func (a annotations) ToToolAnnotations() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    a.readOnly,
		DestructiveHint: &a.destructive,
		IdempotentHint:  a.idempotent,
		OpenWorldHint:   &a.openWorld,
	}
}

// NewReadOnlyAnnotations creates annotations for tools that perform inspection
// or query operations without modifying state or executing user code.
func NewReadOnlyAnnotations() annotations {
	return annotations{
		readOnly:    true,
		destructive: false,
		idempotent:  false,
		openWorld:   false,
	}
}

// NewDestructiveAnnotations creates annotations for tools that execute code,
// modify state, or interact with external services.
func NewDestructiveAnnotations() annotations {
	return annotations{
		readOnly:    false,
		destructive: true,
		idempotent:  false,
		openWorld:   true,
	}
}
