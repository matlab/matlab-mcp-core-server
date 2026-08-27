// Copyright 2026 The MathWorks, Inc.

package tools_test

import (
	"context"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	publictypes "github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/definition"
	basetoolmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/mcp/tools/basetool"
	publictypesmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/publictypes"
	toolsmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/tools"
	telemetrymocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnstructured_HappyPath(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsmocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockCallRequest := &publictypesmocks.MockToolCallRequest{}
	defer mockCallRequest.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedInput := toolInput{Query: "test query"}
	expectedTextContent := []string{"response line 1", "response line 2"}

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordToolCallRequest(ctx, "test-tool", telemetry.ToolSourceBuiltin).
		Return().
		Once()

	mockToolCallRequestFactory.EXPECT().
		New(mockLogger.AsMockArg(), mockConfig, mockMessageCatalog).
		Return(mockCallRequest).
		Once()

	tool := tools.NewUnstructured(
		publictypes.ToolDefinition{Name: "test-tool"},
		func(ctx context.Context, request publictypes.ToolCallRequest, input toolInput) (publictypes.RichContent, publictypes.Error) {
			assert.Equal(t, expectedInput, input)
			assert.Equal(t, mockCallRequest, request)
			return publictypes.RichContent{
				TextContent: expectedTextContent,
			}, nil
		},
	)

	internalTool := tool.ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).(basetool.ToolWithUnstructuredContentOutput[toolInput])

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, _, err := internalTool.Handler()(ctx, mcpCallToolRequest, expectedInput)

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

func TestNewUnstructured_HandlerError(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsmocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockCallRequest := &publictypesmocks.MockToolCallRequest{}
	defer mockCallRequest.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedError := anI18nError

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordToolCallRequest(ctx, "test-tool", telemetry.ToolSourceBuiltin).
		Return().
		Once()

	mockToolCallRequestFactory.EXPECT().
		New(mockLogger.AsMockArg(), mockConfig, mockMessageCatalog).
		Return(mockCallRequest).
		Once()

	tool := tools.NewUnstructured(
		publictypes.ToolDefinition{Name: "test-tool"},
		func(ctx context.Context, request publictypes.ToolCallRequest, input toolInput) (publictypes.RichContent, publictypes.Error) {
			return publictypes.RichContent{}, expectedError
		},
	)

	internalTool := tool.ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).(basetool.ToolWithUnstructuredContentOutput[toolInput])

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, _, err := internalTool.Handler()(ctx, mcpCallToolRequest, toolInput{})

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Nil(t, result)
}

func TestNewUnstructured_HandlerReceivesLogger(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsmocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockCallRequest := &publictypesmocks.MockToolCallRequest{}
	defer mockCallRequest.AssertExpectations(t)

	mockPkgLogger := &publictypesmocks.MockLogger{}
	defer mockPkgLogger.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedMessage := "handler received logger"

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordToolCallRequest(ctx, "test-tool", telemetry.ToolSourceBuiltin).
		Return().
		Once()

	mockToolCallRequestFactory.EXPECT().
		New(mockLogger.AsMockArg(), mockConfig, mockMessageCatalog).
		Return(mockCallRequest).
		Once()

	mockCallRequest.EXPECT().
		Logger().
		Return(mockPkgLogger).
		Once()

	mockPkgLogger.EXPECT().
		Info(expectedMessage).
		Return().
		Once()

	tool := tools.NewUnstructured(
		publictypes.ToolDefinition{Name: "test-tool"},
		func(ctx context.Context, request publictypes.ToolCallRequest, input toolInput) (publictypes.RichContent, publictypes.Error) {
			request.Logger().Info(expectedMessage)
			return publictypes.RichContent{}, nil
		},
	)

	internalTool := tool.ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).(basetool.ToolWithUnstructuredContentOutput[toolInput])

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := internalTool.Handler()(ctx, mcpCallToolRequest, toolInput{Query: "test"})

	// Assert
	require.NoError(t, err)
}

func TestNewUnstructured_HandlerReceivesConfig(t *testing.T) {
	// Arrange
	ctx := t.Context()

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsmocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockCallRequest := &publictypesmocks.MockToolCallRequest{}
	defer mockCallRequest.AssertExpectations(t)

	mockPkgConfig := &publictypesmocks.MockConfig{}
	defer mockPkgConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSession := &mcp.ServerSession{}
	expectedKey := "test-key"
	expectedValue := "test-value"

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockLogger, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockTelemetry.EXPECT().
		RecordToolCallRequest(ctx, "test-tool", telemetry.ToolSourceBuiltin).
		Return().
		Once()

	mockToolCallRequestFactory.EXPECT().
		New(mockLogger.AsMockArg(), mockConfig, mockMessageCatalog).
		Return(mockCallRequest).
		Once()

	mockCallRequest.EXPECT().
		Config().
		Return(mockPkgConfig).
		Once()

	mockPkgConfig.EXPECT().
		Get(expectedKey, "").
		Return(expectedValue, nil).
		Once()

	tool := tools.NewUnstructured(
		publictypes.ToolDefinition{Name: "test-tool"},
		func(ctx context.Context, request publictypes.ToolCallRequest, input toolInput) (publictypes.RichContent, publictypes.Error) {
			result, err := request.Config().Get(expectedKey, "")
			require.NoError(t, err)
			assert.Equal(t, expectedValue, result)
			return publictypes.RichContent{}, nil
		},
	)

	internalTool := tool.ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).(basetool.ToolWithUnstructuredContentOutput[toolInput])

	mcpCallToolRequest := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := internalTool.Handler()(ctx, mcpCallToolRequest, toolInput{Query: "test"})

	// Assert
	require.NoError(t, err)
}

func TestNewUnstructured_DefinitionFieldsForwarded(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsmocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	expectedName := "my-tool"
	expectedTitle := "My Tool"
	expectedDescription := "A tool that does things"
	expectedAnnotations := tools.NewReadOnlyAnnotations()

	tool := tools.NewUnstructured(
		publictypes.ToolDefinition{
			Name:        expectedName,
			Title:       expectedTitle,
			Description: expectedDescription,
			Annotations: expectedAnnotations,
		},
		func(ctx context.Context, request publictypes.ToolCallRequest, input toolInput) (publictypes.RichContent, publictypes.Error) {
			return publictypes.RichContent{}, nil
		},
	)

	// Act
	internalTool := tool.ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).(basetool.ToolWithUnstructuredContentOutput[toolInput])

	// Assert
	require.Equal(t, expectedName, internalTool.Name())
	require.Equal(t, expectedTitle, internalTool.Title())
	require.Equal(t, expectedDescription, internalTool.Description())
	require.Equal(t, expectedAnnotations, internalTool.Annotations())
}

type toolInput struct {
	Query string `json:"query"`
}
