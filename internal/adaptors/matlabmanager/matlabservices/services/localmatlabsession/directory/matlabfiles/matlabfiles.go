// Copyright 2025-2026 The MathWorks, Inc.

package matlabfiles

import _ "embed"

//go:embed assets/+matlab_mcp/initializeMCP.m
var initializeMCP []byte

//go:embed assets/+matlab_mcp/mcpEval.m
var mcpEval []byte

//go:embed assets/+matlab_mcp/getOrStashExceptions.m
var getOrStashExceptions []byte

type MATLABFiles struct{}

func New() MATLABFiles {
	return MATLABFiles{}
}

func (g MATLABFiles) GetAll() map[string][]byte {
	return map[string][]byte{
		"initializeMCP.m":        initializeMCP,
		"mcpEval.m":              mcpEval,
		"getOrStashExceptions.m": getOrStashExceptions,
	}
}
