// Copyright 2025-2026 The MathWorks, Inc.

package evalmatlabcode

const (
	name        = "eval_in_matlab_session"
	title       = "Evaluate MATLAB Code in a MATLAB Session"
	description = "Evaluate a string of MATLAB code (`code`) in an existing MATLAB session, given its session ID (`session_id`). Optionally specify a project folder (`project_path`) to set as the working folder before execution. WARNING: Do not use `restoredefaultpath` as it will remove the MCP server functions from the path and break this tool's ability to communicate with MATLAB."
)

type Args struct {
	SessionID   int    `json:"session_id"             jsonschema:"The ID of the MATLAB session in which to evaluate the code."`
	ProjectPath string `json:"project_path,omitempty" jsonschema:"(Optional) Absolute path to the project folder. When provided, MATLAB sets this as the current working folder. If omitted, code runs in MATLAB's current working folder. Example: C:\\Users\\username\\matlab-project or /home/user/research."`
	Code        string `json:"code"                   jsonschema:"The MATLAB code to evaluate."`
}
