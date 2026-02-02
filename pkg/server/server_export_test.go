// Copyright 2026 The MathWorks, Inc.

package server

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
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

func (m *MockTool) toInternal(lf loggerFactory) tools.Tool {
	args := m.Called(lf)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(tools.Tool)
}
