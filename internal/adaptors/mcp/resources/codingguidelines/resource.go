// Copyright 2025-2026 The MathWorks, Inc.

package codingguidelines

import (
	"context"
	_ "embed"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/baseresource"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

//go:embed assets/codingguidelines.md
var codingGuidelines string

type Resource struct {
	*baseresource.Resource
}

func New(loggerFactory baseresource.LoggerFactory) *Resource {
	return &Resource{
		Resource: baseresource.New(
			name,
			title,
			description,
			mimeType,
			estimatedSize,
			uri,
			loggerFactory,
			Handler(),
		),
	}
}

func Handler() baseresource.ResourceHandler {
	return func(_ context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		logger.Info("Returning MATLAB coding guidelines resource")

		return &baseresource.ReadResourceResult{
			Contents: []baseresource.ResourceContents{
				{
					MIMEType: mimeType,
					Text:     codingGuidelines,
				},
			},
		}, nil
	}
}
