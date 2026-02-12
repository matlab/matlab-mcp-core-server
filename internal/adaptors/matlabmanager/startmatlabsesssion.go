// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

func (m *MATLABManager) StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error) {
	var zeroValue entities.SessionID
	var client matlabsessionstore.MATLABSessionClientWithCleanup

	switch request := startRequest.(type) {
	case entities.LocalSessionDetails:
		localSessionLogger := sessionLogger.With("matlab-root", request.MATLABRoot)
		// For now, we return embedded connector details, to decouple the session start logic from the client creation.
		embeddedConnectorEndpoint, sessionCleanup, err := m.matlabServices.StartLocalMATLABSession(
			ctx,
			localSessionLogger,
			datatypes.LocalSessionDetails{
				MATLABRoot:             request.MATLABRoot,
				IsStartingDirectorySet: request.IsStartingDirectorySet,
				StartingDirectory:      request.StartingDirectory,
				ShowMATLABDesktop:      request.ShowMATLABDesktop,
			},
		)
		if err != nil {
			return zeroValue, err
		}
		embeddedConnectorClient, err := m.clientFactory.New(embeddedConnectorEndpoint)
		if err != nil {
			return zeroValue, err
		}
		client = newMATLABSessionClientWithCleanup(embeddedConnectorClient, sessionCleanup)
	default:
		return zeroValue, fmt.Errorf("unknown request type: %T", request)
	}

	return m.sessionStore.Add(client), nil
}
