// Copyright 2025-2026 The MathWorks, Inc.

package stopmatlabsession_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	// Act
	tool := stopmatlabsession.New(mockLoggerFactory, mockUsecase)

	// Assert
	assert.NotNil(t, tool)
}

func TestTool_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const sessionID = 3
	args := stopmatlabsession.Args{
		SessionID: sessionID,
	}

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(nil).
		Once()

	// Act
	result, err := stopmatlabsession.Handler(mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.NotEmpty(t, result.ResponseText, "Response text should not be empty")
}

func TestTool_Handler_UsecaseReturnsError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const sessionID = 3
	expectedError := assert.AnError
	args := stopmatlabsession.Args{
		SessionID: sessionID,
	}

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), entities.SessionID(sessionID)).
		Return(expectedError).
		Once()

	// Act
	result, err := stopmatlabsession.Handler(mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Empty(t, result.ResponseText, "Response text should be empty when there's an error")
}

func TestStopMATLABSession_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewDestructiveAnnotations()

	// Act
	tool := stopmatlabsession.New(mockLoggerFactory, mockUsecase)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have destructive annotations")
}
