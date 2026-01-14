// Copyright 2025-2026 The MathWorks, Inc.

package localmatlabsession_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/datatypes"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession"
	directorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStarter_HappyPath(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	// Act
	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	// Assert
	assert.NotNil(t, starter)
}

func TestStarter_StartLocalMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDirPath := filepath.Join("tmp", "matlab-session-12345")
	expectedCertificateFile := filepath.Join("tmp", "matlab-session-12345", "cert.pem")
	expectedCertificateKeyFile := filepath.Join("tmp", "matlab-session-12345", "cert.key")
	expectedAPIKey := "test-api-key-12345"
	expectedMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")
	expectedSecurePort := "9999"
	expectedCertificatePEM := []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----")
	expectedEnv := []string{"MATLAB_MCP_API_KEY=" + expectedAPIKey}
	expectedStartupCode := "sessionPath = '" + expectedSessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"
	showDestop := false
	expectedStartupFlags := []string{"-r", expectedStartupCode}
	expectedProcessID := 12345
	processCleanupCalled := false
	processCleanup := func() {
		processCleanupCalled = true
	}

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		Path().
		Return(expectedSessionDirPath).
		Once()

	mockProcessDetails.EXPECT().
		NewAPIKey().
		Return(expectedAPIKey).
		Once()

	mockDirectory.EXPECT().
		CertificateFile().
		Return(expectedCertificateFile).
		Once()

	mockDirectory.EXPECT().
		CertificateKeyFile().
		Return(expectedCertificateKeyFile).
		Once()

	mockProcessDetails.EXPECT().
		EnvironmentVariables(expectedSessionDirPath, expectedAPIKey, expectedCertificateFile, expectedCertificateKeyFile).
		Return(expectedEnv).
		Once()

	mockProcessDetails.EXPECT().
		StartupFlag(runtime.GOOS, showDestop, expectedStartupCode).
		Return(expectedStartupFlags).
		Once()

	mockMATLABProcessLauncher.EXPECT().
		Launch(mockLogger.AsMockArg(), expectedSessionDirPath, expectedMATLABRoot, expectedSessionDirPath, expectedStartupFlags, expectedEnv).
		Return(expectedProcessID, processCleanup, nil).
		Once()

	mockWatchdog.EXPECT().
		RegisterProcessPIDWithWatchdog(expectedProcessID).
		Return(nil).
		Once()

	mockDirectory.EXPECT().
		GetEmbeddedConnectorDetails().
		Return(expectedSecurePort, expectedCertificatePEM, nil).
		Once()

	mockDirectory.EXPECT().
		Cleanup().
		Return(nil).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		IsStartingDirectorySet: false,
		MATLABRoot:             expectedMATLABRoot,
	}

	// Act
	connectionDetails, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cleanup)
	assert.Equal(t, "localhost", connectionDetails.Host)
	assert.Equal(t, expectedSecurePort, connectionDetails.Port)
	assert.Equal(t, expectedAPIKey, connectionDetails.APIKey)
	assert.Equal(t, expectedCertificatePEM, connectionDetails.CertificatePEM)

	assert.False(t, processCleanupCalled)
	err = cleanup()
	require.NoError(t, err)
	assert.True(t, processCleanupCalled)
}

