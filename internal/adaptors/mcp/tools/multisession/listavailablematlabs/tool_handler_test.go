// Copyright 2025-2026 The MathWorks, Inc.

package listavailablematlabs_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/multisession/listavailablematlabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	// Act
	tool := listavailablematlabs.New(mockLoggerFactory, mockUsecase)

	// Assert
	assert.NotNil(t, tool)
}

func TestTool_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	mockEnvironments := []entities.EnvironmentInfo{
		{
			MATLABRoot: "/path/to/matlab/R2023a",
			Version:    "R2023a",
		},
		{
			MATLABRoot: "/path/to/matlab/R2022b",
			Version:    "R2022b",
		},
	}
	ctx := t.Context()
	inputs := listavailablematlabs.Args{}

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg()).
		Return(mockEnvironments).
		Once()

	// Act
	result, err := listavailablematlabs.Handler(mockUsecase)(ctx, mockLogger, inputs)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.AvailableMATLABs, 2, "Should return 2 environments")

	assert.Equal(t, mockEnvironments[0].MATLABRoot, result.AvailableMATLABs[0].MATLABRoot, "First environment should have correct MATLAB root")
	assert.Equal(t, mockEnvironments[0].Version, result.AvailableMATLABs[0].Version, "First environment should have correct version")

	assert.Equal(t, mockEnvironments[1].MATLABRoot, result.AvailableMATLABs[1].MATLABRoot, "Second environment should have correct MATLAB root")
	assert.Equal(t, mockEnvironments[1].Version, result.AvailableMATLABs[1].Version, "Second environment should have correct version")
	assert.Len(t, mockLogger.InfoLogs(), 2, "Bounding info logs should be creates")
}

func TestTool_Handler_EmptyList(t *testing.T) {
	// Arrange
	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	mockEnvironments := []entities.EnvironmentInfo{}

	mockUsecase.EXPECT().
		Execute(mock.Anything, mockLogger.AsMockArg()).
		Return(mockEnvironments).
		Once()

	ctx := t.Context()
	inputs := listavailablematlabs.Args{}

	// Act
	result, err := listavailablematlabs.Handler(mockUsecase)(ctx, mockLogger, inputs)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Empty(t, result.AvailableMATLABs, "Result should be an empty slice")

	assert.Len(t, mockLogger.InfoLogs(), 2, "Bounding info logs should be creates")
}

func TestListAvailableMATLABs_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewReadOnlyAnnotations()

	// Act
	tool := listavailablematlabs.New(mockLoggerFactory, mockUsecase)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have read-only annotations")
}
