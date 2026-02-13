// Copyright 2026 The MathWorks, Inc.

package tools

import (
	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type Definition struct {
	Name        string
	Title       string
	Description string
	Annotations annotations
}

type CallRequest interface {
	Logger() logger.Logger
	Config() config.Config
}

type RichContent struct {
	TextContent []string
}

func NewDefinition(name, title, description string, annotations annotations) Definition {
	if annotations == nil {
		annotations = NewDefaultAnnotation()
	}

	return Definition{
		Name:        name,
		Title:       title,
		Description: description,
		Annotations: annotations,
	}
}
