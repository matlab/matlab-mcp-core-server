// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDirectory_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockMarkerFile := &osfacademocks.MockFile{}
	defer mockMarkerFile.AssertExpectations(t)

	logDir := filepath.Join("tmp", "matlab-mcp-core-server-12345")
	expectedMarkerFileBase := filepath.Join(logDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	markerFileName := filepath.Join(logDir, ".matlab-mcp-core-server-123")
	suffix := "1337"

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
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, expectedMarkerExtension).
		Return(markerFileName, suffix, nil).
		Once()

	// Act
	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, directoryInstance, "Directory instance should not be nil")
}

func TestNewDirectory_MkdirTempError(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedError := messages.New_StartupErrors_FailedToCreateSubdirectory_Error(os.TempDir()) //nolint:usetesting // This is the expected error

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return("", assert.AnError).
		Once()

	// Act
	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.Equal(t, expectedError, err, "New should return FailedToCreateSubdirectory error")
	assert.Nil(t, directoryInstance, "Directory instance should be nil when error occurs")
}

func TestNewDirectory_MkdirAllError(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("logs", "subdir")
	expectedError := messages.New_StartupErrors_FailedToCreateDirectory_Error(expectedLogDir)

	mockConfig.EXPECT().
		BaseDir().
		Return(expectedLogDir).
		Once()

	mockOSLayer.EXPECT().
		MkdirAll(expectedLogDir, os.FileMode(0o700)).
		Return(assert.AnError).
		Once()

	// Act
	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.Equal(t, expectedError, err, "New should return FailedToCreateDirectory error")
	assert.Nil(t, directoryInstance, "Directory instance should be nil when error occurs")
}

func TestNewDirectory_CreateFileWithUniqueSuffixError(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	logDir := filepath.Join("tmp", "logs")
	expectedMarkerFileBaseName := filepath.Join(logDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	expectedError := messages.New_StartupErrors_FailedToCreateFile_Error(expectedMarkerFileBaseName)

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
		CreateFileWithUniqueSuffix(expectedMarkerFileBaseName, expectedMarkerExtension).
		Return("", "", assert.AnError).
		Once()

	// Act
	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.Equal(t, expectedError, err, "New should return FailedToCreateFile error")
	assert.Nil(t, directoryInstance, "Directory instance should be nil when error occurs")
}

func TestNewDirectory_SuppliedServerInstanceID(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	logDir := filepath.Join("tmp", "logs")
	id := "1337"

	mockConfig.EXPECT().
		BaseDir().
		Return(logDir).
		Once()

	mockConfig.EXPECT().
		ServerInstanceID().
		Return(id).
		Once()

	mockOSLayer.EXPECT().
		MkdirAll(logDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	// Act
	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, directoryInstance, "Directory instance should not be nil")
}

func TestDirectory_BaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("tmp", "matlab-mcp-core-server-67890")
	expectedMarkerFileBase := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	markerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	suffix := "1337"

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
		Return(expectedLogDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, expectedMarkerExtension).
		Return(markerFileName, suffix, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	baseDir := directoryInstance.BaseDir()

	// Assert
	assert.Equal(t, expectedLogDir, baseDir, "BaseDir should return the expected log directory")
}

func TestDirectory_BaseDir_SuppliedBaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("logs", "subdir")
	expectedMarkerFileBase := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	markerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	suffix := "123"

	mockConfig.EXPECT().
		BaseDir().
		Return(expectedLogDir).
		Once()

	mockConfig.EXPECT().
		ServerInstanceID().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirAll(expectedLogDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, expectedMarkerExtension).
		Return(markerFileName, suffix, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	baseDir := directoryInstance.BaseDir()

	// Assert
	assert.Equal(t, expectedLogDir, baseDir, "BaseDir should return the expected log directory")
}

func TestDirectory_ID_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	logDir := filepath.Join("tmp", "matlab-mcp-core-server-12345")
	expectedMarkerFileBase := filepath.Join(logDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	markerFileName := filepath.Join(logDir, ".matlab-mcp-core-server")
	expectedSuffix := "1337"

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
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, expectedMarkerExtension).
		Return(markerFileName, expectedSuffix, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	id := directoryInstance.ID()

	// Assert
	require.NoError(t, err, "ID should not return an error")
	assert.Equal(t, expectedSuffix, id, "ID should return the expected ID")
}

func TestDirectory_ID_SuppliedID_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedID := "1337"
	logDir := filepath.Join("tmp", "matlab-mcp-core-server-12345")

	mockConfig.EXPECT().
		BaseDir().
		Return("").
		Once()

	mockConfig.EXPECT().
		ServerInstanceID().
		Return(expectedID).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp("", directory.DefaultLogDirPattern).
		Return(logDir, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	id := directoryInstance.ID()

	// Assert
	require.NoError(t, err, "ID should not return an error")
	assert.Equal(t, expectedID, id, "ID should return the expected ID")
}

func TestDirectory_CreateSubDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("tmp", "matlab-mcp-core-server-11111")
	expectedMarkerBaseName := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	pattern := "test-pattern-"
	expectedTempDir := filepath.Join("tmp", "matlab-mcp-core-server-11111", "test-pattern-22222")
	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "1337"
	expectedPattern := pattern + expectedSuffix + "-"

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
		Return(expectedLogDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerBaseName, expectedMarkerExtension).
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	actualTempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, actualTempDir, "MkdirTemp should return the expected temp directory")
}

