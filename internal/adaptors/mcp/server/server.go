// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AdditionalToolsProvider interface {
	Tools(loggerFactory definition.LoggerFactory) []tools.Tool
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
	NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error)
}

type LifecycleSignaler interface {
	AddShutdownFunction(shutdownFcn func() error)
}

type MCPSDKServerFactory interface {
	NewServer() (*mcp.Server, messages.Error)
}

type MCPServerConfigurator interface {
	GetToolsToAdd() ([]tools.Tool, error)
	GetResourcesToAdd() []resources.Resource
}

type Server struct {
	additionalToolsProvider AdditionalToolsProvider
	mcpSDKServerFactory     MCPSDKServerFactory
	loggerFactory           LoggerFactory
	lifecycleSignaler       LifecycleSignaler
	configurator            MCPServerConfigurator
	serverTransport         mcp.Transport
}

func New(
	additionalToolsProvider AdditionalToolsProvider,
	mcpSDKServerfactory MCPSDKServerFactory,
	loggerFactory LoggerFactory,
	lifecycleSignaler LifecycleSignaler,
	configurator MCPServerConfigurator,
) *Server {
	return &Server{
		additionalToolsProvider: additionalToolsProvider,
		mcpSDKServerFactory:     mcpSDKServerfactory,
		loggerFactory:           loggerFactory,
		lifecycleSignaler:       lifecycleSignaler,
		configurator:            configurator,
		serverTransport:         &mcp.StdioTransport{},
	}
}

func (s *Server) Run() error {
	logger, messagesErr := s.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	mcpServer, messagesErr := s.mcpSDKServerFactory.NewServer()
	if messagesErr != nil {
		return messagesErr
	}

	toolsToAdd, err := s.configurator.GetToolsToAdd()
	if err != nil {
		return err
	}

	for _, tool := range toolsToAdd {
		if err := tool.AddToServer(mcpServer); err != nil {
			return err
		}
	}
	logger.With("count", len(toolsToAdd)).Info("Added tools to MCP SDK server")

	additionalToolsToAdd := s.additionalToolsProvider.Tools(s.loggerFactory)

	for _, tool := range additionalToolsToAdd {
		if err := tool.AddToServer(mcpServer); err != nil {
			return err
		}
	}
	logger.With("count", len(additionalToolsToAdd)).Info("Added additional tools to MCP SDK server")

	resourcesToAdd := s.configurator.GetResourcesToAdd()
	for _, resource := range resourcesToAdd {
		if err := resource.AddToServer(mcpServer); err != nil {
			return err
		}
	}
	logger.With("count", len(resourcesToAdd)).Info("Added resources to MCP SDK server")

	logger.Debug("Starting MCP server")

	ctx, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	serverShutdownC := make(chan struct{})
	defer close(serverShutdownC)

	serverErrC := make(chan error)
	go func() {
		serverErrC <- mcpServer.Run(ctx, s.serverTransport)
	}()
	logger.Debug("Started MCP server")

	s.lifecycleSignaler.AddShutdownFunction(func() error {
		logger.Debug("Stopping MCP server")
		stopServer()
		<-serverShutdownC
		logger.Debug("Stopped MCP server")
		return nil
	})

	if err := <-serverErrC; err != nil && err != context.Canceled {
		logger.WithError(err).Error("MCP server run method returned an unexpected error")
		return err
	}

	return nil
}
