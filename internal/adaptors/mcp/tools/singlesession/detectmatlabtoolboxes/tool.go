// Copyright 2025-2026 The MathWorks, Inc.

package detectmatlabtoolboxes

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/usecases/detectmatlabtoolboxes"
)

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient) (detectmatlabtoolboxes.ReturnArgs, error)
}

type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

func New(
	loggerFactory basetool.LoggerFactory,
	telemetryFactory basetool.TelemetryFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(name, title, description, annotations.NewReadOnlyAnnotations(), loggerFactory, telemetryFactory, Handler(usecase, globalMATLAB)),
	}
}

func (Tool) Name() string {
	return name
}

func (Tool) Description() string {
	return description
}

func Handler(usecase Usecase, globalMATLAB entities.GlobalMATLAB) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionLogger.Info("Executing detect MATLAB toolboxes tool")
		defer sessionLogger.Info("Done - Executing detect MATLAB toolboxes tool")

		client, err := globalMATLAB.Client(ctx, sessionLogger)
		if err != nil {
			return ReturnArgs{}, err
		}

		tbxInfo, err := usecase.Execute(ctx, sessionLogger, client)

		if err != nil {
			return ReturnArgs{}, err
		}

		return ReturnArgs{
			InstallationInfo: tbxInfo.Toolboxes,
		}, nil
	}
}
