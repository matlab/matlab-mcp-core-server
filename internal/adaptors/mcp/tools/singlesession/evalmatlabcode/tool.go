// Copyright 2025-2026 The MathWorks, Inc.

package evalmatlabcode

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/utils/responseconverter"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalmatlabcode"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request evalmatlabcode.Args) (entities.EvalResponse, error)
}

type Tool struct {
	basetool.ToolWithUnstructuredContentOutput[Args]
}

func New(
	loggerFactory basetool.LoggerFactory,
	configFactory ConfigFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
) *Tool {
	return &Tool{
		ToolWithUnstructuredContentOutput: basetool.NewToolWithUnstructuredContent(name, title, description, annotations.NewDestructiveAnnotations(), loggerFactory, Handler(configFactory, usecase, globalMATLAB)),
	}
}

func Handler(configFactory ConfigFactory, usecase Usecase, globalMATLAB entities.GlobalMATLAB) basetool.HandlerWithUnstructuredContentOutput[Args] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (tools.RichContent, error) {
		sessionLogger.Info("Executing Eval tool")
		defer sessionLogger.Info("Done - Executing Eval tool")

		config, messagesErr := configFactory.Config()
		if messagesErr != nil {
			return tools.RichContent{}, messagesErr
		}

		client, err := globalMATLAB.Client(ctx, sessionLogger)
		if err != nil {
			return tools.RichContent{}, err
		}

		response, err := usecase.Execute(ctx, sessionLogger, client, evalmatlabcode.Args{
			Code:          inputs.Code,
			ProjectPath:   inputs.ProjectPath,
			CaptureOutput: !config.ShouldShowMATLABDesktop(),
		})
		if err != nil {
			return tools.RichContent{}, err
		}

		return responseconverter.ConvertEvalResponseToRichContent(response), nil
	}
}
