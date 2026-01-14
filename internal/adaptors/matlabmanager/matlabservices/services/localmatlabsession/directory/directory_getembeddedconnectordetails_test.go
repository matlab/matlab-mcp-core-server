// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"os"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectory_GetEmbeddedConnectorDetails_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForSecurePortFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForCertificateFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(nil, os.ErrNotExist).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForNotEmptyPortFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(""), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_WaitsForNotEmptyCertificateFile(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	expectedPort := "9999"
	expectedCertificate := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return([]byte(""), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(expectedCertificate, nil).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPort, port)
	assert.Equal(t, expectedCertificate, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_TimesoutWaitingForFilesToExists(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(nil, os.ErrNotExist)

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_TimesoutWaitingForFileContent(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil) // Will be called multiple times in wait loop

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(""), nil) // Will be called multiple times in wait loop

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_ReadSecurePortFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read secure port file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}

func TestDirectory_GetEmbeddedConnectorDetails_ReadCertificateFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFileInfo := &osfacademocks.MockFileInfo{}
	defer mockFileInfo.AssertExpectations(t)

	sessionDir := "/tmp/matlab-session-12345"

	dir := directory.NewDirectory(sessionDir, mockOSLayer)
	dir.SetEmbeddedConnectorDetailsTimeout(100 * time.Millisecond)
	dir.SetEmbeddedConnectorDetailsRetry(10 * time.Millisecond)

	securePortFile := dir.SecurePortFile()
	certificateFile := dir.CertificateFile()

	expectedPort := "9999"

	mockOSLayer.EXPECT().
		Stat(securePortFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		Stat(certificateFile).
		Return(mockFileInfo, nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(securePortFile).
		Return([]byte(expectedPort), nil).
		Once()

	mockOSLayer.EXPECT().
		ReadFile(certificateFile).
		Return(nil, assert.AnError).
		Once()

	// Act
	port, certificate, err := dir.GetEmbeddedConnectorDetails()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read certificate path file")
	assert.Empty(t, port)
	assert.Empty(t, certificate)
}
