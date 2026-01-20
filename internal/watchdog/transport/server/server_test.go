// Copyright 2025-2026 The MathWorks, Inc.

package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server"
	httpservermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/http/server"
	servermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server"
	handlermocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport/server/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServer_Start_HappyPath(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	makeHTTPServerServeReturn := make(chan struct{})
	socketPath := filepath.Join(t.TempDir(), "test.sock")
	mockLogger := testutils.NewInspectableLogger()

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Return(mockHTTPServer, nil).
		Once()

	mockHTTPServer.EXPECT().
		Serve(socketPath).
		Run(func(_ string) {
			<-makeHTTPServerServeReturn
		}).
		Return(nil).
		Once()

	serverInstance, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	// Act
	errC := make(chan error)
	go func() {
		errC <- serverInstance.Start(socketPath)
	}()

	select {
	case <-errC:
		t.Fatal("Server stopped unexpectedly")
	case <-time.After(10 * time.Millisecond):
		// Happy path
	}

	close(makeHTTPServerServeReturn)

	// Assert
	require.NoError(t, <-errC)
}

func TestServer_Stop_HappyPath(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Return(mockHTTPServer, nil).
		Once()

	mockHTTPServer.EXPECT().
		Shutdown(mock.Anything).
		Return(nil).
		Once()

	serverInstance, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	// Act
	err = serverInstance.Stop()

	// Assert
	require.NoError(t, err)
}

func TestServer_HandleProcessToKill(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedRequest := messages.ProcessToKillRequest{
		PID: 1234,
	}

	var capturedHandlers map[string]http.HandlerFunc

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Run(func(handlers map[string]http.HandlerFunc) {
			capturedHandlers = handlers
		}).
		Return(mockHTTPServer, nil).
		Once()

	mockHandler.EXPECT().
		HandleProcessToKill(expectedRequest).
		Return(messages.ProcessToKillResponse{}, nil).
		Once()

	_, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	handler, ok := capturedHandlers["POST "+messages.ProcessToKillPath]
	require.True(t, ok, "Handler for POST "+messages.ProcessToKillPath+" should be registered")
	require.NotNil(t, handler)

	// Act
	response := sendJSONPayload(handler, messages.ProcessToKillPath, `{"pid":1234}`)

	// Assert
	require.Equal(t, http.StatusOK, response.Code)
}

func TestServer_HandleProcessToKill_Error(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedRequest := messages.ProcessToKillRequest{
		PID: 1234,
	}

	var capturedHandlers map[string]http.HandlerFunc

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Run(func(handlers map[string]http.HandlerFunc) {
			capturedHandlers = handlers
		}).
		Return(mockHTTPServer, nil).
		Once()

	mockHandler.EXPECT().
		HandleProcessToKill(expectedRequest).
		Return(messages.ProcessToKillResponse{}, assert.AnError).
		Once()

	_, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	handler, ok := capturedHandlers["POST "+messages.ProcessToKillPath]
	require.True(t, ok, "Handler for POST "+messages.ProcessToKillPath+" should be registered")
	require.NotNil(t, handler)

	// Act
	response := sendJSONPayload(handler, messages.ProcessToKillPath, `{"pid":1234}`)

	// Assert
	require.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestServer_HandleShutdown(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedRequest := messages.ShutdownRequest{}

	var capturedHandlers map[string]http.HandlerFunc

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Run(func(handlers map[string]http.HandlerFunc) {
			capturedHandlers = handlers
		}).
		Return(mockHTTPServer, nil).
		Once()

	mockHandler.EXPECT().
		HandleShutdown(expectedRequest).
		Return(messages.ShutdownResponse{}, nil).
		Once()

	_, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	handler := capturedHandlers["POST "+messages.ShutdownPath]
	require.NotNil(t, handler)

	// Act
	response := sendJSONPayload(handler, messages.ShutdownPath, `{}`)

	// Assert
	require.Equal(t, http.StatusOK, response.Code)
}

func TestServer_HandleShutdown_Error(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedRequest := messages.ShutdownRequest{}

	var capturedHandlers map[string]http.HandlerFunc

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Run(func(handlers map[string]http.HandlerFunc) {
			capturedHandlers = handlers
		}).
		Return(mockHTTPServer, nil).
		Once()

	mockHandler.EXPECT().
		HandleShutdown(expectedRequest).
		Return(messages.ShutdownResponse{}, assert.AnError).
		Once()

	_, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	handler := capturedHandlers["POST "+messages.ShutdownPath]
	require.NotNil(t, handler)

	// Act
	response := sendJSONPayload(handler, messages.ShutdownPath, `{}`)

	// Assert
	require.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestServer_AnyHandler_IOReadAllError(t *testing.T) {
	// Arrange
	mockHTTPServerFactory := &servermocks.MockHTTPServerFactory{}
	defer mockHTTPServerFactory.AssertExpectations(t)

	mockHTTPServer := &httpservermocks.MockHttpServer{}
	defer mockHTTPServer.AssertExpectations(t)

	mockHandler := &handlermocks.MockHandler{}
	defer mockHandler.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	var capturedHandlers map[string]http.HandlerFunc

	mockHTTPServerFactory.EXPECT().
		NewServerOverUDS(mock.AnythingOfType("map[string]http.HandlerFunc")).
		Run(func(handlers map[string]http.HandlerFunc) {
			capturedHandlers = handlers
		}).
		Return(mockHTTPServer, nil).
		Once()

	_, err := server.NewServer(
		mockHTTPServerFactory,
		mockLogger,
		mockHandler,
	)
	require.NoError(t, err)

	handler := capturedHandlers["POST "+messages.ShutdownPath]
	require.NotNil(t, handler)

	// Act	req
	reqBody := &IOReaderThatErrorsOnRead{}
	req := httptest.NewRequest(http.MethodPost, messages.ShutdownPath, reqBody)
	response := httptest.NewRecorder()

	handler(response, req)

	// Assert
	require.Equal(t, http.StatusInternalServerError, response.Code)
}

func sendJSONPayload(handler http.HandlerFunc, path string, payload string) *httptest.ResponseRecorder {
	reqBody := bytes.NewBufferString(payload)
	req := httptest.NewRequest(http.MethodPost, path, reqBody)
	rr := httptest.NewRecorder()

	handler(rr, req)

	return rr
}

type IOReaderThatErrorsOnRead struct{}

func (r *IOReaderThatErrorsOnRead) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}