func TestStarter_StartLocalMATLABSession_WithStartingDirectory(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDirPath := filepath.Join("tmp", "matlab-session-12345")
	expectedStartingDir := filepath.Join("home", "somewhere")
	expectedCertificateFile := filepath.Join("tmp", "matlab-session-12345", "cert.pem")
	expectedCertificateKeyFile := filepath.Join("tmp", "matlab-session-12345", "cert.key")
	expectedAPIKey := "test-api-key-12345"
	expectedMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")
	expectedSecurePort := "9999"
	expectedCertificatePEM := []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----")
	expectedEnv := []string{"MATLAB_MCP_API_KEY=" + expectedAPIKey}
	expectedStartupCode := "sessionPath = '" + expectedSessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"
	showDesktop := false
	expectedStartupFlags := []string{"-r", expectedStartupCode}
	expectedProcessID := 12345
	processCleanup := func() {}

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		Path().
		Return(expectedSessionDirPath).
		Once()

	mockProcessDetails.EXPECT().
		NewAPIKey().
		Return(expectedAPIKey).
		Once()

	mockDirectory.EXPECT().
		CertificateFile().
		Return(expectedCertificateFile).
		Once()

	mockDirectory.EXPECT().
		CertificateKeyFile().
		Return(expectedCertificateKeyFile).
		Once()

	mockProcessDetails.EXPECT().
		EnvironmentVariables(expectedSessionDirPath, expectedAPIKey, expectedCertificateFile, expectedCertificateKeyFile).
		Return(expectedEnv).
		Once()

	mockProcessDetails.EXPECT().
		StartupFlag(runtime.GOOS, showDesktop, expectedStartupCode).
		Return(expectedStartupFlags).
		Once()

	// Note: When starting directory is empty, it should use sessionDirPath
	mockMATLABProcessLauncher.EXPECT().
		Launch(mockLogger.AsMockArg(), expectedSessionDirPath, expectedMATLABRoot, expectedStartingDir, expectedStartupFlags, expectedEnv).
		Return(expectedProcessID, processCleanup, nil).
		Once()

	mockWatchdog.EXPECT().
		RegisterProcessPIDWithWatchdog(expectedProcessID).
		Return(nil).
		Once()

	mockDirectory.EXPECT().
		GetEmbeddedConnectorDetails().
		Return(expectedSecurePort, expectedCertificatePEM, nil).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		StartingDirectory:      expectedStartingDir,
		IsStartingDirectorySet: true,
		ShowMATLABDesktop:      showDesktop,
	}

	// Act
	connectionDetails, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cleanup)
	assert.Equal(t, "localhost", connectionDetails.Host)
	assert.Equal(t, expectedSecurePort, connectionDetails.Port)
	assert.Equal(t, expectedAPIKey, connectionDetails.APIKey)
	assert.Equal(t, expectedCertificatePEM, connectionDetails.CertificatePEM)
}

func TestStarter_StartLocalMATLABSession_DirectoryFactoryNewError(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := assert.AnError

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		MATLABRoot:             filepath.Join("usr", "local", "MATLAB", "R2024b"),
		StartingDirectory:      filepath.Join("home", "user", "workspace"),
		IsStartingDirectorySet: true,
		ShowMATLABDesktop:      false,
	}

	// Act
	connectionDetails, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, cleanup)
	assert.Equal(t, embeddedconnector.ConnectionDetails{}, connectionDetails)
}

func TestStarter_StartLocalMATLABSession_MATLABProcessLauncherError(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDirPath := filepath.Join("tmp", "matlab-session-12345")
	expectedCertificateFile := filepath.Join("tmp", "matlab-session-12345", "cert.pem")
	expectedCertificateKeyFile := filepath.Join("tmp", "matlab-session-12345", "cert.key")
	expectedAPIKey := "test-api-key-12345"
	expectedMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")
	expectedEnv := []string{"MATLAB_MCP_API_KEY=" + expectedAPIKey}
	expectedStartupCode := "sessionPath = '" + expectedSessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"
	expectedStartupFlags := []string{"-r", expectedStartupCode}
	expectedError := assert.AnError

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		Path().
		Return(expectedSessionDirPath).
		Once()

	mockProcessDetails.EXPECT().
		NewAPIKey().
		Return(expectedAPIKey).
		Once()

	mockDirectory.EXPECT().
		CertificateFile().
		Return(expectedCertificateFile).
		Once()

	mockDirectory.EXPECT().
		CertificateKeyFile().
		Return(expectedCertificateKeyFile).
		Once()

	mockProcessDetails.EXPECT().
		EnvironmentVariables(expectedSessionDirPath, expectedAPIKey, expectedCertificateFile, expectedCertificateKeyFile).
		Return(expectedEnv).
		Once()

	mockProcessDetails.EXPECT().
		StartupFlag(runtime.GOOS, false, expectedStartupCode).
		Return(expectedStartupFlags).
		Once()

	mockMATLABProcessLauncher.EXPECT().
		Launch(mockLogger.AsMockArg(), expectedSessionDirPath, expectedMATLABRoot, expectedSessionDirPath, expectedStartupFlags, expectedEnv).
		Return(0, nil, expectedError).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	// Act
	connectionDetails, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, cleanup)
	assert.Equal(t, embeddedconnector.ConnectionDetails{}, connectionDetails)
}

