// Copyright 2025 The MathWorks, Inc.

package testdata

import (
	"strings"
	"testing"

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
func AssertProblematicCodeIssues(t testing.TB, messages []string) {
	t.Helper()
	require.NotEmpty(t, ProblematicCodeIssues, "expected issues must be defined")
	allMessages := strings.Join(messages, "\n")
	assert.NotContains(t, allMessages, CleanCodeMessage, "problematic code should have issues")
	for _, expectedIssue := range ProblematicCodeIssues {
		assert.Contains(t, allMessages, expectedIssue)
	}
}

// AssertCleanCode validates that checkcode found no issues.
func AssertCleanCode(t testing.TB, messages []string) {
	t.Helper()
	assert.Contains(t, messages, CleanCodeMessage, "clean code should have no issues")
}
