// Copyright 2025 The MathWorks, Inc.

package globalmatlab

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type MATLABManager interface {
	StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error)
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type MATLABRootSelector interface {
	SelectFirstMATLABVersionOnPath(ctx context.Context, logger entities.Logger) (string, error)
}

type MATLABStartingDirSelector interface {
	SelectMatlabStartingDir() (string, error)
}

type GlobalMATLAB struct {
	matlabManager             MATLABManager
	matlabRootSelector        MATLABRootSelector
	matlabStartingDirSelector MATLABStartingDirSelector

	lock              *sync.Mutex
	matlabRoot        string
	matlabStartingDir string
	sessionID         entities.SessionID
	cachedStartErr    error
	isReady           bool
}

func New(
	matlabManager MATLABManager,
	matlabRootSelector MATLABRootSelector,
	matlabStartingDirSelector MATLABStartingDirSelector,
) *GlobalMATLAB {
	return &GlobalMATLAB{
		matlabManager:             matlabManager,
		matlabRootSelector:        matlabRootSelector,
		matlabStartingDirSelector: matlabStartingDirSelector,

		lock: &sync.Mutex{},
	}
}

func (g *GlobalMATLAB) Initialize(ctx context.Context, logger entities.Logger) error {
	logger = logger.With("mcp_server_pid", os.Getpid())
	logger.Debug("GlobalMATLAB.Initialize called")

	var err error
	g.matlabRoot, err = g.matlabRootSelector.SelectFirstMATLABVersionOnPath(ctx, logger)
	if err != nil {
		return err
	}

	g.matlabStartingDir, err = g.matlabStartingDirSelector.SelectMatlabStartingDir()
	if err != nil {
		logger.WithError(err).Warn("failed to determine MATLAB starting directory, proceeding without one")
	}

	err = g.ensureMATLABClientIsValid(ctx, logger)
	if err != nil {
		return err
	}

	logger.Debug("GlobalMATLAB.Initialize completed successfully")
	return nil
}

func (g *GlobalMATLAB) Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	// Retry logic with exponential backoff for transient connection failures
	var lastErr error
	maxRetries := 5
	baseDelay := 200 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 200ms, 400ms, 800ms, 1600ms, 3200ms
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			logger.With("attempt", attempt+1).With("delay_ms", delay.Milliseconds()).Debug("Retrying MATLAB client connection")
			time.Sleep(delay)
		}

		if err := g.ensureMATLABClientIsValid(ctx, logger); err != nil {
			lastErr = err
			// Don't retry on permanent errors (cached start errors)
			if g.cachedStartErr != nil {
				return nil, err
			}
			continue
		}

		client, err := g.matlabManager.GetMATLABSessionClient(ctx, logger, g.sessionID)
		if err != nil {
			lastErr = err
			continue
		}

		// Test the connection with a simple eval to ensure MATLAB is ready
		g.lock.Lock()
		needsReadyCheck := !g.isReady
		g.lock.Unlock()

		if needsReadyCheck {
			if err := g.waitForMATLABReady(ctx, logger, client); err != nil {
				lastErr = err
				logger.WithError(err).With("attempt", attempt+1).Debug("MATLAB not ready yet, will retry")
				continue
			}
			g.lock.Lock()
			g.isReady = true
			g.lock.Unlock()
			logger.Debug("MATLAB connection verified and ready")
		}

		return client, nil
	}

	return nil, fmt.Errorf("failed to get MATLAB client after %d attempts: %w", maxRetries, lastErr)
}

func (g *GlobalMATLAB) ensureMATLABClientIsValid(ctx context.Context, logger entities.Logger) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.cachedStartErr != nil {
		logger.Debug("ensureMATLABClientIsValid: returning cached error")
		return g.cachedStartErr
	}

	var sessionIDZeroValue entities.SessionID
	if g.sessionID == sessionIDZeroValue {
		logger.With("matlab_root", g.matlabRoot).Debug("ensureMATLABClientIsValid: starting new MATLAB session")
		sessionID, err := g.matlabManager.StartMATLABSession(ctx, logger, entities.LocalSessionDetails{
			MATLABRoot:        g.matlabRoot,
			StartingDirectory: g.matlabStartingDir,
			ShowMATLABDesktop: true,
		})
		if err != nil {
			g.cachedStartErr = err
			logger.WithError(err).Error("ensureMATLABClientIsValid: failed to start MATLAB session")
			return err
		}

		g.sessionID = sessionID
		g.isReady = false // Reset readiness when starting a new session
		logger.With("session_id", sessionID).Debug("ensureMATLABClientIsValid: MATLAB session started, waiting for connection to be ready")
	} else {
		logger.With("session_id", g.sessionID).Debug("ensureMATLABClientIsValid: reusing existing MATLAB session")
	}

	return nil
}

// waitForMATLABReady tests the MATLAB connection with a simple eval to ensure it's ready
// This gives MATLAB time to fully initialize the Embedded Connector before accepting requests
func (g *GlobalMATLAB) waitForMATLABReady(ctx context.Context, logger entities.Logger, client entities.MATLABSessionClient) error {
	// Create a timeout context for the readiness check (max 30 seconds)
	readyCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Try a simple eval with retries
	maxAttempts := 10
	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			if delay > 5*time.Second {
				delay = 5 * time.Second // Cap at 5 seconds
			}
			logger.With("attempt", attempt+1).With("delay_ms", delay.Milliseconds()).Debug("Waiting for MATLAB to be ready")
			
			select {
			case <-readyCtx.Done():
				return fmt.Errorf("timeout waiting for MATLAB to be ready: %w", readyCtx.Err())
			case <-time.After(delay):
			}
		}

		// Test connection with a simple eval
		_, err := client.Eval(readyCtx, logger, entities.EvalRequest{
			Code: "1+1",
		})

		if err == nil {
			logger.Debug("MATLAB connection test successful")
			return nil
		}

		logger.WithError(err).With("attempt", attempt+1).Debug("MATLAB connection test failed, will retry")
	}

	return fmt.Errorf("MATLAB connection not ready after %d attempts", maxAttempts)
}
