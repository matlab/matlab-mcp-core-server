// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"net/http"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpserverfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
)

type HTTPServerFactory interface {
	NewServerOverUDS(handlers map[string]http.HandlerFunc) (httpserverfactory.HttpServer, error)
}

type HandlerFactory interface {
	Handler() (handler.Handler, error)
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type Factory struct {
	httpServerFactory HTTPServerFactory
	loggerFactory     LoggerFactory
	handlerFactory    HandlerFactory
}

func NewFactory(
	httpServerFactory HTTPServerFactory,
	loggerFactory LoggerFactory,
	handlerFactory HandlerFactory,
) *Factory {
	return &Factory{
		httpServerFactory: httpServerFactory,
		loggerFactory:     loggerFactory,
		handlerFactory:    handlerFactory,
	}
}

func (f *Factory) New() (transport.Server, error) {
	logger, messagesErr := f.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return nil, messagesErr
	}

	handler, err := f.handlerFactory.Handler()
	if err != nil {
		return nil, err
	}

	return newServer(
		f.httpServerFactory,
		logger,
		handler,
	)
}
