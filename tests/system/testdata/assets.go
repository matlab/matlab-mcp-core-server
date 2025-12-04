// Copyright 2025 The MathWorks, Inc.

package testdata

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// MATLABFiles holds the embedded test assets
//
//go:embed matlab/*
var matlabFiles embed.FS

const matlabFilesRoot = "matlab"

// CopyToDir extracts the embedded test assets to the specified directory.
// The targetDir will be created if it does not exist.
func CopyToDir(targetDir string) error {
	if err := os.MkdirAll(targetDir, 0o0750); err != nil {
		return err
	}

	return fs.WalkDir(
		matlabFiles,
		matlabFilesRoot,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(matlabFilesRoot, path)
			if err != nil {
				return err
			}

			targetPath := filepath.Join(targetDir, relPath)

			if d.IsDir() {
				return os.MkdirAll(targetPath, 0o0750)
			}

			content, err := matlabFiles.ReadFile(path)
			if err != nil {
				return err
			}

			return os.WriteFile(targetPath, content, 0o0600)
		},
	)
}
