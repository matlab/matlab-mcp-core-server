// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"errors"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	adaptormocks "github.com/matlab/matlab-mcp-core-server/mocks/wire/adaptor"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	serverDefinition := server.Definition[struct{}]{
		Name:         "test-server",
		Title:        "Test Server",
		Instructions: "Test instructions",
	}

	// Act
	s := server.New(serverDefinition)

	// Assert
	require.NotNil(t, s)
}

func TestServer_StartAndWaitForCompletion_HappyPath(t *testing.T) {
	// Arrange
	mockApplicationFactory := &adaptormocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedName := "test-server"
	expectedTitle := "Test Server"
	expectedInstructions := "Test instructions"
	expectedDefinition := definition.New(expectedName, expectedTitle, expectedInstructions)

	mockApplicationFactory.EXPECT().
		New(expectedDefinition).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(nil).
		Once()

	serverDefinition := server.Definition[struct{}]{
		Name:         expectedName,
		Title:        expectedTitle,
		Instructions: expectedInstructions,
	}
	s := server.New(serverDefinition)
	s.SetApplicationFactory(mockApplicationFactory)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 0, exitCode)
}

func TestServer_StartAndWaitForCompletion_KnownError(t *testing.T) {
	// Arrange
	mockApplicationFactory := &adaptormocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockMessageCatalog := &adaptormocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockErrorWriter := &entitiesmocks.MockWriter{}
	defer mockErrorWriter.AssertExpectations(t)

	ctx := t.Context()
	expectedError := errors.New("known error")
	expectedErrorMessage := "A known error occurred"

	mockApplicationFactory.EXPECT().
		New(mock.Anything).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(expectedError).
		Once()

	mockApplication.EXPECT().
		MessageCatalog().
		Return(mockMessageCatalog).
		Once()

	mockMessageCatalog.EXPECT().
		GetFromGeneralError(expectedError).
		Return(expectedErrorMessage, true).
		Once()

	mockErrorWriter.EXPECT().
		Write([]byte(expectedErrorMessage+"\n")).
		Return(len(expectedErrorMessage)+1, nil).
		Once()

	serverDefinition := server.Definition[struct{}]{
		Name:         "test-server",
		Title:        "Test Server",
		Instructions: "Test instructions",
	}
	s := server.New(serverDefinition)
	s.SetApplicationFactory(mockApplicationFactory)
	s.SetErrorWriter(mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 1, exitCode)
}

func TestServer_StartAndWaitForCompletion_UnknownError(t *testing.T) {
	// Arrange
	mockApplicationFactory := &adaptormocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockMessageCatalog := &adaptormocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	mockErrorWriter := &entitiesmocks.MockWriter{}
	defer mockErrorWriter.AssertExpectations(t)

	ctx := t.Context()
	expectedError := errors.New("unknown error")
	expectedFallbackMessage := "Some generic failure message."
	expectedWrittenOutput := expectedFallbackMessage + "\n"

	mockApplicationFactory.EXPECT().
		New(mock.Anything).
		Return(mockApplication).
		Once()

	mockApplication.EXPECT().
		ModeSelector().
		Return(mockModeSelector).
		Once()

	mockModeSelector.EXPECT().
		StartAndWaitForCompletion(ctx).
		Return(expectedError).
		Once()

	mockApplication.EXPECT().
		MessageCatalog().
		Return(mockMessageCatalog).
		Times(2)

	mockMessageCatalog.EXPECT().
		GetFromGeneralError(expectedError).
		Return("", false).
		Once()

	mockMessageCatalog.EXPECT().
		Get(messages.StartupErrors_GenericInitializeFailure).
		Return(expectedFallbackMessage).
		Once()

	mockErrorWriter.EXPECT().
		Write([]byte(expectedWrittenOutput)).
		Return(len(expectedWrittenOutput), nil).
		Once()

	serverDefinition := server.Definition[struct{}]{
		Name:         "test-server",
		Title:        "Test Server",
		Instructions: "Test instructions",
	}
	s := server.New(serverDefinition)
	s.SetApplicationFactory(mockApplicationFactory)
	s.SetErrorWriter(mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 1, exitCode)
}
