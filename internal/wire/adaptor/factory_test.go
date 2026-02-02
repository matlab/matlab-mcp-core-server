// Copyright 2026 The MathWorks, Inc.

package adaptor_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
	adaptormocks "github.com/matlab/matlab-mcp-core-server/mocks/wire/adaptor"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	// No dependencies needed

	// Act
	factory := adaptor.NewFactory()

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockDefinition := &adaptormocks.MockApplicationDefinition{}
	mockDefinition.AssertExpectations(t)

	factory := adaptor.NewFactory()

	// Act
	application := factory.New(mockDefinition)

	// Assert
	assert.NotNil(t, application)
}
