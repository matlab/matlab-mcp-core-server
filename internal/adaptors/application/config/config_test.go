// Copyright 2025 The MathWorks, Inc.

package config_test

import (
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedConfig struct {
	versionMode                      bool
	disableTelemetry                 bool
	useSingleMATLABSession           bool
	logLevel                         entities.LogLevel
	preferredLocalMATLABRoot         string
	preferredMATLABStartingDirectory string
	baseDirectory                    string
	watchdogMode                     bool
	serverInstanceID                 string
	initializeMATLABOnStartup        bool
}

func TestNew_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected expectedConfig
	}{
		{
			name: "default values",
			args: []string{},
			expected: expectedConfig{
				versionMode:                      false,
				disableTelemetry:                 false,
				useSingleMATLABSession:           true,
				logLevel:                         entities.LogLevelInfo,
				preferredLocalMATLABRoot:         "",
				preferredMATLABStartingDirectory: "",
				baseDirectory:                    "",
				watchdogMode:                     false,
				serverInstanceID:                 "",
				initializeMATLABOnStartup:        false,
			},
		},
		{
			name: "custom values",
			args: []string{
				"--version=true",
				"--disable-telemetry",
				"--use-single-matlab-session=false",
				"--log-level=debug",
				"--matlab-root=" + filepath.Join("tmp", "root"),
				"--initial-working-folder=" + filepath.Join("tmp", "pref"),
				"--log-folder=" + filepath.Join("tmp", "logs"),
				"--watchdog=true",
				"--server-instance-id=1337",
				"--initialize-matlab-on-startup=false",
			},
			expected: expectedConfig{
				versionMode:                      true,
				disableTelemetry:                 true,
				useSingleMATLABSession:           false,
				logLevel:                         entities.LogLevelDebug,
				preferredLocalMATLABRoot:         filepath.Join("tmp", "root"),
				preferredMATLABStartingDirectory: filepath.Join("tmp", "pref"),
				baseDirectory:                    filepath.Join("tmp", "logs"),
				watchdogMode:                     true,
				serverInstanceID:                 "1337",
				initializeMATLABOnStartup:        false,
			},
		},
		{
			name: "single session disabled forces initialize on startup false",
			args: []string{
				"--use-single-matlab-session=false",
				"--initialize-matlab-on-startup=true",
			},
			expected: expectedConfig{
				versionMode:                      false,
				disableTelemetry:                 false,
				useSingleMATLABSession:           false,
				logLevel:                         entities.LogLevelInfo,
				preferredLocalMATLABRoot:         "",
				preferredMATLABStartingDirectory: "",
				baseDirectory:                    "",
				watchdogMode:                     false,
				initializeMATLABOnStartup:        false,
			},
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			// Act
			cfg, err := config.New(mockOSLayer)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, cfg, "Config should not be nil")

			assert.Equal(t, testConfig.expected.versionMode, cfg.VersionMode())
			assert.Equal(t, testConfig.expected.disableTelemetry, cfg.DisableTelemetry())
			assert.Equal(t, testConfig.expected.useSingleMATLABSession, cfg.UseSingleMATLABSession())
			assert.Equal(t, testConfig.expected.logLevel, cfg.LogLevel())
			assert.Equal(t, testConfig.expected.preferredLocalMATLABRoot, cfg.PreferredLocalMATLABRoot())
			assert.Equal(t, testConfig.expected.preferredMATLABStartingDirectory, cfg.PreferredMATLABStartingDirectory())
			assert.Equal(t, testConfig.expected.baseDirectory, cfg.BaseDir())
			assert.Equal(t, testConfig.expected.watchdogMode, cfg.WatchdogMode())
			assert.Equal(t, testConfig.expected.serverInstanceID, cfg.ServerInstanceID())
			assert.Equal(t, testConfig.expected.initializeMATLABOnStartup, cfg.InitializeMATLABOnStartup())
		})
	}
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
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return([]string{"testprocess"}).
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

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			version := cfg.Version()

			require.Equal(t, tc.expectedVersion, version)
		})
	}
}

