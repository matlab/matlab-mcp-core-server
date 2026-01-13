// Copyright 2025-2026 The MathWorks, Inc.

package watchdog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog"
	transportmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport"
	handlermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server/handler"
	socketmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/socket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	// Act
	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Assert
	assert.NotNil(t, watchdogInstance, "Watchdog instance should not be nil")
}

func TestWatchdog_StartAndWaitForCompletion_GracefulShutdown(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerHandler := &handlermocks.MockHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	shutdownFuncC := make(chan func())
	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockServerHandlerFactory.EXPECT().
		Handler().
		Return(mockServerHandler, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Run(func(fn func()) {
			shutdownFuncC <- fn
		}).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC, nil).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockServerHandler.EXPECT().
		TerminateAllProcesses().
		Return().
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted
	shutdownFcn := <-shutdownFuncC

	shutdownFcn()

	// Assert
	require.NoError(t, <-errC)
}

func TestWatchdog_StartAndWaitForCompletion_ParentProcessTermination(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerHandler := &handlermocks.MockHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockServerHandlerFactory.EXPECT().
		Handler().
		Return(mockServerHandler, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC, nil).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockServerHandler.EXPECT().
		TerminateAllProcesses().
		Return().
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted
	close(parentTerminationC)

	// Assert
	require.NoError(t, <-errC)
}

func TestWatchdog_StartAndWaitForCompletion_OSSignalInterrupt(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerHandler := &handlermocks.MockHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockServerHandlerFactory.EXPECT().
		Handler().
		Return(mockServerHandler, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC, nil).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockServerHandler.EXPECT().
		TerminateAllProcesses().
		Return().
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted
	interruptSignalC <- os.Interrupt

	// Assert
	require.NoError(t, <-errC)
}

func TestWatchdog_StartAndWaitForCompletion_SocketFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestWatchdog_StartAndWaitForCompletion_ServerFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerHandler := &handlermocks.MockHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockServerHandlerFactory.EXPECT().
		Handler().
		Return(mockServerHandler, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestWatchdog_StartAndWaitForCompletion_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestWatchdog_StartAndWaitForCompletion_ServerHandlerFactoryError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockServerHandlerFactory.EXPECT().
		Handler().
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	err := watchdogInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestWatchdog_StartAndWaitForCompletion_WatchProcessAndGetTerminationChanError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockServerHandlerFactory := &mocks.MockServerHandlerFactory{}
	defer mockServerHandlerFactory.AssertExpectations(t)

	mockServerHandler := &handlermocks.MockHandler{}
	defer mockServerHandler.AssertExpectations(t)

	mockServerFactory := &mocks.MockServerFactory{}
	defer mockServerFactory.AssertExpectations(t)

	mockSocketFactory := &mocks.MockSocketFactory{}
	defer mockSocketFactory.AssertExpectations(t)

	mockServer := &transportmocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockSocket := &socketmocks.MockSocket{}
	defer mockSocket.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	socketPath := filepath.Join(t.TempDir(), "test.sock")
	serverStarted := make(chan struct{})
	expectedParentPID := 1234
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockServerHandlerFactory.EXPECT().
		Handler().
		Return(mockServerHandler, nil).
		Once()

	mockSocketFactory.EXPECT().
		Socket().
		Return(mockSocket, nil).
		Once()

	mockSocket.EXPECT().
		Path().
		Return(socketPath).
		Once()

	mockServerFactory.EXPECT().
		New().
		Return(mockServer, nil).
		Once()

	mockServer.EXPECT().
		Start(socketPath).
		Run(func(_ string) {
			close(serverStarted)
		}).
		Return(nil).
		Once()

	mockServer.EXPECT().
		Stop().
		Return(nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockServerHandler.EXPECT().
		RegisterShutdownFunction(mock.AnythingOfType("func()")).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockServerHandlerFactory,
		mockServerFactory,
		mockSocketFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	<-serverStarted

	// Assert
	require.ErrorIs(t, <-errC, expectedError)
}
