// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactory_Directory_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &directorymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	logDir := filepath.Join("tmp", "matlab-mcp-core-server-12345")
	expectedMarkerFileBase := filepath.Join(logDir, directory.MarkerFileName)
	markerFileName := filepath.Join(logDir, ".matlab-mcp-core-server-123")
	suffix := "1337"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockConfig.EXPECT().
		ServerInstanceID().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(logDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, "").
		Return(markerFileName, suffix, nil).
		Once()

	factory := directory.NewFactory(mockConfigFactory, mockFileNameFactory, mockOSLayer)

	// Act
	directoryInstance, err := factory.Directory()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, directoryInstance)
}

func TestFactory_Directory_ReturnsSameInstance(t *testing.T) {
	// Arrange
	mockConfigFactory := &directorymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	logDir := filepath.Join("tmp", "matlab-mcp-core-server-12345")
	expectedMarkerFileBase := filepath.Join(logDir, directory.MarkerFileName)
	markerFileName := filepath.Join(logDir, ".matlab-mcp-core-server-123")
	suffix := "1337"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockConfig.EXPECT().
		ServerInstanceID().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(logDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, "").
		Return(markerFileName, suffix, nil).
		Once()

	factory := directory.NewFactory(mockConfigFactory, mockFileNameFactory, mockOSLayer)

	// Act
	firstCall, err1 := factory.Directory()
	require.NoError(t, err1)

	secondCall, err2 := factory.Directory()
	require.NoError(t, err2)

	// Assert
	assert.Same(t, firstCall, secondCall, "Factory should return the same Directory instance on subsequent calls")
}

func TestFactory_Directory_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &directorymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	mockOSLayer := &directorymocks.MockOSLayer{}

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := directory.NewFactory(mockConfigFactory, mockFileNameFactory, mockOSLayer)

	// Act
	directoryInstance, err := factory.Directory()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, directoryInstance)
}

func TestFactory_Directory_CachesError(t *testing.T) {
	// Arrange
	mockConfigFactory := &directorymocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := directory.NewFactory(mockConfigFactory, mockFileNameFactory, mockOSLayer)

	// Act
	firstCall, err1 := factory.Directory()
	secondCall, err2 := factory.Directory()

	// Assert
	require.ErrorIs(t, err1, expectedError)
	require.ErrorIs(t, err2, expectedError)
	assert.Nil(t, firstCall)
	assert.Nil(t, secondCall)
}
