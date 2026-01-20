// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_integration_test

import (
	"encoding/base64"
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

func TestClient_EvalWithCapture_HappyPath(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "disp('Hello World')"
	const expectedOutput = "Hello World"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type:     "execute_result",
				MimeType: []string{"text/plain"},
				Value:    []json.RawMessage{json.RawMessage(`"` + expectedOutput + `"`)},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_ReturnImages(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "plot(1:10)"

	expectedImageData := []byte("image data")
	expectedImageBase64 := base64.StdEncoding.EncodeToString(expectedImageData)

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type:     "execute_result",
				MimeType: []string{"image/png"},
				Value:    []json.RawMessage{json.RawMessage(`"` + expectedImageBase64 + `"`)},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, response.ConsoleOutput)
	assert.Equal(t, [][]byte{expectedImageData}, response.Images)
}

func TestClient_EvalWithCapture_ReturnStreams(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "undefined_function"
	const expectedOutput = "some error"
	const expectedName = "stderr"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: expectedOutput,
					Name: expectedName,
				},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MultipleStreams_SameName(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "disp('line1'); disp('line2')"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "line1\n",
					Name: "stdout",
				},
			},
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "line2\n",
					Name: "stdout",
				},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "line1\nline2\n", response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MultipleStreams_DifferentNames(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "fprintf('output'); warning('warning message')"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "output",
					Name: "stdout",
				},
			},
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "Warning: warning message",
					Name: "stderr",
				},
			},
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: " continued",
					Name: "stderr",
				},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "output\nWarning: warning message continued", response.ConsoleOutput)
	assert.Nil(t, response.Images)
}

func TestClient_EvalWithCapture_MixedStreamsAndResults(t *testing.T) {
	// Arrange
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "x = 5; disp('calculating'); y = x * 2; plot(1:y)"

	expectedImageData := []byte("plot_image_data")
	expectedImageBase64 := base64.StdEncoding.EncodeToString(expectedImageData)

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "calculating\n",
					Name: "stdout",
				},
			},
			{
				Type:     "execute_result",
				MimeType: []string{"text/plain"},
				Value:    []json.RawMessage{json.RawMessage(`"y = 10"`)},
			},
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "more output",
					Name: "stdout",
				},
			},
			{
				Type:     "execute_result",
				MimeType: []string{"image/png"},
				Value:    []json.RawMessage{json.RawMessage(`"` + expectedImageBase64 + `"`)},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "calculating\n\ny = 10\nmore output", response.ConsoleOutput)
	assert.Equal(t, [][]byte{expectedImageData}, response.Images)
}

func TestClient_EvalWithCapture_StreamsWithInterruptionByExecuteResult(t *testing.T) {
	httpClientFactory := client.NewFactory()
	mockLogger := testutils.NewInspectableLogger()

	const expectedCode = "warning('first'); x = 1; warning('second')"

	connectionDetails := startTestServerForEvaluation(t, func(responseWriter http.ResponseWriter, request *http.Request) {
		assertFevalMessage(t, request, "matlab_mcp.mcpEval", []string{expectedCode}, 1)

		liveEditorResponseEntries := []embeddedconnector.LiveEditorResponseEntry{
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "Warning: first",
					Name: "stderr",
				},
			},
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: " warning message",
					Name: "stderr",
				},
			},
			{
				Type:     "execute_result",
				MimeType: []string{"text/plain"},
				Value:    []json.RawMessage{json.RawMessage(`"x = 1"`)},
			},
			{
				Type: "stream",
				Content: struct {
					Text string `json:"text"`
					Name string `json:"name"`
				}{
					Text: "Warning: second warning",
					Name: "stderr",
				},
			},
		}
		data, err := json.Marshal(liveEditorResponseEntries)
		assert.NoError(t, err)

		response := embeddedconnector.ConnectorPayload{
			Messages: embeddedconnector.ConnectorMessage{
				FevalResponse: []embeddedconnector.FevalResponseMessage{
					{
						IsError: false,
						Results: []interface{}{
							string(data),
						},
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
	response, err := client.EvalWithCapture(ctx, mockLogger, evalRequest)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Warning: first warning message\nx = 1\nWarning: second warning", response.ConsoleOutput)
	assert.Nil(t, response.Images)
}
