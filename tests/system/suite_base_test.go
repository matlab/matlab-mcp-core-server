// Copyright 2025-2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/system/testdata"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/matlablocator"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpserver"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
	"github.com/stretchr/testify/suite"
)

// SystemTestSuite provides common setup for system tests
type SystemTestSuite struct {
	suite.Suite
	mcpServerPath string
	matlabPath    string
	testDataDir   string
	defaultEnv    []string
}

// SetupSuite runs once before all tests in a suite
func (s *SystemTestSuite) SetupSuite() {
	// Get MCP Server binary path
	mcpServerPath, err := mcpserver.NewLocator().GetPath()
	s.Require().NoError(err, "Failed to get MCP Server binary path")
	s.Require().NotEmpty(mcpServerPath, "MCP Server binary path cannot be empty")
	s.mcpServerPath = mcpServerPath

	// Get MATLAB path
	matlabPath, err := matlablocator.GetPath()
	s.Require().NoError(err, "Failed to get MATLAB path")
	s.Require().NotEmpty(matlabPath, "A non empty MATLAB path is required to run the system tests")
	s.matlabPath = matlabPath

	// Extract test assets to a temporary directory
	// This allows the system tests to be self-contained and compiled with the binary
	testDataDir := s.T().TempDir()

	// Extract files from embedded FS
	err = testdata.CopyToDir(testDataDir)
	s.Require().NoError(err, "Failed to extract test assets")

	s.testDataDir = testDataDir
}

// SetupTest runs before each test
func (s *SystemTestSuite) SetupTest() {
	// Ensure MATLAB is on PATH for each test by constructing a specific environment
	path := os.Getenv("PATH")
	path = pathcontrol.RemoveAllMATLABsFromPath(path)
	path = pathcontrol.AddToPath(path, []string{s.matlabPath})
	s.Require().Contains(path, s.matlabPath, "MATLAB directory should be in the PATH environment variable")

	// Set as the default environment for tests to use
	s.defaultEnv = pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
}

// CreateMCPSession creates an MCP client session with debug logging enabled.
// It returns the session and a function that logs MCP server stderr if the test failed.
//
// The caller is responsible for closing the session.
//
// Usage:
//
//	session, dumpLogs := s.CreateMCPSession(ctx, nil)
//	defer dumpLogs(s.T())
//	defer session.Close()
//
// If env is nil, the suite's defaultEnv is used.
func (s *SystemTestSuite) CreateMCPSession(ctx context.Context, env []string, args ...string) (*mcpclient.MCPClientSession, func(*testing.T)) {
	if env == nil {
		env = s.defaultEnv
	}

	base, err := os.MkdirTemp("", "mcp-logs-")
	s.Require().NoError(err, "should create log temp dir")
	logFolderLocation := filepath.Join(base, "logs")
	s.Require().NoError(os.MkdirAll(logFolderLocation, 0750), "should create log folder")
	args = append(args, "--log-level=debug", "--log-folder="+logFolderLocation)

	client := mcpclient.NewClient(ctx, s.mcpServerPath, env, args...)
	session, err := client.CreateSession(ctx)
	s.Require().NoError(err, "should create MCP session")

	dumpLogs := func(t *testing.T) {
		stderr := session.Stderr()
		if stderr != "" {
			t.Logf("=== MCP Server Logs (stderr) ===\n%s\n=== End MCP Server Logs ===", stderr)
		}
		serverLogPattern := filepath.Join(logFolderLocation, "server-*.log")
		serverLogFiles, err := filepath.Glob(serverLogPattern)
		if err != nil || len(serverLogFiles) == 0 {
			t.Logf("Server log file not found: %s", serverLogPattern)
		} else {
			for _, logFile := range serverLogFiles {
				serverLog, err := os.ReadFile(logFile) //#nosec G304 - log path constructed from known test directory
				if err != nil {
					t.Logf("Failed to read server log file: %s", err.Error())
				} else {
					t.Logf("=== MCP Server Log File (%s) ===\n%s\n=== End MCP Server Log File ===", filepath.Base(logFile), serverLog)
				}
			}
		}
		watchdogLogPattern := filepath.Join(logFolderLocation, "watchdog-*.log")
		watchdogLogFiles, err := filepath.Glob(watchdogLogPattern)
		if err != nil || len(watchdogLogFiles) == 0 {
			t.Logf("Watchdog log file not found: %s", watchdogLogPattern)
		} else {
			for _, logFile := range watchdogLogFiles {
				watchdogLog, err := os.ReadFile(logFile) //#nosec G304 - log path constructed from known test directory
				if err != nil {
					t.Logf("Failed to read watchdog log file: %s", err.Error())
				} else {
					t.Logf("=== MCP Watchdog Log File (%s) ===\n%s\n=== End MCP Watchdog Log File ===", filepath.Base(logFile), watchdogLog)
				}
			}
		}
	}
	return session, dumpLogs
}

// Test file paths
func (s *SystemTestSuite) problematicCodePath() string {
	return filepath.Join(s.testDataDir, "problematic_code.m")
}

func (s *SystemTestSuite) testScriptPath() string {
	return filepath.Join(s.testDataDir, "test_script.m")
}

func (s *SystemTestSuite) testMathFunctionsPath() string {
	return filepath.Join(s.testDataDir, "test_math_functions.m")
}

// matlabRoot returns the MATLAB root directory (without /bin)
// This is needed for the --matlab-root flag
func (s *SystemTestSuite) matlabRoot() string {
	// matlabPath is expected to be the bin directory, so get its parent
	return filepath.Dir(s.matlabPath)
}
