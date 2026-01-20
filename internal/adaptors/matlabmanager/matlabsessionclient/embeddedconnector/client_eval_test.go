// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	httpclientmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_Eval_DoErrors(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(nil, assert.AnError).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	ctx := t.Context()
	evalRequest := entities.EvalRequest{
		Code: "ver",
	}

	// Act
	response, err := client.Eval(ctx, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Empty(t, response)
}

func TestClient_Eval_ContextPropagation(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientmocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "uniqueKey"
	const contextKeyValue = "uniqueValue"

	expectedContext := context.WithValue(t.Context(), contextKey, contextKeyValue)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(request *http.Request) bool {
			return request.Context().Value(contextKey) == contextKeyValue
		})).
		Return(nil, assert.AnError).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)

	evalRequest := entities.EvalRequest{
		Code: "ver",
	}

	// Act
	response, err := client.Eval(expectedContext, mockLogger, evalRequest)

	// Assert
	require.Error(t, err)
	assert.Empty(t, response)
}
