// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabtestfile

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/utils/responseconverter"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/usecases/runmatlabtestfile"
)

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request runmatlabtestfile.Args) (entities.EvalResponse, error)
}

type Tool struct {
	basetool.ToolWithUnstructuredContentOutput[Args]
}

func New(
	loggerFactory basetool.LoggerFactory,
	telemetryFactory basetool.TelemetryFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
) *Tool {
	return &Tool{
		ToolWithUnstructuredContentOutput: basetool.NewToolWithUnstructuredContent(name, title, description, annotations.NewDestructiveAnnotations(), loggerFactory, telemetryFactory, Handler(usecase, globalMATLAB)),
	}
}

func Handler(usecase Usecase, globalMATLAB entities.GlobalMATLAB) basetool.HandlerWithUnstructuredContentOutput[Args] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (tools.RichContent, error) {
		sessionLogger.Info("Executing Run MATLAB Test File tool")
		defer sessionLogger.Info("Done - Executing Run MATLAB Test File tool")

		client, err := globalMATLAB.Client(ctx, sessionLogger)
		if err != nil {
			return tools.RichContent{}, err
		}

		response, err := usecase.Execute(ctx, sessionLogger, client, runmatlabtestfile.Args{
			ScriptPath: inputs.ScriptPath,
		})
		if err != nil {
			return tools.RichContent{}, err
		}

		return responseconverter.ConvertEvalResponseToRichContent(response), nil
	}
}
