// Copyright 2026 The MathWorks, Inc.

package connectionindicator_test

import (
	"fmt"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/connectionindicator"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/matlabmanager/connectionindicator"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func greetingTemplate() string {
	return messages.NewCatalog(messages.Locale_en_US).Get(messages.CLIMessages_GreetingOnConnect)
}

func titleTemplate() string {
	return messages.NewCatalog(messages.Locale_en_US).Get(messages.CLIMessages_TitleOnConnect)
}

func versionCheckRequest() entities.FEvalRequest {
	return entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{"R2025a"},
		NumOutputs: 1,
	}
}

func jsDesktopTitleReadRequest() entities.FEvalRequest {
	return entities.FEvalRequest{
		Function:   "eval",
		Arguments:  []string{`matlab.ui.container.internal.RootApp.getInstance().Title`},
		NumOutputs: 1,
	}
}

func jsDesktopTitleWriteRequest(title string) entities.FEvalRequest {
	return entities.FEvalRequest{
		Function:   "eval",
		Arguments:  []string{fmt.Sprintf(`[~] = subsasgn(matlab.ui.container.internal.RootApp.getInstance(), substruct('.','Title'), '%s');`, title)},
		NumOutputs: 0,
	}
}

func swingTitleReadRequest() entities.FEvalRequest {
	return entities.FEvalRequest{
		Function:   "eval",
		Arguments:  []string{`string(com.mathworks.mlservices.MatlabDesktopServices.getDesktop.getMainFrame.getTitle)`},
		NumOutputs: 1,
	}
}

func swingTitleWriteRequest(title string) entities.FEvalRequest {
	return entities.FEvalRequest{
		Function:   "eval",
		Arguments:  []string{fmt.Sprintf(`com.mathworks.mlservices.MatlabDesktopServices.getDesktop.getMainFrame.setTitle('%s');`, title)},
		NumOutputs: 0,
	}
}

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	// Act
	connectionIndicator := connectionindicator.New(mockCatalog)

	// Assert
	assert.NotNil(t, connectionIndicator)
}

func TestConnectionIndicator_ShowGreeting_UsesTitle(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedRequest := entities.FEvalRequest{
		Function:   "disp",
		Arguments:  []string{"\n<strong>Visual Studio Code is connected to MATLAB via the MATLAB MCP server.\nNow you can use Visual Studio Code to work with MATLAB.</strong>\n"},
		NumOutputs: 0,
	}

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_GreetingOnConnect).
		Return(greetingTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	err := connectionIndicator.ShowGreeting(ctx, mockLogger, mockClient, entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true})

	// Assert
	require.NoError(t, err)
}

func TestConnectionIndicator_ShowGreeting_NoDesktop_UsesCaptureWithMarkup(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedRequest := entities.EvalRequest{
		Code: "disp(''); disp('<strong>Visual Studio Code is connected to MATLAB via the MATLAB MCP server.'); disp('Now you can use Visual Studio Code to work with MATLAB.</strong>'); disp('')",
	}

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_GreetingOnConnect).
		Return(greetingTemplate()).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	// Act
	err := connectionIndicator.ShowGreeting(ctx, mockLogger, mockClient, entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: false})

	// Assert
	require.NoError(t, err)
}

func TestConnectionIndicator_ShowGreeting_FallsBackToNameWhenTitleEmpty(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedRequest := entities.FEvalRequest{
		Function:   "disp",
		Arguments:  []string{"\n<strong>claude is connected to MATLAB via the MATLAB MCP server.\nNow you can use claude to work with MATLAB.</strong>\n"},
		NumOutputs: 0,
	}

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "", Name: "claude"}

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_GreetingOnConnect).
		Return(greetingTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	err := connectionIndicator.ShowGreeting(ctx, mockLogger, mockClient, entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true})

	// Assert
	require.NoError(t, err)
}

func TestConnectionIndicator_ShowGreeting_SkipsWhenTitleAndNameEmpty(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "", Name: ""}

	// Act
	err := connectionIndicator.ShowGreeting(ctx, mockLogger, mockClient, entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true})

	// Assert
	require.NoError(t, err)
	_, hasWarnLog := mockLogger.WarnLogs()["skipping MCP connection greeting: client provided neither a title nor a name"]
	assert.True(t, hasWarnLog, "should log a Warn when neither title nor name is available")
}

func TestConnectionIndicator_ShowGreeting_EscapesHTMLInClientNameAndWarns(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	// The client name is not trusted and the MATLAB Command Window renders a subset of HTML, so any
	// markup characters in the name are escaped to inert entities before being displayed.
	clientName := "Cody <beta>"
	escapedName := "Cody &lt;beta&gt;"
	expectedRequest := entities.FEvalRequest{
		Function:   "disp",
		Arguments:  []string{"\n<strong>" + escapedName + " is connected to MATLAB via the MATLAB MCP server.\nNow you can use " + escapedName + " to work with MATLAB.</strong>\n"},
		NumOutputs: 0,
	}

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: clientName, Name: "vscode"}

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_GreetingOnConnect).
		Return(greetingTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	err := connectionIndicator.ShowGreeting(ctx, mockLogger, mockClient, entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true})

	// Assert
	require.NoError(t, err)
	warnFields, hasWarnLog := mockLogger.WarnLogs()["escaped unusual characters in MCP client name before displaying connection greeting"]
	require.True(t, hasWarnLog, "should log a Warn when the client name contained characters that had to be escaped")
	assert.Equal(t, clientName, warnFields["client-name"], "the warn log should record the original unescaped client name")
}

