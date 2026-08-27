// Copyright 2026 The MathWorks, Inc.

package telemetry_test

import (
	"regexp"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	telemetrymocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry"
	instrumentsmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry/otel/instruments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/noop"
)

func TestNewOTELTelemetry_HappyPath(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(mockInt64Counter, nil).
		Once()

	// Act
	result, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, nil, mockOSLayer, mockOSVersionProvider, mockServerDefinition)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestNewOTELTelemetry_InstrumentCreationFails(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")
	instrumentError := assert.AnError
	expectedError := messages.New_StartupErrors_TelemetryInitializationFailed_Error()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(nil, instrumentError).
		Once()

	// Act
	result, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, nil, mockOSLayer, mockOSVersionProvider, mockServerDefinition)

	// Assert
	require.Nil(t, result)
	require.Equal(t, expectedError, err)
}

func TestNewOTELTelemetry_ClientConnectionCounterCreationFails(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")
	instrumentError := assert.AnError
	expectedError := messages.New_StartupErrors_TelemetryInitializationFailed_Error()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(nil, instrumentError).
		Once()

	// Act
	result, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, nil, mockOSLayer, mockOSVersionProvider, mockServerDefinition)

	// Assert
	require.Nil(t, result)
	require.Equal(t, expectedError, err)
}

func TestNewOTELTelemetry_ToolCallRequestCounterCreationFails(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockServerStartCounter := &instrumentsmocks.MockInt64Counter{}
	defer mockServerStartCounter.AssertExpectations(t)

	mockClientConnectionCounter := &instrumentsmocks.MockInt64Counter{}
	defer mockClientConnectionCounter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")
	instrumentError := assert.AnError
	expectedError := messages.New_StartupErrors_TelemetryInitializationFailed_Error()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockServerStartCounter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockClientConnectionCounter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(nil, instrumentError).
		Once()

	// Act
	result, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, nil, mockOSLayer, mockOSVersionProvider, mockServerDefinition)

	// Assert
	require.Nil(t, result)
	require.Equal(t, expectedError, err)
}

func TestOTELTelemetry_RecordServerStart_HappyPath(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockDirectory := &telemetrymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")

	expectedVersion := "v1.2.3"
	expectedName := "matlab-mcp-server"
	expectedOS := "linux"
	expectedSpecifiedParameters := []string{"disable-telemetry", "log-level"}
	expectedOSVersion := "Debian GNU/Linux 12"
	expectedInstanceID := "test-instance-id"
	expectedConfigDetails := `{"key":"value"}`

	expectedAttributes := []attribute.KeyValue{
		attribute.String("server.instance_id", expectedInstanceID),
		attribute.String("server.name", expectedName),
		attribute.String("server.version", expectedVersion),
		attribute.StringSlice("server.specified_parameters", expectedSpecifiedParameters),
		attribute.String("server.config_details", expectedConfigDetails),
		attribute.String("server.os", expectedOS),
		attribute.String("server.os_version", expectedOSVersion),
	}

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(mockInt64Counter, nil).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedInstanceID).
		Once()

	mockServerDefinition.EXPECT().
		Name().
		Return(expectedName).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockConfig.EXPECT().
		SpecifiedParameters().
		Return(expectedSpecifiedParameters).
		Once()

	mockConfig.EXPECT().
		AsPIISafeJSONString().
		Return(expectedConfigDetails).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return(expectedOS).
		Once()

	mockOSVersionProvider.EXPECT().
		Version().
		Return(expectedOSVersion, nil).
		Once()

	mockInt64Counter.EXPECT().
		Add(mock.Anything, int64(1), expectedAttributes).
		Once()

	otelTelemetry, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, mockDirectory, mockOSLayer, mockOSVersionProvider, mockServerDefinition)
	require.NoError(t, err)

	// Act
	otelTelemetry.RecordServerStart(t.Context())

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestOTELTelemetry_RecordServerStart_WatchdogModeSkips(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockDirectory := &telemetrymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(mockInt64Counter, nil).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	otelTelemetry, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, mockDirectory, mockOSLayer, mockOSVersionProvider, mockServerDefinition)
	require.NoError(t, err)

	// Act
	otelTelemetry.RecordServerStart(t.Context())

	// Assert
	// No counter Add() call should occur - verified via deferred mock expectations.
}

