// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// ServerWithLoggingTestSuite tests SDK logging functionalities.
type ServerWithLoggingTestSuite struct {
	suite.Suite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithLoggingTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithLogging(s.T())
}

func TestServerWithLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithLoggingTestSuite))
}

func (s *ServerWithLoggingTestSuite) TestSDK_Logging_DependenciesAndToolsProviderLogToFile() {
	// Arrange
	logFolder, err := os.MkdirTemp("", "server_session") // Can't use s.T().Tempdir() because too long for socket path
	s.Require().NoError(err)
	defer s.Require().NoError(os.RemoveAll(logFolder))

	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil,
		"--log-level=debug",
		"--log-folder="+logFolder,
	)

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	// Act
	_, err = session.CallTool(s.T().Context(), "tool-that-logs", map[string]any{"name": "World"})
	s.Require().NoError(err, "should call tool successfully")

	// Assert
	serverLogPattern := filepath.Join(logFolder, "server-*.log")

	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := readAllServerLogs(serverLogPattern)
		if err != nil {
			return struct{}{}, false, err
		}

		foundDependenciesProviderLog := strings.Contains(logContent, "Creating Dependencies")
		foundToolsProviderLog := strings.Contains(logContent, "Creating Tools")

		return struct{}{}, foundDependenciesProviderLog && foundToolsProviderLog, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}

func (s *ServerWithLoggingTestSuite) TestSDK_Logging_ToolHandlerLogsToFile() {
	// Arrange
	logFolder, err := os.MkdirTemp("", "server_session") // Can't use s.T().Tempdir() because too long for socket path
	s.Require().NoError(err)
	defer s.Require().NoError(os.RemoveAll(logFolder))

	name := "World"

	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil,
		"--log-level=debug",
		"--log-folder="+logFolder,
	)

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	// Act
	_, err = session.CallTool(s.T().Context(), "tool-that-logs", map[string]any{"name": name})
	s.Require().NoError(err, "should call unstructured tool successfully")

	_, err = session.CallTool(s.T().Context(), "structured-tool-that-logs", map[string]any{"name": name})
	s.Require().NoError(err, "should call structured tool successfully")

	// Assert
	serverLogPattern := filepath.Join(logFolder, "server-*.log")

	ctx, cancel := context.WithTimeout(s.T().Context(), 2*time.Second) // Timeout for the logs to write to disk
	defer cancel()

	_, err = retry.Retry(ctx, func() (struct{}, bool, error) {
		logContent, err := readAllServerLogs(serverLogPattern)
		if err != nil {
			return struct{}{}, false, err
		}

		foundUnstructuredLogEntry := strings.Contains(logContent, "Logging from unstructured tool: "+name)
		foundStructuredLogEntry := strings.Contains(logContent, "Logging from structured tool: "+name)

		return struct{}{}, foundUnstructuredLogEntry && foundStructuredLogEntry, nil
	}, retry.NewLinearRetryStrategy(200*time.Millisecond))

	s.Require().NoError(err)
}

func readAllServerLogs(pattern string) (string, error) {
	logFiles, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	var combined strings.Builder
	for _, logFile := range logFiles {
		content, err := os.ReadFile(logFile) //nolint:gosec // G304: logFile is a controlled test path
		if err != nil {
			return "", err
		}
		combined.Write(content)
	}

	return combined.String(), nil
}
