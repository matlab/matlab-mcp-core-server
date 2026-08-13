// Copyright 2026 The MathWorks, Inc.

package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	internalannotations "github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
)

type ConvertibleAnnotation interface {
	publictypes.Annotations
	ToToolAnnotations() *mcp.ToolAnnotations
}

type ReadOnlyAnnotation struct {
	publictypes.AnnotationSeal
}

var _ ConvertibleAnnotation = ReadOnlyAnnotation{}

func NewReadOnlyAnnotations() ReadOnlyAnnotation {
	return ReadOnlyAnnotation{}
}

func (a ReadOnlyAnnotation) ToToolAnnotations() *mcp.ToolAnnotations {
	return internalannotations.NewReadOnlyAnnotations().ToToolAnnotations()
}

type DestructiveAnnotation struct {
	publictypes.AnnotationSeal
}

var _ ConvertibleAnnotation = DestructiveAnnotation{}

func NewDestructiveAnnotations() DestructiveAnnotation {
	return DestructiveAnnotation{}
}

func (a DestructiveAnnotation) ToToolAnnotations() *mcp.ToolAnnotations {
	return internalannotations.NewDestructiveAnnotations().ToToolAnnotations()
}

type IdempotentWriteAnnotation struct {
	publictypes.AnnotationSeal
}

var _ ConvertibleAnnotation = IdempotentWriteAnnotation{}

func NewIdempotentWriteAnnotations() IdempotentWriteAnnotation {
	return IdempotentWriteAnnotation{}
}

func (a IdempotentWriteAnnotation) ToToolAnnotations() *mcp.ToolAnnotations {
	return internalannotations.NewIdempotentWriteAnnotations().ToToolAnnotations()
}

type ReadOnlyOpenWorldAnnotation struct {
	publictypes.AnnotationSeal
}

var _ ConvertibleAnnotation = ReadOnlyOpenWorldAnnotation{}

func NewReadOnlyOpenWorldAnnotations() ReadOnlyOpenWorldAnnotation {
	return ReadOnlyOpenWorldAnnotation{}
}

func (a ReadOnlyOpenWorldAnnotation) ToToolAnnotations() *mcp.ToolAnnotations {
	return internalannotations.NewReadOnlyOpenWorldAnnotations().ToToolAnnotations()
}

func newDefaultAnnotation() ReadOnlyAnnotation {
	return NewReadOnlyAnnotations()
}
