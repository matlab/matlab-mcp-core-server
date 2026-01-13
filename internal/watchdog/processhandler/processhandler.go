// Copyright 2025-2026 The MathWorks, Inc.

package processhandler

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type OSWrapper interface {
	WaitForProcessToComplete(processPid int)
	FindProcess(processPid int) osfacade.Process
}

type ProcessHandler struct {
	loggerFactory LoggerFactory
	osWrapper     OSWrapper
}

func New(
	loggerFactory LoggerFactory,
	osWrapper OSWrapper,
) *ProcessHandler {
	return &ProcessHandler{
		loggerFactory: loggerFactory,
		osWrapper:     osWrapper,
	}
}

func (f *ProcessHandler) WatchProcessAndGetTerminationChan(processPid int) (<-chan struct{}, error) {
	logger, err := f.loggerFactory.GetGlobalLogger()
	if err != nil {
		return nil, err
	}

	logger = logger.With("process-pid", processPid)
	logger.Debug("Watching process and notifying if it terminates")

	parentTerminatedC := make(chan struct{})

	go func() {
		f.osWrapper.WaitForProcessToComplete(processPid)
		logger.Debug("Process terminated")
		close(parentTerminatedC)
	}()

	return parentTerminatedC, nil
}

func (f *ProcessHandler) KillProcess(processPid int) error {
	if process := f.osWrapper.FindProcess(processPid); process != nil {
		return process.Kill()
	}
	return nil
}
