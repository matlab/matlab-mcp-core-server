// Copyright 2026 The MathWorks, Inc.

package testdata

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed extensions/magic_square.json
var magicSquareExtension string

func WriteMagicSquareExtension(t testing.TB) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "magic_square.json")
	require.NoError(t, os.WriteFile(path, []byte(magicSquareExtension), 0o600))
	return path
}
