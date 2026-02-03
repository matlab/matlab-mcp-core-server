// Copyright 2025-2026 The MathWorks, Inc.

package runmatlabfile_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabfile"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/usecases/runmatlabfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	// Act
	usecase := runmatlabfile.New(mockPathValidator)

	// Assert
	assert.NotNil(t, usecase, "Usecase should not be nil")
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	expectedCdRequest := entities.EvalRequest{
		Code: fmt.Sprintf("cd('%s')", scriptDir),
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(expectedResponse, nil).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_ValidateMATLABScriptError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(scriptPath).
		Return("", expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}

func TestUsecase_Execute_CdEvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	expectedCdRequest := entities.EvalRequest{
		Code: fmt.Sprintf("cd('%s')", scriptDir),
	}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}

func TestUsecase_Execute_RunMATLABFileEvalError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{ScriptPath: scriptPath}

	expectedCdRequest := entities.EvalRequest{
		Code: fmt.Sprintf("cd('%s')", scriptDir),
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}

func TestUsecase_Execute_CaptureOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	captureOutput := true

	usecaseRequest := runmatlabfile.Args{
		ScriptPath:    scriptPath,
		CaptureOutput: captureOutput,
	}

	expectedCdRequest := entities.EvalRequest{
		Code: fmt.Sprintf("cd('%s')", scriptDir),
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	expectedResponse := entities.EvalResponse{
		ConsoleOutput: "Hello, World!",
		Images:        nil,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(expectedResponse, nil).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedResponse, response, "Response should match expected value")
}

func TestUsecase_Execute_CaptureOutput_EvalWithCaptureError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &mocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	ctx := t.Context()
	fileName := "file"
	scriptDir := filepath.Join("some", "path", "to")
	scriptPath := filepath.Join(scriptDir, fileName+".m")
	captureOutput := true
	expectedError := assert.AnError

	usecaseRequest := runmatlabfile.Args{
		ScriptPath:    scriptPath,
		CaptureOutput: captureOutput,
	}

	expectedCdRequest := entities.EvalRequest{
		Code: fmt.Sprintf("cd('%s')", scriptDir),
	}

	expectedEvalRequest := entities.EvalRequest{
		Code: fileName,
	}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(scriptPath).
		Return(scriptPath, nil).
		Once()

	mockClient.EXPECT().
		Eval(ctx, mockLogger.AsMockArg(), expectedCdRequest).
		Return(entities.EvalResponse{}, nil).
		Once()

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := runmatlabfile.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, usecaseRequest)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Empty(t, response, "Response should be empty")
}
