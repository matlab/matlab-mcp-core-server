// Copyright 2025-2026 The MathWorks, Inc.

package client_test

import (
	"encoding/pem"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientFactory_NewClientForSelfSignedTLSServer_HappyPath(t *testing.T) {
	// Arrange
	expectedStatusCode := http.StatusOK

	server := newTestHTTPSServer(t)
	defer server.Close()

	serverCert := server.Certificate()
	certPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCert.Raw,
	})

	factory := newClientFactory()
	client, err := factory.NewClientForSelfSignedTLSServer(certPEMBytes)
	require.NoError(t, err)

	request, err := http.NewRequest("GET", "https://"+server.Listener.Addr().String(), nil)
	require.NoError(t, err)

	// Act
	response, err := client.Do(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedStatusCode, response.StatusCode)
	require.NoError(t, response.Body.Close())
}

func TestHTTPClientFactory_NewClientOverUDS_HappyPath(t *testing.T) {
	// Arrange
	expectedStatusCode := http.StatusOK

	server := newTestUDSServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatusCode)
	})

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve()
	}()
	defer func() {
		server.Close(t)
		require.ErrorIs(t, <-serveErr, http.ErrServerClosed)
	}()

	factory := newClientFactory()
	client := factory.NewClientOverUDS(server.SocketPath)

	request, err := http.NewRequest("GET", "http://uds/test", nil)
	require.NoError(t, err)

	// Act
	response, err := client.Do(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedStatusCode, response.StatusCode)
	require.NoError(t, response.Body.Close())
}

func newClientFactory() *httpclientfactory.HTTPClientFactory {
	application := wire.Initialize()
	return application.HTTPClientFactory
}

func newTestHTTPSServer(t *testing.T) *httptest.Server {
	t.Helper()

	expectedStatusCode := http.StatusOK

	return httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(expectedStatusCode)
	}))
}

func newTestUDSServer(t *testing.T, handler http.HandlerFunc) *testUDSServer {
	t.Helper()

	socketPath := filepath.Join(t.TempDir(), "test.sock")

	listener, err := net.Listen("unix", socketPath)
	require.NoError(t, err, "Failed to create unix socket listener")

	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}

	return &testUDSServer{
		listener:   listener,
		server:     server,
		SocketPath: socketPath,
	}
}

type testUDSServer struct {
	listener   net.Listener
	server     *http.Server
	SocketPath string
}

func (s *testUDSServer) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *testUDSServer) Close(t *testing.T) {
	t.Helper()

	err := s.server.Close()
	require.NoError(t, err, "Failed to close server")
}
