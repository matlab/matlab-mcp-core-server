// Copyright 2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/asyncrunner"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

func (m *MATLABManager) SetMATLABSessionConnectionRetryInterval(matlabSessionConnectionRetryInterval time.Duration) {
	m.matlabSessionConnectionRetryInterval = matlabSessionConnectionRetryInterval
}

func NewGreetingSessionClient(client entities.MATLABSessionClient, greetingTask *asyncrunner.Task) *greetingSessionClient {
	return newGreetingSessionClient(client, greetingTask)
}

func NewCleanupSessionClient(client entities.MATLABSessionClient, stop func(ctx context.Context, sessionLogger entities.Logger) error) *cleanupSessionClient {
	return newCleanupSessionClient(client, stop)
}
