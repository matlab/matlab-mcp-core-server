// Copyright 2025-2026 The MathWorks, Inc.

package client

import (
	httpclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
)

type OSLayer interface {
	Stat(name string) (osfacade.FileInfo, error)
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type HTTPClientFactory interface {
	NewClientOverUDS(socketPath string) httpclient.HttpClient
}

type Factory struct {
	osLayer           OSLayer
	loggerFactory     LoggerFactory
	httpClientFactory HTTPClientFactory
}

func NewFactory(
	osLayer OSLayer,
	loggerFactory LoggerFactory,
	httpClientFactory HTTPClientFactory,
) *Factory {
	return &Factory{
		osLayer:           osLayer,
		loggerFactory:     loggerFactory,
		httpClientFactory: httpClientFactory,
	}
}

func (f *Factory) New() transport.Client {
	return newClient(
		f.osLayer,
		f.httpClientFactory,
		f.loggerFactory,
	)
}
