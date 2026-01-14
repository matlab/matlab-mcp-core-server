// Copyright 2025-2026 The MathWorks, Inc.

package plaintextlivecodegeneration_test

import (
	_ "embed"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	baseresourcemocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/resources/baseresource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := baseresourcemocks.NewMockLoggerFactory(t)

	// Act
	resource := plaintextlivecodegeneration.New(mockLoggerFactory)

	// Assert
	require.NotNil(t, resource)
	assert.Equal(t, "plain_text_live_code_guidelines", resource.Name())
	assert.Equal(t, "Plain Text Live Code Generation", resource.Title())
	assert.Equal(t, "text/markdown", resource.MimeType())
	assert.Equal(t, "guidelines://plain-text-live-code", resource.URI())
}

func TestHandler_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	handler := plaintextlivecodegeneration.Handler()

	// Act
	result, err := handler(t.Context(), mockLogger)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)
	assert.NotNil(t, result.Contents[0].MIMEType)
	assert.NotNil(t, result.Contents[0].Text)
}