func TestStarter_StartLocalMATLABSession_RegisterProcessPIDWithWatchdogError(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedStartingDir := filepath.Join("somewhere")
	expectedSessionDirPath := filepath.Join("tmp", "matlab-session-12345")
	expectedCertificateFile := filepath.Join("tmp", "matlab-session-12345", "cert.pem")
	expectedCertificateKeyFile := filepath.Join("tmp", "matlab-session-12345", "cert.key")
	expectedAPIKey := "test-api-key-12345"
	expectedMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")
	expectedEnv := []string{"MATLAB_MCP_API_KEY=" + expectedAPIKey}
	expectedStartupCode := "sessionPath = '" + expectedSessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"
	expectedStartupFlags := []string{"-r", expectedStartupCode}
	expectedError := assert.AnError
	expectedProcessID := 12345
	processCleanup := func() {}
	expectedCertificatePEM := []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----")
	expectedSecurePort := "9999"
	showDesktop := false

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		Path().
		Return(expectedSessionDirPath).
		Once()

	mockProcessDetails.EXPECT().
		NewAPIKey().
		Return(expectedAPIKey).
		Once()

	mockDirectory.EXPECT().
		CertificateFile().
		Return(expectedCertificateFile).
		Once()

	mockDirectory.EXPECT().
		CertificateKeyFile().
		Return(expectedCertificateKeyFile).
		Once()

	mockProcessDetails.EXPECT().
		EnvironmentVariables(expectedSessionDirPath, expectedAPIKey, expectedCertificateFile, expectedCertificateKeyFile).
		Return(expectedEnv).
		Once()

	mockProcessDetails.EXPECT().
		StartupFlag(runtime.GOOS, false, expectedStartupCode).
		Return(expectedStartupFlags).
		Once()

	mockMATLABProcessLauncher.EXPECT().
		Launch(mockLogger.AsMockArg(), expectedSessionDirPath, expectedMATLABRoot, expectedStartingDir, expectedStartupFlags, expectedEnv).
		Return(expectedProcessID, processCleanup, nil).
		Once()

	mockWatchdog.EXPECT().
		RegisterProcessPIDWithWatchdog(expectedProcessID).
		Return(expectedError).
		Once()

	mockDirectory.EXPECT().
		GetEmbeddedConnectorDetails().
		Return(expectedSecurePort, expectedCertificatePEM, nil).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		StartingDirectory:      expectedStartingDir,
		IsStartingDirectorySet: true,
		ShowMATLABDesktop:      showDesktop,
	}

	// Act
	connectionDetails, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cleanup)
	assert.Equal(t, "localhost", connectionDetails.Host)
	assert.Equal(t, expectedSecurePort, connectionDetails.Port)
	assert.Equal(t, expectedAPIKey, connectionDetails.APIKey)
	assert.Equal(t, expectedCertificatePEM, connectionDetails.CertificatePEM)

	logs := mockLogger.WarnLogs()

	fields, found := logs["Failed to register process with watchdog"]
	require.True(t, found, "Failed to register process with watchdog")

	errField, found := fields["error"]
	require.True(t, found, "Expected an error field in the warning log")

	err, ok := errField.(error)
	require.True(t, ok, "Error field should be of type error")
	require.ErrorIs(t, err, expectedError, "Logged error should match the RegisterProcessPIDWithWatchdog error")

}

func TestStarter_StartLocalMATLABSession_GetEmbeddedConnectorDetailsError(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDirPath := filepath.Join("tmp", "matlab-session-12345")
	expectedCertificateFile := filepath.Join("tmp", "matlab-session-12345", "cert.pem")
	expectedCertificateKeyFile := filepath.Join("tmp", "matlab-session-12345", "cert.key")
	expectedAPIKey := "test-api-key-12345"
	expectedMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")
	expectedEnv := []string{"MATLAB_MCP_API_KEY=" + expectedAPIKey}
	expectedStartupCode := "sessionPath = '" + expectedSessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"
	expectedStartupFlags := []string{"-r", expectedStartupCode}
	expectedProcessID := 12345
	processCleanup := func() {}
	expectedError := assert.AnError

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		Path().
		Return(expectedSessionDirPath).
		Once()

	mockProcessDetails.EXPECT().
		NewAPIKey().
		Return(expectedAPIKey).
		Once()

	mockDirectory.EXPECT().
		CertificateFile().
		Return(expectedCertificateFile).
		Once()

	mockDirectory.EXPECT().
		CertificateKeyFile().
		Return(expectedCertificateKeyFile).
		Once()

	mockProcessDetails.EXPECT().
		EnvironmentVariables(expectedSessionDirPath, expectedAPIKey, expectedCertificateFile, expectedCertificateKeyFile).
		Return(expectedEnv).
		Once()

	mockProcessDetails.EXPECT().
		StartupFlag(runtime.GOOS, false, expectedStartupCode).
		Return(expectedStartupFlags).
		Once()

	mockMATLABProcessLauncher.EXPECT().
		Launch(mockLogger.AsMockArg(), expectedSessionDirPath, expectedMATLABRoot, expectedSessionDirPath, expectedStartupFlags, expectedEnv).
		Return(expectedProcessID, processCleanup, nil).
		Once()

	mockWatchdog.EXPECT().
		RegisterProcessPIDWithWatchdog(expectedProcessID).
		Return(nil).
		Once()

	mockDirectory.EXPECT().
		GetEmbeddedConnectorDetails().
		Return("", nil, expectedError).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	// Act
	connectionDetails, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, cleanup)
	assert.Equal(t, embeddedconnector.ConnectionDetails{}, connectionDetails)
}

