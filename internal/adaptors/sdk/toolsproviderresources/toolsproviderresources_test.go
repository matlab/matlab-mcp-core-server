// Copyright 2026 The MathWorks, Inc.

package toolsproviderresources_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/toolsproviderresources"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/definition"
	basetoolmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/mcp/tools/basetool"
	publictypesmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/publictypes"
	toolsproviderresourcesmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/sdk/toolsproviderresources"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/require"
)

type TestDependencies struct {
	Value string
}

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &toolsproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	// Act
	factory := toolsproviderresources.NewFactory[*TestDependencies](mockLoggerFactory)

	// Assert
	require.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &toolsproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &configmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockBaseToolLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockBaseToolLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(nil).
		Once()

	internalResources := definition.NewToolsProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
		&TestDependencies{Value: "test"},
		mockBaseToolLoggerFactory,
		mockTelemetryFactory,
	)

	// Act
	resources := toolsproviderresources.NewFactory[*TestDependencies](mockLoggerFactory).New(internalResources)

	// Assert
	require.NotNil(t, resources)
}

func TestFactory_New_Logger(t *testing.T) {
	// Arrange
	mockLoggerFactory := &toolsproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &configmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockBaseToolLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockBaseToolLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	expectedLogger := &publictypesmocks.MockLogger{}
	defer expectedLogger.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(expectedLogger).
		Once()

	internalResources := definition.NewToolsProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
		&TestDependencies{Value: "test"},
		mockBaseToolLoggerFactory,
		mockTelemetryFactory,
	)

	// Act
	resources := toolsproviderresources.NewFactory[*TestDependencies](mockLoggerFactory).New(internalResources)

	// Assert
	require.Equal(t, expectedLogger, resources.Logger())
}

func TestFactory_New_Dependencies_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &toolsproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &configmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockBaseToolLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockBaseToolLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	expectedDependencies := &TestDependencies{Value: "test"}

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(nil).
		Once()

	internalResources := definition.NewToolsProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
		expectedDependencies,
		mockBaseToolLoggerFactory,
		mockTelemetryFactory,
	)

	// Act
	resources := toolsproviderresources.NewFactory[*TestDependencies](mockLoggerFactory).New(internalResources)

	// Assert
	require.Equal(t, expectedDependencies, resources.Dependencies())
}

func TestFactory_New_Dependencies_Nil(t *testing.T) {
	// Arrange
	mockLoggerFactory := &toolsproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &configmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockBaseToolLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockBaseToolLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(nil).
		Once()

	internalResources := definition.NewToolsProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
		nil,
		mockBaseToolLoggerFactory,
		mockTelemetryFactory,
	)

	// Act
	resources := toolsproviderresources.NewFactory[*TestDependencies](mockLoggerFactory).New(internalResources)

	// Assert
	require.NotNil(t, resources)
	require.Nil(t, resources.Dependencies())
}

func TestFactory_New_Dependencies_CastFailure(t *testing.T) {
	// Arrange
	mockLoggerFactory := &toolsproviderresourcesmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockInternalLogger := &entitiesmocks.MockLogger{}
	defer mockInternalLogger.AssertExpectations(t)

	mockInternalConfig := &configmocks.MockGenericConfig{}
	defer mockInternalConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockBaseToolLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockBaseToolLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockInternalLogger.EXPECT().
		Error("Dependencies type cast failed, using zero value").
		Return().
		Once()

	mockLoggerFactory.EXPECT().
		New(mockInternalLogger).
		Return(nil).
		Once()

	internalResources := definition.NewToolsProviderResources(
		mockInternalLogger,
		mockInternalConfig,
		mockMessageCatalog,
		"wrong type",
		mockBaseToolLoggerFactory,
		mockTelemetryFactory,
	)

	// Act
	resources := toolsproviderresources.NewFactory[*TestDependencies](mockLoggerFactory).New(internalResources)

	// Assert
	require.Nil(t, resources.Dependencies())
}
