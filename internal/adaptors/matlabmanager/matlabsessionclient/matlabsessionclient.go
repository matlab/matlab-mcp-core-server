// Copyright 2025-2026 The MathWorks, Inc.

package matlabsessionclient

import (
	httpclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type HttpClientFactory interface {
	NewClientForSelfSignedTLSServer(certificatePEM []byte) (httpclient.HttpClient, error)
}

type Factory struct {
	httpClientFactory HttpClientFactory
}

func NewFactory(
	httpClientFactory HttpClientFactory,
) *Factory {
	return &Factory{
		httpClientFactory: httpClientFactory,
	}
}

func (f *Factory) New(endpoint embeddedconnector.ConnectionDetails) (entities.MATLABSessionClient, error) {
	return embeddedconnector.NewClient(endpoint, f.httpClientFactory)
}
