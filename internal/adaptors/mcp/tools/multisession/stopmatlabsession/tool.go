// Copyright 2025-2026 The MathWorks, Inc.

package stopmatlabsession

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-server/internal/entities"
)

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) error
}

type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

func New(
	loggerFactory basetool.LoggerFactory,
	telemetryFactory basetool.TelemetryFactory,
	usecase Usecase,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(name, title, description, annotations.NewDestructiveAnnotations(), loggerFactory, telemetryFactory, Handler(usecase)),
	}
}

func Handler(usecase Usecase) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		err := usecase.Execute(ctx, sessionLogger, entities.SessionID(inputs.SessionID))
		if err != nil {
			return ReturnArgs{}, err
		}

		return ReturnArgs{
			ResponseText: responseTextIfMATLABSessionStoppedSuccessfully,
		}, nil
	}
}
