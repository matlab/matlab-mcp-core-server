// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	resourcemocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/resources"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server"
	toolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	mockConfigurator := &mocks.MockMCPServerConfigurator{}

	// Act
	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Assert
	assert.NotNil(t, svr, "Server should not be nil")
}

func TestServer_Run_HappyPath(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockResource := &resourcemocks.MockResource{}
	defer mockResource.AssertExpectations(t)

	mockFirstTool := &toolsmocks.MockTool{}
	defer mockFirstTool.AssertExpectations(t)

	mockSecondTool := &toolsmocks.MockTool{}
	defer mockSecondTool.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedMCPServer := &mcp.Server{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockMCPSDKServerFactory.EXPECT().
		NewServer(server.Name(), server.Instructions()).
		Return(expectedMCPServer, nil).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return([]tools.Tool{mockFirstTool, mockSecondTool}, nil).
		Once()

	mockConfigurator.EXPECT().
		GetResourcesToAdd().
		Return([]resources.Resource{mockResource}).
		Once()

	mockFirstTool.EXPECT().
		AddToServer(expectedMCPServer).
		Return(nil).
		Once()

	mockSecondTool.EXPECT().
		AddToServer(expectedMCPServer).
		Return(nil).
		Once()

	mockResource.EXPECT().
		AddToServer(expectedMCPServer).
		Return(nil).
		Once()

	capturedShutdownFuncC := make(chan func() error)
	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFuncC <- shutdownFcn
		}).
		Return().
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	_, serverTransport := mcp.NewInMemoryTransports()
	svr.SetServerTransport(serverTransport)

	errC := make(chan error)
	go func() {
		errC <- svr.Run()
	}()

	capturedShutdownFunc := <-capturedShutdownFuncC

	// Act
	err := capturedShutdownFunc()

	// Assert
	require.NoError(t, err, "Shutdown function should not return an error")
	serverErr := <-errC
	require.NoError(t, serverErr, "Server run should exit without error after shutdown")
}

func TestServer_Run_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Act
	err := svr.Run()

	// Assert
	require.ErrorIs(t, err, expectedError, "Run should return the error from GetGlobalLogger")
}

func TestServer_Run_MCPSDKServerFactoryError(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := &messages.StartupErrors_BadFlag_Error{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockMCPSDKServerFactory.EXPECT().
		NewServer(server.Name(), server.Instructions()).
		Return(nil, expectedError).
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Act
	err := svr.Run()

	// Assert
	require.ErrorIs(t, err, expectedError, "Run should return the error from NewServer")
}

func TestServer_Run_ToolAddToServerReturnsError(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockTool := &toolsmocks.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := assert.AnError
	expectedMCPServer := &mcp.Server{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockMCPSDKServerFactory.EXPECT().
		NewServer(server.Name(), server.Instructions()).
		Return(expectedMCPServer, nil).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return([]tools.Tool{mockTool}, nil).
		Once()

	mockTool.EXPECT().
		AddToServer(expectedMCPServer).
		Return(expectedError).
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Act
	err := svr.Run()

	// Assert
	require.Error(t, err, "Run should return an error")
	assert.Equal(t, expectedError, err, "Error should match expected error")
}

func TestServer_Run_ResourceAddToServerReturnsError(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockResource := &resourcemocks.MockResource{}
	defer mockResource.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := assert.AnError
	expectedMCPServer := &mcp.Server{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockMCPSDKServerFactory.EXPECT().
		NewServer(server.Name(), server.Instructions()).
		Return(expectedMCPServer, nil).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return(nil, nil).
		Once()

	mockConfigurator.EXPECT().
		GetResourcesToAdd().
		Return([]resources.Resource{mockResource}).
		Once()

	mockResource.EXPECT().
		AddToServer(expectedMCPServer).
		Return(expectedError).
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Act
	err := svr.Run()

	// Assert
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestServer_Run_HandlesNoToolsOrResources(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedMCPServer := &mcp.Server{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockMCPSDKServerFactory.EXPECT().
		NewServer(server.Name(), server.Instructions()).
		Return(expectedMCPServer, nil).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return(nil, nil).
		Once()

	mockConfigurator.EXPECT().
		GetResourcesToAdd().
		Return(nil).
		Once()

	capturedShutdownFuncC := make(chan func() error)
	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFuncC <- shutdownFcn
		}).
		Return().
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	_, serverTransport := mcp.NewInMemoryTransports()
	svr.SetServerTransport(serverTransport)

	errC := make(chan error)
	go func() {
		errC <- svr.Run()
	}()

	capturedShutdownFunc := <-capturedShutdownFuncC

	// Act
	err := capturedShutdownFunc()

	// Assert
	require.NoError(t, err, "Shutdown function should not return an error")
	serverErr := <-errC
	require.NoError(t, serverErr, "Server run should exit without error after shutdown")
}

func TestServer_Run_GetToolsToAddError(t *testing.T) {
	// Arrange
	mockMCPSDKServerFactory := &mocks.MockMCPSDKServerFactory{}
	defer mockMCPSDKServerFactory.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedMCPServer := &mcp.Server{}
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockMCPSDKServerFactory.EXPECT().
		NewServer(server.Name(), server.Instructions()).
		Return(expectedMCPServer, nil).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return(nil, expectedError).
		Once()

	svr := server.New(mockMCPSDKServerFactory, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Act
	err := svr.Run()

	// Assert
	require.ErrorIs(t, err, expectedError, "Run should return the error from GetToolsToAdd")
}
