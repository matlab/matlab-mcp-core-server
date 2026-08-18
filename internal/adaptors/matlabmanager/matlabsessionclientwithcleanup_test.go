// Copyright 2025-2026 The MathWorks, Inc.

package matlabmanager_test

import (
	"context"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/testutils"
	entitiesmocks "github.com/matlab/matlab-mcp-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanupSessionClient_New_HappyPath(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}

	// Act
	result := matlabmanager.NewCleanupSessionClient(mockClient, func(context.Context, entities.Logger) error { return nil })

	// Assert
	require.NotNil(t, result)
}

func TestCleanupSessionClient_StopSession_InvokesStop(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()

	stopCalled := false
	var gotCtx context.Context
	stop := func(ctx context.Context, _ entities.Logger) error {
		stopCalled = true
		gotCtx = ctx
		return nil
	}

	client := matlabmanager.NewCleanupSessionClient(mockClient, stop)

	// Act
	err := client.StopSession(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.True(t, stopCalled, "StopSession should invoke the injected stop closure")
	assert.Equal(t, ctx, gotCtx, "StopSession should forward its context to the stop closure")
}

func TestCleanupSessionClient_StopSession_PropagatesError(t *testing.T) {
	// Arrange
	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError

	client := matlabmanager.NewCleanupSessionClient(mockClient, func(context.Context, entities.Logger) error {
		return expectedError
	})

	// Act
	err := client.StopSession(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
}
