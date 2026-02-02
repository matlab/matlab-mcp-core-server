// Copyright 2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/functional/sdk/testbinaries"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// ServerWithCustomToolsTestSuite tests SDK custom tools functionalities.
type ServerWithCustomToolsTestSuite struct {
	suite.Suite

	serverDetails testbinaries.ServerDetails
}

// SetupSuite runs once before all tests in a suite
func (s *ServerWithCustomToolsTestSuite) SetupSuite() {
	s.serverDetails = testbinaries.BuildServerWithCustomTools(s.T())
}

func TestServerWithCustomToolsTestSuite(t *testing.T) {
	suite.Run(t, new(ServerWithCustomToolsTestSuite))
}

func (s *ServerWithCustomToolsTestSuite) TestSDK_CustomTools_UnstructuredContentOutput_HappyPath() {
	// Arrange
	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil, "--log-level=debug")

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	name := "World"
	expectedTextOutput := "Hello " + name

	// Act
	result, err := session.CallTool(s.T().Context(), "greet", map[string]any{"name": "World"})

	// Assert
	s.Require().NoError(err, "should call tool successfully")

	textContent, err := session.GetTextContent(result)
	s.Require().NoError(err, "should get text content")
	s.Require().Equal(expectedTextOutput, textContent, "should return greeting message")
}

func (s *ServerWithCustomToolsTestSuite) TestSDK_CustomTools_StructuredContentOutput_HappyPath() {
	// Arrange
	client := mcpclient.NewClient(s.T().Context(), s.serverDetails.BinaryLocation(), nil, "--log-level=debug")

	session, err := client.CreateSession(s.T().Context())
	s.Require().NoError(err, "should create MCP session")
	defer func() {
		s.Require().NoError(session.Close(), "closing session should not error")
	}()

	name := "World"
	expectedResponse := "Hello " + name

	// Act
	result, err := session.CallTool(s.T().Context(), "greet-structured", map[string]any{"name": "World"})

	// Assert
	s.Require().NoError(err, "should call tool successfully")

	var output struct {
		Response string `json:"response"`
	}
	s.Require().NoError(session.UnmarshalStructuredContent(result, &output), "should unmarshal structured content")
	s.Require().Equal(expectedResponse, output.Response, "should return greeting message")
}
