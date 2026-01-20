// Copyright 2025-2026 The MathWorks, Inc.

package os

import "time"

func (pm *ProcessManager) SetCheckParentAliveInterval(interval time.Duration) {
	pm.checkParentAliveInterval = interval
}
