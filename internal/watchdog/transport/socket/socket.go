// Copyright 2025-2026 The MathWorks, Inc.

package socket

import (
	"errors"
	"path/filepath"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

var ErrSocketPathTooLong = errors.New("socket path is too long")

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

type OSLayer interface {
	GOOS() string
}

type Socket interface {
	Path() string
}

type Factory struct {
	directoryFactory DirectoryFactory
	osLayer          OSLayer

	socketInstance Socket
	socketError    error
}

func NewFactory(
	directoryFactory DirectoryFactory,
	osLayer OSLayer,
) *Factory {
	return &Factory{
		directoryFactory: directoryFactory,
		osLayer:          osLayer,
	}
}

func (f *Factory) Socket() (Socket, error) {
	if f.socketError != nil {
		return nil, f.socketError
	}

	if f.socketInstance == nil {
		directory, messagesErr := f.directoryFactory.Directory()
		if messagesErr != nil {
			f.socketError = messagesErr
			return nil, messagesErr
		}

		socket, err := newSocket(
			directory,
			f.osLayer,
		)
		if err != nil {
			f.socketError = err
			return nil, err
		}

		f.socketInstance = socket
	}

	return f.socketInstance, nil
}

type Directory interface {
	BaseDir() string
	ID() string
}

type socket struct {
	path string
}

func newSocket(
	directory Directory,
	osLayer OSLayer,
) (*socket, error) {
	socketPath := filepath.Join(directory.BaseDir(), "watchdog-"+directory.ID()+".sock")

	// Socket path max length is 108 characters, but for safety using 105
	if len(socketPath) > 105 {
		if osLayer.GOOS() == "darwin" {
			socketPath = filepath.Join("/tmp", "watchdog-"+directory.ID()+".sock")
		} else {
			return nil, ErrSocketPathTooLong
		}
	}

	return &socket{
		path: socketPath,
	}, nil
}

func (s *socket) Path() string {
	return s.path
}
