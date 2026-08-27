// Copyright 2025-2026 The MathWorks, Inc.

package listavailablematlabs

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/usecases/listavailablematlabs"
)

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger) listavailablematlabs.ReturnArgs
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
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(name, title, description, annotations.NewReadOnlyAnnotations(), loggerFactory, telemetryFactory, Handler(usecase)),
	}
}

func Handler(usecase Usecase) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionLogger.Info("Executing list available MATLABs tool")
		defer sessionLogger.Info("Done - Executing list available MATLABs tool")

		environments := usecase.Execute(ctx, sessionLogger)

		return convertToAnnotatedEquivalentType(environments), nil
	}
}

func convertToAnnotatedEquivalentType(environmentInfos listavailablematlabs.ReturnArgs) ReturnArgs {
	convertedEnvironmentInfos := make([]EnvironmentInfo, len(environmentInfos))
	for i, env := range environmentInfos {
		convertedEnvironmentInfos[i] = EnvironmentInfo{
			Version:    env.Version,
			MATLABRoot: env.MATLABRoot,
		}
	}
	return ReturnArgs{
		AvailableMATLABs: convertedEnvironmentInfos,
	}
}
