// Copyright 2025-2026 The MathWorks, Inc.

package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultGlobalLogLevel slog.Level = slog.LevelDebug

const (
	logFileName         = "server"
	watchdogLogFileName = "watchdog"
	logFileExt          = ".log"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Directory interface {
	BaseDir() string
	ID() string
}

type FilenameFactory interface {
	FilenameWithSuffix(fileName string, ext string, suffix string) string
}

type OSLayer interface {
	Create(name string) (osfacade.File, error)
}

type Factory struct {
	configFactory   ConfigFactory
	directory       Directory
	filenameFactory FilenameFactory
	osLayer         OSLayer

	initOnce       sync.Once
	initError      messages.Error
	parsedLogLevel slog.Level
	logFile        entities.Writer

	globalLoggerOnce sync.Once
	globalLogger     *slogLogger
}

func NewFactory(
	configFactory ConfigFactory,
	directory Directory,
	filenameFactory FilenameFactory,
	osLayer OSLayer,
) *Factory {
	return &Factory{
		configFactory:   configFactory,
		directory:       directory,
		filenameFactory: filenameFactory,
		osLayer:         osLayer,

		parsedLogLevel: defaultGlobalLogLevel,
	}
}

func (f *Factory) NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error) {
	// In MCP Server development, special care should be given to logging:
	//
	// https://modelcontextprotocol.io/quickstart/server#logging-in-mcp-servers
	//
	// In essence, you can't log to standard `stdout`, and while you may log to `stderr`, you should log to the client:
	//
	// https://modelcontextprotocol.io/specification/2025-06-18/server/utilities/logging
	if err := f.initialize(); err != nil {
		return nil, err
	}

	sessionHandler := mcp.NewLoggingHandler(session, &mcp.LoggingHandlerOptions{})

	handler := slog.NewJSONHandler(f.logFile, &slog.HandlerOptions{
		Level: f.parsedLogLevel,
	})

	return &slogLogger{
		logger: slog.New(NewMultiHandler(sessionHandler, handler)),
	}, nil
}

func (f *Factory) GetGlobalLogger() (entities.Logger, messages.Error) {
	// There are cases where we want to log, but wo don't have an MCP session yet.
	// In those cases, we must log to stderr, to not affect the stdio transport:
	//
	// https://modelcontextprotocol.io/docs/develop/build-server#best-practices
	if err := f.initialize(); err != nil {
		return nil, err
	}

	f.globalLoggerOnce.Do(func() {
		multiWriter := io.MultiWriter(os.Stderr, f.logFile)

		handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: f.parsedLogLevel,
		})
		f.globalLogger = &slogLogger{
			logger: slog.New(handler),
		}
	})

	return f.globalLogger, nil
}

func (f *Factory) initialize() messages.Error {
	f.initOnce.Do(func() {
		config, configErr := f.configFactory.Config()
		if configErr != nil {
			f.initError = configErr
			return
		}

		var parsedLogLevel slog.Level
		logLevel := config.LogLevel()
		switch logLevel {
		case entities.LogLevelDebug:
			parsedLogLevel = slog.LevelDebug
		case entities.LogLevelInfo:
			parsedLogLevel = slog.LevelInfo
		case entities.LogLevelWarn:
			parsedLogLevel = slog.LevelWarn
		case entities.LogLevelError:
			parsedLogLevel = slog.LevelError
		default:
			f.initError = messages.New_StartupErrors_InvalidLogLevel_Error(string(logLevel))
			return
		}
		f.parsedLogLevel = parsedLogLevel

		baseDir := f.directory.BaseDir()
		id := f.directory.ID()

		logFileBase := filepath.Join(baseDir, logFileName)
		if config.WatchdogMode() {
			logFileBase = filepath.Join(baseDir, watchdogLogFileName)
		}

		logFilePath := f.filenameFactory.FilenameWithSuffix(logFileBase, logFileExt, id)

		logFile, err := f.osLayer.Create(logFilePath)
		if err != nil {
			f.initError = messages.New_StartupErrors_FailedToCreateLogFile_Error(logFilePath)
			return
		}

		f.logFile = logFile
	})

	return f.initError
}
