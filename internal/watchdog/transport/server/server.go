// Copyright 2025-2026 The MathWorks, Inc.

package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	httpserver "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/server"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
)

type Server struct {
	httpServer httpserver.HttpServer
	logger     entities.Logger
}

func newServer(
	httpServerFactory HTTPServerFactory,
	logger entities.Logger,
	handler handler.Handler,
) (*Server, error) {
	handlers := map[string]http.HandlerFunc{
		"POST " + messages.ProcessToKillPath: processToKillHandler(logger, handler),
		"POST " + messages.ShutdownPath:      shutdownHandler(logger, handler),
	}

	httpServer, err := httpServerFactory.NewServerOverUDS(handlers)
	if err != nil {
		return nil, err
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (s *Server) Start(socketPath string) error {
	s.logger.
		With("socketPath", socketPath).
		Info("Server started")

	return s.httpServer.Serve(socketPath)
}

func (s *Server) Stop() error {
	s.logger.Info("Server stopping")
	defer s.logger.Info("Server stopped")

	return s.httpServer.Shutdown(context.Background())
}

func processToKillHandler(logger entities.Logger, handler handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleMessage(w, r, logger, handler.HandleProcessToKill)
	}
}

func shutdownHandler(logger entities.Logger, handler handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleMessage(w, r, logger, handler.HandleShutdown)
	}
}

func handleMessage[RequestType any, ResponseType any](
	w http.ResponseWriter,
	r *http.Request,
	logger entities.Logger,
	messageHandler func(RequestType) (ResponseType, error),
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Error("Failed to read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.Body.Close()
	if err != nil {
		logger.WithError(err).Error("Failed to close request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req RequestType
	if err := json.Unmarshal(body, &req); err != nil {
		logger.WithError(err).Error("Failed to decode request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := messageHandler(req)
	if err != nil {
		logger.WithError(err).Error("Failed to handle request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.WithError(err).Error("Failed to encode response")
	}
}
