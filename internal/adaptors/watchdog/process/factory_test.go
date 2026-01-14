// Copyright 2025-2026 The MathWorks, Inc.

package process_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	processmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/watchdog/process"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &processmocks.MockOSLayer{}
	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	mockDirectoryFactory := &processmocks.MockDirectoryFactory{}
	mockConfigFactory := &processmocks.MockConfigFactory{}

	// Act
	factory := process.New(mockOSLayer, mockLoggerFactory, mockDirectoryFactory, mockConfigFactory)

	// Assert
	require.NotNil(t, factory)
}

func TestFactory_StartNewProcess_GetGlobalLoggerError(t *testing.T) {
	// Arrange
	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectoryFactory := &processmocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, messages.AnError).
		Once()

	processInstance := process.New(mockOSLayer, mockLoggerFactory, mockDirectoryFactory, mockConfigFactory)

	// Act
	err := processInstance.StartNewProcess()

	// Assert
	require.ErrorIs(t, err, messages.AnError)
}

func TestFactory_StartNewProcess_DirectoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectoryFactory := &processmocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(nil, messages.AnError).
		Once()

	processInstance := process.New(mockOSLayer, mockLoggerFactory, mockDirectoryFactory, mockConfigFactory)

	// Act
	err := processInstance.StartNewProcess()

	// Assert
	require.ErrorIs(t, err, messages.AnError)
}

func TestFactory_StartNewProcess_ConfigError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectoryFactory := &processmocks.MockDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockConfigFactory := &processmocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockDirectoryFactory.EXPECT().
		Directory().
		Return(mockDirectory, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, messages.AnError).
		Once()

	processInstance := process.New(mockOSLayer, mockLoggerFactory, mockDirectoryFactory, mockConfigFactory)

	// Act
	err := processInstance.StartNewProcess()

	// Assert
	require.ErrorIs(t, err, messages.AnError)
}
