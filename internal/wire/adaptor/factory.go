// Copyright 2026 The MathWorks, Inc.

package adaptor

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

type ApplicationFactory interface {
	New(definition definition.Definition) Application
}

type adaptorFactory struct{}

func NewFactory() ApplicationFactory {
	return &adaptorFactory{}
}

func (f *adaptorFactory) New(definition definition.Definition) Application {
	return &adaptor{
		Application: wire.Initialize(definition),
	}
}
