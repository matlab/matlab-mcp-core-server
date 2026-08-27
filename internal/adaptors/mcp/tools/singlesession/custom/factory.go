// Copyright 2026 The MathWorks, Inc.

package custom

import (
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
)

type Loader interface {
	Load(filePath string) ([]definition.ValidatedTool, messages.Error)
}

type Factory struct {
	loader           Loader
	loggerFactory    basetool.LoggerFactory
	telemetryFactory basetool.TelemetryFactory
	usecase          Usecase
	globalMATLAB     entities.GlobalMATLAB
	configFactory    ConfigFactory
}

func NewFactory(
	loader Loader,
	loggerFactory basetool.LoggerFactory,
	telemetryFactory basetool.TelemetryFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
	configFactory ConfigFactory,
) *Factory {
	return &Factory{
		loader:           loader,
		loggerFactory:    loggerFactory,
		telemetryFactory: telemetryFactory,
		usecase:          usecase,
		globalMATLAB:     globalMATLAB,
		configFactory:    configFactory,
	}
}

func (f *Factory) LoadTools(filePath string) ([]tools.Tool, messages.Error) {
	validatedTools, err := f.loader.Load(filePath)
	if err != nil {
		return nil, err
	}

	result := make([]tools.Tool, 0, len(validatedTools))
	for _, vt := range validatedTools {
		result = append(result, NewTool(vt, f.loggerFactory, f.configFactory, f.telemetryFactory, f.usecase, f.globalMATLAB))
	}

	return result, nil
}
