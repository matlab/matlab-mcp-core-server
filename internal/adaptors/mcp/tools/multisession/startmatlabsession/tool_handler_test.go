// Copyright 2025-2026 The MathWorks, Inc.

package startmatlabsession_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	startmatlabsessionusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/startmatlabsession"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/multisession/startmatlabsession"
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
	tool := startmatlabsession.New(mockLoggerFactory, mockUsecase)

	// Assert
	assert.NotNil(t, tool)
}

func TestTool_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const matlabRoot = "/path/to/matlab"
	const expectedSessionID = entities.SessionID(123)
	const expectedVerOutput = "MATLAB Version: R2023a"
	const expectedAddOnsOutput = "Installed Add-Ons: Toolbox1, Toolbox2"
	expectedResponse := startmatlabsessionusecase.ReturnArgs{
		SessionID:    expectedSessionID,
		VerOutput:    expectedVerOutput,
		AddOnsOutput: expectedAddOnsOutput,
	}

	localSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             matlabRoot,
		IsStartingDirectorySet: false,
	}
	args := startmatlabsession.Args{
		MATLABRoot: matlabRoot,
	}

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), localSessionDetails).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := startmatlabsession.Handler(mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, int(expectedSessionID), result.SessionID, "Session ID should match")
	assert.Equal(t, expectedVerOutput, result.VerOutput, "Ver output should match")
	assert.Equal(t, expectedAddOnsOutput, result.AddOnsOutput, "AddOns output should match")
}

func TestTool_Handler_UsecaseError(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	const matlabRoot = "/path/to/matlab"
	expectedError := assert.AnError

	localSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             matlabRoot,
		IsStartingDirectorySet: false,
	}
	args := startmatlabsession.Args{
		MATLABRoot: matlabRoot,
	}

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), localSessionDetails).
		Return(startmatlabsessionusecase.ReturnArgs{}, expectedError).
		Once()

	// Act
	result, err := startmatlabsession.Handler(mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Empty(t, result.ResponseText, "Response text should be empty on error")
}

func TestStartMATLABSession_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewReadOnlyAnnotations()

	// Act
	tool := startmatlabsession.New(mockLoggerFactory, mockUsecase)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have read-only annotations")
}
