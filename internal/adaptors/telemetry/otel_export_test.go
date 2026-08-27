// Copyright 2026 The MathWorks, Inc.

package telemetry

import (
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
)

func NewOTELTelemetryForTesting(
	logger entities.Logger,
	meter otel.Meter,
	instrumentFactory InstrumentFactory,
	cfg Config,
	dir Directory,
	osLayer OSLayer,
	osVersionProvider OSVersionProvider,
	serverDefinition Definition,
) (Telemetry, messages.Error) {
	return newOTELTelemetry(logger, meter, instrumentFactory, cfg, dir, osLayer, osVersionProvider, serverDefinition)
}

func SHA256Prefix16ForTesting(name string) string {
	return sha256Prefix16(name)
}
