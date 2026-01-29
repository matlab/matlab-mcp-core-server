// Copyright 2025-2026 The MathWorks, Inc.

package sdk

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Definition interface {
	Name() string
	Title() string
	Instructions() string
}

type Factory struct {
	configFactory ConfigFactory
	definition    Definition
}

func NewFactory(
	configFactory ConfigFactory,
	definition Definition,
) *Factory {
	return &Factory{
		configFactory: configFactory,
		definition:    definition,
	}
}

func (f *Factory) NewServer() (*mcp.Server, messages.Error) {
	config, err := f.configFactory.Config()
	if err != nil {
		return nil, err
	}

	impl := &mcp.Implementation{
		Name:    f.definition.Name(),
		Title:   f.definition.Title(),
		Version: config.Version(),
	}
	options := &mcp.ServerOptions{
		Instructions: f.definition.Instructions(),
	}

	return mcp.NewServer(impl, options), nil
}
