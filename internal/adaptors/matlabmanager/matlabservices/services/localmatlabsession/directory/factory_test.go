// Copyright 2025-2026 The MathWorks, Inc.

package directory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	applicationdirectorymocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/directory"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	// Act
	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(expectedSessionDir, "+matlab_mcp")
	expectedCertificateFile := filepath.Join(expectedSessionDir, "cert.pem")
	expectedCertificateKeyFile := filepath.Join(expectedSessionDir, "cert.key")
	expectedMATLABFiles := map[string][]byte{
		"initializeMCP.m": []byte("some content"),
		"eval.m":          []byte("some other content"),
	}

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return(expectedSessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	mockMATLABFiles.EXPECT().
		GetAll().
		Return(expectedMATLABFiles).
		Once()

	for fileName, fileContent := range expectedMATLABFiles {
		filePath := filepath.Join(packageDir, fileName)
		mockOSLayer.EXPECT().
			WriteFile(filePath, fileContent, os.FileMode(0o600)).
			Return(nil).
			Once()
	}

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, dir)
	assert.Equal(t, expectedSessionDir, dir.Path())
	assert.Equal(t, expectedCertificateFile, dir.CertificateFile())
	assert.Equal(t, expectedCertificateKeyFile, dir.CertificateKeyFile())
}

func TestFactory_New_DirectoryFactoryError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := messages.AnError

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(nil, expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_MkdirTempError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := messages.AnError

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return("", expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_PackageDirectoryMkdirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}

func TestFactory_New_WriteFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectoryFactory := &mocks.MockApplicationDirectoryFactory{}
	defer mockApplicationDirectoryFactory.AssertExpectations(t)

	mockApplicationDirectory := &applicationdirectorymocks.MockDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError

	mockApplicationDirectoryFactory.EXPECT().
		Directory().
		Return(mockApplicationDirectory, nil).
		Once()

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	expectedFailingFileName := "initializeMCP.m"

	expectedMATLABFiles := map[string][]byte{
		expectedFailingFileName: []byte("some other content"),
	}

	mockMATLABFiles.EXPECT().
		GetAll().
		Return(expectedMATLABFiles).
		Once()

	mockOSLayer.EXPECT().
		WriteFile(filepath.Join(packageDir, expectedFailingFileName), expectedMATLABFiles[expectedFailingFileName], os.FileMode(0o600)).
		Return(expectedError).
		Once()

	factory := directory.NewFactory(mockOSLayer, mockApplicationDirectoryFactory, mockMATLABFiles)

	// Act
	dir, err := factory.New(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, dir)
}
