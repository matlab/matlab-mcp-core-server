// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"errors"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	internaltoolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
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
	expectedDefinition := definition.New("test-server", "Test Server", "Test instructions", nil)

	mockApplicationFactory.EXPECT().
		New(matchDefinition(expectedDefinition)).
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
		Name:         expectedDefinition.Name(),
		Title:        expectedDefinition.Title(),
		Instructions: expectedDefinition.Instructions(),
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

	expectedDefinition := definition.New("test-server", "Test Server", "Test instructions", nil)

	mockApplicationFactory.EXPECT().
		New(matchDefinition(expectedDefinition)).
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
		Name:         expectedDefinition.Name(),
		Title:        expectedDefinition.Title(),
		Instructions: expectedDefinition.Instructions(),
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
	expectedDefinition := definition.New("test-server", "Test Server", "Test instructions", nil)

	mockApplicationFactory.EXPECT().
		New(matchDefinition(expectedDefinition)).
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
		Name:         expectedDefinition.Name(),
		Title:        expectedDefinition.Title(),
		Instructions: expectedDefinition.Instructions(),
	}
	s := server.New(serverDefinition)
	s.SetApplicationFactory(mockApplicationFactory)
	s.SetErrorWriter(mockErrorWriter)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 1, exitCode)
}

func TestServer_StartAndWaitForCompletion_NilToolsProvider(t *testing.T) {
	// Arrange
	mockApplicationFactory := &adaptormocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockLoggerFactory := &adaptormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	ctx := t.Context()
	var capturedDefinition definition.Definition

	mockApplicationFactory.EXPECT().
		New(mock.MatchedBy(func(d definition.Definition) bool {
			capturedDefinition = d
			return true
		})).
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
		Name:          "test-server",
		Title:         "Test Server",
		Instructions:  "Test instructions",
		ToolsProvider: nil,
	}
	s := server.New(serverDefinition)
	s.SetApplicationFactory(mockApplicationFactory)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 0, exitCode)
	tools := capturedDefinition.Tools(mockLoggerFactory)
	require.Nil(t, tools)
}

func TestServer_StartAndWaitForCompletion_WithToolsProvider(t *testing.T) {
	// Arrange
	mockApplicationFactory := &adaptormocks.MockApplicationFactory{}
	defer mockApplicationFactory.AssertExpectations(t)

	mockApplication := &adaptormocks.MockApplication{}
	defer mockApplication.AssertExpectations(t)

	mockModeSelector := &adaptormocks.MockModeSelector{}
	defer mockModeSelector.AssertExpectations(t)

	mockLoggerFactory := &adaptormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockTool := &server.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockInternalTool := &internaltoolsmocks.MockTool{}
	defer mockInternalTool.AssertExpectations(t)

	ctx := t.Context()
	var capturedDefinition definition.Definition
	toolsProviderCalled := false

	mockApplicationFactory.EXPECT().
		New(mock.MatchedBy(func(d definition.Definition) bool {
			capturedDefinition = d
			return true
		})).
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

	mockTool.On("toInternal", mockLoggerFactory).
		Return(mockInternalTool).
		Once()

	serverDefinition := server.Definition[struct{}]{
		Name:         "test-server",
		Title:        "Test Server",
		Instructions: "Test instructions",
		ToolsProvider: func(resources server.ToolProviderResources[struct{}]) []server.Tool {
			toolsProviderCalled = true
			return []server.Tool{mockTool}
		},
	}
	s := server.New(serverDefinition)
	s.SetApplicationFactory(mockApplicationFactory)

	// Act
	exitCode := s.StartAndWaitForCompletion(ctx)

	// Assert
	require.Equal(t, 0, exitCode)
	tools := capturedDefinition.Tools(mockLoggerFactory)
	require.True(t, toolsProviderCalled)
	require.Len(t, tools, 1)
	require.Equal(t, mockInternalTool, tools[0])
}

func matchDefinition(expected definition.Definition) any {
	return mock.MatchedBy(func(d definition.Definition) bool {
		return d.Name() == expected.Name() &&
			d.Title() == expected.Title() &&
			d.Instructions() == expected.Instructions()
	})
}
