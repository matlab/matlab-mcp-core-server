// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import (
	"errors"
	"fmt"
	"strings"
)

type matlabError struct {
	message string
}

func newMATLABError(message string) matlabError {
	return matlabError{
		message: message,
	}
}

func (e matlabError) Error() string {
	return fmt.Sprintf("matlab error: %v", e.message)
}

var ErrMCPPackageUnavailable = errors.New("the MATLAB session's path no longer includes the MCP server functions (e.g. due to restoredefaultpath). The MCP server cannot evaluate code with output capture in this state")

func isMCPPackageUnavailableError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "Undefined function") && strings.Contains(msg, "matlab_mcp")
}
