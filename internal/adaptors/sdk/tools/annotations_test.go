// Copyright 2026 The MathWorks, Inc.

package tools_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/sdk/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReadOnlyAnnotations_HappyPath(t *testing.T) {
	// Act
	annotations := tools.NewReadOnlyAnnotations()

	// Assert
	require.IsType(t, tools.ReadOnlyAnnotation{}, annotations)

	result := annotations.ToToolAnnotations()
	require.NotNil(t, result)
	assert.True(t, result.ReadOnlyHint)
	require.NotNil(t, result.DestructiveHint)
	assert.False(t, *result.DestructiveHint)
	assert.False(t, result.IdempotentHint)
	require.NotNil(t, result.OpenWorldHint)
	assert.False(t, *result.OpenWorldHint)
}

func TestNewDestructiveAnnotations_HappyPath(t *testing.T) {
	// Act
	annotations := tools.NewDestructiveAnnotations()

	// Assert
	require.IsType(t, tools.DestructiveAnnotation{}, annotations)

	result := annotations.ToToolAnnotations()
	require.NotNil(t, result)
	assert.False(t, result.ReadOnlyHint)
	require.NotNil(t, result.DestructiveHint)
	assert.True(t, *result.DestructiveHint)
	assert.False(t, result.IdempotentHint)
	require.NotNil(t, result.OpenWorldHint)
	assert.True(t, *result.OpenWorldHint)
}

func TestNewIdempotentWriteAnnotations_HappyPath(t *testing.T) {
	// Act
	annotations := tools.NewIdempotentWriteAnnotations()

	// Assert
	require.IsType(t, tools.IdempotentWriteAnnotation{}, annotations)

	result := annotations.ToToolAnnotations()
	require.NotNil(t, result)
	assert.False(t, result.ReadOnlyHint)
	require.NotNil(t, result.DestructiveHint)
	assert.True(t, *result.DestructiveHint)
	assert.True(t, result.IdempotentHint)
	require.NotNil(t, result.OpenWorldHint)
	assert.False(t, *result.OpenWorldHint)
}

func TestNewReadOnlyOpenWorldAnnotations_HappyPath(t *testing.T) {
	// Act
	annotations := tools.NewReadOnlyOpenWorldAnnotations()

	// Assert
	require.IsType(t, tools.ReadOnlyOpenWorldAnnotation{}, annotations)

	result := annotations.ToToolAnnotations()
	require.NotNil(t, result)
	assert.True(t, result.ReadOnlyHint)
	require.NotNil(t, result.DestructiveHint)
	assert.False(t, *result.DestructiveHint)
	assert.False(t, result.IdempotentHint)
	require.NotNil(t, result.OpenWorldHint)
	assert.True(t, *result.OpenWorldHint)
}
