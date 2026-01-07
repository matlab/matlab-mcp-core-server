// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpserverfactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServerFactory_NewServerOverUDS_HappyPath(t *testing.T) {
	// Arrange
	factory := newServerFactory()

	testDataDir, err := os.MkdirTemp("", "mcp_test") //nolint:usetesting // We can't use t.TempDir() here, as it sometimes creates path that are too long for socket paths
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(testDataDir))
	}()

	socketPath := filepath.Join(testDataDir, "test.sock")

	expectedStatusCode := http.StatusOK
	expectedFirstBody := "first hello world"
	expectedSecondBody := "second hello world"

	handlers := map[string]http.HandlerFunc{
		"GET /first": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedFirstBody))
		},
		"POST /second": func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedSecondBody))
		},
	}

	server, err := factory.NewServerOverUDS(handlers)
	require.NoError(t, err)

	serverStopped := make(chan error, 1)
	go func() {
		serverStopped <- server.Serve(socketPath)
	}()
	defer func() {
		require.NoError(t, server.Shutdown(t.Context()))
	}()

	socketFileExists := make(chan error, 1)
	go func() {
		socketFileExists <- waitForSocketFile(socketPath)
	}()

	select {
	case err := <-serverStopped:
		t.Fatalf("Server stopped unexpectedly: %v", err)
	case err := <-socketFileExists:
		require.NoError(t, err)
	}

	client := newUDSClient(socketPath)

	// Act & Assert
	req, err := http.NewRequest("GET", "http://unix/first", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, expectedFirstBody, string(body))

	req, err = http.NewRequest("POST", "http://unix/second", nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, expectedStatusCode, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	assert.Equal(t, expectedSecondBody, string(body))
}

func newServerFactory() *httpserverfactory.HTTPServerFactory {
	osLayer := osfacade.New()
	return httpserverfactory.New(osLayer)
}

func newUDSClient(socketPath string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
}

func waitForSocketFile(socketPath string) error {
	timeout := time.After(1 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		_, err := os.Stat(socketPath)
		if err == nil {
			return nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		select {
		case <-timeout:
			return fmt.Errorf("Failed to wait for socket file: %v", err)
		case <-tick:
		}
	}
}