func TestStarter_StartLocalMATLABSession_CleanupReturnsSessionCleanupError(t *testing.T) {
	// Arrange
	mockDirectoryFactory := &mocks.MockSessionDirectoryFactory{}
	defer mockDirectoryFactory.AssertExpectations(t)

	mockProcessDetails := &mocks.MockProcessDetails{}
	defer mockProcessDetails.AssertExpectations(t)

	mockMATLABProcessLauncher := &mocks.MockMATLABProcessLauncher{}
	defer mockMATLABProcessLauncher.AssertExpectations(t)

	mockWatchdog := &mocks.MockWatchdog{}
	defer mockWatchdog.AssertExpectations(t)

	mockDirectory := &directorymocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDirPath := filepath.Join("tmp", "matlab-session-12345")
	expectedCertificateFile := filepath.Join("tmp", "matlab-session-12345", "cert.pem")
	expectedCertificateKeyFile := filepath.Join("tmp", "matlab-session-12345", "cert.key")
	expectedAPIKey := "test-api-key-12345"
	expectedMATLABRoot := filepath.Join("usr", "local", "MATLAB", "R2024b")
	expectedSecurePort := "9999"
	expectedCertificatePEM := []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----")
	expectedEnv := []string{"MATLAB_MCP_API_KEY=" + expectedAPIKey}
	expectedStartupCode := "sessionPath = '" + expectedSessionDirPath + "';addpath(sessionPath);matlab_mcp.initializeMCP();clear sessionPath;"
	showDestop := false
	expectedStartupFlags := []string{"-r", expectedStartupCode}
	expectedProcessID := 12345
	processCleanup := func() {}
	expectedError := assert.AnError

	mockDirectoryFactory.EXPECT().
		New(mockLogger.AsMockArg()).
		Return(mockDirectory, nil).
		Once()

	mockDirectory.EXPECT().
		Path().
		Return(expectedSessionDirPath).
		Once()

	mockProcessDetails.EXPECT().
		NewAPIKey().
		Return(expectedAPIKey).
		Once()

	mockDirectory.EXPECT().
		CertificateFile().
		Return(expectedCertificateFile).
		Once()

	mockDirectory.EXPECT().
		CertificateKeyFile().
		Return(expectedCertificateKeyFile).
		Once()

	mockProcessDetails.EXPECT().
		EnvironmentVariables(expectedSessionDirPath, expectedAPIKey, expectedCertificateFile, expectedCertificateKeyFile).
		Return(expectedEnv).
		Once()

	mockProcessDetails.EXPECT().
		StartupFlag(runtime.GOOS, showDestop, expectedStartupCode).
		Return(expectedStartupFlags).
		Once()

	mockMATLABProcessLauncher.EXPECT().
		Launch(mockLogger.AsMockArg(), expectedSessionDirPath, expectedMATLABRoot, expectedSessionDirPath, expectedStartupFlags, expectedEnv).
		Return(expectedProcessID, processCleanup, nil).
		Once()

	mockWatchdog.EXPECT().
		RegisterProcessPIDWithWatchdog(expectedProcessID).
		Return(nil).
		Once()

	mockDirectory.EXPECT().
		GetEmbeddedConnectorDetails().
		Return(expectedSecurePort, expectedCertificatePEM, nil).
		Once()

	mockDirectory.EXPECT().
		Cleanup().
		Return(expectedError).
		Once()

	starter := localmatlabsession.NewStarter(
		mockDirectoryFactory,
		mockProcessDetails,
		mockMATLABProcessLauncher,
		mockWatchdog,
	)

	startRequest := datatypes.LocalSessionDetails{
		MATLABRoot:             expectedMATLABRoot,
		IsStartingDirectorySet: false,
	}

	_, cleanup, err := starter.StartLocalMATLABSession(mockLogger, startRequest)
	require.NoError(t, err)
	require.NotNil(t, cleanup)

	// Act

	err = cleanup()

	// Assert
	require.ErrorIs(t, err, expectedError)
}
