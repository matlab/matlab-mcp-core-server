// Copyright 2026 The MathWorks, Inc.

package basetool

import (
	"context"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-server/internal/messages"
)

func (t tool[_, _]) recordToolCall(ctx context.Context, source telemetry.ToolSource) messages.Error {
	tel, err := t.telemetryFactory.Telemetry()
	if err != nil {
		return err
	}
	tel.RecordToolCallRequest(ctx, t.name, source)
	return nil
}
