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

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	// Act
	factory := sdk.NewFactory(mockConfigFactory, mockDefinition)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_NewServer_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	expectedVersion := "1.0.0"
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "test instructions"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockDefinition.EXPECT().
		Name().
		Return(expectedName).
		Once()

	mockDefinition.EXPECT().
		Title().
		Return(expectedTitle).
		Once()

	mockDefinition.EXPECT().
		Instructions().
		Return(expectedInstructions).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.NoError(t, err, "NewServer should not return an error")
	assert.NotNil(t, server, "Server should not be nil")
}

func TestFactory_NewServer_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, server, "Server should be nil when error occurs")
}
