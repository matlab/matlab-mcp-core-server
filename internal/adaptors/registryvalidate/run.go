// Copyright 2026 The MathWorks, Inc.

package registryvalidate

import (
	"flag"
	"fmt"
	"io"
)

func Run(stdout, stderr io.Writer, args []string) error {
	fs := flag.NewFlagSet("registry-validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage:\n  %s <template>\n  %s --render <path> --version <v> --file-sha256 <hex> --mcpb-filename <name> <template>\n\nFlags:\n", fs.Name(), fs.Name()) //nolint:errcheck // stderr write; nothing actionable on failure
		fs.PrintDefaults()
	}
	renderOut := fs.String("render", "", "if set, render the populated JSON to this path instead of validate-only")
	version := fs.String("version", "", "value for {{.Version}} (used with --render)")
	fileSHA256 := fs.String("file-sha256", "", "value for {{.FileSHA256}} (used with --render)")
	mcpbFilename := fs.String("mcpb-filename", "", "value for {{.MCPBFilename}} (used with --render)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("expected exactly one template argument, got %d", fs.NArg())
	}
	templatePath := fs.Arg(0)

	if *renderOut == "" {
		if err := Validate(templatePath, DummyValues()); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "Registry template OK") //nolint:errcheck // stdout write; nothing actionable on failure
		return nil
	}

	values := Values{Version: *version, FileSHA256: *fileSHA256, MCPBFilename: *mcpbFilename}
	if err := Render(templatePath, *renderOut, values); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Wrote %s\n", *renderOut) //nolint:errcheck // stdout write; nothing actionable on failure
	return nil
}
