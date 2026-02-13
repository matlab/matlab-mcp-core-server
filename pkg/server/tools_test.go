// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	internaltoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/require"
)

func TestToolArray_ToInternal_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockTool1 := &server.MockTool{}
	defer mockTool1.AssertExpectations(t)

	mockTool2 := &server.MockTool{}
	defer mockTool2.AssertExpectations(t)

	mockInternalTool1 := &internaltoolsmocks.MockTool{}
	defer mockInternalTool1.AssertExpectations(t)

	mockInternalTool2 := &internaltoolsmocks.MockTool{}
	defer mockInternalTool2.AssertExpectations(t)

	mockTool1.On("toInternal", mockLoggerFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool1).
		Once()

	mockTool2.On("toInternal", mockLoggerFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool2).
		Once()

	tools := server.ToolArray{mockTool1, mockTool2}

	// Act
	result := tools.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	// Assert
	require.Equal(t, []internaltools.Tool{mockInternalTool1, mockInternalTool2}, result)
}

func TestToolArray_ToInternal_EmptyArray(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	tools := server.ToolArray{}

	// Act
	result := tools.ToInternal(mockLoggerFactory, mockConfig, mockMessageCatalog)

	// Assert
	require.Empty(t, result)
}
