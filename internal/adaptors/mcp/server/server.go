// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type LifecycleSignaler interface {
	AddShutdownFunction(shutdownFcn func() error)
}

type MCPSDKServerFactory interface {
	NewServer(name string, instructions string) (*mcp.Server, messages.Error)
}

type MCPServerConfigurator interface {
	GetToolsToAdd() ([]tools.Tool, error)
	GetResourcesToAdd() []resources.Resource
}

type Server struct {
	mcpServer         *mcp.Server
	serverLogger      entities.Logger
	lifecycleSignaler LifecycleSignaler
	serverTransport   mcp.Transport
}

func New(
	mcpSDKServerfactory MCPSDKServerFactory,
	loggerFactory LoggerFactory,
	lifecycleSignaler LifecycleSignaler,
	configurator MCPServerConfigurator,
) (*Server, error) {
	logger, messagesErr := loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return nil, messagesErr
	}

	mcpserver, messagesErr := mcpSDKServerfactory.NewServer(name, instructions)
	if messagesErr != nil {
		return nil, messagesErr
	}

	toolsToAdd, err := configurator.GetToolsToAdd()
	if err != nil {
		return nil, err
	}

	for _, tool := range toolsToAdd {
		if err := tool.AddToServer(mcpserver); err != nil {
			return nil, err
		}
	}
	logger.With("count", len(toolsToAdd)).Info("Added tools to MCP SDK server")

	resourcesToAdd := configurator.GetResourcesToAdd()
	for _, resource := range resourcesToAdd {
		resource.AddToServer(mcpserver)
	}
	logger.With("count", len(resourcesToAdd)).Info("Added resources to MCP SDK server")

	return &Server{
		mcpServer:         mcpserver,
		serverLogger:      logger,
		lifecycleSignaler: lifecycleSignaler,
		serverTransport:   &mcp.StdioTransport{},
	}, nil
}

func (s *Server) Run() error {
	s.serverLogger.Debug("Starting MCP server")

	ctx, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	// This channel only closes when we exit the Run method
	// This ensures that the mcpServer connections have all been closed and resolved
	serverShutdownC := make(chan struct{})
	defer close(serverShutdownC)

	serverErrC := make(chan error)
	go func() {
		serverErrC <- s.mcpServer.Run(ctx, s.serverTransport)
	}()
	s.serverLogger.Debug("Started MCP server")

	s.lifecycleSignaler.AddShutdownFunction(func() error {
		s.serverLogger.Debug("Stopping MCP server")
		stopServer()
		<-serverShutdownC
		s.serverLogger.Debug("Stopped MCP server")
		return nil
	})

	if err := <-serverErrC; err != nil && err != context.Canceled {
		s.serverLogger.WithError(err).Error("MCP server run method returned an unexpected error")
		return err
	}

	return nil
}
