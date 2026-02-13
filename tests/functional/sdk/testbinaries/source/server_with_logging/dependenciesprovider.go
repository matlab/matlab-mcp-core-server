// Copyright 2026 The MathWorks, Inc.

package main

import (
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/logger"
)

type Dependencies struct{}

type DependenciesProviderResources interface { //nolint:iface // Same interface is happenstance
	Logger() logger.Logger
}

func DependenciesProvider(resources DependenciesProviderResources) (Dependencies, i18n.Error) {
	resources.Logger().Info("Creating Dependencies")

	return Dependencies{}, nil
}
