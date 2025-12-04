// Copyright 2025 The MathWorks, Inc.

package mcpserver_test

import (
	"path/filepath"
	"testing"

	mocks "github.com/matlab/matlab-mcp-core-server/tests/mocks/testutils/mcpserver"
	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetPath_HappyPath(t *testing.T) {
	// Arrange
	mockEnv := mocks.NewMockEnvironment(t)
	defer mockEnv.AssertExpectations(t)

	mockFS := mocks.NewMockFileSystem(t)
	defer mockFS.AssertExpectations(t)

	fakeMCPBaseDir := t.TempDir()
	fakeMCPDir := filepath.Join(fakeMCPBaseDir, testconfig.OSDescriptor)
	fakeFilePath := filepath.Join(fakeMCPDir, testconfig.MATLABMCPCoreServerBinariesFilename)

	mockEnv.EXPECT().
		Getenv("MATLAB_MCP_CORE_SERVER_BUILD_DIR").
		Return(fakeMCPBaseDir)
	mockFS.EXPECT().
		Stat(fakeFilePath).
		Return(nil, nil)

	locator := &mcpserver.Locator{
		Env: mockEnv,
		FS:  mockFS,
	}

	// Act
	path, err := locator.GetPath()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, fakeFilePath, path)
}

func Test_GetPath_EnvVarNotSet(t *testing.T) {
	// Arrange
	mockEnv := mocks.NewMockEnvironment(t)
	defer mockEnv.AssertExpectations(t)

	mockFS := mocks.NewMockFileSystem(t)
	defer mockFS.AssertExpectations(t)

	mockEnv.EXPECT().
		Getenv("MATLAB_MCP_CORE_SERVER_BUILD_DIR").
		Return("")

	locator := &mcpserver.Locator{
		Env: mockEnv,
		FS:  mockFS,
	}

	// Act
	_, err := locator.GetPath()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MATLAB_MCP_CORE_SERVER_BUILD_DIR")
	assert.Contains(t, err.Error(), "is not set")
}

func Test_GetPath_FileDoesNotExist(t *testing.T) {
	// Arrange
	mockEnv := mocks.NewMockEnvironment(t)
	defer mockEnv.AssertExpectations(t)

	mockFS := mocks.NewMockFileSystem(t)
	defer mockFS.AssertExpectations(t)

	fakeMCPBaseDir := t.TempDir()
	fakeMCPDir := filepath.Join(fakeMCPBaseDir, testconfig.OSDescriptor)
	fakeFilePath := filepath.Join(fakeMCPDir, testconfig.MATLABMCPCoreServerBinariesFilename)

	mockEnv.EXPECT().
		Getenv("MATLAB_MCP_CORE_SERVER_BUILD_DIR").
		Return(fakeMCPBaseDir)
	mockFS.EXPECT().
		Stat(fakeFilePath).
		Return(nil, assert.AnError)

	locator := &mcpserver.Locator{
		Env: mockEnv,
		FS:  mockFS,
	}

	// Act
	_, err := locator.GetPath()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	assert.Contains(t, err.Error(), fakeFilePath)
}
