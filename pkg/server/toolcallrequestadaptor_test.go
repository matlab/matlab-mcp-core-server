// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/require"
)

func TestNewToolCallRequestAdaptor_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	// Act
	adaptor := server.NewToolCallRequestAdaptor(mockLogger, mockConfig, mockMessageCatalog)

	// Assert
	require.NotNil(t, adaptor)
}

func TestToolCallRequestAdaptor_Logger_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedMessage := "test info message"

	adaptor := server.NewToolCallRequestAdaptor(mockLogger, mockConfig, mockMessageCatalog)

	// Act
	adaptor.Logger().Info(expectedMessage)

	// Assert
	infoLogs := mockLogger.InfoLogs()
	_, found := infoLogs[expectedMessage]
	require.True(t, found)
}

func TestToolCallRequestAdaptor_Config_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedKey := "test-key"
	expectedValue := "test-value"

	mockConfig.EXPECT().
		Get(expectedKey).
		Return(expectedValue, nil).
		Once()

	adaptor := server.NewToolCallRequestAdaptor(mockLogger, mockConfig, mockMessageCatalog)

	// Act
	result, err := adaptor.Config().Get(expectedKey, "")

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedValue, result)
}
