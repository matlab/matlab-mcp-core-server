// Copyright 2025-2026 The MathWorks, Inc.

package codeanalyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
)

const (
	codeIssuesMethodName    = "codeIssues"
	checkCodeMethodName     = "checkcode"
	minVersionForCodeIssues = "R2022b"
)

// matlabIssue represents a single issue from codeIssues function
type matlabIssue struct {
	Description string `json:"Description"`
	LineStart   int    `json:"LineStart"`
	ColumnStart int    `json:"ColumnStart"`
	ColumnEnd   int    `json:"ColumnEnd"`
	Severity    string `json:"Severity"`
	Fixability  string `json:"Fixability"`
}

// matlabCodeIssuesResponse represents the full response from codeIssues function
type matlabCodeIssuesResponse struct {
	Date                      string        `json:"Date"`
	Release                   string        `json:"Release"`
	Files                     string        `json:"Files"`
	CodeAnalyzerConfiguration string        `json:"CodeAnalyzerConfiguration"`
	Issues                    []matlabIssue `json:"Issues"`
}

// matlabCheckcodeItem represents a single item from checkcode function
type matlabCheckcodeItem struct {
	Message string `json:"message"`
	Fix     int    `json:"fix"`
	Line    int    `json:"line"`
	Column  []int  `json:"column"`
}

// Analyzer provides MATLAB code analysis capabilities.
type Analyzer struct{}

// New creates a new Analyzer instance.
func New() *Analyzer {
	return &Analyzer{}
}

// AnalyzeCode runs MATLAB code analysis on the given script and returns any issues found.
func (a *Analyzer) AnalyzeCode(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient, scriptPath string) ([]checkmatlabcode.CodeIssue, error) {
	method := selectCodeCheckMethod(ctx, logger, client)

	jsonOutput, err := runCodeCheck(ctx, logger, client, method, scriptPath)
	if err != nil {
		return nil, err
	}

	return parseAnalysisOutput(method, jsonOutput)
}

// selectCodeCheckMethod determines which MATLAB function to use based on version
func selectCodeCheckMethod(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) string {
	response, err := client.FEval(ctx, logger, entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{minVersionForCodeIssues},
		NumOutputs: 1,
	})

	if err != nil {
		logger.WithError(err).Warn("Failed to check MATLAB version, defaulting to checkcode method")
		return checkCodeMethodName
	}

	if len(response.Outputs) == 0 {
		logger.Warn("MATLAB version check returned no outputs, defaulting to checkcode method")
		return checkCodeMethodName
	}

	isOlderVersion, ok := response.Outputs[0].(bool)
	if !ok {
		logger.Warn("MATLAB version check returned non-bool, defaulting to checkcode method")
		return checkCodeMethodName
	}

	if isOlderVersion {
		return checkCodeMethodName
	}

	return codeIssuesMethodName
}

// runCodeCheck executes the MATLAB code analysis and returns JSON output
func runCodeCheck(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient, method string, scriptPath string) (string, error) {
	escapedPath := strings.ReplaceAll(scriptPath, "'", "''")
	matlabExpression := fmt.Sprintf("disp(jsonencode(%s('%s')))", method, escapedPath)

	response, err := client.EvalWithCapture(ctx, logger, entities.EvalRequest{
		Code: matlabExpression,
	})
	if err != nil {
		return "", err
	}
	return response.ConsoleOutput, nil
}

// parseAnalysisOutput routes to the appropriate parser based on method
func parseAnalysisOutput(method string, jsonOutput string) ([]checkmatlabcode.CodeIssue, error) {
	if method == codeIssuesMethodName {
		return parseCodeIssuesResponse(jsonOutput)
	}
	return parseCheckcodeResponse(jsonOutput)
}

// unmarshalJSON attempts to unmarshal JSON and returns error on failure
func unmarshalJSON(jsonOutput string, target any, methodName string) error {
	if err := json.Unmarshal([]byte(jsonOutput), target); err != nil {
		return fmt.Errorf("failed to parse %s output: %w", methodName, err)
	}
	return nil
}

// parseCodeIssuesResponse processes output from the codeIssues function (R2022b+)
func parseCodeIssuesResponse(jsonOutput string) ([]checkmatlabcode.CodeIssue, error) {
	var response matlabCodeIssuesResponse
	if err := unmarshalJSON(jsonOutput, &response, codeIssuesMethodName); err != nil {
		return nil, err
	}

	issues := make([]checkmatlabcode.CodeIssue, 0, len(response.Issues))
	for _, matlabIssue := range response.Issues {
		issue := checkmatlabcode.CodeIssue{
			Description: matlabIssue.Description,
			Line:        matlabIssue.LineStart,
			StartColumn: matlabIssue.ColumnStart,
			EndColumn:   matlabIssue.ColumnEnd,
			Severity:    matlabIssue.Severity,
			Fixable:     strings.EqualFold(matlabIssue.Fixability, "auto"),
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

// parseCheckcodeResponse processes output from the checkcode function (legacy)
func parseCheckcodeResponse(jsonOutput string) ([]checkmatlabcode.CodeIssue, error) {
	var items []matlabCheckcodeItem
	if err := unmarshalJSON(jsonOutput, &items, checkCodeMethodName); err != nil {
		return nil, err
	}

	issues := make([]checkmatlabcode.CodeIssue, 0, len(items))
	for _, item := range items {
		startColumn, endColumn := extractColumnRange(item.Column)

		issue := checkmatlabcode.CodeIssue{
			Description: item.Message,
			Line:        item.Line,
			StartColumn: startColumn,
			EndColumn:   endColumn,
			Severity:    "unknown",
			Fixable:     item.Fix == 1,
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

func extractColumnRange(columns []int) (startColumn int, endColumn int) {
	startColumn, endColumn = 1, 1

	switch len(columns) {
	case 1:
		startColumn, endColumn = columns[0], columns[0]
	case 2:
		startColumn, endColumn = columns[0], columns[1]
	}

	return startColumn, endColumn
}
