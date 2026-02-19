// Copyright 2025-2026 The MathWorks, Inc.

package checkmatlabcode_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	checkmatlabcodeusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/checkmatlabcode"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	// Act
	tool := checkmatlabcode.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.NotNil(t, tool)
}

func TestTool_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const scriptPath = "/path/to/script.m"
	usecaseCodeIssues := []checkmatlabcodeusecase.CodeIssue{
		{
			Description: "Warning message",
			Line:        1,
			StartColumn: 1,
			EndColumn:   10,
			Severity:    "warning",
			Fixable:     false,
		},
		{
			Description: "Error message",
			Line:        3,
			StartColumn: 5,
			EndColumn:   15,
			Severity:    "error",
			Fixable:     true,
		},
	}
	usecaseResponse := checkmatlabcodeusecase.ReturnArgs{
		CodeIssues: usecaseCodeIssues,
	}
	expectedCodeIssues := []checkmatlabcode.CodeIssue{
		{
			Description: "Warning message",
			Line:        1,
			StartColumn: 1,
			EndColumn:   10,
			Severity:    "warning",
			Fixable:     false,
		},
		{
			Description: "Error message",
			Line:        3,
			StartColumn: 5,
			EndColumn:   15,
			Severity:    "error",
			Fixable:     true,
		},
	}
	args := checkmatlabcode.Args{
		ScriptPath: scriptPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockMATLABSessionClient, checkmatlabcodeusecase.Args{ScriptPath: scriptPath}).
		Return(usecaseResponse, nil).
		Once()

	// Act
	result, err := checkmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, expectedCodeIssues, result.CodeIssues, "Code issues should match")
}

func TestTool_Handler_EmptyOutput(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const scriptPath = "/path/to/script.m"
	var expectedCodeIssues []checkmatlabcodeusecase.CodeIssue
	expectedResponse := checkmatlabcodeusecase.ReturnArgs{
		CodeIssues: expectedCodeIssues,
	}
	args := checkmatlabcode.Args{
		ScriptPath: scriptPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockMATLABSessionClient, checkmatlabcodeusecase.Args{ScriptPath: scriptPath}).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := checkmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.NotNil(t, result.CodeIssues, "Code issues should not be nil")
	assert.Empty(t, result.CodeIssues, "Code issues should be empty")
}

func TestTool_Handler_ClientError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const scriptPath = "/path/to/script.m"
	expectedError := assert.AnError
	args := checkmatlabcode.Args{
		ScriptPath: scriptPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := checkmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.NotNil(t, result.CodeIssues, "Code issues should not be nil")
	assert.Empty(t, result.CodeIssues, "Code issues should be empty on error")
}

func TestTool_Handler_UsecaseError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const scriptPath = "/path/to/script.m"
	expectedError := assert.AnError
	args := checkmatlabcode.Args{
		ScriptPath: scriptPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), mockMATLABSessionClient, checkmatlabcodeusecase.Args{ScriptPath: scriptPath}).
		Return(checkmatlabcodeusecase.ReturnArgs{}, expectedError).
		Once()

	// Act
	result, err := checkmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.NotNil(t, result.CodeIssues, "Code issues should not be nil")
	assert.Empty(t, result.CodeIssues, "Code issues should be empty on error")
}

func TestCheckMATLABCode_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewReadOnlyAnnotations()

	// Act
	tool := checkmatlabcode.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have read-only annotations")
}
