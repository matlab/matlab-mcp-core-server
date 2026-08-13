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

// NewIdempotentWriteAnnotations creates annotations for tools that write or
// overwrite local state but produce the same result when called repeatedly
// with the same arguments (e.g. running an analysis that overwrites a results
// directory).
func NewIdempotentWriteAnnotations() annotations {
	return annotations{
		readOnly:    false,
		destructive: true,
		idempotent:  true,
		openWorld:   false,
	}
}

// NewReadOnlyOpenWorldAnnotations creates annotations for tools that query
// external services without modifying any state (local or remote). Do not use
// this for tools that mutate remote state: the read-only hint tells hosts they
// may skip user confirmation, so misusing it lets writes through silently.
func NewReadOnlyOpenWorldAnnotations() annotations {
	return annotations{
		readOnly:    true,
		destructive: false,
		idempotent:  false,
		openWorld:   true,
	}
}
