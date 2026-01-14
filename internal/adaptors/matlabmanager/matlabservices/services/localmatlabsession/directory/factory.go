// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"fmt"
	"os"
	"path/filepath"

	applicationdirectory "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ApplicationDirectoryFactory interface {
	Directory() (applicationdirectory.Directory, messages.Error)
}

type OSLayer interface {
	Mkdir(name string, perm os.FileMode) error
	RemoveAll(path string) error

	Stat(name string) (osfacade.FileInfo, error)
	ReadFile(filePath string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type MATLABFiles interface {
	GetAll() map[string][]byte
}

type Directory interface {
	Path() string
	CertificateFile() string
	CertificateKeyFile() string
	GetEmbeddedConnectorDetails() (string, []byte, error)
	Cleanup() error
}

type Factory struct {
	osLayer                     OSLayer
	applicationDirectoryFactory ApplicationDirectoryFactory
	matlabFiles                 MATLABFiles
}

func NewFactory(
	osLayer OSLayer,
	applicationDirectoryFactory ApplicationDirectoryFactory,
	matlabFiles MATLABFiles,
) *Factory {
	return &Factory{
		osLayer:                     osLayer,
		applicationDirectoryFactory: applicationDirectoryFactory,
		matlabFiles:                 matlabFiles,
	}
}

func (f *Factory) New(logger entities.Logger) (Directory, error) {
	applicationDirectory, messagesErr := f.applicationDirectoryFactory.Directory()
	if messagesErr != nil {
		return nil, messagesErr
	}

	sessionDir, messagesErr := applicationDirectory.CreateSubDir("matlab-session-")
	if messagesErr != nil {
		return nil, messagesErr
	}

	matlabMCPPackagePath := filepath.Join(sessionDir, "+matlab_mcp")

	err := f.osLayer.Mkdir(matlabMCPPackagePath, 0o700)
	if err != nil {
		return nil, fmt.Errorf("failed to create package directory: %w", err)
	}

	for fileName, fileContent := range f.matlabFiles.GetAll() {
		filePath := filepath.Join(matlabMCPPackagePath, fileName)
		if err := f.osLayer.WriteFile(filePath, fileContent, 0o600); err != nil {
			return nil, fmt.Errorf("failed to create %s file: %w", fileName, err)
		}
	}

	return newDirectory(sessionDir, f.osLayer), nil
}
