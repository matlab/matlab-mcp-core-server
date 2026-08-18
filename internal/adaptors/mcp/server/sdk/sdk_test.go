// Copyright 2025-2026 The MathWorks, Inc.

package sdk_test

import (
	"context"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/server/sdk"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/mcp/server/sdk"
	telemetrymocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/telemetry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockTelemetryFactory := &mocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	// Act
	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockClientInfoStore, mockLoggerFactory, mockGlobalMATLAB, mockTelemetryFactory)

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

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockTelemetryFactory := &mocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedVersion := "1.0.0"
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "test instructions"

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(mockTelemetry, nil).
		Once()

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	mockDefinition.EXPECT().
		Features().
		Return(definition.Features{}).
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

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockClientInfoStore, mockLoggerFactory, mockGlobalMATLAB, mockTelemetryFactory)

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

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockTelemetryFactory := &mocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockClientInfoStore, mockLoggerFactory, mockGlobalMATLAB, mockTelemetryFactory)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, server, "Server should be nil when error occurs")
}

func TestFactory_NewServer_LoggerError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockTelemetryFactory := &mocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockClientInfoStore, mockLoggerFactory, mockGlobalMATLAB, mockTelemetryFactory)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, server, "Server should be nil when error occurs")
}

func TestFactory_NewServer_TelemetryError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDefinition := &mocks.MockDefinition{}
	defer mockDefinition.AssertExpectations(t)

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockTelemetryFactory := &mocks.MockTelemetryFactory{}
	defer mockTelemetryFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := messages.AnError

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger, nil).
		Once()

	mockTelemetryFactory.EXPECT().
		Telemetry().
		Return(nil, expectedError).
		Once()

	factory := sdk.NewFactory(mockConfigFactory, mockDefinition, mockRootStore, mockClientInfoStore, mockLoggerFactory, mockGlobalMATLAB, mockTelemetryFactory)

	// Act
	server, err := factory.NewServer()

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, server, "Server should be nil when error occurs")
}

func TestHandleInitialized_EagerMATLABInit_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}
	called := make(chan struct{})

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(mock.AnythingOfType("context.withoutCancelCtx"), mockLogger.AsMockArg()).
		Run(func(_ context.Context, _ entities.Logger) {
			close(called)
		}).
		Return(nil, nil).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	select {
	case <-called:
		// Expected: Client was called
	case <-time.After(time.Second):
		t.Fatal("MATLAB eager initialization was not called")
	}
	assert.Empty(t, mockLogger.WarnLogs(), "no warnings should be logged on successful eager initialization")
}

func TestHandleInitialized_NilRequest(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, nil)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestHandleInitialized_NilSession(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, &mcp.InitializedRequest{})

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLogClientDetails_HappyPath(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedClientName := "test-client"
	expectedClientTitle := "Test Client"
	expectedClientURL := "https://example.com"
	expectedClientVersion := "1.2.3"

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			ClientInfo: &mcp.Implementation{
				Name:       expectedClientName,
				Title:      expectedClientTitle,
				WebsiteURL: expectedClientURL,
				Version:    expectedClientVersion,
			},
		}).
		Once()

	logClientDetails := sdk.LogClientDetails(mockLogger)

	// Act
	logClientDetails(mockSession)

	// Assert
	logs := mockLogger.InfoLogs()
	fields, found := logs["New client session"]
	require.True(t, found, "Expected info log for new client session")
	assert.Equal(t, expectedClientName, fields["client-name"])
	assert.Equal(t, expectedClientTitle, fields["client-title"])
	assert.Equal(t, expectedClientURL, fields["client-url"])
	assert.Equal(t, expectedClientVersion, fields["client-version"])
}

func TestLogClientDetails_NilInitializeParams(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockSession.EXPECT().
		InitializeParams().
		Return(nil).
		Once()

	logClientDetails := sdk.LogClientDetails(mockLogger)

	// Act
	logClientDetails(mockSession)

	// Assert
	assert.Empty(t, mockLogger.InfoLogs(), "No info logs should be emitted when InitializeParams is nil")
}

