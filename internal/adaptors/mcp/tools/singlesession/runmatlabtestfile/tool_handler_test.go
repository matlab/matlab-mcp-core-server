// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabtestfile_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	runmatlabtestfileusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabtestfile"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/runmatlabtestfile"
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
	tool := runmatlabtestfile.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

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
	const scriptPath = "/some/script/tofile/testFile.m"
	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Test Results",
		Images:        [][]byte{[]byte("image1"), []byte("image2")},
	}
	args := runmatlabtestfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			runmatlabtestfileusecase.Args{ScriptPath: scriptPath},
		).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := runmatlabtestfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

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
	const scriptPath = "/some/path"
	expectedError := assert.AnError
	args := runmatlabtestfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := runmatlabtestfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

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
	const scriptPath = "/invalid/path.m"
	expectedError := assert.AnError
	args := runmatlabtestfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			runmatlabtestfileusecase.Args{ScriptPath: scriptPath},
		).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	// Act
	result, err := runmatlabtestfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

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
	const scriptPath = "/path/tomepty/testFile.m"

	// Set up mock usecase to return an empty response
	emptyResponse := entities.EvalResponse{
		ConsoleOutput: "",
		Images:        [][]byte{},
	}
	args := runmatlabtestfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			runmatlabtestfileusecase.Args{ScriptPath: scriptPath},
		).
		Return(emptyResponse, nil).
		Once()

	// Act
	result, err := runmatlabtestfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")

	require.Len(t, result.TextContent, 1, "Should have one text content item")
	assert.Empty(t, result.TextContent[0], "Text content should be empty")
	assert.Empty(t, result.ImageContent, "Image content should be empty")
}

func TestRunMATLABTestFile_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := runmatlabtestfile.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have destructive annotations")
}
