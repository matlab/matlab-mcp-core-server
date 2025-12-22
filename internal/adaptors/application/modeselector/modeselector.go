// Copyright 2025 The MathWorks, Inc.

package modeselector

import (
	"context"
	"fmt"
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type Config interface {
	Version() string
	HelpMode() bool
	VersionMode() bool
	WatchdogMode() bool
}

type Parser interface {
	Usage() string
}

type WatchdogProcessFactory interface { //nolint:iface // Intentional interface for deps injection
	Create() (entities.Mode, error)
}

type OrchestratorFactory interface { //nolint:iface // Intentional interface for deps injection
	Create() (entities.Mode, error)
}

type OSLayer interface {
	Stdout() io.Writer
}

// ModeSelector is the top level object of the MATLAB MCP Core Server.
// It will be imported in `main.go` to start the application, and wait for it's completion.
// It will select which mode to run in based on the configuration, and defer construction of the required objects until the mode is known.
type ModeSelector struct {
	config                 Config
	watchdogProcessFactory WatchdogProcessFactory
	orchestratorFactory    OrchestratorFactory
	osLayer                OSLayer
	parser                 Parser
}

func New(
	config Config,
	parser Parser,
	watchdogProcessFactory WatchdogProcessFactory,
	orchestratorFactory OrchestratorFactory,
	osLayer OSLayer,
) *ModeSelector {
	return &ModeSelector{
		config:                 config,
		parser:                 parser,
		watchdogProcessFactory: watchdogProcessFactory,
		orchestratorFactory:    orchestratorFactory,
		osLayer:                osLayer,
	}
}

func (a *ModeSelector) StartAndWaitForCompletion(ctx context.Context) error {
	switch {
	case a.config.HelpMode():
		_, err := fmt.Fprintf(a.osLayer.Stdout(), "%s\n", a.parser.Usage())
		return err
	case a.config.VersionMode():
		_, err := fmt.Fprintf(a.osLayer.Stdout(), "%s\n", a.config.Version())
		return err
	case a.config.WatchdogMode():
		watchdogProcess, err := a.watchdogProcessFactory.Create()
		if err != nil {
			return err
		}

		return watchdogProcess.StartAndWaitForCompletion(ctx)
	default:
		orchestrator, err := a.orchestratorFactory.Create()
		if err != nil {
			return err
		}

		return orchestrator.StartAndWaitForCompletion(ctx)
	}
}
