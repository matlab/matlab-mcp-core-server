// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"os"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type FilenameFactory interface {
	CreateFileWithUniqueSuffix(baseName string, ext string) (string, string, error)
}

type OSLayer interface {
	MkdirTemp(dir string, pattern string) (string, error)
	MkdirAll(name string, perm os.FileMode) error
	Create(name string) (osfacade.File, error)
}

type Directory interface {
	BaseDir() string
	ID() string
	CreateSubDir(pattern string) (string, messages.Error)
	RecordToLogger(logger entities.Logger)
}

type Factory struct {
	configFactory   ConfigFactory
	filenameFactory FilenameFactory
	osFacade        OSLayer

	initOnce          sync.Once
	initError         messages.Error
	directoryInstance Directory
}

func NewFactory(
	configFactory ConfigFactory,
	filenameFactory FilenameFactory,
	osFacade OSLayer,
) *Factory {
	return &Factory{
		configFactory:   configFactory,
		filenameFactory: filenameFactory,
		osFacade:        osFacade,
	}
}

func (f *Factory) Directory() (Directory, messages.Error) {
	f.initOnce.Do(func() {
		config, err := f.configFactory.Config()
		if err != nil {
			f.initError = err
			return
		}

		directoryInstance, err := newDirectory(
			config,
			f.filenameFactory,
			f.osFacade,
		)
		if err != nil {
			f.initError = err
			return
		}

		f.directoryInstance = directoryInstance
	})

	if f.initError != nil {
		return nil, f.initError
	}

	return f.directoryInstance, nil
}
