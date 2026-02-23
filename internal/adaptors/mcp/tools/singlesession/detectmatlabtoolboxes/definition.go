// Copyright 2025-2026 The MathWorks, Inc.

package detectmatlabtoolboxes

const (
	name        = "detect_matlab_toolboxes"
	title       = "Detect MATLAB Toolboxes"
	description = "Returns information about installed MATLAB and toolboxes, including version numbers."
)

type Args struct {
}

type ReturnArgs struct {
	InstallationInfo string `json:"installation_info" jsonschema:"MATLAB installation information including MATLAB version and installed toolboxes with their versions."`
}
