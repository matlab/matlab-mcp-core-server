// Copyright 2025 The MathWorks, Inc.

package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	parsermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/inputs/parser"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
		name     string
		args     []string
		expected parser.SpecifiedArguments
	}{
		{
			name: "default values",
			args: []string{},
			expected: parser.SpecifiedArguments{
				VersionMode:                      false,
				DisableTelemetry:                 false,
				UseSingleMATLABSession:           true,
				LogLevel:                         entities.LogLevelInfo,
				PreferredLocalMATLABRoot:         "",
				PreferredMATLABStartingDirectory: "",
				BaseDirectory:                    "",
				WatchdogMode:                     false,
				ServerInstanceID:                 "",
				InitializeMATLABOnStartup:        false,
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
			expected: parser.SpecifiedArguments{
				VersionMode:                      true,
				DisableTelemetry:                 true,
				UseSingleMATLABSession:           false,
				LogLevel:                         entities.LogLevelDebug,
				PreferredLocalMATLABRoot:         filepath.Join("tmp", "root"),
				PreferredMATLABStartingDirectory: filepath.Join("tmp", "pref"),
				BaseDirectory:                    filepath.Join("tmp", "logs"),
				WatchdogMode:                     true,
				ServerInstanceID:                 "1337",
				InitializeMATLABOnStartup:        false,
			},
		},
		{
			name: "single session disabled forces initialize on startup false",
			args: []string{
				"--use-single-matlab-session=false",
				"--initialize-matlab-on-startup=true",
			},
			expected: parser.SpecifiedArguments{
				VersionMode:                      false,
				DisableTelemetry:                 false,
				UseSingleMATLABSession:           false,
				LogLevel:                         entities.LogLevelInfo,
				PreferredLocalMATLABRoot:         "",
				PreferredMATLABStartingDirectory: "",
				BaseDirectory:                    "",
				WatchdogMode:                     false,
				InitializeMATLABOnStartup:        false,
			},
		},
	}

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Setup
			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			cliParser := parser.New(mockMessageCatalog)

			specifiedArguments, err := cliParser.Parse(testConfig.args)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, specifiedArguments, "Config should not be nil")

			assert.Equal(t, testConfig.expected, specifiedArguments)
		})
	}
}

func TestParser_DisableTelemetry_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
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

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			parser := parser.New(mockMessageCatalog)
			result, err := parser.Parse(testConfig.args)
			require.NoError(t, err)
			// Assert
			assert.Equal(t, testConfig.expected, result.DisableTelemetry)
		})
	}
}

func TestParser_UseSingleMATLABSession_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
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

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			parser := parser.New(mockMessageCatalog)
			result, err := parser.Parse(testConfig.args)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, testConfig.expected, result.UseSingleMATLABSession)
		})
	}
}

func TestParser_PreferredLocalMATLABRoot_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
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

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			programName := "testprocess"
			args := append([]string{programName}, testConfig.args...)

			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			parser := parser.New(mockMessageCatalog)
			result, err := parser.Parse(args)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, testConfig.expected, result.PreferredLocalMATLABRoot)
		})
	}
}

func TestParser_PreferredMATLABStartingDirectory_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
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

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			parser := parser.New(mockMessageCatalog)
			result, err := parser.Parse(testConfig.args)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, testConfig.expected, result.PreferredMATLABStartingDirectory)
		})
	}
}

func TestParser_LogDirectory_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
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

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			parser := parser.New(mockMessageCatalog)
			result, err := parser.Parse(testConfig.args)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, testConfig.expected, result.BaseDirectory)
		})
	}
}

func TestParser_LogLevel_HappyPath(t *testing.T) {
	testSpecifiedArgs := []struct {
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

	for _, testConfig := range testSpecifiedArgs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockMessageCatalog := &parsermocks.MockMessageCatalog{}
			defer mockMessageCatalog.AssertExpectations(t)

			mockMessageCatalog.EXPECT().
				Get(mock.Anything).
				Return("any string")

			// Act
			parser := parser.New(mockMessageCatalog)
			result, err := parser.Parse(testConfig.args)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, testConfig.expected, result.LogLevel)
		})
	}
}

