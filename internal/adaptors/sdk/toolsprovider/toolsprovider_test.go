// Copyright 2026 The MathWorks, Inc.

package toolsprovider_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/publictypes"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/toolsprovider"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/definition"
	internaltoolsmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/mcp/tools"
	basetoolmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/mcp/tools/basetool"
	publictypesmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/publictypes"
	toolsmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/tools"
	toolsprovidermocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/toolsprovider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockResourcesFactory := &toolsprovidermocks.MockResourcesFactory[struct{}]{}
	defer mockResourcesFactory.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsprovidermocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	// Act
	factory := toolsprovider.NewFactory(mockResourcesFactory, mockToolCallRequestFactory)

	// Assert
	require.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	type TestDependencies struct{}

	mockResourcesFactory := &toolsprovidermocks.MockResourcesFactory[*TestDependencies]{}
	defer mockResourcesFactory.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsprovidermocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockResources := &publictypesmocks.MockToolsProviderResources[*TestDependencies]{}
	defer mockResources.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockInternalTool := &internaltoolsmocks.MockTool{}
	defer mockInternalTool.AssertExpectations(t)

	mockTool := &mockConvertibleTool{}
	defer mockTool.AssertExpectations(t)

	expectedInternalResources := definition.ToolsProviderResources{
		Logger:           testutils.NewInspectableLogger(),
		Config:           mockConfig,
		MessageCatalog:   mockMessageCatalog,
		LoggerFactory:    mockLoggerFactory,
		TelemetryFactory: mockTelemetryFactory,
	}

	mockResourcesFactory.EXPECT().
		New(expectedInternalResources).
		Return(mockResources).
		Once()

	mockTool.EXPECT().
		ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool).
		Once()

	provider := toolsprovider.ToolsProvider[*TestDependencies](func(resources publictypes.ToolsProviderResources[*TestDependencies]) []publictypes.Tool {
		assert.Equal(t, mockResources, resources)
		return []publictypes.Tool{mockTool}
	})

	// Act
	internalProvider := toolsprovider.NewFactory(mockResourcesFactory, mockToolCallRequestFactory).New(provider)
	tools := internalProvider(expectedInternalResources)

	// Assert
	require.Len(t, tools, 1)
	require.Equal(t, mockInternalTool, tools[0])
}

func TestFactory_New_EmptyTools(t *testing.T) {
	// Arrange
	type TestDependencies struct{}

	mockResourcesFactory := &toolsprovidermocks.MockResourcesFactory[*TestDependencies]{}
	defer mockResourcesFactory.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsprovidermocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockResources := &publictypesmocks.MockToolsProviderResources[*TestDependencies]{}
	defer mockResources.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedInternalResources := definition.ToolsProviderResources{
		Logger:           testutils.NewInspectableLogger(),
		Config:           mockConfig,
		MessageCatalog:   mockMessageCatalog,
		LoggerFactory:    mockLoggerFactory,
		TelemetryFactory: mockTelemetryFactory,
	}

	mockResourcesFactory.EXPECT().
		New(expectedInternalResources).
		Return(mockResources).
		Once()

	provider := toolsprovider.ToolsProvider[*TestDependencies](func(resources publictypes.ToolsProviderResources[*TestDependencies]) []publictypes.Tool {
		return []publictypes.Tool{}
	})

	// Act
	internalProvider := toolsprovider.NewFactory(mockResourcesFactory, mockToolCallRequestFactory).New(provider)
	tools := internalProvider(expectedInternalResources)

	// Assert
	require.Empty(t, tools)
}

