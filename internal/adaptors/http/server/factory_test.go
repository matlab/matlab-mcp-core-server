// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/server"
	servermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &servermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	factory := server.NewFactory(mockOSLayer)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_NewServerOverUDS_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &servermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	handlers := map[string]http.HandlerFunc{}

	factory := server.NewFactory(mockOSLayer)

	// Act
	httpServer, err := factory.NewServerOverUDS(handlers)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, httpServer)
}
