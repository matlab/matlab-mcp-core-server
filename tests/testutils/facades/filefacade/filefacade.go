// Copyright 2025 The MathWorks, Inc.

package filefacade

import (
	"os"
	"path/filepath"
)

// RealFileSystem implements FileSystem using the os and filepath packages
type RealFileSystem struct{}

func (RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (RealFileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
