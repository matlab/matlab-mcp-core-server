// Copyright 2025-2026 The MathWorks, Inc.

package testdata

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Assert validates that output matches the expectation.
func (e Expectation) Assert(t testing.TB, output string) {
	t.Helper()
	require.False(t, len(e.Contains) == 0 && len(e.NotContains) == 0,
		"expectation must have at least one Contains or NotContains entry")
	for _, expected := range e.Contains {
		assert.Contains(t, output, expected)
	}
	for _, forbidden := range e.NotContains {
		assert.NotContains(t, output, forbidden)
	}
}

// AssertProblematicCodeIssues validates that checkcode found the expected issues.
func AssertProblematicCodeIssues(t testing.TB, issues []mcpclient.CodeIssue) {
	t.Helper()
	require.NotEmpty(t, issues, "problematic code should have issues")
	require.NotEmpty(t, ProblematicCodeExpectedLines, "expected lines must be defined")

	issueLines := make([]int, len(issues))
	for i, issue := range issues {
		issueLines[i] = issue.Line
	}

	for _, expectedLine := range ProblematicCodeExpectedLines {
		assert.Contains(t, issueLines, expectedLine, "should report issue on line %d", expectedLine)
	}
}

// AssertCleanCode validates that checkcode found no issues.
func AssertCleanCode(t testing.TB, issues []mcpclient.CodeIssue) {
	t.Helper()
	assert.Empty(t, issues, "clean code should have no issues")
}