func TestFactory_New_MultipleTools(t *testing.T) {
	// Arrange
	type TestDependencies struct{}

	mockResourcesFactory := &toolsprovidermocks.MockResourcesFactory[*TestDependencies]{}
	defer mockResourcesFactory.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsprovidermocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockResources := &publictypesmocks.MockToolsProviderResources[*TestDependencies]{}
	defer mockResources.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockInternalTool1 := &internaltoolsmocks.MockTool{}
	defer mockInternalTool1.AssertExpectations(t)

	mockInternalTool2 := &internaltoolsmocks.MockTool{}
	defer mockInternalTool2.AssertExpectations(t)

	mockTool1 := &mockConvertibleTool{}
	defer mockTool1.AssertExpectations(t)

	mockTool2 := &mockConvertibleTool{}
	defer mockTool2.AssertExpectations(t)

	expectedInternalResources := definition.ToolsProviderResources{
		Logger:           testutils.NewInspectableLogger(),
		Config:           mockConfig,
		MessageCatalog:   mockMessageCatalog,
		LoggerFactory:    mockLoggerFactory,
		TelemetryFactory: mockTelemetryFactory,
	}

	mockResourcesFactory.EXPECT().
		New(expectedInternalResources).
		Return(mockResources).
		Once()

	mockTool1.EXPECT().
		ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool1).
		Once()

	mockTool2.EXPECT().
		ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool2).
		Once()

	provider := toolsprovider.ToolsProvider[*TestDependencies](func(resources publictypes.ToolsProviderResources[*TestDependencies]) []publictypes.Tool {
		return []publictypes.Tool{mockTool1, mockTool2}
	})

	// Act
	internalProvider := toolsprovider.NewFactory(mockResourcesFactory, mockToolCallRequestFactory).New(provider)
	tools := internalProvider(expectedInternalResources)

	// Assert
	require.Len(t, tools, 2)
	require.Equal(t, mockInternalTool1, tools[0])
	require.Equal(t, mockInternalTool2, tools[1])
}

func TestFactory_New_NilProvider(t *testing.T) {
	// Arrange
	mockResourcesFactory := &toolsprovidermocks.MockResourcesFactory[struct{}]{}
	defer mockResourcesFactory.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsprovidermocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	// Act
	internalProvider := toolsprovider.NewFactory(mockResourcesFactory, mockToolCallRequestFactory).New(nil)
	result := internalProvider(definition.ToolsProviderResources{})

	// Assert
	require.Nil(t, result)
}

func TestFactory_New_NonConvertibleToolSkipped(t *testing.T) {
	// Arrange
	type TestDependencies struct{}

	mockResourcesFactory := &toolsprovidermocks.MockResourcesFactory[*TestDependencies]{}
	defer mockResourcesFactory.AssertExpectations(t)

	mockToolCallRequestFactory := &toolsprovidermocks.MockToolCallRequestFactory{}
	defer mockToolCallRequestFactory.AssertExpectations(t)

	mockResources := &publictypesmocks.MockToolsProviderResources[*TestDependencies]{}
	defer mockResources.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockInternalTool := &internaltoolsmocks.MockTool{}
	defer mockInternalTool.AssertExpectations(t)

	mockTool := &mockConvertibleTool{}
	defer mockTool.AssertExpectations(t)

	nonConvertible := &nonConvertibleTool{}

	expectedInternalResources := definition.ToolsProviderResources{
		Logger:           testutils.NewInspectableLogger(),
		Config:           mockConfig,
		MessageCatalog:   mockMessageCatalog,
		LoggerFactory:    mockLoggerFactory,
		TelemetryFactory: mockTelemetryFactory,
	}

	mockResourcesFactory.EXPECT().
		New(expectedInternalResources).
		Return(mockResources).
		Once()

	mockTool.EXPECT().
		ToInternal(mockToolCallRequestFactory, mockLoggerFactory, mockTelemetryFactory, mockConfig, mockMessageCatalog).
		Return(mockInternalTool).
		Once()

	provider := toolsprovider.ToolsProvider[*TestDependencies](func(resources publictypes.ToolsProviderResources[*TestDependencies]) []publictypes.Tool {
		assert.Equal(t, mockResources, resources)
		return []publictypes.Tool{nonConvertible, mockTool}
	})

	// Act
	internalProvider := toolsprovider.NewFactory(mockResourcesFactory, mockToolCallRequestFactory).New(provider)
	tools := internalProvider(expectedInternalResources)

	// Assert
	require.Len(t, tools, 1)
	require.Equal(t, mockInternalTool, tools[0])
}

// In mockery, we can't struct embed a sealing struct, so we have to do it manually here
type mockConvertibleTool struct {
	toolsmocks.MockConvertibleTool
	publictypes.ToolSeal
}

type nonConvertibleTool struct {
	publictypes.ToolSeal
}
