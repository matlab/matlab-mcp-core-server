// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGlobalMATLAB_Client_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_StartingDirectorySet(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_ShowMATLABDesktopFalse(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	shouldShowMATLABDesktop := false

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_SelectMATLABRootError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return("", expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_MATLABStartingDirSelectionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return("", assert.AnError).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_StartMATLABSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := assert.AnError
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_GetMATLABSessionClientError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedNewSessionID := entities.SessionID(456)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := assert.AnError
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil, expectedError).
		Once()

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedNewSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedNewSessionID).
		Return(nil, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_GetMATLABSessionClientError_RetrySucceeds(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedNewSessionID := entities.SessionID(456)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil, assert.AnError).
		Once()

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedNewSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedNewSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_ReturnsInitializeCachedErrorOnSubsequentClientCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return("", expectedError).
		Once()

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	assert.Nil(t, client1)
	require.ErrorIs(t, err1, expectedError)

	assert.Nil(t, client2)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Client_ReturnsMATLABStartupCachedErrorOnSubsequentClientCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := assert.AnError
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	assert.Nil(t, client1)
	require.ErrorIs(t, err1, expectedError)

	assert.Nil(t, client2)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Client_ConcurrentCallsWaitForCompletion(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedSessionID := entities.SessionID(123)
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	blockStartMATLAB := make(chan struct{})
	startMATLABCalled := make(chan struct{})

	type clientResult struct {
		client entities.MATLABSessionClient
		err    error
	}

	firstCallCompleted := make(chan clientResult)
	secondCallCompleted := make(chan clientResult)

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Run(func(ctx context.Context, logger entities.Logger, details entities.SessionDetails) {
			close(startMATLABCalled)
			<-blockStartMATLAB
		}).
		Return(expectedSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		firstCallCompleted <- clientResult{client: client, err: err}
	}()

	<-startMATLABCalled

	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		secondCallCompleted <- clientResult{client: client, err: err}
	}()

	select {
	case <-secondCallCompleted:
		t.Fatal("Second Client call completed before first call was unblocked")
	case <-time.After(100 * time.Millisecond):
		// Second call is still waiting
	}

	close(blockStartMATLAB)
	result1 := <-firstCallCompleted
	result2 := <-secondCallCompleted

	// Assert
	require.NoError(t, result1.err)
	assert.Equal(t, expectedSessionClient, result1.client)

	require.NoError(t, result2.err)
	assert.Equal(t, expectedSessionClient, result2.client)
}

func TestGlobalMATLAB_Client_RestartOnGetClientFailure(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedInitialSessionID := entities.SessionID(123)
	expectedNewSessionID := entities.SessionID(456)
	expectedMATLABRoot := ""
	expectedMATLABStartingDir := ""
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedInitialSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedInitialSessionID).
		Return(expectedSessionClient, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedInitialSessionID).
		Return(nil, assert.AnError).
		Once()

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedInitialSessionID).
		Return(nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedNewSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedNewSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	require.NoError(t, err1)
	assert.Equal(t, expectedSessionClient, client1)

	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err2)
	assert.Equal(t, expectedSessionClient, client2)
}

func TestGlobalMATLAB_Client_DoesNotErrorIfStopSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedInitialSessionID := entities.SessionID(123)
	expectedNewSessionID := entities.SessionID(456)
	expectedMATLABRoot := ""
	expectedMATLABStartingDir := ""
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedInitialSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedInitialSessionID).
		Return(expectedSessionClient, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedInitialSessionID).
		Return(nil, assert.AnError).
		Once()

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedInitialSessionID).
		Return(assert.AnError).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedNewSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedNewSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	require.NoError(t, err1)
	assert.Equal(t, expectedSessionClient, client1)

	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err2)
	assert.Equal(t, expectedSessionClient, client2)
}

func TestGlobalMATLAB_Client_RestartFailure_OnExistingSession(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSessionID := entities.SessionID(123)
	expectedMATLABRoot := ""
	expectedMATLABStartingDir := ""
	expectedError := assert.AnError
	shouldShowMATLABDesktop := true

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(expectedSessionClient, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil, assert.AnError).
		Once()

	mockMATLABManager.EXPECT().
		StopMATLABSession(ctx, mockLogger.AsMockArg(), expectedSessionID).
		Return(nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(entities.SessionID(0), expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	require.NoError(t, err1)
	assert.Equal(t, expectedSessionClient, client1)

	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err2, expectedError)
	assert.Nil(t, client2)
}

func TestGlobalMATLAB_Client_ConfigError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := messages.New_StartupErrors_BadFlag_Error("flag", "value", "reason")

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}
