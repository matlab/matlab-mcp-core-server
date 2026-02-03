// Copyright 2025-2026 The MathWorks, Inc.

package startmatlabsession

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/startmatlabsession"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, request entities.SessionDetails) (startmatlabsession.ReturnArgs, error)
}

type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

func New(
	loggerFactory basetool.LoggerFactory,
	configFactory ConfigFactory,
	usecase Usecase,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(name, title, description, annotations.NewReadOnlyAnnotations(), loggerFactory, Handler(configFactory, usecase)),
	}
}

func Handler(configFactory ConfigFactory, usecase Usecase) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionLogger.Info("Executing Start MATLAB Session tool")
		defer sessionLogger.Info("Done - Executing Start MATLAB Session tool")

		config, messagesErr := configFactory.Config()
		if messagesErr != nil {
			return ReturnArgs{}, messagesErr
		}

		startSessionRequest := entities.LocalSessionDetails{
			MATLABRoot:             inputs.MATLABRoot,
			IsStartingDirectorySet: false,
			ShowMATLABDesktop:      config.ShouldShowMATLABDesktop(),
		}

		response, err := usecase.Execute(ctx, sessionLogger, startSessionRequest)
		if err != nil {
			return ReturnArgs{}, err
		}

		return convertToAnnotatedEquivalentType(response), nil
	}
}

func convertToAnnotatedEquivalentType(response startmatlabsession.ReturnArgs) ReturnArgs {
	return ReturnArgs{
		ResponseText: responseTextIfMATLABSessionStartedSuccesfully,
		SessionID:    int(response.SessionID),
		VerOutput:    response.VerOutput,
		AddOnsOutput: response.AddOnsOutput,
	}
}
