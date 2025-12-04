// Copyright 2025 The MathWorks, Inc.

package mcpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPClient wraps MCP session with helper methods
// Lifecycle: NewClient -> CreateSession -> [operations] -> Close
type MCPClient struct {
	client    *mcp.Client
	transport *mcp.CommandTransport
	stderr    *bytes.Buffer
}

// MCPClientSession represents an active MCP session
type MCPClientSession struct {
	session *mcp.ClientSession
	stderr  *bytes.Buffer
}

func GetMCPClientImplementation() *mcp.Implementation {
	// Those values don't matter for the system tests, but are required to construct an MCP client.
	return &mcp.Implementation{
		Name:    "test-client-for-matlab-mcp-core-server",
		Version: "test-client-for-matlab-mcp-core-server",
	}
}

// NewClient creates a new MCP client
func NewClient(ctx context.Context, serverPath string, env []string, args ...string) *MCPClient {
	cmd := exec.CommandContext(ctx, serverPath, args...)
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	if env != nil {
		cmd.Env = env
	}
	// Use a longer terminate duration to allow the server to gracefully shut down
	// MATLAB sessions, which can take some time.
	transport := &mcp.CommandTransport{
		Command:           cmd,
		TerminateDuration: 3 * time.Minute,
	}
	client := mcp.NewClient(GetMCPClientImplementation(), nil)

	return &MCPClient{
		client:    client,
		transport: transport,
		stderr:    stderr,
	}
}

func (c *MCPClient) CreateSession(ctx context.Context) (*MCPClientSession, error) {
	session, err := c.client.Connect(ctx, c.transport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client session: %w", err)
	}
	return &MCPClientSession{
		session: session,
		stderr:  c.stderr,
	}, nil
}

func (s *MCPClientSession) Close() error {
	err := s.session.Close()
	if err != nil {
		return fmt.Errorf("session close failed: %w\n\nServer stderr:\n%s", err, s.stderr.String())
	}
	return nil
}

// Stderr returns the captured stderr output from the MCP server process.
// This is useful for debugging test failures.
func (s *MCPClientSession) Stderr() string {
	return s.stderr.String()
}

// CallTool calls an MCP tool and asserts it doesn't error
func (s *MCPClientSession) CallTool(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error) {
	result, err := s.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", name, err)
	}

	if result.IsError {
		// Log the error details for better debugging
		errorMsg := "unknown error"
		if textContent, err := s.GetTextContent(result); err == nil {
			errorMsg = textContent
		}
		return nil, fmt.Errorf("tool %s returned an error: %s", name, errorMsg)
	}
	return result, nil
}

// GetTextContent extracts text content from a tool result
func (*MCPClientSession) GetTextContent(result *mcp.CallToolResult) (string, error) {
	if result.Content == nil {
		return "", fmt.Errorf("result content is nil")
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("result should have content")
	}
	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		return "", fmt.Errorf("content should be TextContent")
	}
	return textContent.Text, nil
}

// UnmarshalStructuredContent unmarshals structured content into a target
func (*MCPClientSession) UnmarshalStructuredContent(result *mcp.CallToolResult, target interface{}) error {
	jsonBytes, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return fmt.Errorf("failed to marshal structured content: %w", err)
	}
	err = json.Unmarshal(jsonBytes, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal structured content: %w", err)
	}
	return nil
}

// EvaluateCode evaluates MATLAB code
func (s *MCPClientSession) EvaluateCode(ctx context.Context, code string, projectPath ...string) (string, error) {
	args := map[string]any{"code": code}
	if len(projectPath) > 0 {
		args["project_path"] = projectPath[0]
	}
	result, err := s.CallTool(ctx, "evaluate_matlab_code", args)
	if err != nil {
		return "", err
	}
	return s.GetTextContent(result)
}

// CheckCode checks MATLAB code
func (s *MCPClientSession) CheckCode(ctx context.Context, scriptPath string) ([]string, error) {
	result, err := s.CallTool(ctx, "check_matlab_code", map[string]any{
		"script_path": scriptPath,
	})
	if err != nil {
		return nil, err
	}
	var output struct {
		CheckCodeMessages []string `json:"checkcode_messages"`
	}
	err = s.UnmarshalStructuredContent(result, &output)
	if err != nil {
		return nil, err
	}
	return output.CheckCodeMessages, nil
}

// RunFile runs a MATLAB file
func (s *MCPClientSession) RunFile(ctx context.Context, scriptPath string) (string, error) {
	result, err := s.CallTool(ctx, "run_matlab_file", map[string]any{
		"script_path": scriptPath,
	})
	if err != nil {
		return "", err
	}
	return s.GetTextContent(result)
}

// RunTestFile runs a MATLAB test file
func (s *MCPClientSession) RunTestFile(ctx context.Context, scriptPath string) (string, error) {
	result, err := s.CallTool(ctx, "run_matlab_test_file", map[string]any{
		"script_path": scriptPath,
	})
	if err != nil {
		return "", err
	}
	return s.GetTextContent(result)
}

// DetectToolboxes detects installed MATLAB toolboxes
func (s *MCPClientSession) DetectToolboxes(ctx context.Context) (string, error) {
	result, err := s.CallTool(ctx, "detect_matlab_toolboxes", map[string]any{})
	if err != nil {
		return "", err
	}
	var output struct {
		InstallationInfo string `json:"installation_info"`
	}
	err = s.UnmarshalStructuredContent(result, &output)
	if err != nil {
		return "", err
	}
	return output.InstallationInfo, nil
}

// NewSessionManager creates a new session manager for multi-session workflows
func (s *MCPClientSession) NewSessionManager() *SessionManager {
	return &SessionManager{
		session: s,
	}
}
