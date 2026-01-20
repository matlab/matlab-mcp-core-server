// Copyright 2025-2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/sdk"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	// Act
	factory := sdk.NewFactory(mockConfigFactory)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_NewServer_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	expectedVersion := "1.0.0"
	expectedName := "test-server"
	expectedInstructions := "test instructions"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	factory := sdk.NewFactory(mockConfigFactory)

	// Act
	server, err := factory.NewServer(expectedName, expectedInstructions)

	// Assert
	require.NoError(t, err, "NewServer should not return an error")
	assert.NotNil(t, server, "Server should not be nil")
}

func TestFactory_NewServer_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	expectedName := "test-server"
	expectedInstructions := "test instructions"
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory)

	// Act
	server, err := factory.NewServer(expectedName, expectedInstructions)

	// Assert
	require.ErrorIs(t, err, expectedError, "NewServer should return the error from Config")
	assert.Nil(t, server, "Server should be nil when error occurs")
}
