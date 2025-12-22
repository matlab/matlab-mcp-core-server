// Copyright 2025 The MathWorks, Inc.

package modeselector_test

import (
	"fmt"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector"
	modeselectormocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/modeselector"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	// Act
	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Assert
	assert.NotNil(t, modeSelectorInstance, "ModeSelector instance should not be nil")
}

func TestStartAndWaitForCompletion_VersionMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	expectedVersion := "25.6.68"

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(true).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", expectedVersion))).
		Return(len(expectedVersion)+1, nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in version mode")
}

func TestStartAndWaitForCompletion_VersionMode_WriteError(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	expectedVersion := "25.6.68"
	expectedError := assert.AnError

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(true).
		Once()

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", expectedVersion))).
		Return(0, expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Write")
}

func TestStartAndWaitForCompletion_WatchdogMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockWatchdogProcess := &entitiesmocks.MockMode{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	ctx := t.Context()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	mockWatchdogProcessFactory.EXPECT().
		Create().
		Return(mockWatchdogProcess, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in watchdog mode")
}

func TestStartAndWaitForCompletion_WatchdogMode_CreateError(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	expectedError := assert.AnError

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	mockWatchdogProcessFactory.EXPECT().
		Create().
		Return(nil, expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Create")
}

func TestStartAndWaitForCompletion_WatchdogMode_StartAndWaitError(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockWatchdogProcess := &entitiesmocks.MockMode{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	expectedError := assert.AnError
	ctx := t.Context()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	mockWatchdogProcessFactory.EXPECT().
		Create().
		Return(mockWatchdogProcess, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from StartAndWaitForCompletion")
}

func TestStartAndWaitForCompletion_DefaultMode_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockOrchestrator := &entitiesmocks.MockMode{}
	defer mockOrchestrator.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	ctx := t.Context()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockOrchestratorFactory.EXPECT().
		Create().
		Return(mockOrchestrator, nil).
		Once()

	mockOrchestrator.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(nil).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.NoError(t, err, "StartAndWaitForCompletion should not return an error in default mode")
}

func TestStartAndWaitForCompletion_DefaultMode_CreateError(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	expectedError := assert.AnError

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockOrchestratorFactory.EXPECT().
		Create().
		Return(nil, expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(t.Context())

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from Create")
}

func TestStartAndWaitForCompletion_DefaultMode_StartAndWaitError(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockOrchestrator := &entitiesmocks.MockMode{}
	defer mockOrchestrator.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	expectedError := assert.AnError
	ctx := t.Context()

	mockConfig.EXPECT().
		HelpMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		VersionMode().
		Return(false).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockOrchestratorFactory.EXPECT().
		Create().
		Return(mockOrchestrator, nil).
		Once()

	mockOrchestrator.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(expectedError).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, expectedError, "StartAndWaitForCompletion should return the error from StartAndWaitForCompletion")
}

func TestStartAndWaitForCompletion_HelpMode_StartAndWaitHappyPath(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockOrchestrator := &entitiesmocks.MockMode{}
	defer mockOrchestrator.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	dummyHelpText := "Help me get my feet back on the ground."
	ctx := t.Context()

	mockConfig.EXPECT().
		HelpMode().
		Return(true).
		Once()

	mockParser.EXPECT().
		Usage().
		Return(dummyHelpText).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", dummyHelpText))).
		Return(len(dummyHelpText)+1, nil).
		Once()

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.NoError(t, err)
}

func TestStartAndWaitForCompletion_HelpMode_StartAndWaitWriteError(t *testing.T) {
	// Arrange
	mockConfig := &modeselectormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogProcessFactory := &modeselectormocks.MockWatchdogProcessFactory{}
	defer mockWatchdogProcessFactory.AssertExpectations(t)

	mockOrchestratorFactory := &modeselectormocks.MockOrchestratorFactory{}
	defer mockOrchestratorFactory.AssertExpectations(t)

	mockOsLayer := &modeselectormocks.MockOSLayer{}
	defer mockOsLayer.AssertExpectations(t)

	mockOrchestrator := &entitiesmocks.MockMode{}
	defer mockOrchestrator.AssertExpectations(t)

	mockParser := &modeselectormocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	dummyHelpText := "Help me get my feet back on the ground."
	dummyError := assert.AnError
	ctx := t.Context()

	mockConfig.EXPECT().
		HelpMode().
		Return(true).
		Once()

	mockParser.EXPECT().
		Usage().
		Return(dummyHelpText).
		Once()

	modeSelectorInstance := modeselector.New(
		mockConfig,
		mockParser,
		mockWatchdogProcessFactory,
		mockOrchestratorFactory,
		mockOsLayer,
	)

	mockOsLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockStdout.EXPECT().
		Write([]byte(fmt.Sprintf("%s\n", dummyHelpText))).
		Return(0, dummyError).
		Once()

	// Act
	err := modeSelectorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.ErrorIs(t, err, dummyError)
}
