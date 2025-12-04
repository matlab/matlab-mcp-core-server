// Copyright 2025 The MathWorks, Inc.

package directorymanager

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

type directoryManager struct {
	sessionDir string
	osLayer    OSLayer

	embeddedConnectorDetailsTimeout time.Duration
	embeddedConnectorDetailsRetry   time.Duration
	cleanupTimeout                  time.Duration
	cleanupRetry                    time.Duration
}

func newDirectoryManager(sessionDir string, osLayer OSLayer) *directoryManager {
	return &directoryManager{
		sessionDir: sessionDir,
		osLayer:    osLayer,

		embeddedConnectorDetailsTimeout: defaultEmbeddedConnectorDetailsTimeout,
		embeddedConnectorDetailsRetry:   defaultEmbeddedConnectorDetailsRetry,
		cleanupTimeout:                  defaultCleanupTimeout,
		cleanupRetry:                    defaultCleanupRetry,
	}
}

func (m *directoryManager) Path() string {
	return m.sessionDir
}

func (m *directoryManager) CertificateFile() string {
	return filepath.Join(m.sessionDir, certificateFile)
}

func (m *directoryManager) CertificateKeyFile() string {
	return filepath.Join(m.sessionDir, certificateKeyFile)
}

func (m *directoryManager) GetEmbeddedConnectorDetails() (string, []byte, error) {
	securePortFileFullPath := m.securePortFile()
	certificateFileFullPath := m.CertificateFile()

	timeout := time.After(m.embeddedConnectorDetailsTimeout)
	tick := time.Tick(m.embeddedConnectorDetailsRetry)

	for {
		select {
		case <-timeout:
			return "", nil, fmt.Errorf("timeout waiting for worker to start")
		case <-tick:
			if _, err := m.osLayer.Stat(securePortFileFullPath); err != nil {
				continue
			}
			if _, err := m.osLayer.Stat(certificateFileFullPath); err != nil {
				continue
			}
			securePort, err := m.osLayer.ReadFile(securePortFileFullPath)
			if err != nil {
				return "", nil, fmt.Errorf("failed to read secure port file: %w", err)
			}
			if string(securePort) == "" {
				// File was made, but content is empty, wait for next tick
				continue
			}
			certificatePEM, err := m.osLayer.ReadFile(certificateFileFullPath)
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

func (m *directoryManager) Cleanup() error {
	if m.sessionDir == "" {
		return nil
	}

	timeout := time.After(m.cleanupTimeout)
	tick := time.Tick(m.cleanupRetry)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout trying to delete session directory %s", m.sessionDir)
		case <-tick:
			err := m.osLayer.RemoveAll(m.sessionDir)
			if err == nil {
				return nil
			}
		}
	}
}

func (m *directoryManager) securePortFile() string {
	return filepath.Join(m.sessionDir, securePortFile)
}
