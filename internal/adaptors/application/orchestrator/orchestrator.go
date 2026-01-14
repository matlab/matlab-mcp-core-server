// Copyright 2025-2026 The MathWorks, Inc.

package orchestrator

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type LifecycleSignaler interface {
	RequestShutdown()
	WaitForShutdownToComplete() error
}

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Server interface {
	Run() error
}

type WatchdogClient interface {
	Start() error
	Stop() error
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type GlobalMATLAB interface {
	Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error)
}

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

// Orchestrator
type Orchestrator struct {
	lifecycleSignaler LifecycleSignaler
	configFactory     ConfigFactory
	server            Server
	watchdogClient    WatchdogClient
	loggerFactory     LoggerFactory
	osSignaler        OSSignaler
	globalMATLAB      GlobalMATLAB
	directoryFactory  DirectoryFactory
}

func New(
	lifecycleSignaler LifecycleSignaler,
	configFactory ConfigFactory,
	server Server,
	watchdogClient WatchdogClient,
	loggerFactory LoggerFactory,
	osSignaler OSSignaler,
	globalMATLAB GlobalMATLAB,
	directoryFactory DirectoryFactory,
) *Orchestrator {
	orchestrator := &Orchestrator{
		lifecycleSignaler: lifecycleSignaler,
		configFactory:     configFactory,
		server:            server,
		watchdogClient:    watchdogClient,
		loggerFactory:     loggerFactory,
		osSignaler:        osSignaler,
		globalMATLAB:      globalMATLAB,
		directoryFactory:  directoryFactory,
	}
	return orchestrator
}

func (o *Orchestrator) StartAndWaitForCompletion(ctx context.Context) error {
	config, messagesErr := o.configFactory.Config()
	if messagesErr != nil {
		return messagesErr
	}

	logger, messagesErr := o.loggerFactory.GetGlobalLogger()
	if messagesErr != nil {
		return messagesErr
	}

	defer func() {
		logger.Info("Initiating MATLAB MCP Core Server application shutdown")
		o.lifecycleSignaler.RequestShutdown()

		err := o.lifecycleSignaler.WaitForShutdownToComplete()
		if err != nil {
			logger.WithError(err).Warn("MATLAB MCP Core Server application shutdown failed")
		}

		logger.Debug("Shutdown functions have all completed, stopping the watchdog")
		err = o.watchdogClient.Stop()
		if err != nil {
			logger.WithError(err).Warn("Watchdog shutdown failed")
		}

		logger.Info("MATLAB MCP Core Server application shutdown complete")
	}()

	logger.Info("Initiating MATLAB MCP Core Server application startup")
	config.RecordToLogger(logger)

	directory, messagesErr := o.directoryFactory.Directory()
	if messagesErr != nil {
		return messagesErr
	}
	directory.RecordToLogger(logger)

	err := o.watchdogClient.Start()
	if err != nil {
		return err
	}

	serverErrC := make(chan error, 1)
	go func() {
		serverErrC <- o.server.Run()
	}()

	if config.UseSingleMATLABSession() && config.InitializeMATLABOnStartup() {
		_, err := o.globalMATLAB.Client(ctx, logger)
		if err != nil {
			logger.WithError(err).Warn("MATLAB global initialization failed")
		}
	}

	logger.Info("MATLAB MCP Core Server application startup complete")

	select {
	case <-o.osSignaler.InterruptSignalChan():
		logger.Info("Received termination signal")
		return nil
	case err := <-serverErrC:
		return err
	}
}
