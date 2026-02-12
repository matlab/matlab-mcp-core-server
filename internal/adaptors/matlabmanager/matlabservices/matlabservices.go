// Copyright 2025-2026 The MathWorks, Inc.

package matlabservices

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type MATLABLocator interface {
	ListDiscoveredMatlabInfo(logger entities.Logger) datatypes.ListMatlabInfo
}

type LocalMATLABSessionLauncher interface {
	StartLocalMATLABSession(ctx context.Context, logger entities.Logger, request datatypes.LocalSessionDetails) (embeddedconnector.ConnectionDetails, func() error, error)
}

type MATLABServices struct {
	MATLABLocator
	LocalMATLABSessionLauncher
}

func New(
	matlabLocator MATLABLocator,
	localMATLABSessionLauncher LocalMATLABSessionLauncher,
) *MATLABServices {
	return &MATLABServices{
		MATLABLocator:              matlabLocator,
		LocalMATLABSessionLauncher: localMATLABSessionLauncher,
	}
}
