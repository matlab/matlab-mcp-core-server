// Copyright 2026 The MathWorks, Inc.

package sessiondiscovery_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/sessionselector/sessiondiscovery"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/matlabmanager/sessionselector/sessiondiscovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	result := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Assert
	require.NotNil(t, result)
}

func TestSessionDiscoverer_FromSessionDetails_HappyPath(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCertPath := filepath.Join("path", "to", "cert.pem")
	expectedCertPEM := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")
	expectedPort := "31515"
	expectedAPIKey := "test-api-key"

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        31515,
		"certificate": expectedCertPath,
		"apiKey":      expectedAPIKey,
		"pid":         12345,
	})

	mockOSLayer.EXPECT().
		ReadFile(expectedCertPath).
		Return(expectedCertPEM, nil).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "localhost", result.Host)
	assert.Equal(t, expectedPort, result.Port)
	assert.Equal(t, expectedAPIKey, result.APIKey)
	assert.Equal(t, expectedCertPEM, result.CertificatePEM)
	assert.Empty(t, result.BasePath)
}

func TestSessionDiscoverer_FromSessionDetails_BasePathPassedThrough(t *testing.T) {
	// The base path is stored verbatim; normalisation of slashes happens
	// later when the endpoint URL is built via url.JoinPath.
	testCases := []struct {
		name     string
		basePath string
	}{
		{name: "LeadingAndTrailingSlash", basePath: "/matlab/"},
		{name: "NoLeadingSlash", basePath: "proxy/matlab/"},
		{name: "NoTrailingSlash", basePath: "/user/foo"},
		{name: "BareSlash", basePath: "/"},
		{name: "Empty", basePath: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
			defer mockAppDataDirGetter.AssertExpectations(t)

			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockLogger := testutils.NewInspectableLogger()

			expectedCertPath := filepath.Join("path", "to", "cert.pem")
			expectedCertPEM := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")
			expectedPort := "31515"
			expectedAPIKey := "test-api-key"

			sessionJSON := marshallSessionDetails(t, map[string]any{
				"port":        31515,
				"certificate": expectedCertPath,
				"apiKey":      expectedAPIKey,
				"pid":         12345,
				"basePath":    tc.basePath,
			})

			mockOSLayer.EXPECT().
				ReadFile(expectedCertPath).
				Return(expectedCertPEM, nil).
				Once()

			discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

			// Act
			result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, "localhost", result.Host)
			assert.Equal(t, expectedPort, result.Port)
			assert.Equal(t, expectedAPIKey, result.APIKey)
			assert.Equal(t, expectedCertPEM, result.CertificatePEM)
			assert.Equal(t, tc.basePath, result.BasePath)
		})
	}
}

func TestSessionDiscoverer_FromSessionDetails_InvalidBasePath(t *testing.T) {
	// A base path carrying a scheme, host, query or fragment could redirect
	// requests away from the local session and must be rejected.
	testCases := []struct {
		name     string
		basePath string
	}{
		{name: "WithScheme", basePath: "https://evil.example.com/matlab/"},
		{name: "WithHost", basePath: "//evil.example.com/matlab/"},
		{name: "WithQuery", basePath: "/matlab/?redirect=evil"},
		{name: "WithFragment", basePath: "/matlab/#frag"},
		{name: "WithDotSegment", basePath: "/matlab/./session/"},
		{name: "WithDotDotSegment", basePath: "/user/alice/../bob/"},
		{name: "WithEncodedDotDotSegment", basePath: "/user/alice/%2e%2e/bob/"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
			defer mockAppDataDirGetter.AssertExpectations(t)

			mockOSLayer := &mocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockLogger := testutils.NewInspectableLogger()

			expectedCertPath := filepath.Join("path", "to", "cert.pem")
			expectedCertPEM := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

			sessionJSON := marshallSessionDetails(t, map[string]any{
				"port":        31515,
				"certificate": expectedCertPath,
				"apiKey":      "test-api-key",
				"pid":         12345,
				"basePath":    tc.basePath,
			})

			mockOSLayer.EXPECT().
				ReadFile(expectedCertPath).
				Return(expectedCertPEM, nil).
				Once()

			discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

			// Act
			result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

			// Assert
			require.ErrorIs(t, err, sessiondiscovery.ErrInvalidSessionDetails)
			assert.Empty(t, result.Host)
		})
	}
}

