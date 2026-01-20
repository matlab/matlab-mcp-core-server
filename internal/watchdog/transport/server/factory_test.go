// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server"
	httpservermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/server"
	servermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server"
	handlermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockLoggerFactory := &servermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockHandlerFactory := &servermocks.MockHandlerFactory{}
	defer mockHandlerFactory.AssertExpectations(t)

	// Act
	factory := server.NewFactory(
		mockHTTPServerFactory,
		mockLoggerFactory,
		mockHandlerFactory,
	)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockLoggerFactory := &servermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockHandlerFactory := &servermocks.MockHandlerFactory{}
	defer mockHandlerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockHandlerFactory.EXPECT().
		Handler().
		Return(mockHandler, nil).
		Once()

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Return(mockHTTPServer, nil).
		Once()

	factory := server.NewFactory(
		mockHTTPServerFactory,
		mockLoggerFactory,
		mockHandlerFactory,
	)

	// Act
	serverInstance, err := factory.New()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, serverInstance, "Server should not be nil")
}

func TestFactory_New_HTTPServerFactoryError(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockLoggerFactory := &servermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockHandlerFactory := &servermocks.MockHandlerFactory{}
	defer mockHandlerFactory.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockHandlerFactory.EXPECT().
		Handler().
		Return(mockHandler, nil).
		Once()

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Return(nil, expectedError).
		Once()

	factory := server.NewFactory(
		mockHTTPServerFactory,
		mockLoggerFactory,
		mockHandlerFactory,
	)

	// Act
	serverInstance, err := factory.New()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, serverInstance)
}

func TestFactory_New_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockLoggerFactory := &servermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockHandlerFactory := &servermocks.MockHandlerFactory{}
	defer mockHandlerFactory.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := server.NewFactory(
		mockHTTPServerFactory,
		mockLoggerFactory,
		mockHandlerFactory,
	)

	// Act
	serverInstance, err := factory.New()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, serverInstance)
}
