// Copyright 2025-2026 The MathWorks, Inc.

package matlabsessionstore_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabsessionstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	// Act
	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)

	// Assert
	assert.NotNil(t, store)
}

func TestNew_ShutdownFunctionCallsStopSessionOnAllClients(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient1 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient1.AssertExpectations(t)

	mockClient2 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient2.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	var capturedShutdownFunc func() error

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFunc = shutdownFcn
		}).
		Return().
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockClient1.EXPECT().
		StopSession(mock.AnythingOfType("context.backgroundCtx"), mockLogger.AsMockArg()).
		Return(nil).
		Once()

	mockClient2.EXPECT().
		StopSession(mock.AnythingOfType("context.backgroundCtx"), mockLogger.AsMockArg()).
		Return(nil).
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	require.NotNil(t, capturedShutdownFunc)

	store.Add(mockClient1)
	store.Add(mockClient2)

	// Act
	err := capturedShutdownFunc()

	// Assert
	assert.NoError(t, err)
}

func TestNew_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	expectedError := messages.AnError
	var capturedShutdownFunc func() error

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFunc = shutdownFcn
		}).
		Return().
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	require.NotNil(t, capturedShutdownFunc)

	// Act
	err := capturedShutdownFunc()

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestNew_ShutdownFunctionHandlesEmptyStore(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	var capturedShutdownFunc func() error

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFunc = shutdownFcn
		}).
		Return().
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	require.NotNil(t, capturedShutdownFunc)

	// Act
	err := capturedShutdownFunc()

	// Assert
	require.NoError(t, err)
}

func TestNew_ShutdownFunctionReturnsErrorWhenStopSessionFails(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient1 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient1.AssertExpectations(t)

	mockClient2 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient2.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := assert.AnError
	var capturedShutdownFunc func() error

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFunc = shutdownFcn
		}).
		Return().
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockClient1.EXPECT().
		StopSession(mock.AnythingOfType("context.backgroundCtx"), mockLogger.AsMockArg()).
		Return(nil).
		Once()

	mockClient2.EXPECT().
		StopSession(mock.AnythingOfType("context.backgroundCtx"), mockLogger.AsMockArg()).
		Return(expectedError).
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	require.NotNil(t, capturedShutdownFunc)

	store.Add(mockClient1)
	store.Add(mockClient2)

	// Act
	err := capturedShutdownFunc()

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestStore_Add_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)

	// Act
	sessionID := store.Add(mockClient)

	// Assert
	assert.Equal(t, entities.SessionID(1), sessionID)
}

func TestStore_Add_MultipleClients_ReturnsIncrementingIDs(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient1 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient1.AssertExpectations(t)

	mockClient2 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient2.AssertExpectations(t)

	mockClient3 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient3.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)

	// Act
	sessionID1 := store.Add(mockClient1)
	sessionID2 := store.Add(mockClient2)
	sessionID3 := store.Add(mockClient3)

	// Assert
	assert.Equal(t, entities.SessionID(1), sessionID1)
	assert.Equal(t, entities.SessionID(2), sessionID2)
	assert.Equal(t, entities.SessionID(3), sessionID3)
}

func TestStore_Get_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	sessionID := store.Add(mockClient)

	// Act
	retrievedClient, err := store.Get(sessionID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, mockClient, retrievedClient)
}

func TestStore_Get_NonExistentSession_ReturnsError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	nonExistentSessionID := entities.SessionID(999)

	// Act
	retrievedClient, err := store.Get(nonExistentSessionID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, retrievedClient)
	assert.Contains(t, err.Error(), "session not found")
	assert.Contains(t, err.Error(), "999")
}

func TestStore_Remove_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	sessionID := store.Add(mockClient)

	// Verify client exists before removal
	retrievedClient, err := store.Get(sessionID)
	require.NoError(t, err)
	assert.Equal(t, mockClient, retrievedClient)

	// Act
	store.Remove(sessionID)

	// Assert
	retrievedClient, err = store.Get(sessionID)
	require.Error(t, err)
	assert.Nil(t, retrievedClient)
	assert.Contains(t, err.Error(), "session not found")
}

func TestStore_Remove_NonExistentSession_NoError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)
	nonExistentSessionID := entities.SessionID(999)

	// Act & Assert (should not panic or error)
	store.Remove(nonExistentSessionID)
}

func TestStore_AddGetRemove_MultipleClients(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockClient1 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient1.AssertExpectations(t)

	mockClient2 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient2.AssertExpectations(t)

	mockClient3 := &mocks.MockMATLABSessionClientWithCleanup{}
	defer mockClient3.AssertExpectations(t)

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Return().
		Once()

	store := matlabsessionstore.New(mockLoggerFactory, mockLifecycleSignaler)

	// Act - Add multiple clients
	sessionID1 := store.Add(mockClient1)
	sessionID2 := store.Add(mockClient2)
	sessionID3 := store.Add(mockClient3)

	// Assert - All clients can be retrieved
	retrievedClient1, err := store.Get(sessionID1)
	require.NoError(t, err)
	assert.Equal(t, mockClient1, retrievedClient1)

	retrievedClient2, err := store.Get(sessionID2)
	require.NoError(t, err)
	assert.Equal(t, mockClient2, retrievedClient2)

	retrievedClient3, err := store.Get(sessionID3)
	require.NoError(t, err)
	assert.Equal(t, mockClient3, retrievedClient3)

	// Act - Remove middle client
	store.Remove(sessionID2)

	// Assert - Client 2 is gone, but 1 and 3 remain
	retrievedClient1, err = store.Get(sessionID1)
	require.NoError(t, err)
	assert.Equal(t, mockClient1, retrievedClient1)

	retrievedClient2, err = store.Get(sessionID2)
	require.Error(t, err)
	assert.Nil(t, retrievedClient2)

	retrievedClient3, err = store.Get(sessionID3)
	require.NoError(t, err)
	assert.Equal(t, mockClient3, retrievedClient3)
}
