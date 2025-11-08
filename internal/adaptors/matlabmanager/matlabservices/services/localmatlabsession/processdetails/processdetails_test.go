// Copyright 2025 The MathWorks, Inc.

package processdetails_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processdetails"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processdetails"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	details := processdetails.New(mockOSLayer)

	// Assert
	assert.NotNil(t, details)
}

func TestProcessDetails_NewAPIKey_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	details := processdetails.New(mockOSLayer)

	// Act
	apiKey := details.NewAPIKey()

	// Assert
	assert.NotEmpty(t, apiKey)
	_, err := uuid.Parse(apiKey)
	require.NoError(t, err, "API key should be a valid UUID")
}

func TestProcessDetails_NewAPIKey_ReturnsUniqueValues(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	details := processdetails.New(mockOSLayer)

	// Act
	apiKey1 := details.NewAPIKey()
	apiKey2 := details.NewAPIKey()

	// Assert
	assert.NotEqual(t, apiKey1, apiKey2)
}

func TestProcessDetails_EnvironmentVariables_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	sessionDirPath := "/tmp/matlab-session-12345"
	apiKey := "test-api-key-12345"
	certificateFile := "/tmp/matlab-session-12345/cert.pem"
	certificateKey := "/tmp/matlab-session-12345/cert.key"
	existingEnv := []string{"PATH=/usr/bin", "HOME=/home/user"}

	mockOSLayer.EXPECT().
		Environ().
		Return(existingEnv).
		Once()

	details := processdetails.New(mockOSLayer)

	// Act
	env := details.EnvironmentVariables(sessionDirPath, apiKey, certificateFile, certificateKey)

	// Assert
	expectedEnv := append(
		existingEnv,
		[]string{
			"MATLAB_LOG_DIR=" + sessionDirPath,
			"MW_MCP_SESSION_DIR=" + sessionDirPath,
			`MW_DIAGNOSTIC_DEST="filedir=` + sessionDirPath + `"`,
			"MW_CONTEXT_TAGS=MATLAB:MATLAB_MCP_CORE_SERVER:V1",
			"MWAPIKEY=" + apiKey,
			"MW_CERTFILE=" + certificateFile,
			"MW_PKEYFILE=" + certificateKey,
		}...,
	)
	assert.ElementsMatch(t, expectedEnv, env)
}

func TestProcessDetails_EnvironmentVariables_EmptyExistingEnvironment(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	sessionDirPath := "/tmp/matlab-session-12345"
	apiKey := "test-api-key-12345"
	certificateFile := "/tmp/matlab-session-12345/cert.pem"
	certificateKey := "/tmp/matlab-session-12345/cert.key"
	existingEnv := []string{}

	mockOSLayer.EXPECT().
		Environ().
		Return(existingEnv).
		Once()

	details := processdetails.New(mockOSLayer)

	// Act
	env := details.EnvironmentVariables(sessionDirPath, apiKey, certificateFile, certificateKey)

	// Assert
	expectedEnv := []string{
		"MATLAB_LOG_DIR=" + sessionDirPath,
		"MW_MCP_SESSION_DIR=" + sessionDirPath,
		`MW_DIAGNOSTIC_DEST="filedir=` + sessionDirPath + `"`,
		"MW_CONTEXT_TAGS=MATLAB:MATLAB_MCP_CORE_SERVER:V1",
		"MWAPIKEY=" + apiKey,
		"MW_CERTFILE=" + certificateFile,
		"MW_PKEYFILE=" + certificateKey,
	}
	assert.ElementsMatch(t, expectedEnv, env)
}

func TestProcessDetails_EnvironmentVariables_MWContextTagPropagation(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	sessionDirPath := "/tmp/matlab-session-12345"
	apiKey := "test-api-key-12345"
	certificateFile := "/tmp/matlab-session-12345/cert.pem"
	certificateKey := "/tmp/matlab-session-12345/cert.key"
	expectedUnchangedExistingEnv := []string{
		"PATH=/usr/bin",
		"HOME=/home/user",
	}
	existingMWContextTagEnv := []string{
		"MW_CONTEXT_TAGS=SomeOtherProduct",
	}

	mockOSLayer.EXPECT().
		Environ().
		Return(append(expectedUnchangedExistingEnv, existingMWContextTagEnv...)).
		Once()

	details := processdetails.New(mockOSLayer)

	// Act
	env := details.EnvironmentVariables(sessionDirPath, apiKey, certificateFile, certificateKey)

	// Assert
	expectedEnv := append(
		expectedUnchangedExistingEnv,
		[]string{
			"MATLAB_LOG_DIR=" + sessionDirPath,
			"MW_MCP_SESSION_DIR=" + sessionDirPath,
			`MW_DIAGNOSTIC_DEST="filedir=` + sessionDirPath + `"`,
			"MW_CONTEXT_TAGS=SomeOtherProduct,MATLAB:MATLAB_MCP_CORE_SERVER:V1",
			"MWAPIKEY=" + apiKey,
			"MW_CERTFILE=" + certificateFile,
			"MW_PKEYFILE=" + certificateKey,
		}...,
	)
	assert.ElementsMatch(t, expectedEnv, env)
}

func TestProcessDetails_StartupFlag_HappyPath(t *testing.T) {
	for _, testConfig := range []struct {
		os            string
		showDesktop   bool
		expectedFlags []string
	}{
		{
			os:          "linux",
			showDesktop: true,
			expectedFlags: []string{
				"-desktop",
			},
		},
		{
			os:          "linux",
			showDesktop: false,
			expectedFlags: []string{
				"-nosplash",
				"-softwareopengl",
				"-nodesktop",
				"-minimize",
			},
		},
		{
			os:          "windows",
			showDesktop: true,
			expectedFlags: []string{
				"-desktop",
			},
		},
		{
			os:          "windows",
			showDesktop: false,
			expectedFlags: []string{
				"-nosplash",
				"-softwareopengl",
				"-nodesktop",
				"-noDisplayDesktop",
				"-wait",
				"-log",
				"/minimize",
			},
		},
		{
			os:          "darwin",
			showDesktop: true,
			expectedFlags: []string{
				"-desktop",
			},
		},
		{
			os:          "darwin",
			showDesktop: false,
			expectedFlags: []string{
				"-nosplash",
				"-softwareopengl",
				"-nodesktop",
				"-minimize",
			},
		},
	} {
		t.Run(fmt.Sprintf("%s_desktop_%v", testConfig.os, testConfig.showDesktop), func(t *testing.T) {
			// Arrange
			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			details := processdetails.New(mockOSLayer)
			startupCode := "disp('Hello World');"

			// Act
			flags := details.StartupFlag(testConfig.os, testConfig.showDesktop, startupCode)

			// Assert
			assert.Equal(t,
				append(
					testConfig.expectedFlags,
					"-r",
					startupCode,
				),
				flags,
			)
		})
	}
}
