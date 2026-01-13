// Copyright 2025-2026 The MathWorks, Inc.

package handler

import (
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type ProcessHandler interface {
	KillProcess(processPid int) error
}

type Factory struct {
	loggerFactory  LoggerFactory
	processHandler ProcessHandler

	initOnce        sync.Once
	initError       error
	handlerInstance *handler
}

func NewFactory(
	loggerFactory LoggerFactory,
	processHandler ProcessHandler,
) *Factory {
	return &Factory{
		loggerFactory:  loggerFactory,
		processHandler: processHandler,
	}
}

func (f *Factory) Handler() (Handler, error) {
	f.initOnce.Do(func() {
		logger, err := f.loggerFactory.GetGlobalLogger()
		if err != nil {
			f.initError = err
			return
		}

		f.handlerInstance = newHandler(logger, f.processHandler)
	})

	if f.initError != nil {
		return nil, f.initError
	}

	return f.handlerInstance, nil
}
