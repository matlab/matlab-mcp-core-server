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

type toolInput struct {
	Query string `json:"query"`
}

type testError struct{}

func (e *testError) Error() string { return "test error" }
func (e *testError) MWMarker()     {}

func TestNewToolWithUnstructuredContentOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedInput := toolInput{Query: "test query"}
	expectedTextContent := []string{"response line 1", "response line 2"}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	tool := server.NewToolWithUnstructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input toolInput) (tools.RichContent, i18n.Error) {
			return tools.RichContent{
				TextContent: expectedTextContent,
			}, nil
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, _, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, expectedInput)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Content, 2)

	textContent1, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, expectedTextContent[0], textContent1.Text)

	textContent2, ok := result.Content[1].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, expectedTextContent[1], textContent2.Text)
}

func TestNewToolWithUnstructuredContentOutput_HandlerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedError := &testError{}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	tool := server.NewToolWithUnstructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input toolInput) (tools.RichContent, i18n.Error) {
			return tools.RichContent{}, expectedError
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, _, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, toolInput{})

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Nil(t, result)
}

func TestNewToolWithUnstructuredContentOutput_HandlerReceivesLogger(t *testing.T) {
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

	tool := server.NewToolWithUnstructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input toolInput) (tools.RichContent, i18n.Error) {
			handlerCalled = true

			request.Logger().Info(expectedMessage)

			return tools.RichContent{}, nil
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, toolInput{Query: "test"})

	// Assert
	require.NoError(t, err)
	require.True(t, handlerCalled, "handler should be called")

	infoLogs := mockLogger.InfoLogs()
	_, found := infoLogs[expectedMessage]
	require.True(t, found, "expected log message should be present")
}

func TestNewToolWithUnstructuredContentOutput_HandlerReceivesConfig(t *testing.T) {
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

	tool := server.NewToolWithUnstructuredContentOutput(
		tools.Definition{Name: "test-tool"},
		func(ctx context.Context, request tools.CallRequest, input toolInput) (tools.RichContent, i18n.Error) {
			handlerCalled = true

			result, err := request.Config().Get(expectedKey, "")
			require.NoError(t, err)
			assert.Equal(t, expectedValue, result)

			return tools.RichContent{}, nil
		},
	)

	internalTool := tool.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := internalTool.Handler()(t.Context(), mcpCallToolRequest, toolInput{Query: "test"})

	// Assert
	require.NoError(t, err)
	require.True(t, handlerCalled, "handler should be called")
}
