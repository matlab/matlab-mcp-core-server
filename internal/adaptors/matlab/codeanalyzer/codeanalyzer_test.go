// Copyright 2025-2026 The MathWorks, Inc.

package codeanalyzer_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlab/codeanalyzer"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_AnalyzeCode_CheckcodeMethod(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	checkcodeJSON := `[{"message":"Variable 'x' might be unused.","fix":0,"line":5,"column":[1,10]}]`
	expectedCodeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Variable 'x' might be unused.",
			Line:        5,
			StartColumn: 1,
			EndColumn:   10,
			Severity:    "unknown",
			Fixable:     false,
		},
	}

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Equal(t, expectedCodeIssues, issues, "CodeIssues should match expected value")
}

func TestAnalyzer_AnalyzeCode_CodeIssuesMethod(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{false}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(codeIssues('" + scriptPath + "')))",
	}
	codeIssuesJSON := `{"Date":"2025-01-15","Release":"R2024b","Files":"script.m","CodeAnalyzerConfiguration":"active","Issues":[{"Description":"Variable 'x' might be unused.","LineStart":5,"ColumnStart":1,"ColumnEnd":10,"Severity":"warning","Fixability":"auto"}]}`
	expectedCodeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Variable 'x' might be unused.",
			Line:        5,
			StartColumn: 1,
			EndColumn:   10,
			Severity:    "warning",
			Fixable:     true,
		},
	}

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: codeIssuesJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Equal(t, expectedCodeIssues, issues, "CodeIssues should match expected value")
}

func TestAnalyzer_AnalyzeCode_EmptyOutput(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	emptyCheckcodeJSON := `[]`

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: emptyCheckcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Empty(t, issues, "CodeIssues should be empty")
}

func TestAnalyzer_AnalyzeCode_PathWithSingleQuotes(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("path", "with'quote", "script.m")
	escapedPath := strings.ReplaceAll(scriptPath, "'", "''")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + escapedPath + "')))",
	}
	checkcodeJSON := `[{"message":"Variable 'x' might be unused.","fix":1,"line":5,"column":[1,10]}]`
	expectedCodeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Variable 'x' might be unused.",
			Line:        5,
			StartColumn: 1,
			EndColumn:   10,
			Severity:    "unknown",
			Fixable:     true,
		},
	}

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Equal(t, expectedCodeIssues, issues, "CodeIssues should match expected value")
}

func TestAnalyzer_AnalyzeCode_VersionCheckError_DefaultsToCheckcode(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")
	versionCheckError := assert.AnError

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	checkcodeJSON := `[]`

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(entities.FEvalResponse{}, versionCheckError).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Empty(t, issues, "CodeIssues should be empty")
}

func TestAnalyzer_AnalyzeCode_EvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")
	evalError := assert.AnError

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, evalError).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.Error(t, err, "AnalyzeCode should return an error")
	assert.Empty(t, issues, "CodeIssues should be empty when there's an error")
}

func TestAnalyzer_AnalyzeCode_InvalidJSON(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	invalidJSON := `{invalid json}`

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: invalidJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.Error(t, err, "AnalyzeCode should return an error for invalid JSON")
	assert.Empty(t, issues, "CodeIssues should be empty when there's an error")
}

func TestAnalyzer_AnalyzeCode_InvalidJSON_CodeIssuesMethod(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{false}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(codeIssues('" + scriptPath + "')))",
	}
	invalidJSON := `{invalid json}`

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: invalidJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.Error(t, err, "AnalyzeCode should return an error for invalid JSON")
	assert.Empty(t, issues, "CodeIssues should be empty when there's an error")
}

func TestAnalyzer_AnalyzeCode_VersionCheckEmptyOutputs_DefaultsToCheckcode(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	checkcodeJSON := `[]`

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Empty(t, issues, "CodeIssues should be empty")
}

func TestAnalyzer_AnalyzeCode_VersionCheckNonBool_DefaultsToCheckcode(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{"not a bool"}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	checkcodeJSON := `[]`

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Empty(t, issues, "CodeIssues should be empty")
}

func TestAnalyzer_AnalyzeCode_CheckcodeWithSingleColumnElement(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	checkcodeJSON := `[{"message":"Variable 'x' might be unused.","fix":0,"line":5,"column":[7]}]`
	expectedCodeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Variable 'x' might be unused.",
			Line:        5,
			StartColumn: 7,
			EndColumn:   7,
			Severity:    "unknown",
			Fixable:     false,
		},
	}

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Equal(t, expectedCodeIssues, issues, "CodeIssues should match expected value")
}

func TestAnalyzer_AnalyzeCode_CheckcodeWithEmptyColumnArray(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	scriptPath := filepath.Join("validated", "path", "to", "script.m")

	expectedVersionCheckRequest := entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2022b"},
		NumOutputs: 1,
	}
	versionCheckResponse := entities.FEvalResponse{Outputs: []any{true}}
	expectedEvalRequest := entities.EvalRequest{
		Code: "disp(jsonencode(checkcode('" + scriptPath + "')))",
	}
	checkcodeJSON := `[{"message":"Variable 'x' might be unused.","fix":0,"line":5,"column":[]}]`
	expectedCodeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Variable 'x' might be unused.",
			Line:        5,
			StartColumn: 1,
			EndColumn:   1,
			Severity:    "unknown",
			Fixable:     false,
		},
	}

	mockClient.EXPECT().
		FEval(t.Context(), mockLogger.AsMockArg(), expectedVersionCheckRequest).
		Return(versionCheckResponse, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(t.Context(), mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: checkcodeJSON}, nil).
		Once()

	analyzer := codeanalyzer.New()

	// Act
	issues, err := analyzer.AnalyzeCode(t.Context(), mockLogger, mockClient, scriptPath)

	// Assert
	require.NoError(t, err, "AnalyzeCode should not return an error")
	assert.Equal(t, expectedCodeIssues, issues, "CodeIssues should match expected value")
}
