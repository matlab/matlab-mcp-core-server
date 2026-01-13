// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabfile_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	runmatlabfileusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabfile"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/runmatlabfile"
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
	tool := runmatlabfile.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

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
	const scriptPath = "/some/script/tofile/myfile.m"
	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        [][]byte{[]byte("image1"), []byte("image2")},
	}
	args := runmatlabfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			runmatlabfileusecase.Args{ScriptPath: scriptPath},
		).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := runmatlabfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

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
	const scriptPath = "/some/script/tofile/myfile.m"
	expectedError := assert.AnError
	args := runmatlabfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	// Act
	result, err := runmatlabfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

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
	args := runmatlabfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			runmatlabfileusecase.Args{ScriptPath: scriptPath},
		).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	// Act
	result, err := runmatlabfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
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
	const scriptPath = "/path/tomepty/file.m"

	// Set up mock usecase to return an empty response
	emptyResponse := entities.EvalResponse{
		ConsoleOutput: "",
		Images:        [][]byte{},
	}
	args := runmatlabfile.Args{ScriptPath: scriptPath}

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockLogger.AsMockArg(),
			mockMATLABSessionClient,
			runmatlabfileusecase.Args{ScriptPath: scriptPath},
		).
		Return(emptyResponse, nil).
		Once()

	// Act
	result, err := runmatlabfile.Handler(mockUsecase, mockGlobalMATLAB)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")

	require.Len(t, result.TextContent, 1, "Should have one text content item")
	assert.Empty(t, result.TextContent[0], "Text content should be empty")
	assert.Empty(t, result.ImageContent, "Image content should be empty")
}

func TestRunMATLABFile_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := runmatlabfile.New(mockLoggerFactory, mockUsecase, mockGlobalMATLAB)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have destructive annotations")
}
