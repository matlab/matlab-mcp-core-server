// Copyright 2025-2026 The MathWorks, Inc.

package handler

import (
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/messages"
)

type Handler interface {
	HandleProcessToKill(req messages.ProcessToKillRequest) (messages.ProcessToKillResponse, error)
	HandleShutdown(req messages.ShutdownRequest) (messages.ShutdownResponse, error)
	RegisterShutdownFunction(fn func())
	TerminateAllProcesses()
}

type handler struct {
	logger         entities.Logger
	processHandler ProcessHandler

	lock              *sync.Mutex
	processPIDsToKill map[int]struct{}
	shutdownFuncs     []func()
}

func newHandler(
	logger entities.Logger,
	processHandler ProcessHandler,
) *handler {
	return &handler{
		logger:         logger,
		processHandler: processHandler,

		lock:              &sync.Mutex{},
		processPIDsToKill: make(map[int]struct{}),
		shutdownFuncs:     make([]func(), 0),
	}
}

func (h *handler) HandleProcessToKill(req messages.ProcessToKillRequest) (messages.ProcessToKillResponse, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.logger.
		With("pid", req.PID).
		Info("Adding process to kill")
	h.processPIDsToKill[req.PID] = struct{}{}

	return messages.ProcessToKillResponse{}, nil
}

func (h *handler) RegisterShutdownFunction(fn func()) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.shutdownFuncs = append(h.shutdownFuncs, fn)
}

func (h *handler) HandleShutdown(_ messages.ShutdownRequest) (messages.ShutdownResponse, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, fn := range h.shutdownFuncs {
		fn()
	}
	return messages.ShutdownResponse{}, nil
}

func (h *handler) TerminateAllProcesses() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.logger.
		With("count", len(h.processPIDsToKill)).
		Info("Trying to terminate children")

	for pid := range h.processPIDsToKill {
		h.logger.
			With("pid", pid).
			Debug("Killing process")

		if err := h.processHandler.KillProcess(pid); err != nil {
			h.logger.
				WithError(err).
				With("pid", pid).
				Error("Failed to kill child")
		}
	}
}
