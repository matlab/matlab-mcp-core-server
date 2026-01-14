// Copyright 2025-2026 The MathWorks, Inc.

package orchestrator_test

import (
	"os"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/orchestrator"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	orchestratormocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/orchestrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	//Act
	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Assert
	assert.NotNil(t, orchestratorInstance, "Orchestrator instance should not be nil")
}

func TestOrchestrator_StartAndWaitForCompletion_ConfigError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Config")
}

func TestOrchestrator_StartAndWaitForCompletion_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from GetGlobalLogger")
}

func TestOrchestrator_StartAndWaitForCompletion_DirectoryError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(nil, expectedError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Directory")
}

func TestOrchestrator_StartAndWaitForCompletion_WatchdogStartError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(expectedError).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from watchdogClient.Start")
}

func TestOrchestrator_StartAndWaitForCompletion_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	serverStarted := make(chan struct{})
	stopServer := make(chan struct{})
	defer close(stopServer)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	// Server should run indefinitely (simulate with a blocking channel)
	mockServer.EXPECT().
		Run().
		RunAndReturn(func() error {
			close(serverStarted)
			<-stopServer
			return nil
		}).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	<-serverStarted

	sendInterruptSignal(interruptC)

	// Assert
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error on signal interrupt")
}

func TestOrchestrator_StartAndWaitForCompletion_InitializeMATLABOnStartup_False(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	ctx := t.Context()
	interruptC := getInterruptChannel()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	serverStarted := make(chan struct{})
	stopServer := make(chan struct{})
	defer close(stopServer)

	mockServer.EXPECT().
		Run().
		RunAndReturn(func() error {
			close(serverStarted)
			<-stopServer
			return nil
		}).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(false).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	<-serverStarted
	sendInterruptSignal(interruptC)

	// Assert
	require.NoError(t, <-errC)
}

func TestOrchestrator_StartAndWaitForCompletion_ServerError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		Return(expectedError).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the server error")
}

func TestOrchestrator_StartAndWaitForCompletion_InitializeMATLABErrorDoesNotTriggerShutdown(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	closeServerRoutine := make(chan struct{})
	isShutdownCalled := make(chan struct{})

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		RunAndReturn(func() error {
			<-closeServerRoutine
			return nil
		}).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(getInterruptChannel()).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Run(func() {
			close(isShutdownCalled)
		}).
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(expectedError).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	// Assert
	select {
	case <-isShutdownCalled:
		t.Fatal("RequestShutdown should not be called when MATLAB initialization fails")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: no shutdown request should occur
	}

	close(closeServerRoutine)
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error")
}

func TestOrchestrator_StartAndWaitForCompletion_WaitForShutdownToCompleteError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		Return(nil).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(expectedError).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	// Assert
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error on signal interrupt")

	// This is mostly optional
	logs := mockLogger.WarnLogs()

	fields, found := logs["MATLAB MCP Core Server application shutdown failed"]
	require.True(t, found, "Expected a warning log about shutdown failure")

	errField, found := fields["error"]
	require.True(t, found, "Expected an error field in the warning log")

	err, ok := errField.(error)
	require.True(t, ok, "Error field should be of type error")
	require.ErrorIs(t, err, expectedError, "Logged error should match the shutdown error")
}

func TestOrchestrator_StartAndWaitForCompletion_WatchdogStopError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	serverStarted := make(chan struct{})
	stopServer := make(chan struct{})
	defer close(stopServer)

	mockServer.EXPECT().
		Run().
		RunAndReturn(func() error {
			close(serverStarted)
			<-stopServer
			return nil
		}).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	// Signal
	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	// Shutdown sequence
	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	// Watchdog Stop Fails
	mockWatchdogClient.EXPECT().
		Stop().
		Return(expectedError).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler, mockConfigFactory, mockServer, mockWatchdogClient,
		mockLoggerFactory, mockSignalLayer, mockGlobalMATLABManager, mockDirectoryFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	<-serverStarted
	sendInterruptSignal(interruptC)
	err := <-errC

	// Assert
	require.NoError(t, err)

	// Verify Log
	logs := mockLogger.WarnLogs()
	fields, found := logs["Watchdog shutdown failed"]
	require.True(t, found, "Expected warning log for watchdog failure")
	assert.Equal(t, expectedError, fields["error"])
}

func TestOrchestrator_StartAndWaitForCompletion_MultipleSession_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigFactory := &orchestratormocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t) // Implicit assertion here, Initialize should not be called

	mockDirectoryFactory := &orchestratormocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		Return(nil).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(false).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfigFactory,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectoryFactory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.NoError(t, err)
}

func getInterruptChannel() chan os.Signal {
	return make(chan os.Signal, 1)
}

func sendInterruptSignal(interruptC chan os.Signal) {
	interruptC <- os.Interrupt
}
