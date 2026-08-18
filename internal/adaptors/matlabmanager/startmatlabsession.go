// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/asyncrunner"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

var ErrMATLABSessionNotAlive = errors.New("session is not alive")

const greetingTimeout = 3 * time.Second

const titleTimeout = 20 * time.Second

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
			if cleanupErr := sessionCleanup(); cleanupErr != nil {
				sessionLogger.WithError(cleanupErr).Error("Failed to clean up session after client factory error")
			}
			return zeroValue, err
		}
		gateClient, _, _ := m.indicateConnection(sessionLogger, embeddedConnectorClient, request.ShowMATLABDesktop)
		client = newCleanupSessionClient(gateClient, func(ctx context.Context, sessionLogger entities.Logger) error {
			if _, err := embeddedConnectorClient.Eval(ctx, sessionLogger, entities.EvalRequest{Code: "exit()"}); err != nil {
				return err
			}
			return sessionCleanup()
		})
	case entities.AttachToExistingSession:
		sessionLogger.Info("Attaching to existing session")

		connectionDetails, err := m.sessionSelector.SelectSessionToAttachTo(sessionLogger)
		if err != nil {
			return zeroValue, err
		}

		embeddedConnectorClient, err := m.clientFactory.New(connectionDetails)
		if err != nil {
			return zeroValue, err
		}

		response := embeddedConnectorClient.Ping(ctx, sessionLogger)
		if !response.IsAlive {
			return zeroValue, ErrMATLABSessionNotAlive
		}

		const showMATLABDesktop = true
		gateClient, titleTask, originalTitle := m.indicateConnection(sessionLogger, embeddedConnectorClient, showMATLABDesktop)
		client = newCleanupSessionClient(gateClient, func(ctx context.Context, sessionLogger entities.Logger) error {
			if titleTask.Wait(ctx) && *originalTitle != "" {
				if err := m.connectionIndicator.RestoreTitle(ctx, sessionLogger, embeddedConnectorClient, *originalTitle); err != nil {
					sessionLogger.WithError(err).Warn("failed to restore MATLAB desktop title on detach")
				}
			}
			sessionLogger.Debug("Skipping session stop for externally managed MATLAB session")
			return nil
		})
	default:
		return zeroValue, fmt.Errorf("unknown request type: %T", request)
	}

	sessionID := m.sessionStore.Add(client)

	return sessionID, nil
}

func (m *MATLABManager) indicateConnection(sessionLogger entities.Logger, client entities.MATLABSessionClient, showMATLABDesktop bool) (*greetingSessionClient, *asyncrunner.Task, *string) {
	greetingInfo := m.greetingInfo(showMATLABDesktop)

	greetingTask := asyncrunner.Run(context.Background(), sessionLogger, greetingTimeout,
		"failed to display MCP connection greeting in MATLAB command window",
		func(ctx context.Context) error {
			return m.connectionIndicator.ShowGreeting(ctx, sessionLogger, client, greetingInfo)
		})

	var originalTitle string
	titleTask := asyncrunner.Run(context.Background(), sessionLogger, titleTimeout,
		"failed to set MATLAB desktop title on connect",
		func(ctx context.Context) error {
			title, err := m.connectionIndicator.ApplyConnectedTitle(ctx, sessionLogger, client, greetingInfo)
			if err != nil {
				return err
			}
			originalTitle = title
			return nil
		})

	return newGreetingSessionClient(client, greetingTask), titleTask, &originalTitle
}

func (m *MATLABManager) greetingInfo(showMATLABDesktop bool) entities.GreetingInfo {
	return entities.GreetingInfo{
		ClientInfo:        m.clientInfoProvider.GetClientInfo(),
		ShowMATLABDesktop: showMATLABDesktop,
	}
}
