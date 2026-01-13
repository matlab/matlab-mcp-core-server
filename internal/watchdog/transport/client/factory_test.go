// Copyright 2025-2026 The MathWorks, Inc.

package client_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/client"
	clientmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/client"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &clientmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	// Act
	factory := client.NewFactory(
		mockOSLayer,
		mockLoggerFactory,
		mockHTTPClientFactory,
	)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &clientmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &clientmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockHTTPClientFactory := &clientmocks.MockHTTPClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	factory := client.NewFactory(
		mockOSLayer,
		mockLoggerFactory,
		mockHTTPClientFactory,
	)

	// Act
	clientInstance := factory.New()

	// Assert
	assert.NotNil(t, clientInstance, "Client should not be nil")
}
