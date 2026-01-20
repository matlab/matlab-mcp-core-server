// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	httpclientmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Ping_Retries(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	okResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"messages":{"pingResponse":[{}]}}`)),
	}

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(nil, assert.AnError).
		Once()

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(okResponse, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)
	client.SetPingRetry(10 * time.Millisecond)
	client.SetPingTimeout(100 * time.Millisecond)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.True(t, response.IsAlive)
}

func TestClient_Ping_ContextPropagation(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "uniqueKey"
	const contextKeyValue = "uniqueValue"

	expectedContext := context.WithValue(t.Context(), contextKey, contextKeyValue)
	expectedResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"messages":{"pingResponse":[{}]}}`)),
	}

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(request *http.Request) bool {
			return request.Context().Value(contextKey) == contextKeyValue
		})).
		Return(expectedResponse, nil).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)
	client.SetPingRetry(10 * time.Millisecond)
	client.SetPingTimeout(100 * time.Millisecond)

	// Act
	response := client.Ping(expectedContext, mockLogger)

	// Assert
	assert.True(t, response.IsAlive)
}

func TestClient_Ping_Timeout(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(nil, assert.AnError)

	pingTimeout := 100 * time.Millisecond
	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)
	client.SetPingRetry(10 * time.Millisecond)
	client.SetPingTimeout(pingTimeout)

	ctx := t.Context()

	// Act
	start := time.Now()
	response := client.Ping(ctx, mockLogger)
	duration := time.Since(start)

	// Assert
	assert.False(t, response.IsAlive, "Should return not alive after timeout")
	assert.GreaterOrEqual(t, duration, pingTimeout, "Should have waited for at least the timeout duration")
}
