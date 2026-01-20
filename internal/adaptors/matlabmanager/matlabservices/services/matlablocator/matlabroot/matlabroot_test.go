// Copyright 2025-2026 The MathWorks, Inc.

package matlabroot_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabroot"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabroot"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMATLABRootGetter_GetAll_Success(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	// Capture warning log count as they are expected to be created
	var warnLogs = 0
	defer func() {
		assert.Len(t, mockLogger.WarnLogs(), warnLogs, "Should log expected number of warning messages")
	}()

	inputPaths := []string{
		"/valid/path1",
		"/valid/path2",
	}

	pathEnv := strings.Join(inputPaths, string(os.PathListSeparator))

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return(pathEnv).
		Once()

	for _, path := range inputPaths {
		addSuccessfulPathCheck(path, mockOSLayer, mockFileLayer, mockFileInfo)
	}

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	results := service.GetAll(mockLogger)

	// Assert
	// Trim down expected paths to expected length (i.e. back on directory to go from /bin to /matlabroot)
	expectedValidPaths := make([]string, len(inputPaths))
	for i := range inputPaths {
		expectedValidPaths[i] = filepath.Dir(inputPaths[i])
	}

	assert.Len(t, results, len(expectedValidPaths))
	for _, result := range results {
		assert.Contains(t, expectedValidPaths, result)
	}
}

// TestMATLABRootGetter_GetAll_FirstEvalSymlinksError checks that if the first EvalSymlink call errors
// that path is skipped and later paths are still checked
func TestMATLABRootGetter_GetAll_FirstEvalSymlinksError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	// Setup paths
	validPath := filepath.FromSlash("/some/path/valid")
	addSuccessfulPathCheck(validPath, mockOSLayer, mockFileLayer, mockFileInfo)

	errorPath := filepath.FromSlash("/path/with/stat/error")

	mockFileLayer.EXPECT().
		EvalSymlinks(errorPath).
		Return("", assert.AnError).
		Once()

	pathEnv := strings.Join([]string{validPath, errorPath}, string(os.PathListSeparator))

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return(pathEnv).
		Once()

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	result := service.GetAll(mockLogger)

	// Assert
	require.Len(t, result, 1)
	assert.Equal(t, filepath.Dir(validPath), result[0])

	assert.Len(t, mockLogger.WarnLogs(), 1, "An error when evaluating a symlink should trigger a warning log")
}

// TestMATLABRootGetter_GetAll_StatError checks that if the first Stat call errors
// that path is skipped and later paths are still checked
func TestMATLABRootGetter_GetAll_StatError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	// Setup paths
	validPath := filepath.FromSlash("/some/path/valid")
	addSuccessfulPathCheck(validPath, mockOSLayer, mockFileLayer, mockFileInfo)

	errorPath := filepath.FromSlash("/path/with/stat/error")

	mockFileLayer.EXPECT().
		EvalSymlinks(errorPath).
		Return(errorPath, nil).
		Once()

	matlabExePath := filepath.Join(errorPath, config.MATLABExeName)
	mockOSLayer.EXPECT().
		Stat(matlabExePath).
		Return(nil, assert.AnError).
		Once()

	pathEnv := strings.Join([]string{validPath, errorPath}, string(os.PathListSeparator))

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return(pathEnv).
		Once()

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	result := service.GetAll(mockLogger)

	// Assert
	require.Len(t, result, 1)
	assert.Equal(t, filepath.Dir(validPath), result[0])

	assert.Len(t, mockLogger.WarnLogs(), 1, "os.Stat failing should trigger a warning log")
}

