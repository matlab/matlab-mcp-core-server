// Copyright 2025-2026 The MathWorks, Inc.

package baseresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const UnexpectedErrorPrefix = "unexpected error occurred: "

type LoggerFactory interface {
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

type ResourceContents struct {
	MIMEType string
	Text     string
}

type ReadResourceResult struct {
	Contents []ResourceContents
}

type ResourceHandler func(ctx context.Context, logger entities.Logger) (*ReadResourceResult, error)

func New(
	name string,
	title string,
	description string,
	mimeType string,
	size int64,
	uri string,
	loggerFactory LoggerFactory,
	handler ResourceHandler,
) *Resource {
	return &Resource{
		name:          name,
		title:         title,
		description:   description,
		mimeType:      mimeType,
		size:          size,
		uri:           uri,
		loggerFactory: loggerFactory,
		handler:       handler,
	}
}

type Resource struct {
	name          string
	title         string
	description   string
	mimeType      string
	size          int64
	uri           string
	loggerFactory LoggerFactory
	handler       ResourceHandler
}

func (r *Resource) AddToServer(server resources.Server) error {
	if err := validateMIMEType(r.mimeType); err != nil {
		return err
	}

	server.AddResource(
		&mcp.Resource{
			Name:        r.name,
			Title:       r.title,
			Description: r.description,
			MIMEType:    r.mimeType,
			Size:        r.size,
			URI:         r.uri,
		},
		r.resourceHandler(),
	)

	return nil
}

func (r *Resource) resourceHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		logger, messagesErr := r.loggerFactory.NewMCPSessionLogger(req.Session)
		if messagesErr != nil {
			return nil, messagesErr
		}

		logger = logger.With("resource-name", r.name)
		logger.Debug("Handling resource request")
		defer logger.Debug("Handled resource request")

		if r.handler == nil {
			err := fmt.Errorf(UnexpectedErrorPrefix + "no resource handler available")
			logger.WithError(err).Warn("Resource handler is nil")
			return nil, err
		}

		result, err := r.handler(ctx, logger)
		if err != nil {
			logger.WithError(err).Warn("Resource handler returned an error")
			return nil, err
		}

		mcpContents := make([]*mcp.ResourceContents, len(result.Contents))
		for i, c := range result.Contents {
			mcpContents[i] = &mcp.ResourceContents{
				MIMEType: c.MIMEType,
				Text:     c.Text,
			}
		}

		return &mcp.ReadResourceResult{
			Contents: mcpContents,
		}, nil
	}
}

func (r *Resource) Name() string {
	return r.name
}

func (r *Resource) Title() string {
	return r.title
}

func (r *Resource) Description() string {
	return r.description
}

func (r *Resource) MimeType() string {
	return r.mimeType
}

func (r *Resource) Size() int64 {
	return r.size
}

func (r *Resource) URI() string {
	return r.uri
}

func validateMIMEType(mimeType string) error {
	if mimeType == "" {
		return fmt.Errorf("invalid MIME type: empty string")
	}

	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid MIME type %q: must be in format type/subtype", mimeType)
	}

	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid MIME type %q: type and subtype cannot be empty", mimeType)
	}

	return nil
}
