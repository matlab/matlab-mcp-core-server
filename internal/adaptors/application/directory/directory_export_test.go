// Copyright 2025-2026 The MathWorks, Inc.

package directory

import "github.com/matlab/matlab-mcp-core-server/internal/messages"

const (
	DefaultLogDirPattern = defaultLogDirPattern
	MarkerFileName       = markerFileName
)

func NewDirectory(
	config Config,
	filenameFactory FilenameFactory,
	osFacade OSLayer,
) (Directory, messages.Error) {
	return newDirectory(
		config,
		filenameFactory,
		osFacade,
	)
}
