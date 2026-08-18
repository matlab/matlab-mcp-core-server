// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type cleanupSessionClient struct {
	entities.MATLABSessionClient
	stop func(ctx context.Context, sessionLogger entities.Logger) error
}

func newCleanupSessionClient(client entities.MATLABSessionClient, stop func(ctx context.Context, sessionLogger entities.Logger) error) *cleanupSessionClient {
	return &cleanupSessionClient{
		MATLABSessionClient: client,
		stop:                stop,
	}
}

func (c *cleanupSessionClient) StopSession(ctx context.Context, sessionLogger entities.Logger) error {
	return c.stop(ctx, sessionLogger)
}
