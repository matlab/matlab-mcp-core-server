// Copyright 2025-2026 The MathWorks, Inc.

package evalmatlabcode

const (
	name        = "evaluate_matlab_code"
	title       = "Evaluate MATLAB Code"
	description = "Evaluate a string of MATLAB code (`code`) in an existing MATLAB session. Optionally specify a project folder (`project_path`) to set as the current working folder before execution. Returns the command window output from code execution. WARNING: Do not use `restoredefaultpath` as it will remove the MCP server functions from the path and break this tool's ability to communicate with MATLAB."
)

type Args struct {
	ProjectPath string `json:"project_path,omitempty" jsonschema:"(Optional) Absolute path to the project folder. When provided, MATLAB sets this as the current working folder. If omitted, code runs in MATLAB's current working folder. Example: C:\\Users\\username\\matlab-project or /home/user/research."`
	Code        string `json:"code"                   jsonschema:"The MATLAB code to evaluate."`
}
