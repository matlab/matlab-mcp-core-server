// Copyright 2026 The MathWorks, Inc.

package sdk

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func HandleInitialized(cfg config.Config, logger entities.Logger, features definition.Features, rs RootStore, gm GlobalMATLAB, tel telemetry.Telemetry) func(context.Context, *mcp.InitializedRequest) {
	s := &serverCallbackHandler{config: cfg, logger: logger, features: features, rootStore: rs, globalMATLAB: gm, telemetry: tel}
	return s.handleInitialized
}

func HandleRootsListChanged(cfg config.Config, logger entities.Logger, features definition.Features, rs RootStore, gm GlobalMATLAB, tel telemetry.Telemetry) func(context.Context, *mcp.RootsListChangedRequest) {
	s := &serverCallbackHandler{config: cfg, logger: logger, features: features, rootStore: rs, globalMATLAB: gm, telemetry: tel}
	return s.handleRootsListChanged
}

func UpdateRoots(cfg config.Config, logger entities.Logger, features definition.Features, rs RootStore, gm GlobalMATLAB, tel telemetry.Telemetry) func(context.Context, MCPSession) error {
	s := &serverCallbackHandler{config: cfg, logger: logger, features: features, rootStore: rs, globalMATLAB: gm, telemetry: tel}
	return s.updateRoots
}

func LogClientDetails(logger entities.Logger) func(MCPSession) {
	s := &serverCallbackHandler{logger: logger}
	return s.logClientDetails
}

func RecordClientConnection(logger entities.Logger, tel telemetry.Telemetry, store ClientInfoStore) func(context.Context, MCPSession) {
	s := &serverCallbackHandler{logger: logger, telemetry: tel, clientInfoStore: store}
	return s.recordClientConnection
}
