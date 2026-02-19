// Copyright 2025-2026 The MathWorks, Inc.

package checkmatlabcode

import (
	"context"
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

// CodeIssue represents a single code issue found by the code analysis
type CodeIssue struct {
	Description string
	Line        int
	StartColumn int
	EndColumn   int
	Severity    string
	Fixable     bool
}

type Args struct {
	ScriptPath string
}

type ReturnArgs struct {
	CodeIssues []CodeIssue
}

type PathValidator interface {
	ValidateMATLABScript(filePath string) (string, error)
}

type CodeAnalyzer interface {
	AnalyzeCode(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient, scriptPath string) ([]CodeIssue, error)
}

type Usecase struct {
	pathValidator PathValidator
	codeAnalyzer  CodeAnalyzer
}

func New(
	pathValidator PathValidator,
	codeAnalyzer CodeAnalyzer,
) *Usecase {
	return &Usecase{
		pathValidator: pathValidator,
		codeAnalyzer:  codeAnalyzer,
	}
}

func (u *Usecase) Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request Args) (ReturnArgs, error) {
	sessionLogger.Debug("Entering CheckMATLABCode Usecase")
	defer sessionLogger.Debug("Exiting CheckMATLABCode Usecase")

	validatedPath, err := u.pathValidator.ValidateMATLABScript(request.ScriptPath)
	if err != nil {
		return ReturnArgs{}, fmt.Errorf("path validation failed: %w", err)
	}

	issues, err := u.codeAnalyzer.AnalyzeCode(ctx, sessionLogger, client, validatedPath)
	if err != nil {
		return ReturnArgs{}, err
	}

	return ReturnArgs{CodeIssues: issues}, nil
}
