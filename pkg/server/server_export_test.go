// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
)

func (s *Server[Dependencies]) SetApplicationFactory(factory adaptor.ApplicationFactory) {
	s.applicationFactory = factory
}

func (s *Server[Dependencies]) SetErrorWriter(writer entities.Writer) {
	s.errorWriter = writer
}
