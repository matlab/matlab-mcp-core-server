// Copyright 2026 The MathWorks, Inc.

package provider_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel/meter/provider"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	otelmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel"
	providermocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel/meter/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	// Act
	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	expectedCollectionInterval := 1 * time.Minute

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		TelemetryCollectionInterval().
		Return(expectedCollectionInterval).
		Once()

	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Act
	result, err := factory.New(mockExporter, "matlab-mcp-server", "1.0.0")

	// Assert
	require.NoError(t, err, "Error should be nil")
	require.NotNil(t, result, "MeterProvider should not be nil")
}

func TestFactory_New_LoggerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Act
	result, err := factory.New(mockExporter, "matlab-mcp-server", "1.0.0")

	// Assert
	require.Nil(t, result)
	require.ErrorIs(t, err, expectedError)
}

func TestFactory_New_ConfigError(t *testing.T) {
	// Arrange
	testLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &providermocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &providermocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockLifecycleSignaler := &providermocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockExporter := &otelmocks.MockMetricExporter{}
	defer mockExporter.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(testLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := provider.NewFactory(mockLoggerFactory, mockConfigFactory, mockLifecycleSignaler)

	// Act
	result, err := factory.New(mockExporter, "matlab-mcp-server", "1.0.0")

	// Assert
	require.Nil(t, result, "MeterProvider should be nil when config fails")
	require.ErrorIs(t, err, expectedError)
}
