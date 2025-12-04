// Copyright 2025 The MathWorks, Inc.

package fakematlab

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
	"github.com/stretchr/testify/require"
)

// Create creates a fake MATLAB executable in a temp directory
// Returns the directory containing the executable and the full path to the executable
func Create(t *testing.T) (string, string) {
	tempDir := t.TempDir()
	matlabDir := filepath.Join(tempDir, "bin")
	err := os.MkdirAll(matlabDir, 0o700)
	require.NoError(t, err)

	matlabPath := filepath.Join(matlabDir, testconfig.MATLABExeName)
	err = os.WriteFile(matlabPath, []byte("fake matlab"), 0o700) //nolint:gosec // Create fake executable file
	require.NoError(t, err)

	return matlabDir, matlabPath
}
