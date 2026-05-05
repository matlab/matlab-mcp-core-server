// Copyright 2026 The MathWorks, Inc.

package functional_test

import (
	"os"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/pathcontrol"
)

// MockMATLABTestSuite extends the base with an environment that has mock MATLAB
// on PATH, suitable for testing local-install mode.
type MockMATLABTestSuite struct {
	MockMATLABBaseSuite
	defaultEnv []string
}

func (s *MockMATLABTestSuite) SetupSuite() {
	s.MockMATLABBaseSuite.SetupSuite()

	mockMATLABBinDir := filepath.Join(s.installation.MATLABRoot, "bin")
	path := pathcontrol.RemoveAllMATLABsFromPath(os.Getenv("PATH"))
	path = pathcontrol.AddToPath(path, []string{mockMATLABBinDir})

	env := pathcontrol.UpdateEnvEntry(os.Environ(), "PATH", path)
	env = pathcontrol.UpdateEnvEntry(env, "MW_MCP_SERVER_EMBEDDED_CONNECTOR_DETAILS_TIMEOUT", "10s")

	s.defaultEnv = env
}

// CreateSession creates a mock MATLAB session with debug logging enabled.
// Additional CLI args (e.g. "--extension-file=...") are forwarded to the server.
// Usage:
//
//	session, err := s.CreateSession(mockmatlab.HappyConfig())
//	s.Require().NoError(err)
//	defer s.CleanupSession(session, true)
func (s *MockMATLABTestSuite) CreateSession(cfg mockmatlab.Config, args ...string) (*mcpclient.LoggedSession, error) {
	value, err := cfg.ToEnvValue()
	s.Require().NoError(err, "failed to serialize mock config")
	env := pathcontrol.UpdateEnvEntry(s.defaultEnv, mockmatlab.EnvMockMATLABConfig, value)

	return s.createLoggedSession(s.T().Context(), env, args...)
}
