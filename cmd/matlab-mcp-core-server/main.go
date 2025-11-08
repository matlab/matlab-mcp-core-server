// Copyright 2025 The MathWorks, Inc.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/utils/instancelock"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

func main() {
	// Check for existing instance before doing anything else
	instanceLock, err := instancelock.New()
	if err != nil {
		slog.With("error", err).Error("Failed to create instance lock.")
		os.Exit(1)
	}

	acquired, err := instanceLock.TryLock()
	if err != nil {
		slog.With("error", err).Error("Failed to check for existing instance.")
		os.Exit(1)
	}

	if !acquired {
		// Another instance is already running
		fmt.Fprintf(os.Stderr, "MATLAB MCP Core Server is already running. Only one instance is allowed.\n")
		os.Exit(0)
	}

	// Ensure lock is released on exit
	defer func() {
		if err := instanceLock.Unlock(); err != nil {
			slog.With("error", err).Warn("Failed to release instance lock on exit.")
		}
	}()

	modeSelector, err := wire.InitializeModeSelector()
	if err != nil {
		// As we failed to even initialize, we cannot use a LoggerFactory,
		// and we can't assume whatever failed had a logger factory to log the error either.
		// In this case, we use the default slog.
		slog.With("error", err).Error("Failed to initialize MATLAB MCP Core Server.")
		os.Exit(1)
	}

	ctx := context.Background()
	err = modeSelector.StartAndWaitForCompletion(ctx)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
