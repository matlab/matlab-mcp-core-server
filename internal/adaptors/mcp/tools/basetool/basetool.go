// Copyright 2025-2026 The MathWorks, Inc.

package basetool

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const UnexpectedErrorPrefixForLLM = "unexpected error occurred: "

type AnnotationProvider interface {
	ToToolAnnotations() *mcp.ToolAnnotations
}

type LoggerFactory interface {
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

type TelemetryFactory interface {
	Telemetry() (telemetry.Telemetry, messages.Error)
}

type ToolAdder[ToolInput, ToolOutput any] interface {
	AddTool(server *mcp.Server, tool *mcp.Tool, handler mcp.ToolHandlerFor[ToolInput, ToolOutput])
}

type tool[ToolInput any, ToolOutput any] struct {
	name             string
	title            string
	description      string
	annotations      AnnotationProvider
	loggerFactory    LoggerFactory
	telemetryFactory TelemetryFactory
	toolAdder        ToolAdder[ToolInput, ToolOutput]
}

func (t tool[_, _]) Name() string {
	return t.name
}

func (t tool[_, _]) Title() string {
	return t.title
}

func (t tool[_, _]) Description() string {
	return t.description
}

func (t tool[_, _]) Annotations() AnnotationProvider {
	return t.annotations
}

func (_ tool[ToolInput, _]) GetInputSchema() (any, error) {
	return jsonschema.For[ToolInput](&jsonschema.ForOptions{})
}
