// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/asyncrunner"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type greetingSessionClient struct {
	entities.MATLABSessionClient
	greetingTask *asyncrunner.Task
}

func newGreetingSessionClient(underlying entities.MATLABSessionClient, greetingTask *asyncrunner.Task) *greetingSessionClient {
	return &greetingSessionClient{
		MATLABSessionClient: underlying,
		greetingTask:        greetingTask,
	}
}

func (c *greetingSessionClient) awaitGreeting(ctx context.Context) {
	c.greetingTask.Wait(ctx)
}

func (c *greetingSessionClient) Eval(ctx context.Context, logger entities.Logger, req entities.EvalRequest) (entities.EvalResponse, error) {
	c.awaitGreeting(ctx)
	return c.MATLABSessionClient.Eval(ctx, logger, req)
}

func (c *greetingSessionClient) EvalWithCapture(ctx context.Context, logger entities.Logger, req entities.EvalRequest) (entities.EvalResponse, error) {
	c.awaitGreeting(ctx)
	return c.MATLABSessionClient.EvalWithCapture(ctx, logger, req)
}

func (c *greetingSessionClient) FEval(ctx context.Context, logger entities.Logger, req entities.FEvalRequest) (entities.FEvalResponse, error) {
	c.awaitGreeting(ctx)
	return c.MATLABSessionClient.FEval(ctx, logger, req)
}