func TestConfig_DisableTelemetry_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: false,
		},
		{
			name:     "implicitly true",
			args:     []string{"--disable-telemetry"},
			expected: true,
		},
		{
			name:     "explicitly true",
			args:     []string{"--disable-telemetry=true"},
			expected: true,
		},
		{
			name:     "explicitly false",
			args:     []string{"--disable-telemetry=false"},
			expected: false,
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.DisableTelemetry()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_UseSingleMATLABSession_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: true,
		},
		{
			name:     "explicitly true",
			args:     []string{"--use-single-matlab-session=true"},
			expected: true,
		},
		{
			name:     "explicitly false",
			args:     []string{"--use-single-matlab-session=false"},
			expected: false,
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.UseSingleMATLABSession()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_PreferredLocalMATLABRoot_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: "",
		},
		{
			name:     "custom path",
			args:     []string{"--matlab-root=" + filepath.Join("path", "to", "matlab")},
			expected: filepath.Join("path", "to", "matlab"),
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.PreferredLocalMATLABRoot()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_PreferredMATLABStartingDirectory_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: "",
		},
		{
			name:     "custom project path",
			args:     []string{"--initial-working-folder=" + filepath.Join("path", "to", "project")},
			expected: filepath.Join("path", "to", "project"),
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.PreferredMATLABStartingDirectory()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_LogDirectory_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: "",
		},
		{
			name:     "Supplied log directory",
			args:     []string{"--log-folder=" + filepath.Join("tmp", "logs")},
			expected: filepath.Join("tmp", "logs"),
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.BaseDir()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_LogLevel_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected entities.LogLevel
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: entities.LogLevelInfo,
		},
		{
			name:     "debug level",
			args:     []string{"--log-level=debug"},
			expected: entities.LogLevelDebug,
		},
		{
			name:     "info level",
			args:     []string{"--log-level=info"},
			expected: entities.LogLevelInfo,
		},
		{
			name:     "warn level",
			args:     []string{"--log-level=warn"},
			expected: entities.LogLevelWarn,
		},
		{
			name:     "error level",
			args:     []string{"--log-level=error"},
			expected: entities.LogLevelError,
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.LogLevel()

			// Assert
			require.NoError(t, err)
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_LogLevel_Invalid(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	programName := "testprocess"
	args := append([]string{programName}, "--log-level=invalid")

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	// Act
	cfg, err := config.New(mockOSLayer)

	// Assert
	require.Errorf(t, err, "invalid log level")
	assert.Empty(t, cfg)
}

func TestConfig_LogLevel_EmptyIsInvalid(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	programName := "testprocess"
	args := append([]string{programName}, "--log-level=")

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	// Act
	cfg, err := config.New(mockOSLayer)

	// Assert
	require.Errorf(t, err, "invalid log level")
	assert.Empty(t, cfg)
}

func TestConfig_Log_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name                string
		args                []string
		expectedLogMessage  string
		expectedConfigField map[string]any
	}{
		{
			name:               "default configuration",
			args:               []string{},
			expectedLogMessage: "Configuration state",
			expectedConfigField: map[string]any{
				"disable-telemetry":         false,
				"initial-working-folder":    "",
				"log-level":                 entities.LogLevelInfo,
				"matlab-root":               "",
				"use-single-matlab-session": true,
			},
		},
		{
			name: "custom configuration",
			args: []string{
				"--disable-telemetry",
				"--use-single-matlab-session=false",
				"--log-level=debug",
				"--initial-working-folder=" + filepath.Join("home", "user"),
				"--matlab-root=" + filepath.Join("home", "matlab"),
			},
			expectedLogMessage: "Configuration state",
			expectedConfigField: map[string]any{
				"disable-telemetry":         true,
				"initial-working-folder":    filepath.Join("home", "user"),
				"log-level":                 entities.LogLevelDebug,
				"matlab-root":               filepath.Join("home", "matlab"),
				"use-single-matlab-session": false,
			},
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockOSLayer.EXPECT().
				Args().
				Return(args).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			testLogger := testutils.NewInspectableLogger()

			// Act
			cfg.RecordToLogger(testLogger)

			// Assert
			infoLogs := testLogger.InfoLogs()
			require.Len(t, infoLogs, 1)

			fields, found := infoLogs[testConfig.expectedLogMessage]
			require.True(t, found, "Expected log message not found")

			for expectedField, expectedValue := range testConfig.expectedConfigField {
				actualValue, exists := fields[expectedField]
				require.True(t, exists, "%s field not found in log", expectedField)
				assert.Equal(t, expectedValue, actualValue, "%s field has incorrect value", expectedField)
			}
		})
	}
}
