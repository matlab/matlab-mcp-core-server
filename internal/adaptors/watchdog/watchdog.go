// Copyright 2025-2026 The MathWorks, Inc.

package watchdog

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

type WatchdogProcess interface {
	Start() error
}

type ClientFactory interface {
	New() transport.Client
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type SocketFactory interface {
	Socket() (socket.Socket, error)
}

type Watchdog struct {
	loggerFactory LoggerFactory
	logger        entities.Logger

	watchdogProcess WatchdogProcess
	socketFactory   SocketFactory

	client transport.Client

	startedC chan struct{}
}

func New(
	watchdogProcess WatchdogProcess,
	clientFactory ClientFactory,
	loggerFactory LoggerFactory,
	socketFactory SocketFactory,
) *Watchdog {
	return &Watchdog{
		loggerFactory: loggerFactory,

		watchdogProcess: watchdogProcess,
		socketFactory:   socketFactory,

		client: clientFactory.New(),

		startedC: make(chan struct{}),
	}
}

func (w *Watchdog) Start() error {
	logger, messagesErr := w.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	w.logger = logger
	w.logger.Debug("Starting watchdog")

	socket, err := w.socketFactory.Socket()
	if err != nil {
		w.logger.WithError(err).Error("Failed to get socket")
		return err
	}

	err = w.watchdogProcess.Start()
	if err != nil {
		w.logger.WithError(err).Error("Failed to start watchdog process")
		return err
	}

	if err := w.client.Connect(socket.Path()); err != nil {
		w.logger.WithError(err).Error("Failed to connect to watchdog socket")
		return err
	}

	close(w.startedC)

	w.logger.Debug("Started watchdog")

	return nil
}

func (w *Watchdog) RegisterProcessPIDWithWatchdog(processPID int) error {
	<-w.startedC

	w.logger.With("pid", processPID).Debug("Adding child process to watchdog")
	_, err := w.client.SendProcessPID(processPID)
	return err
}

func (w *Watchdog) Stop() error {
	<-w.startedC

	w.logger.Debug("Sending graceful shutdown signal to watchdog")
	_, err := w.client.SendStop()
	return err
}
