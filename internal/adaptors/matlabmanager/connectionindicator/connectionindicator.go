// Copyright 2026 The MathWorks, Inc.

package connectionindicator

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/usecases/utils/matlabstring"
)

const connectedTitleSeparator = "  |  "

// JavaScript Desktop API (R2025a+).
const jsDesktopTitleReadCode = `matlab.ui.container.internal.RootApp.getInstance().Title`

const jsDesktopTitleWriteCodeFormat = `[~] = subsasgn(matlab.ui.container.internal.RootApp.getInstance(), substruct('.','Title'), '%s');`

// Java/Swing desktop API (before R2025a).
const swingDesktopTitleReadCode = `string(com.mathworks.mlservices.MatlabDesktopServices.getDesktop.getMainFrame.getTitle)`

const swingDesktopTitleWriteCodeFormat = `com.mathworks.mlservices.MatlabDesktopServices.getDesktop.getMainFrame.setTitle('%s');`

const minVersionForJSDesktopTitle = "R2025a"

type MessageCatalog interface {
	Get(key messages.MessageKey) string
}

type ConnectionIndicator struct {
	messageCatalog MessageCatalog
}

type desktopTitleCodes struct {
	readCode        string
	writeCodeFormat string
}

func New(messageCatalog MessageCatalog) *ConnectionIndicator {
	return &ConnectionIndicator{messageCatalog: messageCatalog}
}

func (i *ConnectionIndicator) ShowGreeting(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, greetingInfo entities.GreetingInfo) error {
	name := clientName(greetingInfo.ClientInfo)
	if name == "" {
		sessionLogger.Warn("skipping MCP connection greeting: client provided neither a title nor a name")
		return nil
	}

	escapedName := html.EscapeString(name)
	if escapedName != name {
		sessionLogger.With("client-name", name).Warn("escaped unusual characters in MCP client name before displaying connection greeting")
	}

	message := messages.New_CLIMessages_GreetingOnConnect_Message_FromCatalog(i.messageCatalog, escapedName, escapedName)
	greetingText := "\n<strong>" + message + "</strong>\n"

	var greeting string
	var err error
	if greetingInfo.ShowMATLABDesktop {
		greeting = greetingText
		_, err = client.FEval(ctx, sessionLogger, entities.FEvalRequest{
			Function:   "disp",
			Arguments:  []string{greeting},
			NumOutputs: 0,
		})
	} else {
		greeting = dispCode(greetingText)
		_, err = client.EvalWithCapture(ctx, sessionLogger, entities.EvalRequest{
			Code: greeting,
		})
	}

	if err != nil {
		sessionLogger.WithError(err).With("greeting", greeting).Debug("MCP connection greeting failed")
		return err
	}
	return nil
}

func (i *ConnectionIndicator) ApplyConnectedTitle(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, greetingInfo entities.GreetingInfo) (string, error) {
	if !greetingInfo.ShowMATLABDesktop {
		sessionLogger.Debug("skipping MCP connection title: MATLAB desktop not shown")
		return "", nil
	}

	name := clientName(greetingInfo.ClientInfo)
	if name == "" {
		sessionLogger.Debug("skipping MCP connection title: client provided neither a title nor a name")
		return "", nil
	}

	titleCodes := i.resolveTitleCodes(ctx, sessionLogger, client)

	response, err := client.FEval(ctx, sessionLogger, entities.FEvalRequest{
		Function:   "eval",
		Arguments:  []string{titleCodes.readCode},
		NumOutputs: 1,
	})
	if err != nil {
		return "", err
	}

	if len(response.Outputs) == 0 {
		sessionLogger.Debug("skipping MCP connection title: no output returned when reading desktop title")
		return "", nil
	}
	originalTitle, _ := response.Outputs[0].(string)
	originalTitle = strings.TrimRight(originalTitle, "\r\n")
	if originalTitle == "" {
		sessionLogger.Debug("skipping MCP connection title: desktop title read returned empty")
		return "", nil
	}

	suffix := i.connectedTitleSuffix()
	for strings.HasSuffix(originalTitle, suffix) {
		originalTitle = strings.TrimSuffix(originalTitle, suffix)
	}
	if originalTitle == "" {
		sessionLogger.Debug("skipping MCP connection title: desktop title contained only the connected suffix")
		return "", nil
	}

	if err := i.writeDesktopTitle(ctx, sessionLogger, client, titleCodes.writeCodeFormat, originalTitle+suffix); err != nil {
		return "", err
	}

	return originalTitle, nil
}

func (i *ConnectionIndicator) connectedTitleSuffix() string {
	return connectedTitleSeparator + messages.New_CLIMessages_TitleOnConnect_Message_FromCatalog(i.messageCatalog)
}

func (i *ConnectionIndicator) RestoreTitle(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, originalTitle string) error {
	titleCodes := i.resolveTitleCodes(ctx, sessionLogger, client)
	return i.writeDesktopTitle(ctx, sessionLogger, client, titleCodes.writeCodeFormat, originalTitle)
}

func (i *ConnectionIndicator) writeDesktopTitle(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, writeCodeFormat string, value string) error {
	code := fmt.Sprintf(writeCodeFormat, matlabstring.EscapeSingleQuotes(value))
	_, err := client.FEval(ctx, sessionLogger, entities.FEvalRequest{
		Function:   "eval",
		Arguments:  []string{code},
		NumOutputs: 0,
	})
	return err
}

func jsDesktopTitleCodes() desktopTitleCodes {
	return desktopTitleCodes{
		readCode:        jsDesktopTitleReadCode,
		writeCodeFormat: jsDesktopTitleWriteCodeFormat,
	}
}

func swingDesktopTitleCodes() desktopTitleCodes {
	return desktopTitleCodes{
		readCode:        swingDesktopTitleReadCode,
		writeCodeFormat: swingDesktopTitleWriteCodeFormat,
	}
}

func (i *ConnectionIndicator) resolveTitleCodes(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient) desktopTitleCodes {
	response, err := client.FEval(ctx, sessionLogger, entities.FEvalRequest{
		Function:   "isMATLABReleaseOlderThan",
		Arguments:  []string{minVersionForJSDesktopTitle},
		NumOutputs: 1,
	})
	if err != nil {
		sessionLogger.WithError(err).Warn("failed to check MATLAB release for desktop title API, defaulting to JavaScript Desktop API")
		return jsDesktopTitleCodes()
	}

	if len(response.Outputs) == 0 {
		sessionLogger.Warn("MATLAB release check returned no outputs, defaulting to JavaScript Desktop title API")
		return jsDesktopTitleCodes()
	}

	isOlderRelease, ok := response.Outputs[0].(bool)
	if !ok {
		sessionLogger.Warn("MATLAB release check returned a non-bool, defaulting to JavaScript Desktop title API")
		return jsDesktopTitleCodes()
	}

	// <R2025a: use swing
	if isOlderRelease {
		return swingDesktopTitleCodes()
	}
	return jsDesktopTitleCodes()
}

func clientName(clientInfo entities.MCPClientInfo) string {
	if clientInfo.Title != "" {
		return clientInfo.Title
	}
	return clientInfo.Name
}

func dispCode(message string) string {
	lines := strings.Split(message, "\n")
	statements := make([]string, len(lines))
	for index, line := range lines {
		statements[index] = "disp('" + matlabstring.EscapeSingleQuotes(line) + "')"
	}
	return strings.Join(statements, "; ")
}