func TestSessionDiscoverer_FromSessionDetails_InvalidJSON(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	invalidJSON := []byte("not valid json")

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, invalidJSON)

	// Assert
	require.Error(t, err)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_MissingPort(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"certificate": "cert.pem",
		"apiKey":      "key",
		"pid":         100,
	})

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.ErrorIs(t, err, sessiondiscovery.ErrInvalidSessionDetails)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_PortNotANumber(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionJSON := []byte(`{"port":"not-a-number","certificate":"cert.pem","apiKey":"key","pid":100}`)

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.Error(t, err)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_PortNotAnInt(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionJSON := []byte(`{"port":"1.5","certificate":"cert.pem","apiKey":"key","pid":100}`)

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.Error(t, err)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_PortBelowRange(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        0,
		"certificate": filepath.Join("path", "to", "cert.pem"),
		"apiKey":      "key",
		"pid":         100,
	})

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.ErrorIs(t, err, sessiondiscovery.ErrInvalidSessionDetails)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_PortAboveRange(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        65536,
		"certificate": filepath.Join("path", "to", "cert.pem"),
		"apiKey":      "key",
		"pid":         100,
	})

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.ErrorIs(t, err, sessiondiscovery.ErrInvalidSessionDetails)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_EmptyAPIKey(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        31515,
		"certificate": "",
		"apiKey":      "",
		"pid":         100,
	})

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.ErrorIs(t, err, sessiondiscovery.ErrInvalidSessionDetails)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_CertificateReadError(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCertPath := filepath.Join("path", "to", "cert.pem")

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        31515,
		"certificate": expectedCertPath,
		"apiKey":      "key",
		"pid":         100,
	})

	mockOSLayer.EXPECT().
		ReadFile(expectedCertPath).
		Return(nil, assert.AnError).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_FromSessionDetails_EmptyCertificatePEM(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedCertPath := filepath.Join("path", "to", "cert.pem")

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        31515,
		"certificate": expectedCertPath,
		"apiKey":      "test-api-key",
		"pid":         100,
	})

	mockOSLayer.EXPECT().
		ReadFile(expectedCertPath).
		Return([]byte(""), nil).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	result, err := discoverer.FromSessionDetails(mockLogger, sessionJSON)

	// Assert
	require.ErrorIs(t, err, sessiondiscovery.ErrInvalidSessionDetails)
	assert.Empty(t, result.Host)
}

func TestSessionDiscoverer_DiscoverSessions_HappyPath(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedAppDataDir := filepath.Join("home", "user", "MATLABMCPServer")
	expectedCertPath := filepath.Join("path", "to", "cert.pem")
	expectedCertPEM := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")
	expectedPort := 31515
	expectedAPIKey := "test-api-key"

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        expectedPort,
		"certificate": expectedCertPath,
		"apiKey":      expectedAPIKey,
		"pid":         12345,
	})

	mockAppDataDirGetter.EXPECT().
		AppDataDir().
		Return(expectedAppDataDir, nil).
		Once()

	expectedSessionFile := filepath.Join(expectedAppDataDir, "v1", "sessionDetails.json")
	mockOSLayer.EXPECT().
		ReadFile(expectedSessionFile).
		Return(sessionJSON, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(expectedCertPath).
		Return(expectedCertPEM, nil).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	sessions := discoverer.DiscoverSessions(mockLogger)

	// Assert
	require.Len(t, sessions, 1)
	assert.Equal(t, "localhost", sessions[0].Host)
	assert.Equal(t, fmt.Sprintf("%d", expectedPort), sessions[0].Port)
	assert.Equal(t, expectedAPIKey, sessions[0].APIKey)
	assert.Equal(t, expectedCertPEM, sessions[0].CertificatePEM)
}

