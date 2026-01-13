// Copyright 2025-2026 The MathWorks, Inc.

package process_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	processmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/watchdog/process"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	expectedProgramPath := filepath.Join("path", "to", "program")
	expectedBaseDir := filepath.Join("tmp", "base", "dir")
	expectedServerID := "server-id"
	expectedLogLevel := entities.LogLevelInfo

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		Executable().
		Return(expectedProgramPath, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedServerID).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(expectedLogLevel).
		Once()

	mockOSLayer.EXPECT().
		Command(expectedProgramPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, expectedBaseDir,
			"--" + flags.ServerInstanceID, expectedServerID,
			"--" + flags.LogLevel, string(expectedLogLevel),
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory, mockConfigFactory)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, processInstance, "Process instance should not be nil")
}

func TestNew_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory, mockConfigFactory)

	// Assert
	require.ErrorIs(t, err, expectedError, "New should return the error from GetGlobalLogger")
	assert.Nil(t, processInstance, "Process instance should be nil on error")
}

func TestNew_ConfigError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedError := &messages.StartupErrors_BadFlag_Error{}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory, mockConfigFactory)

	// Assert
	require.ErrorIs(t, err, expectedError, "New should return the error from Config")
	assert.Nil(t, processInstance, "Process instance should be nil on error")
}

func TestNew_ExecutableError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	expectedError := errors.New("executable error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		Executable().
		Return("", expectedError).
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory, mockConfigFactory)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the executable error")
	assert.Nil(t, processInstance, "Process instance should be nil on error")
}

func TestProcess_Start_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	expectedProgramPath := filepath.Join("path", "to", "program")
	expectedBaseDir := filepath.Join("tmp", "base", "dir")
	expectedServerID := "server-id"
	expectedLogLevel := entities.LogLevelInfo

	// Setup mocks for New
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		Executable().
		Return(expectedProgramPath, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedServerID).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(expectedLogLevel).
		Once()

	mockOSLayer.EXPECT().
		Command(expectedProgramPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, expectedBaseDir,
			"--" + flags.ServerInstanceID, expectedServerID,
			"--" + flags.LogLevel, string(expectedLogLevel),
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	mockCmd.EXPECT().
		Start().
		Return(nil).
		Once()

	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory, mockConfigFactory)
	require.NoError(t, err, "New should not return an error")

	// Act
	err = processInstance.Start()

	// Assert
	assert.NoError(t, err, "Start should not return an error")
}

func TestProcess_Start_Error(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	expectedProgramPath := filepath.Join("path", "to", "program")
	expectedBaseDir := filepath.Join("tmp", "base", "dir")
	expectedServerID := "server-id"
	expectedLogLevel := entities.LogLevelInfo
	expectedError := errors.New("start process error")

	// Setup mocks for New
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockOSLayer.EXPECT().
		Executable().
		Return(expectedProgramPath, nil).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedServerID).
		Once()

	mockConfig.EXPECT().
		LogLevel().
		Return(expectedLogLevel).
		Once()

	mockOSLayer.EXPECT().
		Command(expectedProgramPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, expectedBaseDir,
			"--" + flags.ServerInstanceID, expectedServerID,
			"--" + flags.LogLevel, string(expectedLogLevel),
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	mockCmd.EXPECT().
		Start().
		Return(expectedError).
		Once()

	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory, mockConfigFactory)
	require.NoError(t, err, "New should not return an error")

	// Act
	err = processInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the start process error")
}
