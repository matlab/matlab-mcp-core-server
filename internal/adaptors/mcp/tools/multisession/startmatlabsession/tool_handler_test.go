// Copyright 2025-2026 The MathWorks, Inc.

package startmatlabsession_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	startmatlabsessionusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/startmatlabsession"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	basetoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	// Act
	tool := startmatlabsession.New(mockLoggerFactory, mockConfigFactory, mockUsecase)

	// Assert
	assert.NotNil(t, tool)
}

func TestTool_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	matlabRoot := filepath.Join("path", "to", "matlab")
	shouldShowMATLABDesktop := true
	expectedSessionID := entities.SessionID(123)
	expectedVerOutput := "MATLAB Version: R2023a"
	expectedAddOnsOutput := "Installed Add-Ons: Toolbox1, Toolbox2"
	expectedResponse := startmatlabsessionusecase.ReturnArgs{
		SessionID:    expectedSessionID,
		VerOutput:    expectedVerOutput,
		AddOnsOutput: expectedAddOnsOutput,
	}

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             matlabRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}
	args := startmatlabsession.Args{
		MATLABRoot: matlabRoot,
	}

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := startmatlabsession.Handler(mockConfigFactory, mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, int(expectedSessionID), result.SessionID, "Session ID should match")
	assert.Equal(t, expectedVerOutput, result.VerOutput, "Ver output should match")
	assert.Equal(t, expectedAddOnsOutput, result.AddOnsOutput, "AddOns output should match")
}

func TestTool_Handler_UsecaseError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	matlabRoot := filepath.Join("path", "to", "matlab")
	shouldShowMATLABDesktop := true
	expectedError := assert.AnError

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             matlabRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}
	args := startmatlabsession.Args{
		MATLABRoot: matlabRoot,
	}

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(startmatlabsessionusecase.ReturnArgs{}, expectedError).
		Once()

	// Act
	result, err := startmatlabsession.Handler(mockConfigFactory, mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Empty(t, result.ResponseText, "Response text should be empty on error")
}

func TestTool_Handler_ConfigError(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	matlabRoot := filepath.Join("path", "to", "matlab")
	expectedError := messages.New_StartupErrors_BadFlag_Error("flag", "value", "reason")

	args := startmatlabsession.Args{
		MATLABRoot: matlabRoot,
	}

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	// Act
	result, err := startmatlabsession.Handler(mockConfigFactory, mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return config error")
	assert.Empty(t, result.ResponseText, "Response text should be empty on error")
}

func TestTool_Handler_ShowMATLABDesktopFalse(t *testing.T) {
	// Arrange
	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	matlabRoot := filepath.Join("path", "to", "matlab")
	shouldShowMATLABDesktop := false
	expectedSessionID := entities.SessionID(123)
	expectedVerOutput := "MATLAB Version: R2023a"
	expectedAddOnsOutput := "Installed Add-Ons: Toolbox1, Toolbox2"
	expectedResponse := startmatlabsessionusecase.ReturnArgs{
		SessionID:    expectedSessionID,
		VerOutput:    expectedVerOutput,
		AddOnsOutput: expectedAddOnsOutput,
	}

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             matlabRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      shouldShowMATLABDesktop,
	}
	args := startmatlabsession.Args{
		MATLABRoot: matlabRoot,
	}

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockUsecase.EXPECT().
		Execute(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(expectedResponse, nil).
		Once()

	// Act
	result, err := startmatlabsession.Handler(mockConfigFactory, mockUsecase)(ctx, mockLogger, args)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, int(expectedSessionID), result.SessionID, "Session ID should match")
	assert.Equal(t, expectedVerOutput, result.VerOutput, "Ver output should match")
	assert.Equal(t, expectedAddOnsOutput, result.AddOnsOutput, "AddOns output should match")
}

func TestStartMATLABSession_Annotations(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolsmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockUsecase := &mocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	expectedAnnotations := annotations.NewReadOnlyAnnotations()

	// Act
	tool := startmatlabsession.New(mockLoggerFactory, mockConfigFactory, mockUsecase)

	// Assert
	assert.Equal(t, expectedAnnotations, tool.Annotations(), "Tool should have read-only annotations")
}
