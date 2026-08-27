// Copyright 2026 The MathWorks, Inc.

package telemetry

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry/otel/instruments"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
)

type Directory interface {
	ID() string
}

type Config interface {
	Version() string
	WatchdogMode() bool
	SpecifiedParameters() []string
	AsPIISafeJSONString() string
}

type otelTelemetry struct {
	logger            entities.Logger
	meter             otel.Meter
	instrumentFactory InstrumentFactory

	config            Config
	directory         Directory
	osLayer           OSLayer
	osVersionProvider OSVersionProvider
	definition        Definition

	// Instruments
	serverStartCounter      instruments.Int64Counter
	clientConnectionCounter instruments.Int64Counter
	toolCallRequestCounter  instruments.Int64Counter
}

func newOTELTelemetry(
	logger entities.Logger,
	meter otel.Meter,
	instrumentFactory InstrumentFactory,
	cfg Config,
	dir Directory,
	osLayer OSLayer,
	osVersionProvider OSVersionProvider,
	definition Definition,
) (*otelTelemetry, messages.Error) {
	telemetry := &otelTelemetry{
		logger:            logger,
		meter:             meter,
		instrumentFactory: instrumentFactory,
		config:            cfg,
		directory:         dir,
		osLayer:           osLayer,
		osVersionProvider: osVersionProvider,
		definition:        definition,
	}

	err := telemetry.createInstruments(logger)
	if err != nil {
		return nil, err
	}

	return telemetry, nil
}

func (t *otelTelemetry) RecordServerStart(ctx context.Context) {
	if t.config.WatchdogMode() {
		t.logger.Debug("Skipping server start metric in watchdog mode")
		return
	}

	t.logger.Debug("Recording server start metric")
	attributes := NewAttributes(t.logger)

	attributes.AddString("server.instance_id", t.directory.ID())
	attributes.AddString("server.name", t.definition.Name())
	attributes.AddString("server.version", t.config.Version())
	attributes.AddStringSlice("server.specified_parameters", t.config.SpecifiedParameters())
	attributes.AddString("server.config_details", t.config.AsPIISafeJSONString())
	attributes.AddString("server.os", t.osLayer.GOOS())
	osVersion, err := t.osVersionProvider.Version()
	if err != nil {
		t.logger.WithError(err).Debug("Failed to get OS version")
		osVersion = ""
	}
	attributes.AddString("server.os_version", osVersion)

	t.serverStartCounter.Add(ctx, 1, attributes.AsOTEL())
}

func (t *otelTelemetry) RecordToolCallRequest(ctx context.Context, toolName string, source ToolSource) {
	t.logger.With("tool-name", toolName).With("tool-source", source).Debug("Recording tool call request metric")
	attributes := NewAttributes(t.logger)

	attributes.AddString("server.instance_id", t.directory.ID())
	attributes.AddString("tool.name", toolNameAttribute(toolName, source))

	t.toolCallRequestCounter.Add(ctx, 1, attributes.AsOTEL())
}

func toolNameAttribute(toolName string, source ToolSource) string {
	if source == ToolSourceExtension {
		return sha256Prefix16(toolName)
	}
	return toolName
}

func sha256Prefix16(name string) string {
	sum := sha256.Sum256([]byte(name))
	return hex.EncodeToString(sum[:])[:16]
}

func (t *otelTelemetry) RecordClientConnection(ctx context.Context, info ClientConnectionInfo) {
	t.logger.Debug("Recording client connection metric")
	attributes := NewAttributes(t.logger)

	attributes.AddString("server.instance_id", t.directory.ID())
	attributes.AddString("client.name", info.Name)
	attributes.AddString("client.title", info.Title)
	attributes.AddString("client.website_url", info.WebsiteURL)
	attributes.AddString("client.version", info.Version)
	attributes.AddStringSlice("client.capabilities", info.Capabilities)
	attributes.AddString("client.capabilities_details", info.CapabilitiesJSON)

	t.clientConnectionCounter.Add(ctx, 1, attributes.AsOTEL())
}

func (t *otelTelemetry) createInstruments(logger entities.Logger) messages.Error {
	logger.Debug("Creating telemetry instruments")
	defer logger.Debug("Done creating telemetry instruments")

	serverStartCounter, err := t.instrumentFactory.NewInt64Counter(t.meter, "server.starts", "Number of times the server has started", "{start}")
	if err != nil {
		logger.WithError(err).Error("Failed to create server start counter")
		return messages.New_StartupErrors_TelemetryInitializationFailed_Error()
	}
	t.serverStartCounter = serverStartCounter

	clientConnectionCounter, err := t.instrumentFactory.NewInt64Counter(t.meter, "server.client_connections", "Number of times a client connected to a server", "{connection}")
	if err != nil {
		logger.WithError(err).Error("Failed to create client connection counter")
		return messages.New_StartupErrors_TelemetryInitializationFailed_Error()
	}
	t.clientConnectionCounter = clientConnectionCounter

	toolCallRequestCounter, err := t.instrumentFactory.NewInt64Counter(t.meter, "tool.calls_request", "Number of tool invocations", "{call}")
	if err != nil {
		logger.WithError(err).Error("Failed to create tool call request counter")
		return messages.New_StartupErrors_TelemetryInitializationFailed_Error()
	}
	t.toolCallRequestCounter = toolCallRequestCounter

	return nil
}
