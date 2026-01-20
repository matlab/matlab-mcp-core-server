// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pingRetry = 10 * time.Millisecond
const pingTimeout = 40 * time.Millisecond

func TestClient_Ping_HappyPath(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForState(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertPingMessage(t, request)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				PingResponse: []embeddedconnector.PingResponseMessage{
					{
						MessageFaults: []json.RawMessage{},
					},
				},
			},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		assert.NoError(t, json.NewEncoder(responseWriter).Encode(response))
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.True(t, response.IsAlive)
}

func TestClient_Ping_MATLABNotAvailable(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedErrorMessage = "MATLAB is not available"
	fault := json.RawMessage(`{"message":"` + expectedErrorMessage + `","faultCode":"MATLAB.PingError"}`)
	connectionDetails := startTestServerForState(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertPingMessage(t, request)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				PingResponse: []embeddedconnector.PingResponseMessage{
					{
						MessageFaults: []json.RawMessage{fault},
					},
				},
			},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		assert.NoError(t, json.NewEncoder(responseWriter).Encode(response))
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	client.SetPingRetry(pingRetry)
	client.SetPingTimeout(pingTimeout)

	require.NoError(t, err)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_HTTPError(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForState(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	client.SetPingRetry(pingRetry)
	client.SetPingTimeout(pingTimeout)

	require.NoError(t, err)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_NoResponseMessages(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForState(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				PingResponse: []embeddedconnector.PingResponseMessage{},
			},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		err := json.NewEncoder(responseWriter).Encode(response)
		assert.NoError(t, err)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	client.SetPingRetry(pingRetry)
	client.SetPingTimeout(pingTimeout)

	require.NoError(t, err)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_InvalidJSONResponse(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForState(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		_, err := responseWriter.Write([]byte("invalid json"))
		assert.NoError(t, err)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	client.SetPingRetry(pingRetry)
	client.SetPingTimeout(pingTimeout)

	require.NoError(t, err)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_ContextCancellation(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		t.Error("Handler should not be called when context is cancelled")
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	cancel() // Cancel immediately

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_RetriesOnError(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const failuresBeforeSuccess = 2
	var requestCount int32

	connectionDetails := startTestServerForState(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)

		if count <= failuresBeforeSuccess {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		assertPingMessage(t, request)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				PingResponse: []embeddedconnector.PingResponseMessage{
					{
						MessageFaults: []json.RawMessage{},
					},
				},
			},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		assert.NoError(t, json.NewEncoder(responseWriter).Encode(response))
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	client.SetPingRetry(1 * time.Millisecond)
	client.SetPingTimeout(5 * time.Second)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.True(t, response.IsAlive, "Ping should succeed after retries")
	finalCount := atomic.LoadInt32(&requestCount)
	assert.Equal(t, int32(failuresBeforeSuccess+1), finalCount, "Expected exactly %d attempts", failuresBeforeSuccess+1)
}

func assertPingMessage(t *testing.T, request *http.Request) {
	var requestPayload embeddedconnector.ConnectorPayload
	err := json.NewDecoder(request.Body).Decode(&requestPayload)
	require.NoError(t, err)

	require.NotNil(t, requestPayload.Messages)
	require.Len(t, requestPayload.Messages.Ping, 1, "expected exactly one Ping message")
}
