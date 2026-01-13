// Copyright 2025-2026 The MathWorks, Inc.

package evalmatlabcode_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	evalmatlabcodeusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/evalmatlabcode"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/evalmatlabcode"
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
	tool := evalmatlabcode.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

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
	const code = "disp('Hello, World!')"
	const projectPath = "/some/path"
	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        [][]byte{[]byte("image1"), []byte("image2")},
	}
	args := evalmatlabcode.Args{
		Code:        code,
		ProjectPath: projectPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			evalmatlabcodeusecase.Args{
				Code:        code,
				ProjectPath: projectPath,
			},
		).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := evalmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")

	require.Len(t, result.TextContent, 1, "Should have one text content item")
	assert.Equal(t, expectedResponse.ConsoleOutput, result.TextContent[0], "Text content should match")

	require.Len(t, result.ImageContent, 2, "Should have two image content items")
	assert.Equal(t, "image1", string(result.ImageContent[0]), "First image should match")
	assert.Equal(t, "image2", string(result.ImageContent[1]), "Second image should match")
}

func TestTool_Handler_ClientReturnsError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const code = "invalid code"
	const projectPath = "/some/path"
	expectedError := assert.AnError
	args := evalmatlabcode.Args{
		Code:        code,
		ProjectPath: projectPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := evalmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Empty(t, result, "Result should be empty in an error case")
}

func TestTool_Handler_UsecaseReturnsError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const code = "invalid code"
	const projectPath = "/some/path"
	expectedError := assert.AnError
	args := evalmatlabcode.Args{
		Code:        code,
		ProjectPath: projectPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			evalmatlabcodeusecase.Args{
				Code:        code,
				ProjectPath: projectPath,
			},
		).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	// Act
	result, err := evalmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Empty(t, result, "Result should be empty in an error case")
}

func TestTool_Handler_UsecaseReturnsEmptyResponse(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const code = "% Empty comment"
	const projectPath = "/some/path"

	emptyResponse := entities.EvalResponse{
		ConsoleOutput: "",
		Images:        nil,
	}
	args := evalmatlabcode.Args{
		Code:        code,
		ProjectPath: projectPath,
	}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			evalmatlabcodeusecase.Args{
				Code:        code,
				ProjectPath: projectPath,
			},
		).
		Return(emptyResponse, nil).
		Once()

	// Act
	result, err := evalmatlabcode.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")

	require.Len(t, result.TextContent, 1, "Should have one text content item")
	assert.Empty(t, result.TextContent[0], "Text content should be empty")
	assert.Empty(t, result.ImageContent, "Image content should be empty")
}

func TestEvaluateMATLABCode_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := evalmatlabcode.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have destructive annotations")
}
