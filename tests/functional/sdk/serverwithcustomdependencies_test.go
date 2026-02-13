// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// ServerWithCustomDependenciesTestSuite tests SDK custom dependencies functionalities.
type ServerWithCustomDependenciesTestSuite struct {
	suite.Suite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithCustomDependenciesTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithCustomDependencies(s.T())
}

func TestServerWithCustomDependenciesTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithCustomDependenciesTestSuite))
}

func (s *ServerWithCustomDependenciesTestSuite) TestSDK_CustomDependencies_ToolUsesDependency_HappyPath() {
	// Connect to a session
	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil,
		"--log-level=debug",
	)

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	name := "World"
	expectedTextOutput := "Service Hello " + name

	// Call the tool
	result, err := session.CallTool(s.T().Context(), "greet", map[string]any{"name": name})
	s.Require().NoError(err, "should call tool successfully")

	textContent, err := session.GetTextContent(result)
	s.Require().NoError(err, "should get text content")
	s.Require().Equal(expectedTextOutput, textContent, "should return greeting message using dependency")
}
