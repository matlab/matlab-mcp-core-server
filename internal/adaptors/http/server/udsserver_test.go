// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/server"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
	servermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUDSServer_Serve_Shutdown_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &servermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")

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

	udsServer := server.NewUDSServer(mockOSLayer, nil)

	// Act
	errC := make(chan error)
	go func() {
		errC <- udsServer.Serve(socketPath)
	}()

	select {
	case <-errC:
		t.Fatal("Serve should be blocking")
	case <-time.After(10 * time.Millisecond):
		// Normal behaviour
	}

	err := udsServer.Shutdown(t.Context())

	// Assert
	require.NoError(t, err)
	require.NoError(t, <-errC)
}

func TestUDSServer_Serve_PathTooLong(t *testing.T) {
	// Arrange
	mockOSLayer := &servermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	socketPath := filepath.Join("tmp", strings.Repeat("a", 200)+".sock")

	udsServer := server.NewUDSServer(mockOSLayer, nil)

	// Act
	err := udsServer.Serve(socketPath)

	// Assert
	require.ErrorIs(t, err, socket.ErrSocketPathTooLong)
}

func TestUDSServer_Serve_RemoveAllError(t *testing.T) {
	// Arrange
	mockOSLayer := &servermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	socketPath := filepath.Join("tmp", "test.sock")
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		RemoveAll(socketPath).
		Return(expectedError).
		Once()

	udsServer := server.NewUDSServer(mockOSLayer, nil)

	// Act
	err := udsServer.Serve(socketPath)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedError)
}
