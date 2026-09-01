// Copyright 2026 The MathWorks, Inc.

package registryvalidate_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/registryvalidate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minimalSchema = `{
  "$id": "https://example.test/schemas/v1/server.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["name", "version"],
  "properties": {
    "$schema":    { "type": "string" },
    "name":       { "type": "string" },
    "version":    { "type": "string", "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9.-]+)?$" },
    "identifier": { "type": "string" },
    "fileSha256": { "type": "string" }
  },
  "additionalProperties": false
}`

const goodTemplate = `{
  "$schema": "https://example.test/schemas/v1/server.schema.json",
  "name": "io.example/server",
  "version": "{{.Version}}",
  "identifier": "matlab/matlab-mcp-server/releases/download/v{{.Version}}/{{.MCPBFilename}}",
  "fileSha256": "{{.FileSHA256}}"
}`

const realSchemaTemplate = `{
  "$schema": "https://static.modelcontextprotocol.io/schemas/2025-12-11/server.schema.json",
  "name": "io.github.matlab/matlab-mcp-server",
  "title": "MATLAB MCP Server",
  "description": "Connect AI coding agents to MATLAB. Run code, tests, and analysis via MCP.",
  "version": "{{.Version}}",
  "repository": {
    "url": "https://github.com/matlab/matlab-mcp-server",
    "source": "github"
  },
  "packages": [
    {
      "registryType": "mcpb",
      "registryBaseUrl": "https://github.com",
      "identifier": "matlab/matlab-mcp-server/releases/download/v{{.Version}}/{{.MCPBFilename}}",
      "fileSha256": "{{.FileSHA256}}",
      "version": "{{.Version}}",
      "transport": {
        "type": "stdio"
      }
    }
  ]
}`

const validFileSHA256 = "0000000000000000000000000000000000000000000000000000000000000000"

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path) //nolint:gosec // test controls the path
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func TestValidateHappyPath(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", goodTemplate)

	// Act
	err := registryvalidate.ValidateWithSchema(tmpl, []byte(minimalSchema), registryvalidate.DummyValues())

	// Assert
	assert.NoError(t, err)
}

func TestValidateRealEmbeddedSchema(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", realSchemaTemplate)

	// Act
	err := registryvalidate.Validate(tmpl, registryvalidate.DummyValues())

	// Assert
	assert.NoError(t, err)
}

func TestValidateMissingRequiredField(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", `{
  "$schema": "https://example.test/schemas/v1/server.schema.json",
  "version": "{{.Version}}"
}`)

	// Act
	err := registryvalidate.ValidateWithSchema(tmpl, []byte(minimalSchema), registryvalidate.DummyValues())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "schema validation failed")
}

func TestValidateMalformedJSON(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", `{ "name": "broken", `)

	// Act
	err := registryvalidate.ValidateWithSchema(tmpl, []byte(minimalSchema), registryvalidate.DummyValues())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "valid JSON")
}

func TestValidateSchemaURLDriftFails(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", `{
  "$schema": "https://example.test/schemas/OLD/server.schema.json",
  "name": "io.example/server",
  "version": "{{.Version}}"
}`)

	// Act
	err := registryvalidate.ValidateWithSchema(tmpl, []byte(minimalSchema), registryvalidate.DummyValues())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "schema URL drift")
}

func TestValidateUnknownPlaceholderFails(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	unknownField := "V" + "erison"
	tmpl := writeFile(t, dir, "tmpl.json", `{
  "$schema": "https://example.test/schemas/v1/server.schema.json",
  "name": "io.example/server",
  "version": "{{.`+unknownField+`}}"
}`)

	// Act
	err := registryvalidate.ValidateWithSchema(tmpl, []byte(minimalSchema), registryvalidate.DummyValues())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expand template")
}

func TestRenderWritesPopulatedJSON(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", goodTemplate)
	out := filepath.Join(dir, "rendered.json")
	values := registryvalidate.Values{
		Version:      "1.2.3",
		FileSHA256:   validFileSHA256,
		MCPBFilename: "matlab-mcp-server.mcpb",
	}

	// Act
	err := registryvalidate.RenderWithSchema(tmpl, []byte(minimalSchema), out, values)

	// Assert
	require.NoError(t, err)
	rendered, readErr := readFile(out)
	require.NoError(t, readErr)
	assert.Contains(t, rendered, `"version": "1.2.3"`)
	assert.Contains(t, rendered, values.FileSHA256)
	assert.Contains(t, rendered, "matlab/matlab-mcp-server/releases/download/v1.2.3/matlab-mcp-server.mcpb")
	assert.Contains(t, rendered, `"name": "io.example/server"`)
	assert.NotContains(t, rendered, "{{.Version}}")
}

func TestRenderDoesNotWriteOnValidationFailure(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", goodTemplate)
	out := filepath.Join(dir, "should-not-exist.json")
	values := registryvalidate.Values{
		Version:      "garbage",
		FileSHA256:   validFileSHA256,
		MCPBFilename: "matlab-mcp-server.mcpb",
	}

	// Act
	err := registryvalidate.RenderWithSchema(tmpl, []byte(minimalSchema), out, values)

	// Assert
	require.Error(t, err)
	exists, _ := fileExists(out)
	assert.False(t, exists)
}

func TestValidateRejectsInvalidValues(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	tmpl := writeFile(t, dir, "tmpl.json", goodTemplate)

	tests := []struct {
		name   string
		values registryvalidate.Values
		want   string
	}{
		{
			name: "bad version suffix",
			values: registryvalidate.Values{
				Version:      "1.2.3bad",
				FileSHA256:   validFileSHA256,
				MCPBFilename: "matlab-mcp-server.mcpb",
			},
			want: "version",
		},
		{
			name: "bad sha",
			values: registryvalidate.Values{
				Version:      "1.2.3",
				FileSHA256:   "DEADBEEF",
				MCPBFilename: "matlab-mcp-server.mcpb",
			},
			want: "file-sha256",
		},
		{
			name: "bad filename path traversal",
			values: registryvalidate.Values{
				Version:      "1.2.3",
				FileSHA256:   validFileSHA256,
				MCPBFilename: "../matlab-mcp-server.mcpb",
			},
			want: "mcpb filename",
		},
		{
			name: "bad filename shell character",
			values: registryvalidate.Values{
				Version:      "1.2.3",
				FileSHA256:   validFileSHA256,
				MCPBFilename: "matlab;rm.mcpb",
			},
			want: "mcpb filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := registryvalidate.ValidateWithSchema(tmpl, []byte(minimalSchema), tt.values)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
}
