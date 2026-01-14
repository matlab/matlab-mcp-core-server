// Copyright 2025-2026 The MathWorks, Inc.

package localmatlabsession

import (
	"runtime"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type SessionDirectoryFactory interface {
	New(logger entities.Logger) (directory.Directory, error)
}

type ProcessDetails interface {
	NewAPIKey() string
	EnvironmentVariables(sessionDirPath string, apiKey string, certificateFile string, certificateKey string) []string
	StartupFlag(os string, showMATLAB bool, startupCode string) []string
}

type MATLABProcessLauncher interface {
	Launch(logger entities.Logger, sessionRoot string, matlabRoot string, workingDir string, args []string, env []string) (int, func(), error)
}

type Watchdog interface {
	RegisterProcessPIDWithWatchdog(processPID int) error
}

type Starter struct {
	directoryFactory      SessionDirectoryFactory
	processDetails        ProcessDetails
	matlabProcessLauncher MATLABProcessLauncher
	watchdog              Watchdog
}

func NewStarter(
	directoryFactory SessionDirectoryFactory,
	procesDetails ProcessDetails,
	matlabProcessLauncher MATLABProcessLauncher,
	watchdog Watchdog,
) *Starter {
	return &Starter{
		directoryFactory:      directoryFactory,
		processDetails:        procesDetails,
		matlabProcessLauncher: matlabProcessLauncher,
		watchdog:              watchdog,
	}
}

func (m *Starter) StartLocalMATLABSession(logger entities.Logger, request datatypes.LocalSessionDetails) (embeddedconnector.ConnectionDetails, func() error, error) {
	logger.Debug("Starting a local MATLAB session")

	sessionDir, err := m.directoryFactory.New(logger)
	if err != nil {
		return embeddedconnector.ConnectionDetails{}, nil, err
	}

	sessionDirPath := sessionDir.Path()

	logger = logger.With("session_dir", sessionDirPath)
	logger.Debug("Created session directory")

	if !request.IsStartingDirectorySet {
		request.StartingDirectory = sessionDirPath
	}

	uniqueAPIKey := m.processDetails.NewAPIKey()

	env := m.processDetails.EnvironmentVariables(
		sessionDirPath,
		uniqueAPIKey,
		sessionDir.CertificateFile(),
		sessionDir.CertificateKeyFile(),
	)

	startupCode := "sessionPath = '" + sessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"

	startupFlags := m.processDetails.StartupFlag(runtime.GOOS, request.ShowMATLABDesktop, startupCode)

	processID, processCleanup, err := m.matlabProcessLauncher.Launch(logger, sessionDirPath, request.MATLABRoot, request.StartingDirectory, startupFlags, env)
	if err != nil {
		return embeddedconnector.ConnectionDetails{}, nil, err
	}

	if err = m.watchdog.RegisterProcessPIDWithWatchdog(processID); err != nil {
		logger.WithError(err).Warn("Failed to register process with watchdog")
	}

	securePort, certificatePEM, err := sessionDir.GetEmbeddedConnectorDetails()
	if err != nil {
		return embeddedconnector.ConnectionDetails{}, nil, err
	}

	return embeddedconnector.ConnectionDetails{
			Host:           "localhost",
			Port:           securePort,
			APIKey:         uniqueAPIKey,
			CertificatePEM: certificatePEM,
		}, func() error {
			processCleanup()
			return sessionDir.Cleanup()
		}, nil
}
