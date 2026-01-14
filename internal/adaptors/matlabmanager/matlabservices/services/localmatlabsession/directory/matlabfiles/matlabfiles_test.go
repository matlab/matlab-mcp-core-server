// Copyright 2025-2026 The MathWorks, Inc.

package matlabfiles_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory/matlabfiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed assets/+matlab_mcp
var expectedMATLABFiles embed.FS

func TestMATLABFiles_GetAll_HappyPath(t *testing.T) {
	// Arrange
	matlabFiles := matlabfiles.New()
	subFS, err := fs.Sub(expectedMATLABFiles, "assets/+matlab_mcp")
	require.NoError(t, err)

	// Act
	files := matlabFiles.GetAll()

	// Assert
	for fileName, fileContent := range files {
		expectedFileContent, err := fs.ReadFile(subFS, fileName)
		require.NoError(t, err)
		assert.Equal(t, expectedFileContent, fileContent)
	}
}

func TestMATLABFiles_GetAll_ReturnAllFiles(t *testing.T) {
	// Arrange
	matlabFiles := matlabfiles.New()
	subFS, err := fs.Sub(expectedMATLABFiles, "assets/+matlab_mcp")
	require.NoError(t, err)

	// Act
	files := matlabFiles.GetAll()

	// Assert
	entries, err := fs.ReadDir(subFS, ".")
	require.NoError(t, err)
	for _, entry := range entries {
		assert.Contains(t, files, entry.Name())
	}
}
