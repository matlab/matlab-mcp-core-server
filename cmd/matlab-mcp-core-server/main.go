// Copyright 2025 The MathWorks, Inc.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/messagecatalog"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

func main() {
	modeSelector, err := wire.InitializeModeSelector()
	if err != nil {
		globalMessageCatalog := messagecatalog.New()
		errorMessage, ok := globalMessageCatalog.GetFromGeneralError(err)
		if ok {
			fmt.Fprintf(os.Stderr, "%s\n", errorMessage)
			os.Exit(1)
		}

		// As we failed to even initialize, we cannot use a LoggerFactory. Output error to stderr
		fallbackMessage := globalMessageCatalog.Get(messages.StartupErrors_GenericInitializeFailure)
		fmt.Fprintf(os.Stderr, fallbackMessage, err)
		os.Exit(1)
	}

	ctx := context.Background()
	err = modeSelector.StartAndWaitForCompletion(ctx)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
