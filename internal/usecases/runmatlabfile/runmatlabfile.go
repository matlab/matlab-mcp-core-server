// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabfile

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/utils/pathextractor"
)

type Args struct {
	ScriptPath    string
	CaptureOutput bool
}

type PathValidator interface {
	ValidateMATLABScript(filePath string) (string, error)
}

type Usecase struct {
	pathValidator PathValidator
}

func New(
	pathValidator PathValidator,
) *Usecase {
	return &Usecase{
		pathValidator: pathValidator,
	}
}

func (u *Usecase) Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request Args) (entities.EvalResponse, error) {
	sessionLogger.Debug("Entering RunMATLABFile Usecase")
	defer sessionLogger.Debug("Exiting RunMATLABFile Usecase")

	validatedPath, err := u.pathValidator.ValidateMATLABScript(request.ScriptPath)
	if err != nil {
		return entities.EvalResponse{}, err
	}

	scriptDir, scriptName := pathextractor.ExtractPathComponents(validatedPath)

	_, err = client.Eval(ctx, sessionLogger, entities.EvalRequest{
		Code: fmt.Sprintf("cd('%s')", scriptDir),
	})
	if err != nil {
		return entities.EvalResponse{}, err
	}

	runCodeRequest := entities.EvalRequest{
		Code: scriptName,
	}

	if request.CaptureOutput {
		return client.EvalWithCapture(ctx, sessionLogger, runCodeRequest)
	}
	return client.Eval(ctx, sessionLogger, runCodeRequest)
}
