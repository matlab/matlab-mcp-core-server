// Copyright 2026 The MathWorks, Inc.

package registryvalidate_test

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/registryvalidate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunRenderHappyPath(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", realSchemaTemplate)
	out := filepath.Join(dir, "server.json")

	// Act
	err := registryvalidate.Run(io.Discard, io.Discard, []string{
		"--render", out,
		"--version", "1.2.3",
		"--file-sha256", validFileSHA256,
		"--mcpb-filename", "matlab-mcp-server.mcpb",
		tmpl,
	})

	// Assert
	require.NoError(t, err)
	rendered, readErr := readFile(out)
	require.NoError(t, readErr)
	assert.Contains(t, rendered, `"version": "1.2.3"`)
	assert.Contains(t, rendered, `"fileSha256": "`+validFileSHA256+`"`)
	assert.Contains(t, rendered, "matlab/matlab-mcp-server/releases/download/v1.2.3/matlab-mcp-server.mcpb")
}

func TestRunValidateOnly(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", realSchemaTemplate)
	var stdout strings.Builder

	// Act
	err := registryvalidate.Run(&stdout, io.Discard, []string{tmpl})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Registry template OK\n", stdout.String())
}

func TestRunRequiresTemplateArg(t *testing.T) {
	// Arrange
	var stderr strings.Builder

	// Act
	err := registryvalidate.Run(io.Discard, &stderr, []string{})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected exactly one template argument")
	assert.Contains(t, stderr.String(), "Usage:")
}
