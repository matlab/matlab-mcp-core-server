// Copyright 2025 The MathWorks, Inc.

package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/spf13/pflag"
)

type MessageCatalog interface {
	Get(message messages.MessageKey) string
}

type SpecifiedArguments struct {
	VersionMode                      bool
	HelpMode                         bool
	DisableTelemetry                 bool
	UseSingleMATLABSession           bool
	LogLevel                         entities.LogLevel
	PreferredLocalMATLABRoot         string
	PreferredMATLABStartingDirectory string
	BaseDirectory                    string
	WatchdogMode                     bool
	ServerInstanceID                 string
	InitializeMATLABOnStartup        bool
}

type Parser struct {
	usageText string
	flagSet   *pflag.FlagSet
}

func New(messageCatalog MessageCatalog) *Parser {
	flagSet := pflag.NewFlagSet(pflag.CommandLine.Name(), pflag.ContinueOnError)
	setupFlags(messageCatalog, flagSet)
	usageText := generateUsageText(flagSet)
	return &Parser{flagSet: flagSet, usageText: usageText}
}

func (p *Parser) Usage() string {
	return p.usageText
}

func (p *Parser) Parse(args []string) (SpecifiedArguments, messages.Error) {
	err := p.flagSet.Parse(args)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	helpMode, err := p.flagSet.GetBool(flags.HelpMode)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	versionMode, err := p.flagSet.GetBool(flags.VersionMode)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	disableTelemetry, err := p.flagSet.GetBool(flags.DisableTelemetry)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	useSingleMATLABSession, err := p.flagSet.GetBool(flags.UseSingleMATLABSession)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	logLevel, err := p.flagSet.GetString(flags.LogLevel)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	switch logLevel {
	case string(entities.LogLevelDebug), string(entities.LogLevelInfo), string(entities.LogLevelWarn), string(entities.LogLevelError):
		break
	default:
		return SpecifiedArguments{}, messages.New_StartupErrors_InvalidLogLevel_Error(logLevel)
	}

	preferredLocalMATLABRoot, err := p.flagSet.GetString(flags.PreferredLocalMATLABRoot)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	preferredMATLABStartingDirectory, err := p.flagSet.GetString(flags.PreferredMATLABStartingDirectory)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	baseDir, err := p.flagSet.GetString(flags.BaseDir)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	watchdogMode, err := p.flagSet.GetBool(flags.WatchdogMode)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	serverInstanceID, err := p.flagSet.GetString(flags.ServerInstanceID)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	initializeMATLABOnStartup, err := p.flagSet.GetBool(flags.InitializeMATLABOnStartup)
	if err != nil {
		return SpecifiedArguments{}, p.convertToUserFacingError(err)
	}

	if !useSingleMATLABSession {
		initializeMATLABOnStartup = false
	}

	return SpecifiedArguments{
		VersionMode:                      versionMode,
		HelpMode:                         helpMode,
		DisableTelemetry:                 disableTelemetry,
		UseSingleMATLABSession:           useSingleMATLABSession,
		LogLevel:                         entities.LogLevel(logLevel),
		PreferredLocalMATLABRoot:         preferredLocalMATLABRoot,
		PreferredMATLABStartingDirectory: preferredMATLABStartingDirectory,
		BaseDirectory:                    baseDir,
		WatchdogMode:                     watchdogMode,
		ServerInstanceID:                 serverInstanceID,
		InitializeMATLABOnStartup:        initializeMATLABOnStartup,
	}, nil
}

func setHiddenBoolFlag(flagSet *pflag.FlagSet, flagName string, flagDefaultValue bool, flagDescription string) {
	flagSet.Bool(flagName, flagDefaultValue, flagDescription)
	//nolint:errcheck,gosec // Logically impossible to hit NotExistError
	flagSet.MarkHidden(flagName)
}

func setHiddenStringFlag(flagSet *pflag.FlagSet, flagName string, flagDefaultValue string, flagDescription string) {
	flagSet.String(flagName, flagDefaultValue, flagDescription)
	//nolint:errcheck,gosec // Logically impossible to hit NotExistError
	flagSet.MarkHidden(flagName)
}

