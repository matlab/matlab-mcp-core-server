// Copyright 2026 The MathWorks, Inc.

package definition_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	toolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	"github.com/stretchr/testify/require"
)

func TestDefinition_Name_HappyPath(t *testing.T) {
	// Arrange
	expectedName := "my-definition"
	def := definition.New(expectedName, "", "", nil)

	// Act
	result := def.Name()

	// Assert
	require.Equal(t, expectedName, result)
}

func TestDefinition_Title_HappyPath(t *testing.T) {
	// Arrange
	expectedTitle := "My Definition Title"
	def := definition.New("", expectedTitle, "", nil)

	// Act
	result := def.Title()

	// Assert
	require.Equal(t, expectedTitle, result)
}

func TestDefinition_Instructions_HappyPath(t *testing.T) {
	// Arrange
	expectedInstructions := "These are the instructions"
	def := definition.New("", "", expectedInstructions, nil)

	// Act
	result := def.Instructions()

	// Assert
	require.Equal(t, expectedInstructions, result)
}

func TestDefinition_Tools_HappyPath(t *testing.T) {
	// Arrange
	mockTool := &toolsmocks.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockLoggerFactory := &definitionmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedTools := []tools.Tool{mockTool}

	toolsProvider := func(loggerFactory definition.LoggerFactory) []tools.Tool {
		require.Equal(t, mockLoggerFactory, loggerFactory)
		return expectedTools
	}

	def := definition.New("", "", "", toolsProvider)

	// Act
	result := def.Tools(mockLoggerFactory)

	// Assert
	require.Equal(t, expectedTools, result)
}

func TestDefinition_Tools_NilProvider(t *testing.T) {
	// Arrange
	mockLoggerFactory := &definitionmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	def := definition.New("", "", "", nil)

	// Act
	result := def.Tools(mockLoggerFactory)

	// Assert
	require.Nil(t, result)
}