func TestSessionDiscoverer_DiscoverSessions_AppDataDirError(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockAppDataDirGetter.EXPECT().
		AppDataDir().
		Return("", assert.AnError).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	sessions := discoverer.DiscoverSessions(mockLogger)

	// Assert
	assert.Empty(t, sessions)
	assert.Len(t, mockLogger.DebugLogs(), 1)
}

func TestSessionDiscoverer_DiscoverSessions_SessionFileNotFound(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedAppDataDir := filepath.Join("home", "user", "MATLABMCPServer")

	mockAppDataDirGetter.EXPECT().
		AppDataDir().
		Return(expectedAppDataDir, nil).
		Once()

	expectedSessionFile := filepath.Join(expectedAppDataDir, "v1", "sessionDetails.json")
	mockOSLayer.EXPECT().
		ReadFile(expectedSessionFile).
		Return(nil, assert.AnError).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	sessions := discoverer.DiscoverSessions(mockLogger)

	// Assert
	assert.Empty(t, sessions)
	assert.Len(t, mockLogger.DebugLogs(), 1)
}

func TestSessionDiscoverer_DiscoverSessions_InvalidJSON(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedAppDataDir := filepath.Join("home", "user", "MATLABMCPServer")

	mockAppDataDirGetter.EXPECT().
		AppDataDir().
		Return(expectedAppDataDir, nil).
		Once()

	expectedSessionFile := filepath.Join(expectedAppDataDir, "v1", "sessionDetails.json")
	mockOSLayer.EXPECT().
		ReadFile(expectedSessionFile).
		Return([]byte("not valid json"), nil).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	sessions := discoverer.DiscoverSessions(mockLogger)

	// Assert
	assert.Empty(t, sessions)
	assert.Len(t, mockLogger.DebugLogs(), 1)
}

func TestSessionDiscoverer_DiscoverSessions_InvalidPort(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedAppDataDir := filepath.Join("home", "user", "MATLABMCPServer")

	sessionJSON := []byte(`{"port":"not-a-number","certificate":"cert.pem","apiKey":"key","pid":100}`)

	mockAppDataDirGetter.EXPECT().
		AppDataDir().
		Return(expectedAppDataDir, nil).
		Once()

	expectedSessionFile := filepath.Join(expectedAppDataDir, "v1", "sessionDetails.json")
	mockOSLayer.EXPECT().
		ReadFile(expectedSessionFile).
		Return(sessionJSON, nil).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	sessions := discoverer.DiscoverSessions(mockLogger)

	// Assert
	assert.Empty(t, sessions)
	assert.Len(t, mockLogger.DebugLogs(), 1)
}

func TestSessionDiscoverer_DiscoverSessions_CertificateReadError(t *testing.T) {
	// Arrange
	mockAppDataDirGetter := &mocks.MockAppDataDirGetter{}
	defer mockAppDataDirGetter.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedAppDataDir := filepath.Join("home", "user", "MATLABMCPServer")
	expectedCertPath := filepath.Join("path", "to", "cert.pem")

	sessionJSON := marshallSessionDetails(t, map[string]any{
		"port":        31515,
		"certificate": expectedCertPath,
		"apiKey":      "key",
		"pid":         100,
	})

	mockAppDataDirGetter.EXPECT().
		AppDataDir().
		Return(expectedAppDataDir, nil).
		Once()

	expectedSessionFile := filepath.Join(expectedAppDataDir, "v1", "sessionDetails.json")
	mockOSLayer.EXPECT().
		ReadFile(expectedSessionFile).
		Return(sessionJSON, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(expectedCertPath).
		Return(nil, assert.AnError).
		Once()

	discoverer := sessiondiscovery.New(mockAppDataDirGetter, mockOSLayer)

	// Act
	sessions := discoverer.DiscoverSessions(mockLogger)

	// Assert
	assert.Empty(t, sessions)
	assert.Len(t, mockLogger.DebugLogs(), 1)
}

func marshallSessionDetails(t *testing.T, rawData map[string]any) []byte {
	t.Helper()

	data, err := json.Marshal(rawData)
	require.NoError(t, err, "Failed to marshall session details")

	return data
}
