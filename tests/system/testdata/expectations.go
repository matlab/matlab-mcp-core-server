// Copyright 2025-2026 The MathWorks, Inc.

package testdata

// Expectation defines expected output patterns for a test file.
type Expectation struct {
	Contains    []string // Output should contain these strings
	NotContains []string // Output should NOT contain these strings
}

// TestScript expectations for test_script.m
var TestScript = Expectation{
	Contains:    []string{"test_script.m: ALL TESTS PASSED"},
	NotContains: []string{"SOME TESTS FAILED"},
}

// TestMathFunctions expectations for test_math_functions.m
var TestMathFunctions = Expectation{
	Contains: []string{"7 Passed", "0 Failed"},
}

// CheckCode expectations

// ProblematicCodeExpectedLines contains line numbers that should appear in checkcode output
// for problematic_code.m across all supported MATLAB versions.
var ProblematicCodeExpectedLines = []int{
	8,  // unused variable
	17, // missing semicolon / unused variable
	31, // preallocating warning
}
