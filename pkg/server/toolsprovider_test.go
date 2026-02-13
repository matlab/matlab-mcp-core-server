// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	internaltools "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	internaltoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolsProvider_toInternal_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTool := &server.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockInternalTool := &internaltoolsmocks.MockTool{}
	defer mockInternalTool.AssertExpectations(t)

	expectedMessage := "test message"

	mockLogger.EXPECT().
		Info(expectedMessage).
		Once()

	type TestDependencies struct{}
	expectedDependencies := &TestDependencies{}

	provider := server.ToolsProvider[*TestDependencies](func(resources server.ToolsProviderResources[*TestDependencies]) []server.Tool {
		resources.Logger().Info(expectedMessage)

		assert.Equal(t, expectedDependencies, resources.Dependencies())
		return []server.Tool{mockTool}
	})

	mockTool.On("toInternal", mockLoggerFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool).
		Once()

	// Act
	internalProvider := provider.ToInternal()
	tools := internalProvider(definition.NewToolsProviderResources(
		mockLogger,
		mockConfig,
		mockMessageCatalog,
		expectedDependencies,
		mockLoggerFactory,
	))

	// Assert
	require.Equal(t, []internaltools.Tool{mockInternalTool}, tools)
}

func TestToolsProvider_toInternal_NilProvider(t *testing.T) {
	// Arrange
	var provider server.ToolsProvider[struct{}]

	// Act
	internalProvider := provider.ToInternal()
	tools := internalProvider(definition.ToolsProviderResources{})

	// Assert
	require.Nil(t, tools)
}

func TestToolsProvider_toInternal_DependenciesCastFailure(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	type TestDependencies struct{}

	provider := server.ToolsProvider[*TestDependencies](func(resources server.ToolsProviderResources[*TestDependencies]) []server.Tool {
		assert.Nil(t, resources.Dependencies()) // nil when cast fails for pointer type
		return nil
	})

	// Act
	internalProvider := provider.ToInternal()
	tools := internalProvider(definition.NewToolsProviderResources(
		mockLogger,
		mockConfig,
		mockMessageCatalog,
		"wrong type", // not *TestDependencies
		mockLoggerFactory,
	))

	// Assert
	require.Empty(t, tools)
}
