// Copyright 2026 The MathWorks, Inc.

package server

import (
	"context"
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
)

type Definition[Dependencies any] struct {
	Name         string
	Title        string
	Instructions string
}

type Server[Dependencies any] struct {
	serverDefinition Definition[Dependencies]

	applicationFactory adaptor.ApplicationFactory
	errorWriter        entities.Writer
}

func New[Dependencies any](thisDefinition Definition[Dependencies]) *Server[Dependencies] {
	return &Server[Dependencies]{
		applicationFactory: adaptor.NewFactory(),

		serverDefinition: thisDefinition,
		errorWriter:      os.Stderr,
	}
}

func (s *Server[Dependencies]) StartAndWaitForCompletion(ctx context.Context) int {
	serverDefinition := definition.New(
		s.serverDefinition.Name,
		s.serverDefinition.Title,
		s.serverDefinition.Instructions,
	)
	application := s.applicationFactory.New(serverDefinition)

	if err := application.ModeSelector().StartAndWaitForCompletion(ctx); err != nil {
		errorMessage, ok := application.MessageCatalog().GetFromGeneralError(err)
		if ok {
			fmt.Fprintf(s.errorWriter, "%s\n", errorMessage) //nolint:errcheck // Nothing we can do then
			return 1
		}

		fallbackMessage := application.MessageCatalog().Get(messages.StartupErrors_GenericInitializeFailure)
		fmt.Fprintf(s.errorWriter, "%s\n", fallbackMessage) //nolint:errcheck // Nothing we can do then
		return 1
	}

	return 0
}
