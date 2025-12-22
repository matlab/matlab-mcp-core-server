// Copyright 2025 The MathWorks, Inc.

package flags

const (
	HelpMode             = "help"
	HelpModeDefaultValue = false

	VersionMode             = "version"
	VersionModeDefaultValue = false

	DisableTelemetry             = "disable-telemetry"
	DisableTelemetryDefaultValue = false

	PreferredLocalMATLABRoot             = "matlab-root"
	PreferredLocalMATLABRootDefaultValue = ""

	PreferredMATLABStartingDirectory             = "initial-working-folder"
	PreferredMATLABStartingDirectoryDefaultValue = ""

	BaseDir             = "log-folder"
	BaseDirDefaultValue = ""

	LogLevel             = "log-level"
	LogLevelDefaultValue = "info"

	InitializeMATLABOnStartup             = "initialize-matlab-on-startup"
	InitializeMATLABOnStartupDefaultValue = false

	// Hidden
	UseSingleMATLABSession             = "use-single-matlab-session"
	UseSingleMATLABSessionDefaultValue = true

	WatchdogMode             = "watchdog"
	WatchdogModeDefaultValue = false

	ServerInstanceID             = "server-instance-id"
	ServerInstanceIDDefaultValue = ""
)
