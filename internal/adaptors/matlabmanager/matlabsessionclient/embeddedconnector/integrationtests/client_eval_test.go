// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Eval_HappyPath(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "disp('Hello World')"
	const expectedOutput = "Hello World\n"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertEvalMessage(t, request, expectedCode)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				EvalResponse: []embeddedconnector.EvalResponseMessage{
					{
						IsError:     false,
						ResponseStr: expectedOutput,
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
	evalRequest := entities.EvalRequest{
		Code: expectedCode,
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_MATLABError(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedCode := "invalid_function()"
	expectedErrorMessage := "Undefined function 'invalid_function'"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertEvalMessage(t, request, expectedCode)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				EvalResponse: []embeddedconnector.EvalResponseMessage{
					{
						IsError:     true,
						ResponseStr: expectedErrorMessage,
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
	evalRequest := entities.EvalRequest{
		Code: expectedCode,
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrorMessage)
	assert.Empty(t, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_HTTPError(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "disp('test')",
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
	assert.Empty(t, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_NoResponseMessages(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				EvalResponse: []embeddedconnector.EvalResponseMessage{},
			},
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		err := json.NewEncoder(responseWriter).Encode(response)
		assert.NoError(t, err)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "disp('test')",
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no response messages received")
	assert.Empty(t, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_InvalidJSONResponse(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		_, err := responseWriter.Write([]byte("invalid json"))
		assert.NoError(t, err)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "disp('test')",
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
	assert.Empty(t, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_Eval_ContextCancellation(t *testing.T) {
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

	evalRequest := entities.EvalRequest{
		Code: "disp('test')",
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
	assert.Empty(t, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func assertEvalMessage(t *testing.T, request *http.Request, expectedCode string) {
	var requestPayload embeddedconnector.ConnectorPayload
	err := json.NewDecoder(request.Body).Decode(&requestPayload)
	require.NoError(t, err)

	require.NotNil(t, requestPayload.Messages)
	require.Len(t, requestPayload.Messages.Eval, 1, "expected exactly one Eval message")
	assert.Equal(t, expectedCode, requestPayload.Messages.Eval[0].Code, "eval code does not match expected code")
}
