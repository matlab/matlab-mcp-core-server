// Copyright 2026 The MathWorks, Inc.

package registryvalidate

import (
	"fmt"
	"os"
)

func ValidateWithSchema(templatePath string, schemaBytes []byte, values Values) error {
	_, err := renderAndValidate(templatePath, schemaBytes, "schema fixture", values)
	return err
}

func RenderWithSchema(templatePath string, schemaBytes []byte, outPath string, values Values) error {
	expanded, err := renderAndValidate(templatePath, schemaBytes, "schema fixture", values)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outPath, []byte(expanded), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	return nil
}
