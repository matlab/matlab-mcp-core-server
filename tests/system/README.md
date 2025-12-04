# System Tests

---

This README is intended for MathWorks® developers only.

---

System tests validate the MATLAB MCP Core Server end-to-end through realistic workflows. These tests use actual MATLAB installations and exercise the full server stack.

## Test Philosophy

System tests focus on **realistic user scenarios** rather than exhaustive feature testing:

- **Scenario-driven**: Tests simulate actual use cases that users encounter
- **Workflow-oriented**: Multiple features tested together as they would be used in practice

Think of system tests as answering: "Does this work for a real user trying to solve a real problem?"

## Running System Tests

### Prerequisites
```bash
# Build the server
make build
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MCP_MATLAB_PATH` | Directory containing the MATLAB executable (e.g., `/usr/local/MATLAB/R2024b/bin`) | Yes |
| `MATLAB_MCP_CORE_SERVER_BUILD_DIR` | Base directory for built binaries. The test locates the server at `<base>/<os>/matlab-mcp-core-server` (e.g., `.bin/glnxa64/matlab-mcp-core-server`) | Only for `go test` |

The Makefile automatically sets `MATLAB_MCP_CORE_SERVER_BUILD_DIR` to `.bin/`, so you only need to set it when running `go test` directly.

### Run Tests
```bash
# Using make (recommended) - automatically sets binary path
export MCP_MATLAB_PATH=/path/to/matlab/bin
make system-tests

# Using go test directly - must set both variables
export MCP_MATLAB_PATH=/path/to/matlab/bin
export MATLAB_MCP_CORE_SERVER_BUILD_DIR=$(pwd)/.bin
go test -v ./tests/system/...

# Run specific suite
go test -v ./tests/system/... -run TestWorkflowSuite

# Run specific scenario
go test -v ./tests/system/... -run TestInteractiveDevelopmentWorkflow
go test -v ./tests/system/... -run TestParallelExperimentationWorkflow
```

## Adding New Tests

When adding new features, think in terms of **user scenarios** rather than individual tools:

**Ask yourself:**
- "What user problem does this solve?"
- "When would someone actually use this?"
- "Does this fit into an existing scenario or require a new one?"

**Guidelines:**
1. **Fits existing scenario**: Add to the appropriate workflow test
2. **New user scenario**: Create a new test method describing the scenario 

---

Copyright 2025 The MathWorks, Inc.

---