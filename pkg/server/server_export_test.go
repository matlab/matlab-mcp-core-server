// Copyright 2026 The MathWorks, Inc.

package server

import (
	internalconfig "github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
	"github.com/stretchr/testify/mock"
)

func (s *Server[Dependencies]) SetApplicationFactory(factory adaptor.ApplicationFactory) {
	s.applicationFactory = factory
}

func (s *Server[Dependencies]) SetErrorWriter(writer entities.Writer) {
	s.errorWriter = writer
}

type MockTool struct {
	mock.Mock
}

func (m *MockTool) toInternal(lf basetool.LoggerFactory, config internalconfig.GenericConfig, messageCatalog definition.MessageCatalog) tools.Tool {
	args := m.Called(lf, config, messageCatalog)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(tools.Tool)
}
