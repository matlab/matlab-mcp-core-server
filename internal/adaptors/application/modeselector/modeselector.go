// Copyright 2025-2026 The MathWorks, Inc.

package modeselector

import (
	"context"
	"fmt"
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Parser interface {
	Usage() string
}

type WatchdogProcess interface { //nolint:iface // Intentional interface for deps injection
	StartAndWaitForCompletion(ctx context.Context) error
}

type Orchestrator interface { //nolint:iface // Intentional interface for deps injection
	StartAndWaitForCompletion(ctx context.Context) error
}

type OSLayer interface {
	Stdout() io.Writer
}

type ModeSelector struct {
	configFactory   ConfigFactory
	watchdogProcess WatchdogProcess
	orchestrator    Orchestrator
	osLayer         OSLayer
	parser          Parser
}

func New(
	configFactory ConfigFactory,
	parser Parser,
	watchdogProcess WatchdogProcess,
	orchestrator Orchestrator,
	osLayer OSLayer,
) *ModeSelector {
	return &ModeSelector{
		configFactory:   configFactory,
		parser:          parser,
		watchdogProcess: watchdogProcess,
		orchestrator:    orchestrator,
		osLayer:         osLayer,
	}
}

func (m *ModeSelector) StartAndWaitForCompletion(ctx context.Context) error {
	config, err := m.configFactory.Config()
	if err != nil {
		return err
	}

	switch {
	case config.HelpMode():
		_, err := fmt.Fprintf(m.osLayer.Stdout(), "%s\n", m.parser.Usage())
		return err
	case config.VersionMode():
		_, err := fmt.Fprintf(m.osLayer.Stdout(), "%s\n", config.Version())
		return err
	case config.WatchdogMode():
		return m.watchdogProcess.StartAndWaitForCompletion(ctx)
	default:
		return m.orchestrator.StartAndWaitForCompletion(ctx)
	}
}
