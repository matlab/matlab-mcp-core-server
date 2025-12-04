// Copyright 2025 The MathWorks, Inc.

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
	args = append(args, "--log-level=debug")
	client := mcpclient.NewClient(ctx, s.mcpServerPath, env, args...)
	session, err := client.CreateSession(ctx)
	s.Require().NoError(err, "should create MCP session")

	dumpLogs := func(t *testing.T) {
		if t.Failed() {
			stderr := session.Stderr()
			if stderr != "" {
				t.Logf("=== MCP Server Logs (stderr) ===\n%s\n=== End MCP Server Logs ===", stderr)
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