func TestDirectory_CreateSubDir_EnforcesDashSuffix(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("tmp", "matlab-mcp-core-server-11111")
	expectedMarkerBaseName := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	pattern := "test-pattern"
	expectedTempDir := filepath.Join("tmp", "matlab-mcp-core-server-11111", "test-pattern-22222")
	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "1337"
	expectedPattern := pattern + "-" + expectedSuffix + "-"

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
		Return(expectedLogDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerBaseName, expectedMarkerExtension).
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	tempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, tempDir, "MkdirTemp should return the expected temp directory")
}

func TestDirectory_CreateSubDir_MkdirTempError(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("tmp", "matlab-mcp-core-server-33333")
	expectedMarkerBaseName := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	pattern := "test-pattern-"
	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "1337"
	expectedPattern := pattern + expectedSuffix + "-"
	expectedError := messages.New_StartupErrors_FailedToCreateSubdirectory_Error(expectedLogDir)

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
		Return(expectedLogDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerBaseName, expectedMarkerExtension).
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedPattern).
		Return("", assert.AnError).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	actualTempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.Equal(t, expectedError, err, "CreateSubDir should return FailedToCreateSubdirectory error")
	assert.Empty(t, actualTempDir, "CreateSubDir should return empty string when error occurs")
}

func TestDirectory_CreateSubDir_SuppliedBaseDir_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("logs", "subdir")
	expectedMarkerFileBase := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	markerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	suffix := "1337"
	pattern := "test-pattern-"
	expectedDirPattern := pattern + suffix + "-"
	expectedTempDir := filepath.Join("logs", "subdir", "test-pattern-22222")

	mockConfig.EXPECT().
		BaseDir().
		Return(expectedLogDir).
		Once()

	mockConfig.EXPECT().
		ServerInstanceID().
		Return("").
		Once()

	mockOSLayer.EXPECT().
		MkdirAll(expectedLogDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, expectedMarkerExtension).
		Return(markerFileName, suffix, nil).
		Once()

	mockOSLayer.EXPECT().
		MkdirTemp(expectedLogDir, expectedDirPattern).
		Return(expectedTempDir, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	// Act
	actualTempDir, err := directoryInstance.CreateSubDir(pattern)

	// Assert
	require.NoError(t, err, "MkdirTemp should not return an error")
	assert.Equal(t, expectedTempDir, actualTempDir, "MkdirTemp should return the expected temp directory")
}

func TestDirectory_RecordToLogger_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &directorymocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockFileNameFactory := &directorymocks.MockFilenameFactory{}
	defer mockFileNameFactory.AssertExpectations(t)

	mockOSLayer := &directorymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	expectedLogDir := filepath.Join("tmp", "matlab-mcp-core-server-33333")
	expectedMarkerFileBase := filepath.Join(expectedLogDir, directory.MarkerFileName)
	expectedMarkerExtension := ""
	expectedMarkerFileName := filepath.Join(expectedLogDir, ".matlab-mcp-core-server")
	expectedSuffix := "1337"

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
		Return(expectedLogDir, nil).
		Once()

	mockFileNameFactory.EXPECT().
		CreateFileWithUniqueSuffix(expectedMarkerFileBase, expectedMarkerExtension).
		Return(expectedMarkerFileName, expectedSuffix, nil).
		Once()

	directoryInstance, err := directory.NewDirectory(mockConfig, mockFileNameFactory, mockOSLayer)
	require.NoError(t, err)

	testLogger := testutils.NewInspectableLogger()

	// Act
	directoryInstance.RecordToLogger(testLogger)

	// Assert
	infoLogs := testLogger.InfoLogs()
	require.Len(t, infoLogs, 1)

	fields, found := infoLogs["Application directory state"]
	require.True(t, found, "Expected log message not found")

	actualValue, exists := fields["log-dir"]
	require.True(t, exists, "log-dir field not found in log")
	assert.Equal(t, expectedLogDir, actualValue, "log-dir field has incorrect value")

	actualValue, exists = fields["id"]
	require.True(t, exists, "id field not found in log")
	assert.Equal(t, expectedSuffix, actualValue, "id field has incorrect value")
}