func TestOTELTelemetry_RecordServerStart_OSVersionError(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockInt64Counter := &instrumentsmocks.MockInt64Counter{}
	defer mockInt64Counter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockDirectory := &telemetrymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")

	expectedVersion := "v1.2.3"
	expectedName := "matlab-mcp-server"
	expectedOS := "linux"
	expectedSpecifiedParameters := []string{"disable-telemetry", "log-level"}
	expectedInstanceID := "test-instance-id"
	expectedConfigDetails := `{"key":"value"}`

	expectedAttributes := []attribute.KeyValue{
		attribute.String("server.instance_id", expectedInstanceID),
		attribute.String("server.name", expectedName),
		attribute.String("server.version", expectedVersion),
		attribute.StringSlice("server.specified_parameters", expectedSpecifiedParameters),
		attribute.String("server.config_details", expectedConfigDetails),
		attribute.String("server.os", expectedOS),
		attribute.String("server.os_version", ""),
	}

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockInt64Counter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(mockInt64Counter, nil).
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedInstanceID).
		Once()

	mockServerDefinition.EXPECT().
		Name().
		Return(expectedName).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockConfig.EXPECT().
		SpecifiedParameters().
		Return(expectedSpecifiedParameters).
		Once()

	mockConfig.EXPECT().
		AsPIISafeJSONString().
		Return(expectedConfigDetails).
		Once()

	mockOSLayer.EXPECT().
		GOOS().
		Return(expectedOS).
		Once()

	mockOSVersionProvider.EXPECT().
		Version().
		Return("", assert.AnError).
		Once()

	mockInt64Counter.EXPECT().
		Add(mock.Anything, int64(1), expectedAttributes).
		Once()

	otelTelemetry, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, mockDirectory, mockOSLayer, mockOSVersionProvider, mockServerDefinition)
	require.NoError(t, err)

	// Act
	otelTelemetry.RecordServerStart(t.Context())

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestOTELTelemetry_RecordClientConnection_HappyPath(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockServerStartCounter := &instrumentsmocks.MockInt64Counter{}
	defer mockServerStartCounter.AssertExpectations(t)

	mockClientConnectionCounter := &instrumentsmocks.MockInt64Counter{}
	defer mockClientConnectionCounter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockDirectory := &telemetrymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")

	expectedInstanceID := "test-instance-id"
	expectedClientName := "vscode"
	expectedClientTitle := "Visual Studio Code"
	expectedClientURL := "https://code.visualstudio.com"
	expectedClientVersion := "1.0.0"
	expectedCapabilities := []string{"roots", "sampling"}
	expectedCapabilitiesJSON := `{"roots":{"listChanged":true},"sampling":{}}`

	expectedAttributes := []attribute.KeyValue{
		attribute.String("server.instance_id", expectedInstanceID),
		attribute.String("client.name", expectedClientName),
		attribute.String("client.title", expectedClientTitle),
		attribute.String("client.website_url", expectedClientURL),
		attribute.String("client.version", expectedClientVersion),
		attribute.StringSlice("client.capabilities", expectedCapabilities),
		attribute.String("client.capabilities_details", expectedCapabilitiesJSON),
	}

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockServerStartCounter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockClientConnectionCounter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(mockClientConnectionCounter, nil).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedInstanceID).
		Once()

	mockClientConnectionCounter.EXPECT().
		Add(mock.Anything, int64(1), expectedAttributes).
		Once()

	otelTelemetry, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, mockDirectory, mockOSLayer, mockOSVersionProvider, mockServerDefinition)
	require.NoError(t, err)

	info := telemetry.ClientConnectionInfo{
		Name:             expectedClientName,
		Title:            expectedClientTitle,
		WebsiteURL:       expectedClientURL,
		Version:          expectedClientVersion,
		Capabilities:     expectedCapabilities,
		CapabilitiesJSON: expectedCapabilitiesJSON,
	}

	// Act
	otelTelemetry.RecordClientConnection(t.Context(), info)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestOTELTelemetry_RecordClientConnection_EmptyInfo(t *testing.T) {
	// Arrange
	mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
	defer mockInstrumentFactory.AssertExpectations(t)

	mockServerStartCounter := &instrumentsmocks.MockInt64Counter{}
	defer mockServerStartCounter.AssertExpectations(t)

	mockClientConnectionCounter := &instrumentsmocks.MockInt64Counter{}
	defer mockClientConnectionCounter.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockOSLayer := &telemetrymocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockServerDefinition := &telemetrymocks.MockDefinition{}
	defer mockServerDefinition.AssertExpectations(t)

	mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
	defer mockOSVersionProvider.AssertExpectations(t)

	mockDirectory := &telemetrymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	testLogger := testutils.NewInspectableLogger()
	meter := noop.NewMeterProvider().Meter("test")

	expectedInstanceID := "test-instance-id"

	expectedAttributes := []attribute.KeyValue{
		attribute.String("server.instance_id", expectedInstanceID),
		attribute.String("client.name", ""),
		attribute.String("client.title", ""),
		attribute.String("client.website_url", ""),
		attribute.String("client.version", ""),
		attribute.StringSlice("client.capabilities", nil),
		attribute.String("client.capabilities_details", ""),
	}

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
		Return(mockServerStartCounter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
		Return(mockClientConnectionCounter, nil).
		Once()

	mockInstrumentFactory.EXPECT().
		NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
		Return(mockClientConnectionCounter, nil).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(expectedInstanceID).
		Once()

	mockClientConnectionCounter.EXPECT().
		Add(mock.Anything, int64(1), expectedAttributes).
		Once()

	otelTelemetry, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, mockDirectory, mockOSLayer, mockOSVersionProvider, mockServerDefinition)
	require.NoError(t, err)

	info := telemetry.ClientConnectionInfo{}

	// Act
	otelTelemetry.RecordClientConnection(t.Context(), info)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestOTELTelemetry_RecordToolCallRequest_HappyPath(t *testing.T) {
	cases := []struct {
		name         string
		source       telemetry.ToolSource
		inputName    string
		expectedName string
	}{
		{
			name:         "builtin tool name is emitted verbatim",
			source:       telemetry.ToolSourceBuiltin,
			inputName:    "evalMATLABCode",
			expectedName: "evalMATLABCode",
		},
		{
			name:         "extension tool name is emitted as truncated SHA-256 hex",
			source:       telemetry.ToolSourceExtension,
			inputName:    "analyzeData",
			expectedName: "1dcbd0dbb3165d43",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockInstrumentFactory := &telemetrymocks.MockInstrumentFactory{}
			defer mockInstrumentFactory.AssertExpectations(t)

			mockServerStartCounter := &instrumentsmocks.MockInt64Counter{}
			defer mockServerStartCounter.AssertExpectations(t)

			mockClientConnectionCounter := &instrumentsmocks.MockInt64Counter{}
			defer mockClientConnectionCounter.AssertExpectations(t)

			mockToolCallRequestCounter := &instrumentsmocks.MockInt64Counter{}
			defer mockToolCallRequestCounter.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockOSLayer := &telemetrymocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockServerDefinition := &telemetrymocks.MockDefinition{}
			defer mockServerDefinition.AssertExpectations(t)

			mockOSVersionProvider := &telemetrymocks.MockOSVersionProvider{}
			defer mockOSVersionProvider.AssertExpectations(t)

			mockDirectory := &telemetrymocks.MockDirectory{}
			defer mockDirectory.AssertExpectations(t)

			testLogger := testutils.NewInspectableLogger()
			meter := noop.NewMeterProvider().Meter("test")

			expectedInstanceID := "test-instance-id"

			expectedAttributes := []attribute.KeyValue{
				attribute.String("server.instance_id", expectedInstanceID),
				attribute.String("tool.name", tc.expectedName),
			}

			mockInstrumentFactory.EXPECT().
				NewInt64Counter(meter, "server.starts", "Number of times the server has started", "{start}").
				Return(mockServerStartCounter, nil).
				Once()

			mockInstrumentFactory.EXPECT().
				NewInt64Counter(meter, "server.client_connections", "Number of times a client connected to a server", "{connection}").
				Return(mockClientConnectionCounter, nil).
				Once()

			mockInstrumentFactory.EXPECT().
				NewInt64Counter(meter, "tool.calls_request", "Number of tool invocations", "{call}").
				Return(mockToolCallRequestCounter, nil).
				Once()

			mockDirectory.EXPECT().
				ID().
				Return(expectedInstanceID).
				Once()

			mockToolCallRequestCounter.EXPECT().
				Add(mock.Anything, int64(1), expectedAttributes).
				Once()

			otelTelemetry, err := telemetry.NewOTELTelemetryForTesting(testLogger, meter, mockInstrumentFactory, mockConfig, mockDirectory, mockOSLayer, mockOSVersionProvider, mockServerDefinition)
			require.NoError(t, err)

			// Act
			otelTelemetry.RecordToolCallRequest(t.Context(), tc.inputName, tc.source)

			// Assert
			// Assertions are verified via deferred mock expectations.
		})
	}
}

