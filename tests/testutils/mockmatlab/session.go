// Copyright 2026 The MathWorks, Inc.

package mockmatlab

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/sessiondetails"
)

const (
	securePortFile     = "connector.securePort"
	certificateFile    = "cert.pem"
	certificateKeyFile = "cert.key"

	EnvMockMATLABConfig = mockruntime.EnvMockMATLABConfig

	defaultReadyTimeout = 10 * time.Second
	defaultReadyPoll    = 100 * time.Millisecond
)

type Config = mockruntime.Config

func HappyConfig() Config {
	return mockruntime.HappyConfig()
}

func HangBeforeFilesConfig() Config {
	return mockruntime.HangBeforeFilesConfig()
}

func ExitImmediatelyConfig(exitCode int) Config {
	return mockruntime.ExitImmediatelyConfig(exitCode)
}

func SlowStartupConfig(delayMs int) Config {
	return mockruntime.SlowStartupConfig(delayMs)
}

func StartupFailureConfig() Config {
	return mockruntime.StartupFailureConfig()
}

type Session struct {
	cmd               *exec.Cmd
	SessionDir        string
	APIKey            string
	connectionDetails embeddedconnector.ConnectionDetails
}

func StartSession(ctx context.Context, installation *Installation, cfg Config) (*Session, error) {
	sessionDir, err := os.MkdirTemp("", "mock-matlab-session-")
	if err != nil {
		return nil, fmt.Errorf("failed to create session dir: %w", err)
	}

	apiKey := "mock-api-key" //nolint:gosec // Not a real credential
	certFile := filepath.Join(sessionDir, certificateFile)
	keyFile := filepath.Join(sessionDir, certificateKeyFile)

	configJSON, err := cfg.ToEnvValue()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize config: %w", err)
	}

	binaryPath := mockMATLABBinaryPath(installation.MATLABRoot)
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec // Trusted test path
	cmd.Env = append(os.Environ(),
		"MW_MCP_SESSION_DIR="+sessionDir,
		"MWAPIKEY="+apiKey,
		"MW_CERTFILE="+certFile,
		"MW_PKEYFILE="+keyFile,
		mockruntime.EnvMockMATLABConfig+"="+configJSON,
	)

	if err := cmd.Start(); err != nil {
		_ = os.RemoveAll(sessionDir)
		return nil, fmt.Errorf("failed to start mock MATLAB: %w", err)
	}

	return &Session{
		cmd:        cmd,
		SessionDir: sessionDir,
		APIKey:     apiKey,
	}, nil
}

func (s *Session) WaitForReady(ctx context.Context) (embeddedconnector.ConnectionDetails, error) {
	portPath := filepath.Join(s.SessionDir, securePortFile)
	certPath := filepath.Join(s.SessionDir, certificateFile)

	ctx, cancel := context.WithTimeout(ctx, defaultReadyTimeout)
	defer cancel()

	details, err := retry.Retry(ctx, func() (embeddedconnector.ConnectionDetails, bool, error) {
		port, readErr := readNonEmptyFile(portPath)
		if readErr != nil {
			return embeddedconnector.ConnectionDetails{}, false, nil
		}
		certPEM, readErr := readNonEmptyFile(certPath)
		if readErr != nil {
			return embeddedconnector.ConnectionDetails{}, false, nil
		}
		return embeddedconnector.ConnectionDetails{
			Host:           "localhost",
			Port:           string(port),
			APIKey:         s.APIKey,
			CertificatePEM: certPEM,
		}, true, nil
	}, retry.NewLinearRetryStrategy(defaultReadyPoll))

	if err != nil {
		return embeddedconnector.ConnectionDetails{}, fmt.Errorf("timeout waiting for mock MATLAB to become ready")
	}

	s.connectionDetails = details
	return details, nil
}

// ToSessionDetailsJSON returns a JSON string in the format expected by
// --matlab-session-connection-details. The certificate field is the file path
// to the PEM cert in the session directory.
func (s *Session) ToSessionDetailsJSON() (string, error) {
	return sessiondetails.MarshalJSON(
		s.connectionDetails.Port,
		s.CertificatePath(),
		s.APIKey,
		s.cmd.Process.Pid,
	)
}

// CertificatePath returns the path to the PEM certificate file in the session directory.
func (s *Session) CertificatePath() string {
	return filepath.Join(s.SessionDir, certificateFile)
}

// ShareMATLABSession publishes session details to the standard discovery
// location under homeDir, mimicking what a real MATLAB session would do.
func (s *Session) ShareMATLABSession(homeDir string) (string, error) {
	detailsJSON, err := s.ToSessionDetailsJSON()
	if err != nil {
		return "", err
	}
	return sessiondetails.Publish(homeDir, detailsJSON)
}

// ReceivedRequest represents a single request logged by the mock MATLAB process.
type ReceivedRequest struct {
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"messageType"`
	Content     string    `json:"content"`
}

// ReceivedRequests reads all requests logged by this mock MATLAB session.
func (s *Session) ReceivedRequests() ([]ReceivedRequest, error) {
	logPath := filepath.Join(s.SessionDir, "requests.jsonl")
	f, err := os.Open(logPath) //nolint:gosec // Test utility reading from session temp dir
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open request log: %w", err)
	}
	defer f.Close() //nolint:errcheck // Read-only file handle in test utility

	var requests []ReceivedRequest
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var entry ReceivedRequest
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, fmt.Errorf("failed to parse request log entry: %w", err)
		}
		requests = append(requests, entry)
	}
	return requests, scanner.Err()
}

// ReceivedEvals returns only the Eval requests received by this mock, with their mcode content.
func (s *Session) ReceivedEvals() ([]ReceivedRequest, error) {
	all, err := s.ReceivedRequests()
	if err != nil {
		return nil, err
	}
	var evals []ReceivedRequest
	for _, r := range all {
		if r.MessageType == "Eval" || r.MessageType == "FEval" {
			evals = append(evals, r)
		}
	}
	return evals, nil
}

func (s *Session) Stop() error {
	var firstErr error
	recordErr := func(err error) {
		if firstErr == nil {
			firstErr = err
		}
	}

	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			recordErr(fmt.Errorf("failed to kill mock MATLAB process: %w", err))
		}

		if err := s.cmd.Wait(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				recordErr(fmt.Errorf("failed to wait for mock MATLAB process: %w", err))
			}
		}
	}

	if err := os.RemoveAll(s.SessionDir); err != nil {
		recordErr(fmt.Errorf("failed to remove mock MATLAB session directory: %w", err))
	}

	return firstErr
}

func (s *Session) Wait() error {
	return s.cmd.Wait()
}

func (s *Session) ProcessExited() bool {
	return s.cmd.ProcessState != nil
}

func readNonEmptyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path) //nolint:gosec // Trusted test path
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	return data, nil
}
