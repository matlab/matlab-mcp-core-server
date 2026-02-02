// Copyright 2026 The MathWorks, Inc.

package tools

import (
	internalannotations "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type annotations interface {
	ToToolAnnotations() *mcp.ToolAnnotations
	lock() // We want to control the annotation combination for now
}

type ReadOnlyAnnotation struct{}

func NewReadOnlyAnnotations() annotations {
	return ReadOnlyAnnotation{}
}

func (a ReadOnlyAnnotation) ToToolAnnotations() *mcp.ToolAnnotations {
	return internalannotations.NewReadOnlyAnnotations().ToToolAnnotations()
}

func (a ReadOnlyAnnotation) lock() {}

func NewDefaultAnnotation() annotations {
	return NewReadOnlyAnnotations()
}
