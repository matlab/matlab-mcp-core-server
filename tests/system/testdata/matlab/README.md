# Test Data - MATLAB Files

---

This README is intended for MathWorks® developers only.

---

Test assets for validating that MCP server tools execute MATLAB code correctly.

## Self-Validating Test Pattern

MATLAB files validate their own correctness and report final status. Go tests verify the success markers appeared, not the implementation details.

**Why this pattern?**
- **Maintainable**: Change MATLAB logic without breaking Go tests
- **Separation**: MATLAB validates correctness, Go validates integration
- **Clear**: Explicit success/failure reporting
- **Robust**: Survives MATLAB version changes

**How it works:**
1. MATLAB files wrap operations in try/catch blocks
2. MATLAB files use assertions to validate correctness
3. MATLAB files output explicit success markers (e.g., `"ALL TESTS PASSED"`)
4. Go tests check for these markers using assertion helpers

**Example:**
```matlab
% MATLAB script runs code and reports status
try
    x = 1:10;
    mean_val = mean(x);
    assert(mean_val > 0, 'Mean should be positive');
    fprintf('  ✓ Statistics calculation\n');
catch e
    fprintf('  ✗ Statistics calculation: %s\n', e.message);
    failed_tests{end+1} = 'Statistics calculation';
end
fprintf('test_script.m: ALL TESTS PASSED\n');
```

```go
// Go test verifies success marker appeared
testdata.TestScript.Assert(t, output)
```

## Go Package Structure

The parent `testdata/` directory contains Go files that work with these MATLAB assets:

- **`expectations.go`**: Pure data definitions (expected strings, patterns)
- **`assertions.go`**: Reusable assertion helpers that use the expectation data
- **`assets.go`**: Embedded filesystem for compiling MATLAB files into test binaries

## Maintenance Guidelines

**When to modify:**
- Adding new MATLAB capability coverage
- Fixing bugs in test code
- Improving robustness

**Safe changes:**
- ✅ Modify data values/ranges
- ✅ Add new test files
- ✅ Update assertions/validations
- ✅ Improve error messages

**Breaking changes:**
- ❌ Remove or rename success markers (breaks Go assertions)
- ❌ Change output format that Go tests depend on

**Design constraints:**
- No external dependencies (base MATLAB only)
- Fast execution, deterministic results
- MATLAB version-agnostic

---
Copyright 2025 The MathWorks, Inc.
