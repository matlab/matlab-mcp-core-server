// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"context"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/matlabmanager"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMATLABManager_StartMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	sessionCleanupFunc := func() error { return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	var sessionID entities.SessionID
	var err error
	synctest.Test(t, func(t *testing.T) {
		sessionID, err = manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		// The greeting and title run asynchronously; drain both goroutines before asserting.
		synctest.Wait()
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionID, sessionID)
}

func TestMATLABManager_StartMATLABSession_NoDesktop_GreetingReceivesShowMATLABDesktopFalse(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	sessionCleanupFunc := func() error { return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      false,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: false}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      false,
	}

	// Act
	var sessionID entities.SessionID
	var err error
	synctest.Test(t, func(t *testing.T) {
		sessionID, err = manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		// The greeting and title run asynchronously; drain both goroutines before asserting.
		synctest.Wait()
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionID, sessionID)
}

func TestMATLABManager_StartMATLABSession_GreetingError_IsSwallowed(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	sessionCleanupFunc := func() error { return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

		// gate greeting behind title so the Warn write can't race with the title task formatting the logger
	titleArgsFormatted := make(chan struct{})

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Run(func(context.Context, entities.Logger, entities.MATLABSessionClient, entities.GreetingInfo) {
			close(titleArgsFormatted)
		}).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Run(func(context.Context, entities.Logger, entities.MATLABSessionClient, entities.GreetingInfo) {
			<-titleArgsFormatted
		}).
		Return(assert.AnError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	var sessionID entities.SessionID
	var err error
	synctest.Test(t, func(t *testing.T) {
		sessionID, err = manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		// The greeting and title run asynchronously; drain both goroutines before asserting.
		synctest.Wait()
	})

	// Assert -- the connection still succeeds even though the greeting failed.
	require.NoError(t, err)
	assert.Equal(t, expectedSessionID, sessionID)
	_, hasWarnLog := mockLogger.WarnLogs()["failed to display MCP connection greeting in MATLAB command window"]
	assert.True(t, hasWarnLog, "greeting failure should be logged at Warn")
}

func TestMATLABManager_StartMATLABSession_MATLABServicesError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedError := assert.AnError

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(embeddedconnector.ConnectionDetails{}, nil, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_ClientFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "12345",
	}
	cleanupCalled := false
	sessionCleanupFunc := func() error { cleanupCalled = true; return nil }
	expectedError := assert.AnError

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	expectedCtx := t.Context()

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(nil, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, sessionID)
	assert.True(t, cleanupCalled, "session cleanup should be called when client factory fails")
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedSessionID := entities.SessionID(42)
	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}

	clientInfo := entities.MCPClientInfo{Title: "Claude Code", Name: "claude"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		Ping(expectedCtx, mockLogger.AsMockArg()).
		Return(entities.PingResponse{IsAlive: true}).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	// Act
	var sessionID entities.SessionID
	var err error
	synctest.Test(t, func(t *testing.T) {
		sessionID, err = manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})
		// The greeting and title run asynchronously; drain both goroutines before asserting.
		synctest.Wait()
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionID, sessionID)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_SessionSelectorError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(embeddedconnector.ConnectionDetails{}, assert.AnError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_ClientFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "key",
		CertificatePEM: []byte("cert"),
	}
	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(nil, assert.AnError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_PingFailure(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}
	expectedCtx := t.Context()

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		Ping(expectedCtx, mockLogger.AsMockArg()).
		Return(entities.PingResponse{IsAlive: false}).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	// Act
	sessionID, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})

	// Assert
	require.ErrorIs(t, err, matlabmanager.ErrMATLABSessionNotAlive)
	assert.Empty(t, sessionID)
}