// TestMATLABRootGetter_GetAll_FindsDirectoryInsteadOfBinary checks that if the first path check find a directory
// that path is skipped and later paths are still checked
func TestMATLABRootGetter_GetAll_FindsDirectoryInsteadOfBinary(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	// Setup paths
	validPath := filepath.FromSlash("/some/path/valid")
	addSuccessfulPathCheck(validPath, mockOSLayer, mockFileLayer, mockFileInfo)

	errorPath := filepath.FromSlash("/path/with/stat/error")

	mockFileLayer.EXPECT().
		EvalSymlinks(errorPath).
		Return(errorPath, nil).
		Once()

	matlabExePath := filepath.Join(errorPath, config.MATLABExeName)
	mockOSLayer.EXPECT().
		Stat(matlabExePath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(true).
		Once()

	pathEnv := strings.Join([]string{validPath, errorPath}, string(os.PathListSeparator))

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return(pathEnv).
		Once()

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	result := service.GetAll(mockLogger)

	// Assert
	require.Len(t, result, 1)
	assert.Equal(t, filepath.Dir(validPath), result[0])

	//nolint:testifylint // Clearer to use len check for number of errors
	assert.Len(t, mockLogger.WarnLogs(), 0, "No warning logs should be produced by finding a directory")
}

// TestMATLABRootGetter_GetAll_SecondEvalSymlinksError checks that if the second EvalSymlink call errors
// that path is skipped and later paths are still checked
func TestMATLABRootGetter_GetAll_SecondEvalSymlinksError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	// Setup paths
	validPath := filepath.FromSlash("/some/path/valid")
	addSuccessfulPathCheck(validPath, mockOSLayer, mockFileLayer, mockFileInfo)

	errorPath := filepath.FromSlash("/path/with/stat/error")

	mockFileLayer.EXPECT().
		EvalSymlinks(errorPath).
		Return(errorPath, nil).
		Once()

	matlabExePath := filepath.Join(errorPath, config.MATLABExeName)
	mockOSLayer.EXPECT().
		Stat(matlabExePath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(false).
		Once()

	mockFileLayer.EXPECT().
		EvalSymlinks(matlabExePath).
		Return("", assert.AnError).
		Once()

	pathEnv := strings.Join([]string{validPath, errorPath}, string(os.PathListSeparator))

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return(pathEnv).
		Once()

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	result := service.GetAll(mockLogger)

	// Assert
	require.Len(t, result, 1)
	assert.Equal(t, filepath.Dir(validPath), result[0])

	assert.Len(t, mockLogger.WarnLogs(), 1, "An error when evaluating a symlink should trigger a warning log")
}

func TestMATLABRootGetter_GetAll_EmptyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return("").
		Once()

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	result := service.GetAll(mockLogger)

	// Assert
	assert.Nil(t, result, "Should return nil when PATH is empty")
}

func TestMATLABRootGetter_GetAll_FileNotExist(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileLayer := &mocks.MockFileLayer{}
	defer mockFileLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	// Setup paths
	path := filepath.FromSlash("/path/with/no/matlab")
	pathEnv := path

	mockOSLayer.EXPECT().
		Getenv("PATH").
		Return(pathEnv).
		Once()

	mockFileLayer.EXPECT().
		EvalSymlinks(path).
		Return(path, nil).
		Once()

	matlabExePath := filepath.Join(path, config.MATLABExeName)
	mockOSLayer.EXPECT().
		Stat(matlabExePath).
		Return(nil, os.ErrNotExist).
		Once()

	service := matlabroot.New(mockOSLayer, mockFileLayer)

	// Act
	result := service.GetAll(mockLogger)

	// Assert
	assert.Nil(t, result, "Should return nil when no MATLAB executable exists")

	assert.Len(t, mockLogger.WarnLogs(), 0, "No warning logs should be generated for files not existing") //nolint:testifylint // Len check is consistent with other logger checks
}

func addSuccessfulPathCheck(
	path string,
	mockOSLayer *mocks.MockOSLayer,
	mockFileLayer *mocks.MockFileLayer,
	mockFileInfo *osfacademocks.MockFileInfo,
) {
	mockFileLayer.EXPECT().
		EvalSymlinks(path).
		Return(path, nil).
		Once()

	matlabExePath := filepath.Join(path, config.MATLABExeName)

	mockOSLayer.EXPECT().
		Stat(matlabExePath).
		Return(mockFileInfo, nil).
		Once()

	mockFileInfo.EXPECT().
		IsDir().
		Return(false).
		Once()

	mockFileLayer.EXPECT().
		EvalSymlinks(matlabExePath).
		Return(matlabExePath, nil).
		Once()
}
