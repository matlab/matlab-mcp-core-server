// Copyright 2026 The MathWorks, Inc.

package tools

import (
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/tools"
)

type Annotations = publictypes.Annotations

// NewReadOnlyAnnotations returns annotations for tools that only inspect state
// without modifying anything (ReadOnly=true, Destructive=false, Idempotent=false, OpenWorld=false).
func NewReadOnlyAnnotations() Annotations {
	return tools.NewReadOnlyAnnotations()
}

// NewDestructiveAnnotations returns annotations for tools that execute code,
// modify state, or interact with external services
// (ReadOnly=false, Destructive=true, Idempotent=false, OpenWorld=true).
func NewDestructiveAnnotations() Annotations {
	return tools.NewDestructiveAnnotations()
}

// NewIdempotentWriteAnnotations returns annotations for tools that write or
// overwrite local state but produce the same result when called repeatedly
// with the same arguments
// (ReadOnly=false, Destructive=true, Idempotent=true, OpenWorld=false).
func NewIdempotentWriteAnnotations() Annotations {
	return tools.NewIdempotentWriteAnnotations()
}

// NewReadOnlyOpenWorldAnnotations returns annotations for tools that query
// external services without modifying any state (local or remote)
// (ReadOnly=true, Destructive=false, Idempotent=false, OpenWorld=true).
// Do not use this for tools that mutate remote state: the read-only hint tells
// hosts they may skip user confirmation, so misusing it lets writes through
// silently.
func NewReadOnlyOpenWorldAnnotations() Annotations {
	return tools.NewReadOnlyOpenWorldAnnotations()
}
