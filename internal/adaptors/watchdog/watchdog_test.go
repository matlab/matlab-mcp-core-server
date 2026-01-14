// Copyright 2025-2026 The MathWorks, Inc.

package watchdog_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	transportmessages "github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
	watchdogmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/watchdog"
	transportmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport"
	socketmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/socket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	// Act
	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Assert
	assert.NotNil(t, watchdogInstance, "Watchdog instance should not be nil")
}

func TestWatchdog_Start_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	expectedError := messages.AnError

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the GetGlobalLogger error")
}

func TestWatchdog_Start_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	require.NoError(t, err, "Start should not return an error")
}

func TestWatchdog_Start_SocketFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the watchdog process start error")
}

func TestWatchdog_Start_WatchdogProcessStartError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the watchdog process start error")
}

func TestWatchdog_Start_ClientConnectError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the client connect error")
}

func TestWatchdog_RegisterProcessPIDWithWatchdog_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"
	expectedPID := 12345

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	mockClient.EXPECT().
		SendProcessPID(expectedPID).
		Return(transportmessages.ProcessToKillResponse{}, nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Start the watchdog first
	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.RegisterProcessPIDWithWatchdog(expectedPID)

	// Assert
	require.NoError(t, err, "RegisterProcessPIDWithWatchdog should not return an error")
}

func TestWatchdog_RegisterProcessPIDWithWatchdog_WaitsIfNotStarted(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"
	expectedPID := 12345

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	mockClient.EXPECT().
		SendProcessPID(expectedPID).
		Return(transportmessages.ProcessToKillResponse{}, nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act & Assert
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.RegisterProcessPIDWithWatchdog(expectedPID)
	}()

	select {
	case <-errC:
		t.Fatal("RegisterProcessPIDWithWatchdog should block until Start is called")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: RegisterProcessPIDWithWatchdog blocks
	}

	// Start after we've tried to register a PID
	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	require.NoError(t, <-errC, "RegisterProcessPIDWithWatchdog should not return an error")
}

func TestWatchdog_RegisterProcessPIDWithWatchdog_SendProcessPIDError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"
	expectedPID := 12345
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	mockClient.EXPECT().
		SendProcessPID(expectedPID).
		Return(transportmessages.ProcessToKillResponse{}, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Start the watchdog first
	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.RegisterProcessPIDWithWatchdog(expectedPID)

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the SendProcessPID error")
}

func TestWatchdog_Stop_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	mockClient.EXPECT().
		SendStop().
		Return(transportmessages.ShutdownResponse{}, nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.Stop()

	// Assert
	assert.NoError(t, err, "Stop should not return an error")
}

func TestWatchdog_Stop_StopErrors(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	mockClient.EXPECT().
		SendStop().
		Return(transportmessages.ShutdownResponse{}, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.Stop()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the Stop error")
}

func TestWatchdog_Stop_WaitsIfNotStarted(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockClientFactory := &watchdogmocks.MockClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSocketFactory := &watchdogmocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockClient := &transportmocks.MockClient{}
	defer mockClient.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedSocketPath := "socket-path"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartNewProcess().
		Return(nil).
		Once()

	mockClientFactory.EXPECT().
		New().
		Return(mockClient).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(expectedSocketPath).
		Once()

	mockClient.EXPECT().
		Connect(expectedSocketPath).
		Return(nil).
		Once()

	mockClient.EXPECT().
		SendStop().
		Return(transportmessages.ShutdownResponse{}, nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockClientFactory,
		mockLoggerFactory,
		mockSocketFactory,
	)

	// Act & Assert
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.Stop()
	}()

	select {
	case <-errC:
		t.Fatal("Stop should block until started")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: Stop blocks until started
	}

	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	assert.NoError(t, <-errC, "Stop should not return an error")
}
