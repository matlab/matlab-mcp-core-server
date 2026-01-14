// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectory_Cleanup_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetCleanupTimeout(100 * time.Millisecond)
	dir.SetCleanupRetry(10 * time.Millisecond)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	// Act
	err := dir.Cleanup()

	// Assert
	require.NoError(t, err)
}

func TestDirectory_Cleanup_WaitsForRemoveAllToPass(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetCleanupTimeout(100 * time.Millisecond)
	dir.SetCleanupRetry(10 * time.Millisecond)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(assert.AnError).
		Once()

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(nil).
		Once()

	// Act
	err := dir.Cleanup()

	// Assert
	require.NoError(t, err)
}

func TestDirectory_Cleanup_Timesout(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetCleanupTimeout(100 * time.Millisecond)
	dir.SetCleanupRetry(10 * time.Millisecond)

	mockOSLayer.EXPECT().
		RemoveAll(sessionDir).
		Return(assert.AnError) // Will be called many times with retry

	// Act
	err := dir.Cleanup()

	// Assert
	require.Error(t, err)
}

func TestDirectory_Cleanup_EmptySessionDir(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	dir := directory.NewDirectory("", mockOSLayer)
	dir.SetCleanupTimeout(100 * time.Millisecond)
	dir.SetCleanupRetry(10 * time.Millisecond)

	// Act
	err := dir.Cleanup()

	// Assert
	require.NoError(t, err)
}
