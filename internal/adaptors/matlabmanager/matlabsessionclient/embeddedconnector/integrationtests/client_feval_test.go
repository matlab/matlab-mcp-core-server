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

func TestClient_FEval_HappyPath(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedFunction := "sum"
	expectedArguments := []string{"1", "2"}
	expectedNumOutputs := 1
	expectedResults := []interface{}{"3"}

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, expectedFunction, expectedArguments, expectedNumOutputs)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: expectedResults,
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
	fevalRequest := entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResults, response.Outputs)
}

func TestClient_FEval_MultipleOutputs(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedFunction := "size"
	expectedArguments := []string{"a"}
	expectedNumOutputs := 2
	expectedResults := []interface{}{"2", "3"}

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, expectedFunction, expectedArguments, expectedNumOutputs)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: expectedResults,
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
	fevalRequest := entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResults, response.Outputs)
}

func TestClient_FEval_NoArguments(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedFunction := "rand"
	expectedArguments := []string{}
	expectedNumOutputs := 1
	expectedResults := []interface{}{"0.8147"}

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, expectedFunction, expectedArguments, expectedNumOutputs)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: expectedResults,
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
	fevalRequest := entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedResults, response.Outputs)
}

func TestClient_FEval_MATLABError(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedFunction := "invalid_function"
	expectedArguments := []string{}
	expectedNumOutputs := 1
	expectedErrorMessage := "Undefined function 'invalid_function'"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, expectedFunction, expectedArguments, expectedNumOutputs)

		faultMessage := embeddedconnector.Fault{
			Message: expectedErrorMessage,
		}
		faultBytes, err := json.Marshal(faultMessage)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError:       true,
						MessageFaults: []json.RawMessage{faultBytes},
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
	fevalRequest := entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrorMessage)
	assert.Nil(t, response.Outputs)
}

func TestClient_FEval_MATLABErrorWithMultipleFaults(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedFunction := "invalid_function"
	expectedArguments := []string{}
	expectedNumOutputs := 1
	expectedErrorMessage1 := "First error message"
	expectedErrorMessage2 := "Second error message"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, expectedFunction, expectedArguments, expectedNumOutputs)

		fault1 := embeddedconnector.Fault{
			Message: expectedErrorMessage1,
		}
		fault1Bytes, err := json.Marshal(fault1)
		assert.NoError(t, err)

		fault2 := embeddedconnector.Fault{
			Message: expectedErrorMessage2,
		}
		fault2Bytes, err := json.Marshal(fault2)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError:       true,
						MessageFaults: []json.RawMessage{fault1Bytes, fault2Bytes},
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
	fevalRequest := entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrorMessage1)
	assert.Contains(t, err.Error(), expectedErrorMessage2)
	assert.Nil(t, response.Outputs)
}

func TestClient_FEval_MATLABErrorWithNoFaults(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	expectedFunction := "invalid_function"
	expectedArguments := []string{}
	expectedNumOutputs := 1

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, expectedFunction, expectedArguments, expectedNumOutputs)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError:       true,
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
	fevalRequest := entities.FEvalRequest{
		Function:   expectedFunction,
		Arguments:  expectedArguments,
		NumOutputs: expectedNumOutputs,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "response was in error state but no fault messages received")
	assert.Nil(t, response.Outputs)
}

func TestClient_FEval_HTTPError(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	})

	client, err := embeddedconnector.NewClient(connectionDetails, httpClientFactory)
	require.NoError(t, err)

	ctx := t.Context()
	fevalRequest := entities.FEvalRequest{
		Function:   "sum",
		Arguments:  []string{"1", "2"},
		NumOutputs: 1,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
	assert.Nil(t, response.Outputs)
}

func TestClient_FEval_NoResponseMessages(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{},
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
	fevalRequest := entities.FEvalRequest{
		Function:   "sum",
		Arguments:  []string{"1", "2"},
		NumOutputs: 1,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no response messages received")
	assert.Nil(t, response.Outputs)
}

func TestClient_FEval_InvalidJSONResponse(t *testing.T) {
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
	fevalRequest := entities.FEvalRequest{
		Function:   "sum",
		Arguments:  []string{"1", "2"},
		NumOutputs: 1,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal response")
	assert.Nil(t, response.Outputs)
}

func TestClient_FEval_ContextCancellation(t *testing.T) {
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

	fevalRequest := entities.FEvalRequest{
		Function:   "sum",
		Arguments:  []string{"1", "2"},
		NumOutputs: 1,
	}

	// Act
	response, err := client.FEval(ctx, mockLogger, fevalRequest)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
	assert.Nil(t, response.Outputs)
}

func assertFevalMessage(t *testing.T, request *http.Request, expectedFunction string, expectedArgs []string, expectedNumOutputs int) {
	var requestPayload embeddedconnector.ConnectorPayload
	err := json.NewDecoder(request.Body).Decode(&requestPayload)
	require.NoError(t, err)

	require.NotNil(t, requestPayload.Messages)
	require.Len(t, requestPayload.Messages.FEval, 1, "expected exactly one FEval message")
	assert.Equal(t, expectedFunction, requestPayload.Messages.FEval[0].Function, "feval function does not match expected function")
	assert.Equal(t, expectedArgs, requestPayload.Messages.FEval[0].Arguments, "feval args does not match expected args")
	assert.Equal(t, expectedNumOutputs, requestPayload.Messages.FEval[0].Nargout, "feval nargout does not match expected nargout")
}