func TestConnectionIndicator_ShowGreeting_FEvalError(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedError := assert.AnError

	expectedRequest := entities.FEvalRequest{
		Function:   "disp",
		Arguments:  []string{"\n<strong>Visual Studio Code is connected to MATLAB via the MATLAB MCP server.\nNow you can use Visual Studio Code to work with MATLAB.</strong>\n"},
		NumOutputs: 0,
	}

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_GreetingOnConnect).
		Return(greetingTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), expectedRequest).
		Return(entities.FEvalResponse{}, expectedError).
		Once()

	// Act
	err := connectionIndicator.ShowGreeting(ctx, mockLogger, mockClient, entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true})

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestConnectionIndicator_ApplyConnectedTitle_UsesTitle(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_FallsBackToNameWhenTitleEmpty(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "", Name: "claude"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_AlreadyConnected_DoesNotDoubleAppend(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB  |  Agent Connected\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_AlreadyDoubleSuffixed_StripsToSingle(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	// A title already polluted with two suffixes by the earlier bug is healed back to a single suffix.
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB  |  Agent Connected  |  Agent Connected\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_NoDesktop_Skips(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: false}

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, original)
}

func TestConnectionIndicator_ApplyConnectedTitle_EmptyReadOutput_SkipsWrite(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	// The read yields no outputs, so the title write must be skipped without erroring.
	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: nil}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, original)
}

func TestConnectionIndicator_ApplyConnectedTitle_SkipsWhenTitleAndNameEmpty(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "", Name: ""}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, original)
}

func TestConnectionIndicator_ApplyConnectedTitle_ReadError(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	expectedError := assert.AnError

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{}, expectedError).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, original)
}

func TestConnectionIndicator_ApplyConnectedTitle_WriteError(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")
	expectedError := assert.AnError

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, expectedError).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, original)
}

func TestConnectionIndicator_RestoreTitle_WritesOriginal(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	writeRequest := jsDesktopTitleWriteRequest("Original Title")

	connectionIndicator := connectionindicator.New(mockCatalog)

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	err := connectionIndicator.RestoreTitle(ctx, mockLogger, mockClient, "Original Title")

	// Assert
	require.NoError(t, err)
}

func TestConnectionIndicator_RestoreTitle_WriteError(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	writeRequest := jsDesktopTitleWriteRequest("Original Title")
	expectedError := assert.AnError

	connectionIndicator := connectionindicator.New(mockCatalog)

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{false}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, expectedError).
		Once()

	// Act
	err := connectionIndicator.RestoreTitle(ctx, mockLogger, mockClient, "Original Title")

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestConnectionIndicator_ApplyConnectedTitle_OldRelease_UsesSwingAPI(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := swingTitleReadRequest()
	writeRequest := swingTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{true}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_VersionCheckError_DefaultsToJSDesktop(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	// A failed release check must not abort the title change; it falls back to the JavaScript Desktop API.
	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{}, assert.AnError).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_VersionCheckNoOutput_DefaultsToJSDesktop(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: nil}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_ApplyConnectedTitle_VersionCheckNonBool_DefaultsToJSDesktop(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	readRequest := jsDesktopTitleReadRequest()
	writeRequest := jsDesktopTitleWriteRequest("My MATLAB  |  Agent Connected")

	connectionIndicator := connectionindicator.New(mockCatalog)
	clientInfo := entities.MCPClientInfo{Title: "Visual Studio Code", Name: "vscode"}
	greetingInfo := entities.GreetingInfo{ClientInfo: clientInfo, ShowMATLABDesktop: true}

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{"not-a-bool"}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), readRequest).
		Return(entities.FEvalResponse{Outputs: []any{"My MATLAB\n"}}, nil).
		Once()

	mockCatalog.EXPECT().
		Get(messages.CLIMessages_TitleOnConnect).
		Return(titleTemplate()).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	original, err := connectionIndicator.ApplyConnectedTitle(ctx, mockLogger, mockClient, greetingInfo)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "My MATLAB", original)
}

func TestConnectionIndicator_RestoreTitle_OldRelease_UsesSwingAPI(t *testing.T) {
	// Arrange
	mockCatalog := &mocks.MockMessageCatalog{}
	defer mockCatalog.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	writeRequest := swingTitleWriteRequest("Original Title")

	connectionIndicator := connectionindicator.New(mockCatalog)

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), versionCheckRequest()).
		Return(entities.FEvalResponse{Outputs: []any{true}}, nil).
		Once()

	mockClient.EXPECT().
		FEval(ctx, mockLogger.AsMockArg(), writeRequest).
		Return(entities.FEvalResponse{}, nil).
		Once()

	// Act
	err := connectionIndicator.RestoreTitle(ctx, mockLogger, mockClient, "Original Title")

	// Assert
	require.NoError(t, err)
}
