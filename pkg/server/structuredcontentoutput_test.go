// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"context"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/matlab/matlab-mcp-core-server/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type structuredToolInput struct {
	Query string `json:"query"`
}

type structuredToolOutput struct {
	Result string `json:"result"`
}

type structuredTestError struct{}

func (e *structuredTestError) Error() string { return "structured test error" }

func (e *structuredTestError) MWMarker() {}

func TestNewToolWithStructuredContentOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedInput := structuredToolInput{Query: "test query"}
	expectedOutput := structuredToolOutput{Result: "success"}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	tool := server.NewToolWithStructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input structuredToolInput) (structuredToolOutput, i18n.Error) {
			return expectedOutput, nil
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, output, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, expectedInput)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestNewToolWithStructuredContentOutput_HandlerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedError := &structuredTestError{}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	tool := server.NewToolWithStructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input structuredToolInput) (structuredToolOutput, i18n.Error) {
			return structuredToolOutput{}, expectedError
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, output, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, structuredToolInput{})

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Equal(t, structuredToolOutput{}, output)
}

func TestNewToolWithStructuredContentOutput_HandlerReceivesLogger(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedMessage := "handler received logger"
	handlerCalled := false

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	tool := server.NewToolWithStructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input structuredToolInput) (structuredToolOutput, i18n.Error) {
			handlerCalled = true

			request.Logger().Info(expectedMessage)

			return structuredToolOutput{}, nil
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, structuredToolInput{Query: "test"})

	// Assert
	require.NoError(t, err)
	require.True(t, handlerCalled, "handler should be called")

	infoLogs := mockLogger.InfoLogs()
	_, found := infoLogs[expectedMessage]
	require.True(t, found, "expected log message should be present")
}

func TestNewToolWithStructuredContentOutput_HandlerReceivesConfig(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedKey := "test-key"
	expectedValue := "test-value"
	handlerCalled := false

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	mockConfig.EXPECT().
		Get(expectedKey).
		Return(expectedValue, nil).
		Once()

	tool := server.NewToolWithStructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input structuredToolInput) (structuredToolOutput, i18n.Error) {
			handlerCalled = true

			result, err := request.Config().Get(expectedKey, "")
			require.NoError(t, err)
			assert.Equal(t, expectedValue, result)

			return structuredToolOutput{}, nil
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, structuredToolInput{Query: "test"})

	// Assert
	require.NoError(t, err)
	require.True(t, handlerCalled, "handler should be called")
}
