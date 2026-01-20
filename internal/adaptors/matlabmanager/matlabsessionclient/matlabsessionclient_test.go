// Copyright 2025-2026 The MathWorks, Inc.

package matlabsessionclient_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	httpclientmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/client"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabsessionclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockHTTPClientFactory := &mocks.MockHttpClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	// Act
	factory := matlabsessionclient.NewFactory(mockHTTPClientFactory)

	// Assert
	assert.NotNil(t, factory)
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockHTTPClientFactory := &mocks.MockHttpClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	mockHTTPClient := &httpclientmocks.MockHttpClient{}
	defer mockHTTPClient.AssertExpectations(t)

	expectedCertificatePEM := []byte("some cert")
	mockHTTPClientFactory.EXPECT().
		NewClientForSelfSignedTLSServer(expectedCertificatePEM).
		Return(mockHTTPClient, nil).
		Once()

	factory := matlabsessionclient.NewFactory(mockHTTPClientFactory)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "9910",
		APIKey:         "test-api-key",
		CertificatePEM: expectedCertificatePEM,
	}

	// Act
	client, err := factory.New(connectionDetails)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestFactory_New_HttpClientCreationErrors(t *testing.T) {
	// Arrange
	mockHTTPClientFactory := &mocks.MockHttpClientFactory{}
	defer mockHTTPClientFactory.AssertExpectations(t)

	expectedCertificatePEM := []byte("some cert")
	expectedError := assert.AnError

	mockHTTPClientFactory.EXPECT().
		NewClientForSelfSignedTLSServer(expectedCertificatePEM).
		Return(nil, expectedError).
		Once()

	factory := matlabsessionclient.NewFactory(mockHTTPClientFactory)

	connectionDetails := embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           "9910",
		APIKey:         "test-api-key",
		CertificatePEM: expectedCertificatePEM,
	}

	// Act
	client, err := factory.New(connectionDetails)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}
