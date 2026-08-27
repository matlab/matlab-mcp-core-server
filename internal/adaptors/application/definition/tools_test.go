// Copyright 2026 The MathWorks, Inc.

package definition_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/definition"
	basetoolmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/stretchr/testify/require"
)

func TestNewToolsProviderResources_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTelemetryFactory := &basetoolmocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	dependencies := &struct{ Value string }{Value: "test"}

	// Act
	result := definition.NewToolsProviderResources(
		mockLogger,
		mockConfig,
		mockMessageCatalog,
		dependencies,
		mockLoggerFactory,
		mockTelemetryFactory,
	)

	// Assert
	require.Equal(t, mockLogger, result.Logger)
	require.Equal(t, mockConfig, result.Config)
	require.Equal(t, mockMessageCatalog, result.MessageCatalog)
	require.Equal(t, dependencies, result.Dependencies)
	require.Equal(t, mockLoggerFactory, result.LoggerFactory)
}