func TestLogClientDetails_NilClientInfo(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{}).
		Once()

	logClientDetails := sdk.LogClientDetails(mockLogger)

	// Act
	logClientDetails(mockSession)

	// Assert
	assert.Empty(t, mockLogger.InfoLogs(), "No info logs should be emitted when ClientInfo is nil")
}

func TestHandleInitialized_EagerMATLABInit_MATLABFeatureDisabled(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: false}}

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	mockGlobalMATLAB.AssertNotCalled(t, "Client")
}

func TestHandleInitialized_EagerMATLABInit_MultipleSession(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(false).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	mockGlobalMATLAB.AssertNotCalled(t, "Client")
}

func TestHandleInitialized_EagerMATLABInit_InitializeMATLABOnStartupFalse(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(false).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	mockGlobalMATLAB.AssertNotCalled(t, "Client")
}

func TestHandleInitialized_EagerMATLABInit_ErrorLogsWarning(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	features := definition.Features{MATLAB: definition.MATLABFeature{Enabled: true}}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockConfig.EXPECT().
		InitializeMATLABOnStartup().
		Return(true).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(mock.AnythingOfType("context.withoutCancelCtx"), mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	handler := sdk.HandleInitialized(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	handler(ctx, &mcp.InitializedRequest{Session: &mcp.ServerSession{}})

	// Assert
	require.Eventually(t, func() bool {
		_, found := mockLogger.WarnLogs()["MATLAB eager initialization failed"]
		return found
	}, time.Second, time.Millisecond)
	fields := mockLogger.WarnLogs()["MATLAB eager initialization failed"]
	assert.Equal(t, expectedError, fields["error"])
}

func TestUpdateRoots_NilInitializeParams_ReturnsNil(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(nil).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_NilCapabilities_ReturnsNil(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{}).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_NilRootsV2_ReturnsNil(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{},
		}).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_ListRootsError_ReturnsError(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	features := definition.Features{}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{
				RootsV2: &mcp.RootCapabilities{},
			},
		}).
		Once()

	mockSession.EXPECT().
		ListRoots(ctx, (*mcp.ListRootsParams)(nil)).
		Return(nil, expectedError).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.ErrorIs(t, err, expectedError)
	mockRootStore.AssertNotCalled(t, "UpdateRoots")
}

func TestUpdateRoots_HappyPath_UpdatesRootStore(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockRootStore := &mocks.MockRootStore{}
	defer mockRootStore.AssertExpectations(t)

	mockGlobalMATLAB := &mocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	features := definition.Features{}
	expectedRoots := []*mcp.Root{
		{URI: "file:///home/user/project", Name: "project"},
	}

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{
				RootsV2: &mcp.RootCapabilities{},
			},
		}).
		Once()

	mockSession.EXPECT().
		ListRoots(ctx, (*mcp.ListRootsParams)(nil)).
		Return(&mcp.ListRootsResult{Roots: expectedRoots}, nil).
		Once()

	mockRootStore.EXPECT().
		UpdateRoots(expectedRoots).
		Once()

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	updateRoots := sdk.UpdateRoots(mockConfig, mockLogger, features, mockRootStore, mockGlobalMATLAB, mockTelemetry)

	// Act
	err := updateRoots(ctx, mockSession)

	// Assert
	require.NoError(t, err)
}

