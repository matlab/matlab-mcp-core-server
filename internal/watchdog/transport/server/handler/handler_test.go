// Copyright 2025-2026 The MathWorks, Inc.

package handler_test

import (
	"testing"

	internalmessages "github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
	handlermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	// Act
	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)

	// Act
	h, err := factory.Handler()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, h)
}

func TestFactory_Handler_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	expectedError := internalmessages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)

	// Act
	h, err := factory.Handler()

	// Assert
	assert.Nil(t, h)
	assert.Equal(t, expectedError, err)
}

func TestHandler_HandleProcessToKill_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedPID := 12345
	request := messages.ProcessToKillRequest{
		PID: expectedPID,
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	// Act
	response, err := h.HandleProcessToKill(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ProcessToKillResponse{}, response)
}

func TestHandler_HandleShutdown_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	request := messages.ShutdownRequest{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	// Act
	response, err := h.HandleShutdown(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ShutdownResponse{}, response)
}

func TestHandler_HandleShutdown_CallsRegisteredFunctions(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	request := messages.ShutdownRequest{}
	functionCalled := false

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	h.RegisterShutdownFunction(func() {
		functionCalled = true
	})

	// Act
	response, err := h.HandleShutdown(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ShutdownResponse{}, response)
	assert.True(t, functionCalled, "Registered shutdown function should have been called")
}

func TestHandler_HandleShutdown_CallsMultipleRegisteredFunctions(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	request := messages.ShutdownRequest{}
	callOrder := []int{}
	expectedCallOrder := []int{1, 2}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	h.RegisterShutdownFunction(func() {
		callOrder = append(callOrder, expectedCallOrder[0])
	})
	h.RegisterShutdownFunction(func() {
		callOrder = append(callOrder, expectedCallOrder[1])
	})

	// Act
	response, err := h.HandleShutdown(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, messages.ShutdownResponse{}, response)
	assert.Equal(t, expectedCallOrder, callOrder, "Registered shutdown functions should be called in order")
}

func TestHandler_TerminateAllProcesses_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedPID := 12345

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID).
		Return(nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID})
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}

func TestHandler_TerminateAllProcesses_MultiplePIDs(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedPID1 := 12345
	expectedPID2 := 67890

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID1).
		Return(nil).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID2).
		Return(nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID1})
	require.NoError(t, err)
	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID2})
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}

func TestHandler_TerminateAllProcesses_KillProcessError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedPID1 := 12345
	expectedPID2 := 67890
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID1).
		Return(expectedError).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPID2).
		Return(nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID1})
	require.NoError(t, err)
	_, err = h.HandleProcessToKill(messages.ProcessToKillRequest{PID: expectedPID2})
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}

func TestHandler_TerminateAllProcesses_NoPIDs(t *testing.T) {
	// Arrange
	mockLoggerFactory := &handlermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockProcessHandler := &handlermocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	factory := handler.NewFactory(mockLoggerFactory, mockProcessHandler)
	h, err := factory.Handler()
	require.NoError(t, err)

	// Act
	h.TerminateAllProcesses()

	// Assert
}
