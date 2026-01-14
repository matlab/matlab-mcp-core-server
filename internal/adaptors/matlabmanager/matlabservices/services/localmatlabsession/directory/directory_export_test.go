// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"time"
)

func NewDirectory(sessionDir string, osLayer OSLayer) *directory {
	return newDirectory(sessionDir, osLayer)
}

func (d *directory) SecurePortFile() string {
	return d.securePortFile()
}

func (d *directory) SetEmbeddedConnectorDetailsTimeout(timeout time.Duration) {
	d.embeddedConnectorDetailsTimeout = timeout
}

func (d *directory) SetEmbeddedConnectorDetailsRetry(retry time.Duration) {
	d.embeddedConnectorDetailsRetry = retry
}

func (d *directory) SetCleanupTimeout(timeout time.Duration) {
	d.cleanupTimeout = timeout
}

func (d *directory) SetCleanupRetry(retry time.Duration) {
	d.cleanupRetry = retry
}