func TestParser_LogLevel_Invalid(t *testing.T) {
	// Arrange
	badLogLevel := "invalid"
	args := []string{"--log-level=" + badLogLevel}

	mockMessageCatalog := &parsermocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("any string")

	// Act
	parser := parser.New(mockMessageCatalog)
	result, err := parser.Parse(args)

	// Assert
	require.Equal(t, err, messages.New_StartupErrors_InvalidLogLevel_Error(badLogLevel))
	assert.Empty(t, result)
}

func TestParser_LogLevel_EmptyIsInvalid(t *testing.T) {
	// Arrange
	args := []string{"--log-level="}

	mockMessageCatalog := &parsermocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("any string")

	// Act
	parser := parser.New(mockMessageCatalog)
	result, err := parser.Parse(args)

	// Assert
	require.Equal(t, err, messages.New_StartupErrors_InvalidLogLevel_Error(""))
	assert.Empty(t, result)
}

func TestParser_BadFlagResultsInError(t *testing.T) {
	// Arrange
	args := []string{"--notaflag=true"}

	mockMessageCatalog := &parsermocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("any string")

	// Act
	parser := parser.New(mockMessageCatalog)
	result, err := parser.Parse(args)
	dummyUsage := parser.Usage()

	// Assert
	require.Equal(t, messages.New_StartupErrors_BadFlag_Error("notaflag", "\n", dummyUsage), err)
	assert.Empty(t, result)
}

func TestParser_BadValuesResultsInError(t *testing.T) {
	// Arrange
	BadArg := "trur"
	args := []string{"--" + flags.DisableTelemetry + "=" + BadArg}

	mockMessageCatalog := &parsermocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("any string")

	// Act
	parser := parser.New(mockMessageCatalog)
	result, err := parser.Parse(args)

	// Assert
	require.Equal(t, messages.New_StartupErrors_BadValue_Error(BadArg, flags.DisableTelemetry), err)
	assert.Empty(t, result)
}

func TestParser_BadSyntaxResultsInError(t *testing.T) {
	// Arrange
	args := []string{"---" + flags.DisableTelemetry + "="}

	mockMessageCatalog := &parsermocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("any string")

	// Act
	parser := parser.New(mockMessageCatalog)
	result, err := parser.Parse(args)
	dummyUsage := parser.Usage()

	// Assert
	require.Equal(t, messages.New_StartupErrors_BadSyntax_Error(args[0], "\n", dummyUsage), err)
	assert.Empty(t, result)
}

func TestParser_DefaultValues(t *testing.T) {
	// Arrange
	args := []string{}

	mockMessageCatalog := &parsermocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockMessageCatalog.EXPECT().
		Get(mock.Anything).
		Return("any string")

	expectedResult := parser.SpecifiedArguments{
		VersionMode:                      flags.VersionModeDefaultValue,
		HelpMode:                         flags.HelpModeDefaultValue,
		DisableTelemetry:                 flags.DisableTelemetryDefaultValue,
		UseSingleMATLABSession:           flags.UseSingleMATLABSessionDefaultValue,
		LogLevel:                         flags.LogLevelDefaultValue,
		PreferredLocalMATLABRoot:         flags.PreferredLocalMATLABRootDefaultValue,
		PreferredMATLABStartingDirectory: flags.PreferredMATLABStartingDirectoryDefaultValue,
		BaseDirectory:                    flags.BaseDirDefaultValue,
		WatchdogMode:                     flags.WatchdogModeDefaultValue,
		ServerInstanceID:                 flags.ServerInstanceIDDefaultValue,
		InitializeMATLABOnStartup:        flags.InitializeMATLABOnStartupDefaultValue,
	}

	// Act
	parser := parser.New(mockMessageCatalog)
	result, err := parser.Parse(args)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestParser_GenerateUsage(t *testing.T) {
	// Arrange
	flagSet := pflag.NewFlagSet(pflag.CommandLine.Name(), pflag.ContinueOnError)
	flagSet.Bool("short", true, "This flag name is very short")
	flagSet.Bool("really-long-flag", true, "This flag name is extremely verbose")
	flagSet.Bool("hidden-flag", true, "This flag is hidden")
	require.NoError(t, flagSet.MarkHidden("hidden-flag"))

	expectedResult := `Usage:
      --really-long-flag    This flag name is extremely verbose
      --short               This flag name is very short
`
	// Act
	result := parser.GenerateUsageText(flagSet)

	// Assert
	assert.Equal(t, expectedResult, result)
	assert.NotContains(t, result, "hidden-flag")
}
