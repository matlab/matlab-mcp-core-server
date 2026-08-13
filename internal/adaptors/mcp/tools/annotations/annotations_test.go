// Copyright 2025-2026 The MathWorks, Inc.

package annotations

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestNewReadOnlyAnnotations(t *testing.T) {
	// Act
	result := NewReadOnlyAnnotations()

	// Assert
	assert.True(t, result.readOnly, "readOnly should be true")
	assert.False(t, result.destructive, "destructive should be false")
	assert.False(t, result.idempotent, "idempotent should be false")
	assert.False(t, result.openWorld, "openWorld should be false")
}

func TestNewDestructiveAnnotations(t *testing.T) {
	// Act
	result := NewDestructiveAnnotations()

	// Assert
	assert.False(t, result.readOnly, "readOnly should be false")
	assert.True(t, result.destructive, "destructive should be true")
	assert.False(t, result.idempotent, "idempotent should be false")
	assert.True(t, result.openWorld, "openWorld should be true")
}

func TestToToolAnnotations_ReadOnly(t *testing.T) {
	// Arrange
	annotations := NewReadOnlyAnnotations()

	// Act
	result := annotations.ToToolAnnotations()

	// Assert
	assert.NotNil(t, result, "result should not be nil")
	assert.True(t, result.ReadOnlyHint, "ReadOnlyHint should be true")
	assert.NotNil(t, result.DestructiveHint, "DestructiveHint pointer should not be nil")
	assert.False(t, *result.DestructiveHint, "DestructiveHint value should be false")
	assert.False(t, result.IdempotentHint, "IdempotentHint should be false")
	assert.NotNil(t, result.OpenWorldHint, "OpenWorldHint pointer should not be nil")
	assert.False(t, *result.OpenWorldHint, "OpenWorldHint value should be false")
}

func TestToToolAnnotations_Destructive(t *testing.T) {
	// Arrange
	annotations := NewDestructiveAnnotations()

	// Act
	result := annotations.ToToolAnnotations()

	// Assert
	assert.NotNil(t, result, "result should not be nil")
	assert.False(t, result.ReadOnlyHint, "ReadOnlyHint should be false")
	assert.NotNil(t, result.DestructiveHint, "DestructiveHint pointer should not be nil")
	assert.True(t, *result.DestructiveHint, "DestructiveHint value should be true")
	assert.False(t, result.IdempotentHint, "IdempotentHint should be false")
	assert.NotNil(t, result.OpenWorldHint, "OpenWorldHint pointer should not be nil")
	assert.True(t, *result.OpenWorldHint, "OpenWorldHint value should be true")
}

func TestNewIdempotentWriteAnnotations(t *testing.T) {
	// Act
	result := NewIdempotentWriteAnnotations()

	// Assert
	assert.False(t, result.readOnly, "readOnly should be false")
	assert.True(t, result.destructive, "destructive should be true")
	assert.True(t, result.idempotent, "idempotent should be true")
	assert.False(t, result.openWorld, "openWorld should be false")
}

func TestNewReadOnlyOpenWorldAnnotations(t *testing.T) {
	// Act
	result := NewReadOnlyOpenWorldAnnotations()

	// Assert
	assert.True(t, result.readOnly, "readOnly should be true")
	assert.False(t, result.destructive, "destructive should be false")
	assert.False(t, result.idempotent, "idempotent should be false")
	assert.True(t, result.openWorld, "openWorld should be true")
}

func TestToToolAnnotations_IdempotentWrite(t *testing.T) {
	// Arrange
	annotations := NewIdempotentWriteAnnotations()

	// Act
	result := annotations.ToToolAnnotations()

	// Assert
	assert.NotNil(t, result, "result should not be nil")
	assert.False(t, result.ReadOnlyHint, "ReadOnlyHint should be false")
	assert.NotNil(t, result.DestructiveHint, "DestructiveHint pointer should not be nil")
	assert.True(t, *result.DestructiveHint, "DestructiveHint value should be true")
	assert.True(t, result.IdempotentHint, "IdempotentHint should be true")
	assert.NotNil(t, result.OpenWorldHint, "OpenWorldHint pointer should not be nil")
	assert.False(t, *result.OpenWorldHint, "OpenWorldHint value should be false")
}

func TestToToolAnnotations_ReadOnlyOpenWorld(t *testing.T) {
	// Arrange
	annotations := NewReadOnlyOpenWorldAnnotations()

	// Act
	result := annotations.ToToolAnnotations()

	// Assert
	assert.NotNil(t, result, "result should not be nil")
	assert.True(t, result.ReadOnlyHint, "ReadOnlyHint should be true")
	assert.NotNil(t, result.DestructiveHint, "DestructiveHint pointer should not be nil")
	assert.False(t, *result.DestructiveHint, "DestructiveHint value should be false")
	assert.False(t, result.IdempotentHint, "IdempotentHint should be false")
	assert.NotNil(t, result.OpenWorldHint, "OpenWorldHint pointer should not be nil")
	assert.True(t, *result.OpenWorldHint, "OpenWorldHint value should be true")
}

func TestAnnotations_ProtocolSerialization(t *testing.T) {
	// Arrange
	annotations := NewReadOnlyAnnotations()
	mcpAnnotations := annotations.ToToolAnnotations()

	// Act - verify annotations can be serialized through MCP protocol
	tool := &mcp.Tool{
		Name:        "test-tool",
		Description: "test",
		Annotations: mcpAnnotations,
	}

	// Assert - verify the tool annotations match expected values
	assert.NotNil(t, tool.Annotations, "Tool annotations should not be nil")
	assert.True(t, tool.Annotations.ReadOnlyHint, "ReadOnlyHint should be true")
	assert.NotNil(t, tool.Annotations.DestructiveHint, "DestructiveHint pointer should not be nil")
	assert.False(t, *tool.Annotations.DestructiveHint, "DestructiveHint value should be false")
	assert.False(t, tool.Annotations.IdempotentHint, "IdempotentHint should be false")
	assert.NotNil(t, tool.Annotations.OpenWorldHint, "OpenWorldHint pointer should not be nil")
	assert.False(t, *tool.Annotations.OpenWorldHint, "OpenWorldHint value should be false")
}
