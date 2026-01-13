// Copyright 2025-2026 The MathWorks, Inc.

package watchdog

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type OSLayer interface {
	Getppid() int
}

type ProcessHandler interface {
	WatchProcessAndGetTerminationChan(processPid int) (<-chan struct{}, error)
	KillProcess(processPid int) error
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type ServerHandlerFactory interface {
	Handler() (handler.Handler, error)
}

type ServerFactory interface {
	New() (transport.Server, error)
}

type SocketFactory interface {
	Socket() (socket.Socket, error)
}

type Watchdog struct {
	loggerFactory        LoggerFactory
	osLayer              OSLayer
	processHandler       ProcessHandler
	osSignaler           OSSignaler
	serverHandlerFactory ServerHandlerFactory
	serverFactory        ServerFactory
	socketFactory        SocketFactory

	parentPID         int
	shutdownRequestC  chan struct{}
	shutdownResponseC chan struct{}
}

func New(
	loggerFactory LoggerFactory,
	osLayer OSLayer,
	processHandler ProcessHandler,
	osSignaler OSSignaler,
	serverHandlerFactory ServerHandlerFactory,
	serverFactory ServerFactory,
	socketFactory SocketFactory,
) *Watchdog {
	return &Watchdog{
		loggerFactory:        loggerFactory,
		osLayer:              osLayer,
		processHandler:       processHandler,
		osSignaler:           osSignaler,
		serverHandlerFactory: serverHandlerFactory,
		serverFactory:        serverFactory,
		socketFactory:        socketFactory,

		shutdownRequestC:  make(chan struct{}),
		shutdownResponseC: make(chan struct{}),
	}
}

func (w *Watchdog) StartAndWaitForCompletion(_ context.Context) error {
	logger, messagesErr := w.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	socket, err := w.socketFactory.Socket()
	if err != nil {
		return err
	}

	serverHandler, err := w.serverHandlerFactory.Handler()
	if err != nil {
		return err
	}

	server, err := w.serverFactory.New()
	if err != nil {
		return err
	}

	go func() {
		if err := server.Start(socket.Path()); err != nil {
			logger.WithError(err).Error("Server Start method returned an error")
		}
	}()

	defer func() {
		if err := server.Stop(); err != nil {
			logger.WithError(err).Error("Failed to stop server")
		}
	}()

	logger.Info("Watchdog process has started")
	defer logger.Info("Watchdog process has exited")

	w.parentPID = w.osLayer.Getppid()

	serverHandler.RegisterShutdownFunction(func() {
		close(w.shutdownRequestC)
		// Make sure we broke out of the select, before returning
		<-w.shutdownResponseC
	})

	parentTerminatedC, err := w.processHandler.WatchProcessAndGetTerminationChan(w.parentPID)
	if err != nil {
		return err
	}

	select {
	case <-w.shutdownRequestC:
		logger.Debug("Graceful shutdown signal received")
		// Ackownledge shutdown
		close(w.shutdownResponseC)

	case <-parentTerminatedC:
		logger.Debug("Lost connection to parent, shutting down")

	case <-w.osSignaler.InterruptSignalChan():
		logger.Debug("Received unexpected graceful shutdown OS signal")
	}

	serverHandler.TerminateAllProcesses()

	return nil
}
