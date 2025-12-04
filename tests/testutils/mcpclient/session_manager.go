// Copyright 2025 The MathWorks, Inc.

package mcpclient

import (
	"context"
	"fmt"
)

// SessionManager helps manage MATLAB sessions
type SessionManager struct {
	session *MCPClientSession
}

// MATLABInfo represents information about a MATLAB installation
type MATLABInfo struct {
	Path    string `json:"matlab_root"`
	Version string `json:"version"`
}

// ListAvailableMATLABs lists available MATLAB installations
func (sm *SessionManager) ListAvailableMATLABs(ctx context.Context) ([]MATLABInfo, error) {
	result, err := sm.session.CallTool(ctx, "list_available_matlabs", map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("list_available_matlabs tool call failed: %v", err)
	}
	var output struct {
		AvailableMATLABs []MATLABInfo `json:"available_matlabs"`
	}
	err = sm.session.UnmarshalStructuredContent(result, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal list_available_matlabs output: %v", err)
	}
	return output.AvailableMATLABs, nil
}

// StartSession starts a new MATLAB session
func (sm *SessionManager) StartSession(ctx context.Context, matlabRoot ...string) (int, error) {
	args := map[string]any{}
	if len(matlabRoot) > 0 {
		args["matlab_root"] = matlabRoot[0]
	}
	result, err := sm.session.CallTool(ctx, "start_matlab_session", args)
	if err != nil {
		return 0, fmt.Errorf("failed to call start_matlab_session tool: %v", err)
	}
	var output struct {
		SessionID int `json:"session_id"`
	}
	err = sm.session.UnmarshalStructuredContent(result, &output)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal start_matlab_session output: %v", err)
	}
	return output.SessionID, nil
}

// StopSession stops a MATLAB session
func (sm *SessionManager) StopSession(ctx context.Context, sessionID int) error {
	result, err := sm.session.CallTool(ctx, "stop_matlab_session", map[string]any{
		"session_id": sessionID,
	})
	if err != nil {
		return fmt.Errorf("stop_matlab_session tool call failed: %v", err)
	}
	if result != nil && result.IsError {
		textContent, _ := sm.session.GetTextContent(result)
		return fmt.Errorf("stop_matlab_session failed: %s", textContent)
	}
	return nil
}

// CleanupSession stops a MATLAB session for cleanup, ignoring errors if the session is already gone
func (sm *SessionManager) CleanupSession(ctx context.Context, sessionID int) {
	_ = sm.StopSession(ctx, sessionID)
}

// EvaluateInSession evaluates code in a specific session
func (sm *SessionManager) EvaluateInSession(ctx context.Context, sessionID int, code string, projectPath ...string) (string, error) {
	args := map[string]any{
		"session_id": sessionID,
		"code":       code,
	}
	if len(projectPath) > 0 {
		args["project_path"] = projectPath[0]
	}
	result, err := sm.session.CallTool(ctx, "eval_in_matlab_session", args)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate in MATLAB session: %v", err)
	}
	return sm.session.GetTextContent(result)
}
