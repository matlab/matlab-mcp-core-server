// Copyright 2026 The MathWorks, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/registryvalidate"
)

func main() {
	if err := registryvalidate.Run(os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "registry-validate: %v\n", err)
		os.Exit(1)
	}
}
