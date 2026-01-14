// Copyright 2025-2026 The MathWorks, Inc.

package directory

import (
	"fmt"
	"path/filepath"
	"time"
)

const defaultEmbeddedConnectorDetailsTimeout = 5 * time.Minute
const defaultEmbeddedConnectorDetailsRetry = 500 * time.Millisecond

const defaultCleanupTimeout = 2 * time.Minute
const defaultCleanupRetry = 500 * time.Millisecond

const securePortFile = "connector.securePort"
const certificateFile = "cert.pem"
const certificateKeyFile = "cert.key"

type directory struct {
	sessionDir string
	osLayer    OSLayer

	embeddedConnectorDetailsTimeout time.Duration
	embeddedConnectorDetailsRetry   time.Duration
	cleanupTimeout                  time.Duration
	cleanupRetry                    time.Duration
}

func newDirectory(sessionDir string, osLayer OSLayer) *directory {
	return &directory{
		sessionDir: sessionDir,
		osLayer:    osLayer,

		embeddedConnectorDetailsTimeout: defaultEmbeddedConnectorDetailsTimeout,
		embeddedConnectorDetailsRetry:   defaultEmbeddedConnectorDetailsRetry,
		cleanupTimeout:                  defaultCleanupTimeout,
		cleanupRetry:                    defaultCleanupRetry,
	}
}

func (d *directory) Path() string {
	return d.sessionDir
}

func (d *directory) CertificateFile() string {
	return filepath.Join(d.sessionDir, certificateFile)
}

func (d *directory) CertificateKeyFile() string {
	return filepath.Join(d.sessionDir, certificateKeyFile)
}

func (d *directory) GetEmbeddedConnectorDetails() (string, []byte, error) {
	securePortFileFullPath := d.securePortFile()
	certificateFileFullPath := d.CertificateFile()

	timeout := time.After(d.embeddedConnectorDetailsTimeout)
	tick := time.Tick(d.embeddedConnectorDetailsRetry)

	for {
		select {
		case <-timeout:
			return "", nil, fmt.Errorf("timeout waiting for worker to start")
		case <-tick:
			if _, err := d.osLayer.Stat(securePortFileFullPath); err != nil {
				continue
			}
			if _, err := d.osLayer.Stat(certificateFileFullPath); err != nil {
				continue
			}
			securePort, err := d.osLayer.ReadFile(securePortFileFullPath)
			if err != nil {
				return "", nil, fmt.Errorf("failed to read secure port file: %w", err)
			}
			if string(securePort) == "" {
				// File was made, but content is empty, wait for next tick
				continue
			}
			certificatePEM, err := d.osLayer.ReadFile(certificateFileFullPath)
			if err != nil {
				return "", nil, fmt.Errorf("failed to read certificate path file: %w", err)
			}
			if string(certificatePEM) == "" {
				// File was made, but content is empty, wait for next tick
				continue
			}
			return string(securePort), certificatePEM, nil
		}
	}
}

func (d *directory) Cleanup() error {
	if d.sessionDir == "" {
		return nil
	}

	timeout := time.After(d.cleanupTimeout)
	tick := time.Tick(d.cleanupRetry)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout trying to delete session directory %s", d.sessionDir)
		case <-tick:
			err := d.osLayer.RemoveAll(d.sessionDir)
			if err == nil {
				return nil
			}
		}
	}
}

func (d *directory) securePortFile() string {
	return filepath.Join(d.sessionDir, securePortFile)
}
