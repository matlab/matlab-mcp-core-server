// Copyright 2025-2026 The MathWorks, Inc.

package checkmatlabcode_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	checkmatlabcodemocks "github.com/matlab/matlab-mcp-core-server/mocks/usecases/checkmatlabcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockCodeAnalyzer := &checkmatlabcodemocks.MockCodeAnalyzer{}
	defer mockCodeAnalyzer.AssertExpectations(t)

	// Act
	usecase := checkmatlabcode.New(mockPathValidator, mockCodeAnalyzer)

	// Assert
	assert.NotNil(t, usecase, "Usecase should not be nil")
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockCodeAnalyzer := &checkmatlabcodemocks.MockCodeAnalyzer{}
	defer mockCodeAnalyzer.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedScriptPath := filepath.Join("path", "to", "script.m")
	expectedValidatedPath := filepath.Join("validated", "path", "to", "script.m")
	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: expectedScriptPath,
	}
	codeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Variable 'x' might be unused.",
			Line:        5,
			StartColumn: 1,
			EndColumn:   10,
			Severity:    "warning",
			Fixable:     true,
		},
	}
	expectedCodeIssues := codeIssues

	mockPathValidator.EXPECT().
		ValidateMATLABScript(expectedScriptPath).
		Return(expectedValidatedPath, nil).
		Once()

	mockCodeAnalyzer.EXPECT().
		AnalyzeCode(ctx, mockLogger.AsMockArg(), mockClient, expectedValidatedPath).
		Return(codeIssues, nil).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator, mockCodeAnalyzer)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedCodeIssues, response.CodeIssues, "CodeIssues should match expected value")
}

func TestUsecase_Execute_PathValidationError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockCodeAnalyzer := &checkmatlabcodemocks.MockCodeAnalyzer{}
	defer mockCodeAnalyzer.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedScriptPath := filepath.Join("path", "to", "script.m")
	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: expectedScriptPath,
	}
	pathValidationErr := fmt.Errorf("invalid script path")

	mockPathValidator.EXPECT().
		ValidateMATLABScript(expectedScriptPath).
		Return("", pathValidationErr).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator, mockCodeAnalyzer)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.Error(t, err, "Execute should return an error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}

func TestUsecase_Execute_AnalyzeCodeError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockCodeAnalyzer := &checkmatlabcodemocks.MockCodeAnalyzer{}
	defer mockCodeAnalyzer.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedScriptPath := filepath.Join("path", "to", "script.m")
	expectedValidatedPath := filepath.Join("validated", "path", "to", "script.m")
	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: expectedScriptPath,
	}
	analyzeCodeErr := fmt.Errorf("code analysis failed")

	mockPathValidator.EXPECT().
		ValidateMATLABScript(expectedScriptPath).
		Return(expectedValidatedPath, nil).
		Once()

	mockCodeAnalyzer.EXPECT().
		AnalyzeCode(ctx, mockLogger.AsMockArg(), mockClient, expectedValidatedPath).
		Return(nil, analyzeCodeErr).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator, mockCodeAnalyzer)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.Error(t, err, "Execute should return an error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}

func TestUsecase_Execute_EmptyIssues(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockCodeAnalyzer := &checkmatlabcodemocks.MockCodeAnalyzer{}
	defer mockCodeAnalyzer.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedScriptPath := filepath.Join("path", "to", "script.m")
	expectedValidatedPath := filepath.Join("validated", "path", "to", "script.m")
	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: expectedScriptPath,
	}
	codeIssues := []checkmatlabcode.CodeIssue{}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(expectedScriptPath).
		Return(expectedValidatedPath, nil).
		Once()

	mockCodeAnalyzer.EXPECT().
		AnalyzeCode(ctx, mockLogger.AsMockArg(), mockClient, expectedValidatedPath).
		Return(codeIssues, nil).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator, mockCodeAnalyzer)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Empty(t, response.CodeIssues, "CodeIssues should be empty")
}
