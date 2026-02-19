// Copyright 2025-2026 The MathWorks, Inc.

package checkmatlabcode

const (
	name        = "check_matlab_code"
	title       = "Check MATLAB Code"
	description = "Perform static code analysis on a MATLAB script (`script_path`) using MATLAB's built-in Code Analyzer function in an existing MATLAB session. Returns warnings about coding style, potential errors, deprecated functions, performance issues, and best practice violations. It also includes information about where each issue occurs and how it can be fixed in MATLAB. This is a non-destructive, read-only operation that helps identify code quality issues without executing the script."
)

type Args struct {
	ScriptPath string `json:"script_path" jsonschema:"The full absolute path to the MATLAB script file to analyze - Must be a .m file that exists - File is not modified during analysis - Example: C:\\Users\\username\\matlab\\myFunction.m or /home/user/scripts/analysis.m."`
}

type ReturnArgs struct {
	CodeIssues []CodeIssue `json:"code_issues" jsonschema:"Detailed information about each code issue including location, severity, and fixability."`
}

type CodeIssue struct {
	Description string `json:"description"  jsonschema:"Description of the code issue."`
	Line        int    `json:"line"         jsonschema:"Line number where the issue occurs."`
	StartColumn int    `json:"start_column" jsonschema:"Starting column position of the issue."`
	EndColumn   int    `json:"end_column"   jsonschema:"Ending column position of the issue."`
	Severity    string `json:"severity"     jsonschema:"Severity level of the issue (e.g., warning, error)."`
	Fixable     bool   `json:"fixable"      jsonschema:"Whether the issue can be automatically fixed using MATLAB in-built 'fix' method."`
}