func TestRecordClientConnection_HappyPath(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedClientName := "vscode"
	expectedClientTitle := "Visual Studio Code"
	expectedClientURL := "https://code.visualstudio.com"
	expectedClientVersion := "1.0.0"

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			ClientInfo: &mcp.Implementation{
				Name:       expectedClientName,
				Title:      expectedClientTitle,
				WebsiteURL: expectedClientURL,
				Version:    expectedClientVersion,
			},
			Capabilities: &mcp.ClientCapabilities{
				RootsV2:  &mcp.RootCapabilities{},
				Sampling: &mcp.SamplingCapabilities{},
			},
		}).
		Once()

	expectedInfo := telemetry.ClientConnectionInfo{
		Name:             expectedClientName,
		Title:            expectedClientTitle,
		WebsiteURL:       expectedClientURL,
		Version:          expectedClientVersion,
		Capabilities:     []string{"roots", "sampling"},
		CapabilitiesJSON: `{"roots":{},"sampling":{}}`,
	}

	mockTelemetry.EXPECT().
		RecordClientConnection(ctx, expectedInfo).
		Once()

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockClientInfoStore.EXPECT().
		SetClientInfo(entities.MCPClientInfo{
			Name:       expectedClientName,
			Title:      expectedClientTitle,
			WebsiteURL: expectedClientURL,
			Version:    expectedClientVersion,
		}).
		Once()

	recordClientConnection := sdk.RecordClientConnection(mockLogger, mockTelemetry, mockClientInfoStore)

	// Act
	recordClientConnection(ctx, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestRecordClientConnection_NilInitializeParams(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	mockSession.EXPECT().
		InitializeParams().
		Return(nil).
		Once()

	recordClientConnection := sdk.RecordClientConnection(mockLogger, mockTelemetry, mockClientInfoStore)

	// Act
	recordClientConnection(ctx, mockSession)

	// Assert
	mockTelemetry.AssertNotCalled(t, "RecordClientConnection")
}

func TestRecordClientConnection_NilClientInfo(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{}).
		Once()

	mockTelemetry.EXPECT().
		RecordClientConnection(ctx, telemetry.ClientConnectionInfo{}).
		Once()

	recordClientConnection := sdk.RecordClientConnection(mockLogger, mockTelemetry, mockClientInfoStore)

	// Act
	recordClientConnection(ctx, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestRecordClientConnection_NilCapabilities(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	expectedClientName := "cursor"

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			ClientInfo: &mcp.Implementation{
				Name: expectedClientName,
			},
		}).
		Once()

	mockTelemetry.EXPECT().
		RecordClientConnection(ctx, telemetry.ClientConnectionInfo{
			Name: expectedClientName,
		}).
		Once()

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	mockClientInfoStore.EXPECT().
		SetClientInfo(entities.MCPClientInfo{Name: expectedClientName, Title: ""}).
		Once()

	recordClientConnection := sdk.RecordClientConnection(mockLogger, mockTelemetry, mockClientInfoStore)

	// Act
	recordClientConnection(ctx, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestRecordClientConnection_LegacyRootsOnly_Excluded(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{
				Roots: struct {
					ListChanged bool `json:"listChanged,omitempty"`
				}{ListChanged: true},
			},
		}).
		Once()

	expectedInfo := telemetry.ClientConnectionInfo{
		Capabilities:     nil,
		CapabilitiesJSON: `{}`,
	}

	mockTelemetry.EXPECT().
		RecordClientConnection(ctx, expectedInfo).
		Once()

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	recordClientConnection := sdk.RecordClientConnection(mockLogger, mockTelemetry, mockClientInfoStore)

	// Act
	recordClientConnection(ctx, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestRecordClientConnection_RootsV2OverridesLegacyRoots(t *testing.T) {
	// Arrange
	mockSession := &mocks.MockMCPSession{}
	defer mockSession.AssertExpectations(t)

	mockTelemetry := &telemetrymocks.MockTelemetry{}
	defer mockTelemetry.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	mockSession.EXPECT().
		InitializeParams().
		Return(&mcp.InitializeParams{
			Capabilities: &mcp.ClientCapabilities{
				Roots: struct {
					ListChanged bool `json:"listChanged,omitempty"`
				}{ListChanged: false},
				RootsV2: &mcp.RootCapabilities{ListChanged: true},
			},
		}).
		Once()

	expectedInfo := telemetry.ClientConnectionInfo{
		Capabilities:     []string{"roots"},
		CapabilitiesJSON: `{"roots":{"listChanged":true}}`,
	}

	mockTelemetry.EXPECT().
		RecordClientConnection(ctx, expectedInfo).
		Once()

	mockClientInfoStore := &mocks.MockClientInfoStore{}
	defer mockClientInfoStore.AssertExpectations(t)

	recordClientConnection := sdk.RecordClientConnection(mockLogger, mockTelemetry, mockClientInfoStore)

	// Act
	recordClientConnection(ctx, mockSession)

	// Assert
	// Assertions are verified via deferred mock expectations.
}
