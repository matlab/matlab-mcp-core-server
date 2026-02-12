// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type MATLABServices interface {
	ListDiscoveredMatlabInfo(logger entities.Logger) datatypes.ListMatlabInfo
	StartLocalMATLABSession(ctx context.Context, logger entities.Logger, request datatypes.LocalSessionDetails) (embeddedconnector.ConnectionDetails, func() error, error)
}

type MATLABSessionStore interface {
	Add(client matlabsessionstore.MATLABSessionClientWithCleanup) entities.SessionID
	Get(sessionID entities.SessionID) (matlabsessionstore.MATLABSessionClientWithCleanup, error)
	Remove(sessionID entities.SessionID)
}

type MATLABSessionClientFactory interface {
	New(endpoint embeddedconnector.ConnectionDetails) (entities.MATLABSessionClient, error)
}

type MATLABManager struct {
	matlabServices MATLABServices
	sessionStore   MATLABSessionStore
	clientFactory  MATLABSessionClientFactory
}

var _ entities.MATLABManager = (*MATLABManager)(nil)

func New(
	matlabServices MATLABServices,
	sessionStore MATLABSessionStore,
	clientFactory MATLABSessionClientFactory,
) *MATLABManager {
	return &MATLABManager{
		matlabServices: matlabServices,
		sessionStore:   sessionStore,
		clientFactory:  clientFactory,
	}
}
