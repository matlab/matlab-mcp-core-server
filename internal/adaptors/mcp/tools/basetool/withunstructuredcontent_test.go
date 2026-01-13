// Copyright 2025-2026 The MathWorks, Inc.

package basetool_test

import (
	"context"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestUnstructuredInput struct {
	Query string `json:"query"`
}

const (
	testUnstructuredToolName        = "test-unstructured-tool"
	testUnstructuredToolTitle       = "Test Unstructured Tool"
	testUnstructuredToolDescription = "A test tool for unstructured content"
)

func TestNewToolWithUnstructuredContentOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{
			TextContent: []string{"test response"},
		}, nil
	}

	// Act
	tool := basetool.NewToolWithUnstructuredContent(
		testUnstructuredToolName,
		testUnstructuredToolTitle,
		testUnstructuredToolDescription,
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	// Assert
	assert.Equal(t, testUnstructuredToolName, tool.Name(), "Tool name should match")
	assert.Equal(t, testUnstructuredToolTitle, tool.Title(), "Tool title should match")
	assert.Equal(t, testUnstructuredToolDescription, tool.Description(), "Tool description should match")

	expectedInputSchema, err := jsonschema.For[TestUnstructuredInput](&jsonschema.ForOptions{})
	require.NoError(t, err, "Input schema generation should succeed")
	inputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "Input schema generation should succeed")
	require.Equal(t, expectedInputSchema, inputSchema, "Input schema should match expected")
}

func TestToolWithUnstructuredContentOutput_AddToServer_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockAdder := &mocks.MockToolAdder[TestUnstructuredInput, any]{}
	defer mockAdder.AssertExpectations(t)

	expectedServer := mcp.NewServer(&mcp.Implementation{}, &mcp.ServerOptions{})

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{
			TextContent: []string{"test response"},
		}, nil
	}

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	tool := basetool.NewToolWithUnstructuredContent(
		testUnstructuredToolName,
		testUnstructuredToolTitle,
		testUnstructuredToolDescription,
		expectedAnnotations,
		mockLoggerFactory,
		handler,
	)

	expectedToolInputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "GetInputSchema should not return an error")

	mockAdder.EXPECT().AddTool(
		expectedServer,
		&mcp.Tool{
			Name:         testUnstructuredToolName,
			Title:        testUnstructuredToolTitle,
			Description:  testUnstructuredToolDescription,
			Annotations:  expectedAnnotations.ToToolAnnotations(),
			InputSchema:  expectedToolInputSchema,
			OutputSchema: nil,
		},
		mock.Anything,
	)

	tool.SetToolAdder(mockAdder)

	// Act
	err = tool.AddToServer(expectedServer)

	// Assert
	require.NoError(t, err, "AddToServer should not return an error")
}

func TestToolWithUnstructuredContentOutput_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent:  []string{"text response"},
		ImageContent: []tools.PNGImageData{[]byte("image1")},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.Content, 2, "Should have 2 content items")

	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "First content should be text content")
	assert.Equal(t, expectedRichContent.TextContent[0], textContent.Text, "Text content should match")

	imageContent, ok := result.Content[1].(*mcp.ImageContent)
	require.True(t, ok, "Second content should be image content")
	assert.Equal(t, "image/png", imageContent.MIMEType, "Image MIME type should be PNG")
	assert.Equal(t, []byte(expectedRichContent.ImageContent[0]), imageContent.Data, "Image data should match")
}

func TestToolWithUnstructuredContentOutput_Handler_TextContentOnly(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent: []string{"response 1", "response 2"},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.Content, 2, "Should have 2 content items")

	textContent1, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "First content should be text content")
	assert.Equal(t, expectedRichContent.TextContent[0], textContent1.Text, "First text content should match")

	textContent2, ok := result.Content[1].(*mcp.TextContent)
	require.True(t, ok, "Second content should be text content")
	assert.Equal(t, expectedRichContent.TextContent[1], textContent2.Text, "Second text content should match")
}

func TestToolWithUnstructuredContentOutput_Handler_ImageContentOnly(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		ImageContent: []tools.PNGImageData{
			[]byte("image1"),
			[]byte("image2"),
		},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.Content, 2, "Should have 2 content items")

	imageContent1, ok := result.Content[0].(*mcp.ImageContent)
	require.True(t, ok, "First content should be image content")
	assert.Equal(t, "image/png", imageContent1.MIMEType, "First image MIME type should be PNG")
	assert.Equal(t, []byte(expectedRichContent.ImageContent[0]), imageContent1.Data, "First image data should match")

	imageContent2, ok := result.Content[1].(*mcp.ImageContent)
	require.True(t, ok, "Second content should be image content")
	assert.Equal(t, "image/png", imageContent2.MIMEType, "Second image MIME type should be PNG")
	assert.Equal(t, []byte(expectedRichContent.ImageContent[1]), imageContent2.Data, "Second image data should match")
}

func TestToolWithUnstructuredContentOutput_Handler_NoContent(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent:  []string{},
		ImageContent: []tools.PNGImageData{},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	assert.Empty(t, result.Content, "Content should be empty")
}

func TestToolWithUnstructuredContentOutput_Handler_UnstructuredHandlerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedError := assert.AnError
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{}, expectedError
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return the expected error")
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Nil(t, output, "Output should be nil when error occurs")
}

func TestToolWithUnstructuredContentOutput_Handler_NewMCPSessionLoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedError := messages.AnError

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{}, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(nil, expectedError).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return the expected error")
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Nil(t, output, "Output should be nil when error occurs")
}

func TestToolWithUnstructuredContentOutput_Handler_ContextPropagation(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent: []string{"success"},
	}
	mockSessionLogger := testutils.NewInspectableLogger()

	contextReceived := make(chan context.Context, 1) // Buffering to avoid deadlock
	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		contextReceived <- ctx
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		annotations.NewReadOnlyAnnotations(),
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, t.Context(), <-contextReceived, "Context should be propagated to handler")
}

func TestToolWithUnstructuredContent_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{
			TextContent: []string{"test response"},
		}, nil
	}

	expectedAnnotations := annotations.NewReadOnlyAnnotations()

	// Act
	tool := basetool.NewToolWithUnstructuredContent(
		"",
		"",
		"",
		expectedAnnotations,
		mockLoggerFactory,
		handler,
	)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have read-only annotations")
}

func TestToolWithUnstructuredContentOutput_AddToServer_NilAnnotationInterface(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{
			TextContent: []string{"test response"},
		}, nil
	}

	tool := basetool.NewToolWithUnstructuredContent(
		testUnstructuredToolName,
		testUnstructuredToolTitle,
		testUnstructuredToolDescription,
		nil,
		mockLoggerFactory,
		handler,
	)

	// Act
	err := tool.AddToServer(nil)

	// Assert
	require.Error(t, err, "AddToServer should return an error for nil annotations")
	assert.Contains(t, err.Error(), "annotations must not be nil", "Error message should indicate nil annotations")
}