var hexPattern = regexp.MustCompile(`^[0-9a-f]{16}$`)

func TestSHA256Prefix16_ReferenceHashMustNotChange(t *testing.T) {
	got := telemetry.SHA256Prefix16ForTesting("analyzeData")

	assert.Equal(t, "1dcbd0dbb3165d43", got, "reference vector must not drift")
}

func TestSHA256Prefix16_ReturnsSixteenLowercaseHexChars(t *testing.T) {
	inputs := []string{
		"",
		"a",
		"eval_matlab_code",
		"Some Tool With Spaces",
		"emoji_\U0001F389",
		"very-long-name-with-lots-of-content-that-still-truncates-to-sixteen-hex",
	}

	for _, in := range inputs {
		got := telemetry.SHA256Prefix16ForTesting(in)
		assert.Regexp(t, hexPattern, got, "hash of %q must be 16 lowercase hex chars", in)
	}
}

func TestSHA256Prefix16_IsDeterministic(t *testing.T) {
	first := telemetry.SHA256Prefix16ForTesting("consistent")
	second := telemetry.SHA256Prefix16ForTesting("consistent")

	assert.Equal(t, first, second, "same input must produce same output")
}

func TestSHA256Prefix16_IsCaseSensitive(t *testing.T) {
	lower := telemetry.SHA256Prefix16ForTesting("analyzeData")
	upper := telemetry.SHA256Prefix16ForTesting("AnalyzeData")

	assert.NotEqual(t, lower, upper, "case must be significant")
}

func TestSHA256Prefix16_TreatsUTF8AsIs(t *testing.T) {
	nfc := telemetry.SHA256Prefix16ForTesting("café")
	nfd := telemetry.SHA256Prefix16ForTesting("café")

	assert.NotEqual(t, nfc, nfd, "different byte sequences must produce different hashes")
}
