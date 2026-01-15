// Copyright 2025-2026 The MathWorks, Inc.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

func main() {
	application := wire.Initialize()

	ctx := context.Background()

	if err := application.ModeSelector.StartAndWaitForCompletion(ctx); err != nil {
		errorMessage, ok := application.MessageCatalog.GetFromGeneralError(err)
		if ok {
			fmt.Fprintf(os.Stderr, "%s\n", errorMessage)
			os.Exit(1)
		}
		os.Exit(1)
	}

	os.Exit(0)
}