func setupFlags(messageCatalog MessageCatalog, flagSet *pflag.FlagSet) {
	flagSet.Bool(flags.HelpMode, flags.HelpModeDefaultValue,
		messageCatalog.Get(messages.CLIMessages_HelpDescription),
	)

	flagSet.Bool(flags.VersionMode, flags.VersionModeDefaultValue,
		messageCatalog.Get(messages.CLIMessages_VersionDescription),
	)

	flagSet.Bool(flags.DisableTelemetry, flags.DisableTelemetryDefaultValue,
		messageCatalog.Get(messages.CLIMessages_DisableTelemetryDescription),
	)

	flagSet.String(flags.LogLevel, flags.LogLevelDefaultValue,
		messageCatalog.Get(messages.CLIMessages_LogLevelDescription),
	)

	flagSet.String(flags.PreferredLocalMATLABRoot, flags.PreferredLocalMATLABRootDefaultValue,
		messageCatalog.Get(messages.CLIMessages_PreferredLocalMATLABRootDescription),
	)

	flagSet.String(flags.PreferredMATLABStartingDirectory, flags.PreferredMATLABStartingDirectoryDefaultValue,
		messageCatalog.Get(messages.CLIMessages_PreferredMATLABStartingDirectoryDescription),
	)

	flagSet.String(flags.BaseDir, flags.BaseDirDefaultValue,
		messageCatalog.Get(messages.CLIMessages_BaseDirDescription),
	)

	flagSet.Bool(flags.InitializeMATLABOnStartup, flags.InitializeMATLABOnStartupDefaultValue,
		messageCatalog.Get(messages.CLIMessages_InitializeMATLABOnStartupDescription),
	)

	// Hidden flags, for internal use only
	setHiddenBoolFlag(flagSet, flags.UseSingleMATLABSession, flags.UseSingleMATLABSessionDefaultValue,
		messageCatalog.Get(messages.CLIMessages_UseSingleMATLABSessionDescription))

	setHiddenBoolFlag(flagSet, flags.WatchdogMode, flags.WatchdogModeDefaultValue,
		messageCatalog.Get(messages.CLIMessages_InternalUseDescription))

	setHiddenStringFlag(flagSet, flags.ServerInstanceID, flags.ServerInstanceIDDefaultValue,
		messageCatalog.Get(messages.CLIMessages_InternalUseDescription))
}

func (p *Parser) convertToUserFacingError(err error) messages.Error {
	var notExistError *pflag.NotExistError
	var invalidSyntaxError *pflag.InvalidSyntaxError
	var invalidValueError *pflag.InvalidValueError
	var valueRequiredError *pflag.ValueRequiredError

	switch {
	case errors.As(err, &notExistError):
		return messages.New_StartupErrors_BadFlag_Error(notExistError.GetSpecifiedName(), "\n", p.usageText)
	case errors.As(err, &invalidSyntaxError):
		return messages.New_StartupErrors_BadSyntax_Error(invalidSyntaxError.GetSpecifiedFlag(), "\n", p.usageText)
	case errors.As(err, &invalidValueError):
		return messages.New_StartupErrors_BadValue_Error(invalidValueError.GetValue(), invalidValueError.GetFlag().Name)
	case errors.As(err, &valueRequiredError):
		return messages.New_StartupErrors_MissingValue_Error(valueRequiredError.GetSpecifiedName())
	}
	return messages.New_StartupErrors_ParseFailed_Error("\n", p.usageText)
}

func generateUsageText(flagSet *pflag.FlagSet) string {
	usageText := fmt.Sprintf("%s\n", "Usage:")

	// Determine max flag length
	maxFlagLength := 0
	prePadding := 6
	postPadding := 2

	flagSet.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		if len(f.Name) > maxFlagLength {
			maxFlagLength = len(f.Name)
		}
	})

	flagSet.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		padding := maxFlagLength + postPadding + 2 - len(f.Name)
		usageText += fmt.Sprintf("%s--%s%s%s\n", strings.Repeat(" ", prePadding), f.Name, strings.Repeat(" ", padding), f.Usage)
	})
	return usageText
}
