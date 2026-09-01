// Copyright 2026 The MathWorks, Inc.

package registryvalidate

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
)

//go:embed server.schema.json
var embeddedSchema []byte

const (
	dummyVersion      = "0.0.0"
	dummyFileSHA256   = "0000000000000000000000000000000000000000000000000000000000000000"
	dummyMCPBFilename = "matlab-mcp-server.mcpb"
)

var (
	versionPattern      = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	fileSHA256Pattern   = regexp.MustCompile(`^[a-f0-9]{64}$`)
	mcpbFilenamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
)

type Values struct {
	Version      string
	FileSHA256   string
	MCPBFilename string
}

func DummyValues() Values {
	return Values{Version: dummyVersion, FileSHA256: dummyFileSHA256, MCPBFilename: dummyMCPBFilename}
}

func Validate(templatePath string, values Values) error {
	_, err := renderAndValidate(templatePath, embeddedSchema, "embedded schema", values)
	return err
}

func Render(templatePath, outPath string, values Values) error {
	expanded, err := renderAndValidate(templatePath, embeddedSchema, "embedded schema", values)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outPath, []byte(expanded), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	return nil
}

func renderAndValidate(templatePath string, schemaBytes []byte, schemaName string, values Values) (string, error) {
	templateBytes, err := os.ReadFile(templatePath) //nolint:gosec // template path comes from trusted caller (CLI arg)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", templatePath, err)
	}

	return renderAndValidateContent(string(templateBytes), schemaBytes, schemaName, values)
}

func renderAndValidateContent(templateContent string, schemaBytes []byte, schemaName string, values Values) (string, error) {
	if err := validateValues(values); err != nil {
		return "", err
	}

	expanded, err := expandPlaceholders(templateContent, values)
	if err != nil {
		return "", fmt.Errorf("expand template: %w", err)
	}

	var instance map[string]any
	if err := json.Unmarshal([]byte(expanded), &instance); err != nil {
		return "", fmt.Errorf("expanded template is not valid JSON: %w", err)
	}
	templateSchemaURL, _ := instance["$schema"].(string)

	var schema jsonschema.Schema
	if err := json.Unmarshal(schemaBytes, &schema); err != nil {
		return "", fmt.Errorf("parse %s: %w", schemaName, err)
	}

	if templateSchemaURL != "" && schema.ID != "" && templateSchemaURL != schema.ID {
		return "", fmt.Errorf("schema URL drift: template references %q but %s has $id %q. Refresh the vendored schema or fix the template", templateSchemaURL, schemaName, schema.ID)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return "", fmt.Errorf("resolve schema: %w", err)
	}

	if err := resolved.Validate(instance); err != nil {
		return "", fmt.Errorf("schema validation failed: %w", err)
	}

	return expanded, nil
}

func validateValues(v Values) error {
	if !versionPattern.MatchString(v.Version) {
		return fmt.Errorf("version %q must match N.N.N", v.Version)
	}
	if !fileSHA256Pattern.MatchString(v.FileSHA256) {
		return fmt.Errorf("file-sha256 %q must be 64 lowercase hex characters", v.FileSHA256)
	}
	if !mcpbFilenamePattern.MatchString(v.MCPBFilename) {
		return fmt.Errorf("mcpb filename %q must contain only letters, numbers, dots, underscores, and hyphens", v.MCPBFilename)
	}
	return nil
}

func expandPlaceholders(s string, v Values) (string, error) {
	tmpl, err := template.New("server.json").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, v); err != nil {
		return "", err
	}
	return buf.String(), nil
}
