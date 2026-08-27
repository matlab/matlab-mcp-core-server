// Copyright 2026 The MathWorks, Inc.

package tools

import (
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
)

type Definition struct {
	Name        string
	Description string
}

func Definitions() []Definition {
	checkCode := checkmatlabcode.New(nil, nil, nil, nil)
	detectToolboxes := detectmatlabtoolboxes.New(nil, nil, nil, nil)
	evalCode := evalmatlabcode.New(nil, nil, nil, nil, nil)
	runFile := runmatlabfile.New(nil, nil, nil, nil, nil)
	runTestFile := runmatlabtestfile.New(nil, nil, nil, nil)

	return []Definition{
		{Name: checkCode.Name(), Description: checkCode.Description()},
		{Name: detectToolboxes.Name(), Description: detectToolboxes.Description()},
		{Name: evalCode.Name(), Description: evalCode.Description()},
		{Name: runFile.Name(), Description: runFile.Description()},
		{Name: runTestFile.Name(), Description: runTestFile.Description()},
	}
}
