// Copyright 2025-2026 The MathWorks, Inc.

package httpserverfactory_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpserverfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
	httpserverfactorymocks "github.com/matlab/matlab-mcp-core-server/mocks/utils/httpserverfactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &httpserverfactorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	factory := httpserverfactory.New(mockOSLayer)

	// Assert
	assert.NotNil(t, factory)
}

func TestHTTPServerFactory_NewServerOverUDS_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &httpserverfactorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	handlers := map[string]http.HandlerFunc{}

	factory := httpserverfactory.New(mockOSLayer)

	// Act
	server, err := factory.NewServerOverUDS(handlers)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestUDSServer_Serve_Shutdown_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &httpserverfactorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	handlers := map[string]http.HandlerFunc{}

	mockOSLayer.EXPECT().
		RemoveAll(socketPath).
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(socketPath).
		RunAndReturn(func(name string) error {
			return os.RemoveAll(name)
		}).
		Once()

	factory := httpserverfactory.New(mockOSLayer)

	server, err := factory.NewServerOverUDS(handlers)
	require.NoError(t, err)

	// Act
	errC := make(chan error)
	go func() {
		errC <- server.Serve(socketPath)
	}()

	select {
	case <-errC:
		t.Fatal("Serve should be blocking")
	case <-time.After(10 * time.Millisecond):
		// Normal behaviour
	}

	err = server.Shutdown(t.Context())

	// Assert
	require.NoError(t, err)
	require.NoError(t, <-errC)
}

func TestUDSServer_Serve_PathTooLong(t *testing.T) {
	// Arrange
	mockOSLayer := &httpserverfactorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	socketPath := filepath.Join("tmp", strings.Repeat("a", 200)+".sock")
	handlers := map[string]http.HandlerFunc{}

	factory := httpserverfactory.New(mockOSLayer)
	server, err := factory.NewServerOverUDS(handlers)
	require.NoError(t, err)

	// Act
	err = server.Serve(socketPath)

	// Assert
	require.ErrorIs(t, err, socket.ErrSocketPathTooLong)
}

func TestUDSServer_Serve_RemoveAllError(t *testing.T) {
	// Arrange
	mockOSLayer := &httpserverfactorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	socketPath := filepath.Join("tmp", "test.sock")
	expectedError := assert.AnError
	handlers := map[string]http.HandlerFunc{}

	mockOSLayer.EXPECT().
		RemoveAll(socketPath).
		Return(expectedError).
		Once()

	factory := httpserverfactory.New(mockOSLayer)
	server, err := factory.NewServerOverUDS(handlers)
	require.NoError(t, err)

	// Act
	err = server.Serve(socketPath)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedError)
}
