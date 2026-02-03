// Copyright 2025-2026 The MathWorks, Inc.

package config_test

import (
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	// Act
	factory := config.NewFactory(mockParser, mockOSLayer)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_Config_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parser.SpecifiedArguments{}, nil).
		Once()

	factory := config.NewFactory(mockParser, mockOSLayer)

	// Act
	cfg, err := factory.Config()

	// Assert
	require.NoError(t, err, "Config should not return an error")
	assert.NotNil(t, cfg, "Config should not be nil")
}

func TestFactory_Config_ParseError(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	expectedError := &messages.StartupErrors_BadFlag_Error{}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parser.SpecifiedArguments{}, expectedError).
		Once()

	factory := config.NewFactory(mockParser, mockOSLayer)

	// Act
	cfg, err := factory.Config()

	// Assert
	require.ErrorIs(t, err, expectedError, "Config should return the expected error")
	assert.Nil(t, cfg, "Config should be nil")
}

func TestFactory_Config_IsSingleton(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parser.SpecifiedArguments{}, nil).
		Once()

	factory := config.NewFactory(mockParser, mockOSLayer)

	// Act
	cfg1, err1 := factory.Config()
	cfg2, err2 := factory.Config()

	// Assert
	require.NoError(t, err1, "First Config call should not return an error")
	require.NoError(t, err2, "Second Config call should not return an error")
	assert.Same(t, cfg1, cfg2, "Config should return the same instance on multiple calls")
}

func TestFactory_Config_SingletonPreservesError(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	expectedError := &messages.StartupErrors_BadFlag_Error{}

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockParser.EXPECT().
		Parse(args[1:]).
		Return(parser.SpecifiedArguments{}, expectedError).
		Once()

	factory := config.NewFactory(mockParser, mockOSLayer)

	// Act
	cfg1, err1 := factory.Config()
	cfg2, err2 := factory.Config()

	// Assert
	require.Error(t, err1, "First Config call should return an error")
	require.Error(t, err2, "Second Config call should return the same error")
	assert.Equal(t, expectedError, err1, "First error should match expected")
	assert.Equal(t, expectedError, err2, "Second error should match expected")
	assert.Nil(t, cfg1, "First config should be nil")
	assert.Nil(t, cfg2, "Second config should be nil")
}

func TestConfig_Version(t *testing.T) {
	modulePath := "github.com/matlab/matlab-mcp-core-server"

	testCases := []struct {
		name            string
		buildInfoOK     bool
		moduleVersion   string
		expectedVersion string
	}{
		{
			name:            "version from build info",
			buildInfoOK:     true,
			moduleVersion:   "v1.2.3",
			expectedVersion: modulePath + " v1.2.3",
		},
		{
			name:            "devel fallback",
			buildInfoOK:     true,
			moduleVersion:   "(devel)",
			expectedVersion: modulePath + " (devel)",
		},
		{
			name:            "build info unavailable",
			buildInfoOK:     false,
			moduleVersion:   "",
			expectedVersion: "(unknown)",
		},
		{
			name:            "empty version string",
			buildInfoOK:     true,
			moduleVersion:   "",
			expectedVersion: modulePath + " (devel)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}

			mockOSLayer.EXPECT().
				Args().
				Return([]string{"testprocess"}).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return(parser.SpecifiedArguments{}, nil).
				Once()

			var buildInfo *debug.BuildInfo
			if tc.buildInfoOK {
				buildInfo = &debug.BuildInfo{
					Main: debug.Module{
						Path:    modulePath,
						Version: tc.moduleVersion,
					},
				}
			}

			mockOSLayer.EXPECT().
				ReadBuildInfo().
				Return(buildInfo, tc.buildInfoOK).
				Once()

			// Act
			cfg, err := config.New(mockOSLayer, mockParser)
			require.NoError(t, err)

			version := cfg.Version()

			// Assert
			require.Equal(t, tc.expectedVersion, version)
		})
	}
}

func TestConfig_ShouldShowMATLABDesktop(t *testing.T) {
	testCases := []struct {
		name         string
		displayMode  entities.DisplayMode
		expectedShow bool
	}{
		{
			name:         "desktop mode shows MATLAB",
			displayMode:  entities.DisplayModeDesktop,
			expectedShow: true,
		},
		{
			name:         "nodesktop mode hides MATLAB",
			displayMode:  entities.DisplayModeNoDesktop,
			expectedShow: false,
		},
		{
			name:         "unset mode shows MATLAB (default)",
			displayMode:  "",
			expectedShow: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockParser := &configmocks.MockParser{}
			defer mockParser.AssertExpectations(t)

			programName := "testprocess"
			args := []string{programName}
			specifiedArguments := parser.SpecifiedArguments{
				DisplayMode: tc.displayMode,
			}

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			mockParser.EXPECT().
				Parse(args[1:]).
				Return(specifiedArguments, nil).
				Once()

			cfg, err := config.New(mockOSLayer, mockParser)
			require.NoError(t, err)

			// Act
			result := cfg.ShouldShowMATLABDesktop()

			// Assert
			assert.Equal(t, tc.expectedShow, result)
		})
	}
}

func TestConfig_Log_HappyPath(t *testing.T) {
	// Arrange
	specifiedArguments := parser.SpecifiedArguments{
		DisableTelemetry:                 true,
		PreferredMATLABStartingDirectory: filepath.Join("home", "user"),
		LogLevel:                         entities.LogLevelDebug,
		PreferredLocalMATLABRoot:         filepath.Join("home", "matlab"),
		UseSingleMATLABSession:           false,
		DisplayMode:                      entities.DisplayModeNoDesktop,
	}
	expectedLogMessage := "Configuration state"
	expectedConfigField := map[string]any{
		"disable-telemetry":         true,
		"initial-working-folder":    filepath.Join("home", "user"),
		"log-level":                 entities.LogLevelDebug,
		"matlab-root":               filepath.Join("home", "matlab"),
		"use-single-matlab-session": false,
		"matlab-display-mode":       entities.DisplayModeNoDesktop,
	}
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockParser := &configmocks.MockParser{}
	defer mockParser.AssertExpectations(t)

	programName := "testprocess"
	args := []string{programName}
	mockParser.EXPECT().
		Parse(args[1:]).
		Return(specifiedArguments, nil)

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	cfg, err := config.New(mockOSLayer, mockParser)
	require.NoError(t, err)

	testLogger := testutils.NewInspectableLogger()

	// Act
	cfg.RecordToLogger(testLogger)

	// Assert
	infoLogs := testLogger.InfoLogs()
	require.Len(t, infoLogs, 1)

	fields, found := infoLogs[expectedLogMessage]
	require.True(t, found, "Expected log message not found")

	for expectedField, expectedValue := range expectedConfigField {
		actualValue, exists := fields[expectedField]
		require.True(t, exists, "%s field not found in log", expectedField)
		assert.Equal(t, expectedValue, actualValue, "%s field has incorrect value", expectedField)
	}
}
