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
	expectedCheckCodeOutput := []string{"Line 1: Warning message", "Line 3: Error message"}
	expectedResponse := checkmatlabcodeusecase.ReturnArgs{
		CheckCodeOutput: expectedCheckCodeOutput,
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
	expectedCleanedOutput := []string{"Line 1: Warning message", "Line 3: Error message"}
	assert.Equal(t, expectedCleanedOutput, result.CheckCodeOutput, "Check code output should match")
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
	var expectedCheckCodeOutput []string
	expectedResponse := checkmatlabcodeusecase.ReturnArgs{
		CheckCodeOutput: expectedCheckCodeOutput,
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
	assert.NotNil(t, result.CheckCodeOutput, "Check code output should not be nil")
	assert.Empty(t, result.CheckCodeOutput, "Check code output should be empty")
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
	assert.NotNil(t, result.CheckCodeOutput, "Check code output should not be nil")
	assert.Empty(t, result.CheckCodeOutput, "Check code output should be empty on error")
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
	assert.NotNil(t, result.CheckCodeOutput, "Check code output should not be nil")
	assert.Empty(t, result.CheckCodeOutput, "Check code output should be empty on error")
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
