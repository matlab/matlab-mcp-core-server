// Copyright 2025 The MathWorks, Inc.

package pathcontrol

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/tests/testconfig"
)

func RemoveFromPath(currentPath string, pathsToRemove []string) string {
	pathParts := filepath.SplitList(currentPath)
	var newPathParts []string

	for _, part := range pathParts {
		if !slices.Contains(pathsToRemove, part) {
			newPathParts = append(newPathParts, part)
		}
	}

	return strings.Join(newPathParts, string(os.PathListSeparator))
}

func AddToPath(currentPath string, pathsToAdd []string) string {
	newPath := currentPath
	for _, path := range pathsToAdd {
		if newPath == "" {
			newPath = path
		} else {
			newPath = newPath + string(os.PathListSeparator) + path
		}
	}
	return newPath
}

// RemoveAllMATLABsFromPath removes any directory containing a MATLAB executable from the path string
func RemoveAllMATLABsFromPath(currentPath string) string {
	pathParts := filepath.SplitList(currentPath)
	var newPathParts []string

	for _, part := range pathParts {
		matlabExePath := filepath.Join(part, testconfig.MATLABExeName)
		// Check if MATLAB executable exists in this directory
		if _, err := os.Stat(matlabExePath); err == nil {
			continue
		}
		newPathParts = append(newPathParts, part)
	}

	return strings.Join(newPathParts, string(os.PathListSeparator))
}

func UpdateEnvEntry(env []string, key string, value string) []string {
	var newEnv []string
	keyFound := false
	prefix := key + "="

	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			newEnv = append(newEnv, prefix+value)
			keyFound = true
		} else {
			newEnv = append(newEnv, e)
		}
	}

	if !keyFound {
		newEnv = append(newEnv, prefix+value)
	}

	return newEnv
}