func TestMATLABManager_StartMATLABSession_LocalStopSession_EvalsExitThenCleanup(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	cleanupCalled := false
	sessionCleanupFunc := func() error { cleanupCalled = true; return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()

	var capturedClient matlabsessionstore.MATLABSessionClientWithCleanup

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Run(func(client matlabsessionstore.MATLABSessionClientWithCleanup) {
			capturedClient = client
		}).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	// The teardown closure evals exit() through the gate client, then runs cleanup.
	mockSessionClient.EXPECT().
		Eval(mock.Anything, mockLogger.AsMockArg(), entities.EvalRequest{Code: "exit()"}).
		Return(entities.EvalResponse{}, nil).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	var stopErr error
	synctest.Test(t, func(t *testing.T) {
		_, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		require.NoError(t, err)
		// Greeting and title run asynchronously; drain both so the gate is open before teardown.
		synctest.Wait()
		stopErr = capturedClient.StopSession(t.Context(), mockLogger)
	})

	// Assert
	require.NoError(t, stopErr)
	assert.True(t, cleanupCalled, "session cleanup should run after exit() succeeds")
}

func TestMATLABManager_StartMATLABSession_LocalStopSession_EvalError_SkipsCleanup(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	cleanupCalled := false
	sessionCleanupFunc := func() error { cleanupCalled = true; return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()
	expectedError := assert.AnError

	var capturedClient matlabsessionstore.MATLABSessionClientWithCleanup

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Run(func(client matlabsessionstore.MATLABSessionClientWithCleanup) {
			capturedClient = client
		}).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	// exit() fails, so the closure returns early and never reaches cleanup.
	mockSessionClient.EXPECT().
		Eval(mock.Anything, mockLogger.AsMockArg(), entities.EvalRequest{Code: "exit()"}).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	var stopErr error
	synctest.Test(t, func(t *testing.T) {
		_, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		require.NoError(t, err)
		// Greeting and title run asynchronously; drain both so the gate is open before teardown.
		synctest.Wait()
		stopErr = capturedClient.StopSession(t.Context(), mockLogger)
	})

	// Assert
	require.ErrorIs(t, stopErr, expectedError)
	assert.False(t, cleanupCalled, "session cleanup should be skipped when exit() fails")
}

func TestMATLABManager_StartMATLABSession_LocalStopSession_PendingGreetingWithDeadline_StillRunsCleanup(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	cleanupCalled := false
	sessionCleanupFunc := func() error { cleanupCalled = true; return nil }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()

	var capturedClient matlabsessionstore.MATLABSessionClientWithCleanup

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Run(func(client matlabsessionstore.MATLABSessionClientWithCleanup) {
			capturedClient = client
		}).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	// greeting stays pending until the test releases it, modelling a connection greeting that has not finished when teardown starts. The channel is made inside the synctest bubble so a goroutine blocked on it counts as durably blocked.
	var releaseGreeting chan struct{}

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Run(func(_ context.Context, _ entities.Logger, _ entities.MATLABSessionClient, _ entities.GreetingInfo) {
			<-releaseGreeting
		}).
		Return(nil).
		Once()

	// exit() returns the ctx error, mirroring the real HTTP client failing on an expired ctx.
	mockSessionClient.EXPECT().
		Eval(mock.Anything, mockLogger.AsMockArg(), entities.EvalRequest{Code: "exit()"}).
		RunAndReturn(func(ctx context.Context, _ entities.Logger, _ entities.EvalRequest) (entities.EvalResponse, error) {
			return entities.EvalResponse{}, ctx.Err()
		}).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	var stopErr error
	synctest.Test(t, func(t *testing.T) {
		releaseGreeting = make(chan struct{})

		_, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		require.NoError(t, err)

		// The stop context still has time left when teardown starts, but expires within
		// the greeting timeout, so anything that waits on the pending greeting burns it.
		stopCtx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
		defer cancel()
		stopErr = capturedClient.StopSession(stopCtx, mockLogger)

		close(releaseGreeting)
		synctest.Wait()
	})

	// Assert
	assert.True(t, cleanupCalled, "session cleanup must run when teardown does not wait on the pending greeting")
	require.NoError(t, stopErr)
}

func TestMATLABManager_StartMATLABSession_LocalStopSession_CleanupError_IsPropagated(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedMATLABRoot := filepath.Join("path", "to", "matlab", "R2023a")
	expectedSessionID := entities.SessionID(123)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host: "localhost",
		Port: "1234",
	}

	expectedCleanupError := assert.AnError
	sessionCleanupFunc := func() error { return expectedCleanupError }

	expectedLocalSessionDetails := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	expectedCtx := t.Context()

	var capturedClient matlabsessionstore.MATLABSessionClientWithCleanup

	mockMATLABServices.EXPECT().
		StartLocalMATLABSession(expectedCtx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(connectionDetails, sessionCleanupFunc, nil).
		Once()

	mockClientFactory.EXPECT().
		New(connectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Run(func(client matlabsessionstore.MATLABSessionClientWithCleanup) {
			capturedClient = client
		}).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return("", nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	// exit() succeeds, so cleanup runs and its error surfaces from StopSession.
	mockSessionClient.EXPECT().
		Eval(mock.Anything, mockLogger.AsMockArg(), entities.EvalRequest{Code: "exit()"}).
		Return(entities.EvalResponse{}, nil).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	// Act
	var stopErr error
	synctest.Test(t, func(t *testing.T) {
		_, err := manager.StartMATLABSession(expectedCtx, mockLogger, startRequest)
		require.NoError(t, err)
		// Greeting and title run asynchronously; drain both so the gate is open before teardown.
		synctest.Wait()
		stopErr = capturedClient.StopSession(t.Context(), mockLogger)
	})

	// Assert
	require.ErrorIs(t, stopErr, expectedCleanupError)
}

func TestMATLABManager_StartMATLABSession_AttachToExistingSession_StopSession_RestoresTitleAndSkipsTeardown(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABServices := &mocks.MockMATLABServices{}
	defer mockMATLABServices.AssertExpectations(t)

	mockSessionStore := &mocks.MockMATLABSessionStore{}
	defer mockSessionStore.AssertExpectations(t)

	mockClientFactory := &mocks.MockMATLABSessionClientFactory{}
	defer mockClientFactory.AssertExpectations(t)

	mockSessionSelector := &mocks.MockSessionSelector{}
	defer mockSessionSelector.AssertExpectations(t)

	mockClientInfoProvider := &mocks.MockMCPClientInfoProvider{}
	defer mockClientInfoProvider.AssertExpectations(t)

	mockConnectionIndicator := &mocks.MockConnectionIndicator{}
	defer mockConnectionIndicator.AssertExpectations(t)

	mockSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockSessionClient.AssertExpectations(t)

	expectedSessionID := entities.SessionID(42)
	expectedConnectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "31515",
		APIKey:         "test-api-key",
		CertificatePEM: []byte("cert-content"),
	}

	clientInfo := entities.MCPClientInfo{Title: "Claude Code", Name: "claude"}

	expectedGreetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	originalTitle := "Original Title"

	expectedCtx := t.Context()

	var capturedClient matlabsessionstore.MATLABSessionClientWithCleanup

	mockSessionSelector.EXPECT().
		SelectSessionToAttachTo(mockLogger.AsMockArg()).
		Return(expectedConnectionDetails, nil).
		Once()

	mockClientFactory.EXPECT().
		New(expectedConnectionDetails).
		Return(mockSessionClient, nil).
		Once()

	mockSessionClient.EXPECT().
		Ping(expectedCtx, mockLogger.AsMockArg()).
		Return(entities.PingResponse{IsAlive: true}).
		Once()

	mockSessionStore.EXPECT().
		Add(mock.AnythingOfType("*matlabmanager.cleanupSessionClient")).
		Run(func(client matlabsessionstore.MATLABSessionClientWithCleanup) {
			capturedClient = client
		}).
		Return(expectedSessionID).
		Once()

	mockClientInfoProvider.EXPECT().
		GetClientInfo().
		Return(clientInfo).
		Once()

	mockConnectionIndicator.EXPECT().
		ApplyConnectedTitle(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(originalTitle, nil).
		Once()

	mockConnectionIndicator.EXPECT().
		ShowGreeting(mock.MatchedBy(func(ctx context.Context) bool {
			_, hasDeadline := ctx.Deadline()
			return hasDeadline
		}), mockLogger.AsMockArg(), mockSessionClient, expectedGreetingInfo).
		Return(nil).
		Once()

	// Detaching an externally managed session restores the title but never evals exit().
	mockConnectionIndicator.EXPECT().
		RestoreTitle(mock.Anything, mockLogger.AsMockArg(), mockSessionClient, originalTitle).
		Return(nil).
		Once()

	manager := matlabmanager.New(mockConfigFactory, mockMATLABServices, mockSessionStore, mockClientFactory, mockSessionSelector, mockClientInfoProvider, mockConnectionIndicator)

	// Act
	var stopErr error
	synctest.Test(t, func(t *testing.T) {
		_, err := manager.StartMATLABSession(expectedCtx, mockLogger, entities.AttachToExistingSession{})
		require.NoError(t, err)
		// Greeting and title run asynchronously; drain both so the captured title is available before teardown.
		synctest.Wait()
		stopErr = capturedClient.StopSession(t.Context(), mockLogger)
	})

	// Assert
	require.NoError(t, stopErr)
	_, hasDebugLog := mockLogger.DebugLogs()["Skipping session stop for externally managed MATLAB session"]
	assert.True(t, hasDebugLog, "should log that session stop was skipped for the externally managed session")
}
