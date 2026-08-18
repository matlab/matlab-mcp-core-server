// Copyright 2025-2026 The MathWorks, Inc.

package sdk

import (
	"context"
	"encoding/json"
	"maps"
	"slices"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/matlab/matlab-mcp-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Definition interface {
	Name() string
	Title() string
	Instructions() string
	Features() definition.Features
}

type RootStore interface {
	UpdateRoots(roots []*mcp.Root)
}

type ClientInfoStore interface {
	SetClientInfo(info entities.MCPClientInfo)
}

type LoggerFactory interface {
	GetGlobalLogger() (entities.Logger, messages.Error)
}

type GlobalMATLAB interface {
	Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error)
}

type TelemetryFactory interface {
	Telemetry() (telemetry.Telemetry, messages.Error)
}

type MCPSession interface {
	InitializeParams() *mcp.InitializeParams
	ListRoots(ctx context.Context, params *mcp.ListRootsParams) (*mcp.ListRootsResult, error)
}

type Factory struct {
	configFactory    ConfigFactory
	definition       Definition
	rootStore        RootStore
	clientInfoStore  ClientInfoStore
	loggerFactory    LoggerFactory
	globalMATLAB     GlobalMATLAB
	telemetryFactory TelemetryFactory
}

type serverCallbackHandler struct {
	config          config.Config
	logger          entities.Logger
	features        definition.Features
	rootStore       RootStore
	clientInfoStore ClientInfoStore
	globalMATLAB    GlobalMATLAB
	telemetry       telemetry.Telemetry
}

func NewFactory(
	configFactory ConfigFactory,
	definition Definition,
	rootStore RootStore,
	clientInfoStore ClientInfoStore,
	loggerFactory LoggerFactory,
	globalMATLAB GlobalMATLAB,
	telemetryFactory TelemetryFactory,
) *Factory {
	return &Factory{
		configFactory:    configFactory,
		definition:       definition,
		rootStore:        rootStore,
		clientInfoStore:  clientInfoStore,
		loggerFactory:    loggerFactory,
		globalMATLAB:     globalMATLAB,
		telemetryFactory: telemetryFactory,
	}
}

func (f *Factory) NewServer() (*mcp.Server, messages.Error) {
	cfg, err := f.configFactory.Config()
	if err != nil {
		return nil, err
	}

	logger, err := f.loggerFactory.GetGlobalLogger()
	if err != nil {
		return nil, err
	}

	tel, err := f.telemetryFactory.Telemetry()
	if err != nil {
		return nil, err
	}

	s := &serverCallbackHandler{
		config:          cfg,
		logger:          logger,
		features:        f.definition.Features(),
		rootStore:       f.rootStore,
		clientInfoStore: f.clientInfoStore,
		globalMATLAB:    f.globalMATLAB,
		telemetry:       tel,
	}

	impl := &mcp.Implementation{
		Name:    f.definition.Name(),
		Title:   f.definition.Title(),
		Version: cfg.Version(),
	}
	options := &mcp.ServerOptions{
		Instructions:            f.definition.Instructions(),
		InitializedHandler:      s.handleInitialized,
		RootsListChangedHandler: s.handleRootsListChanged,
	}

	return mcp.NewServer(impl, options), nil
}

func (s *serverCallbackHandler) handleInitialized(ctx context.Context, req *mcp.InitializedRequest) {
	if req == nil ||
		req.Session == nil {
		return
	}

	s.logClientDetails(req.Session)
	s.recordClientConnection(ctx, req.Session)

	if err := s.updateRoots(ctx, req.Session); err != nil {
		s.logger.
			WithError(err).
			Warn("failed to update MCP roots, using fallback starting folder")
	}

	matlabEnabled := s.features.MATLAB.Enabled
	if matlabEnabled && s.config.UseSingleMATLABSession() && s.config.InitializeMATLABOnStartup() {
		go func() {
			s.logger.Debug("Eagerly initializing MATLAB")

			startMATLABCtx := context.WithoutCancel(ctx)
			if _, err := s.globalMATLAB.Client(startMATLABCtx, s.logger); err != nil {
				s.logger.
					WithError(err).
					Warn("MATLAB eager initialization failed")
			}
		}()
	}
}

func (s *serverCallbackHandler) handleRootsListChanged(ctx context.Context, req *mcp.RootsListChangedRequest) {
	if err := s.updateRoots(ctx, req.Session); err != nil {
		s.logger.WithError(err).Warn("failed to update MCP roots, using fallback starting folder")
	}
}

func (s *serverCallbackHandler) recordClientConnection(ctx context.Context, session MCPSession) {
	initializeParams := session.InitializeParams()
	if initializeParams == nil {
		return
	}

	info := telemetry.ClientConnectionInfo{}

	if initializeParams.ClientInfo != nil {
		info.Name = initializeParams.ClientInfo.Name
		info.Title = initializeParams.ClientInfo.Title
		info.WebsiteURL = initializeParams.ClientInfo.WebsiteURL
		info.Version = initializeParams.ClientInfo.Version

		// Cache the client identity for the connection indicator to render later.
		s.clientInfoStore.SetClientInfo(entities.MCPClientInfo{
			Name:       info.Name,
			Title:      info.Title,
			WebsiteURL: info.WebsiteURL,
			Version:    info.Version,
		})
	}

	if initializeParams.Capabilities != nil {
		info.Capabilities, info.CapabilitiesJSON = marshalCapabilities(initializeParams.Capabilities)
	}

	s.telemetry.RecordClientConnection(ctx, info)
}

func marshalCapabilities(caps *mcp.ClientCapabilities) ([]string, string) {
	data, err := json.Marshal(caps)
	if err != nil {
		return nil, ""
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, string(data)
	}

	// The legacy "roots" field is a non-pointer struct that always marshals even when empty.
	// Replace it with RootsV2 (the authoritative pointer field) or remove it entirely.
	delete(fields, "roots")
	if caps.RootsV2 != nil {
		rootsJSON, err := json.Marshal(caps.RootsV2)
		if err == nil {
			fields["roots"] = rootsJSON
		}
	}

	names := slices.Collect(maps.Keys(fields))
	slices.Sort(names)

	detailsJSON, err := json.Marshal(fields)
	if err != nil {
		return names, string(data)
	}

	return names, string(detailsJSON)
}

func (s *serverCallbackHandler) logClientDetails(session MCPSession) {
	initializeParams := session.InitializeParams()
	if initializeParams != nil &&
		initializeParams.ClientInfo != nil {
		clientInfo := initializeParams.ClientInfo
		s.logger.
			With("client-name", clientInfo.Name).
			With("client-title", clientInfo.Title).
			With("client-url", clientInfo.WebsiteURL).
			With("client-version", clientInfo.Version).
			Info("New client session")
	}
}

func (s *serverCallbackHandler) updateRoots(ctx context.Context, session MCPSession) error {
	// RootsV2 is the correct pointer field for this check.
	// The legacy Roots field is a value type (not a pointer) due to a go-sdk bug (issue #607),
	// making it impossible to distinguish "no roots support" from "empty roots support".
	params := session.InitializeParams()
	if params == nil || params.Capabilities == nil || params.Capabilities.RootsV2 == nil {
		return nil
	}

	result, err := session.ListRoots(ctx, nil)
	if err != nil {
		return err
	}

	s.rootStore.UpdateRoots(result.Roots)
	s.logger.With("roots", result.Roots).Debug("Updated MCP roots from client")

	return nil
}
