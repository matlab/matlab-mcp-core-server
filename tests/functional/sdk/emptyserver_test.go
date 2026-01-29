// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"os/exec"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// EmptyServerTestSuite tests SDK definition functionnalities, using a mock MATLAB.
type EmptyServerTestSuite struct {
	suite.Suite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *EmptyServerTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildEmptyServer(s.T())
}

func TestEmptyServerTestSuite(t *testing.T) {
	suite.Run(t, new(EmptyServerTestSuite))
}

func (s *EmptyServerTestSuite) TestSDK_CLI_Version() {
	// Arrange

	// Act
	cmd := exec.Command(s.serverDetails.BinaryLocation(), "--version") //nolint:gosec // Trusted test path
	output, err := cmd.CombinedOutput()

	// Assert
	s.Require().NoError(err, "version flag should execute successfully")
	s.Require().Contains(string(output), s.serverDetails.ModuleName(), "should display server package path")
}

func (s *EmptyServerTestSuite) TestSDK_ServerDefinition_NameTitleAndInstruction() {
	// Arrange
	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil, "--log-level=debug")

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	// Act
	result := session.InitializeResult()

	// Assert
	s.Require().NotNil(result)
	s.Require().Equal(s.serverDetails.Name(), result.ServerInfo.Name)
	s.Require().Equal(s.serverDetails.Title(), result.ServerInfo.Title)
	s.Require().Equal(s.serverDetails.Instructions(), result.Instructions)
}
