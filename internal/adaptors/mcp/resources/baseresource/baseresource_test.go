// Copyright 2025-2026 The MathWorks, Inc.

package baseresource_test

import (
	"context"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/baseresource"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/resources"
	baseresourcemocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/resources/baseresource"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	const (
		name        = "test_resource"
		title       = "Test Resource"
		description = "A test resource"
		mimeType    = "text/plain"
		size        = 100
		uri         = "test://resource"
	)

	mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		return &baseresource.ReadResourceResult{}, nil
	}

	// Act
	r := baseresource.New(name, title, description, mimeType, size, uri, mockLoggerFactory, handler)

	// Assert
	assert.NotNil(t, r)
	assert.Equal(t, name, r.Name())
	assert.Equal(t, title, r.Title())
	assert.Equal(t, description, r.Description())
	assert.Equal(t, mimeType, r.MimeType())
	assert.Equal(t, int64(size), r.Size())
	assert.Equal(t, uri, r.URI())
}

func TestResource_AddToServer_InvalidMimeType(t *testing.T) {
	tests := []struct {
		name             string
		mimeType         string
		expectedErrorMsg string
	}{
		{
			name:             "empty string",
			mimeType:         "",
			expectedErrorMsg: "invalid MIME type: empty string",
		},
		{
			name:             "missing slash",
			mimeType:         "invalid-mime-type",
			expectedErrorMsg: "must be in format type/subtype",
		},
		{
			name:             "empty type",
			mimeType:         "/subtype",
			expectedErrorMsg: "type and subtype cannot be empty",
		},
		{
			name:             "empty subtype",
			mimeType:         "type/",
			expectedErrorMsg: "type and subtype cannot be empty",
		},
		{
			name:             "multiple slashes",
			mimeType:         "type/sub/type",
			expectedErrorMsg: "must be in format type/subtype",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
			defer mockLoggerFactory.AssertExpectations(t)

			handler := func(ctx context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
				return &baseresource.ReadResourceResult{}, nil
			}

			r := baseresource.New("test_resource", "Test Resource", "A test resource", tt.mimeType, 100, "test://resource", mockLoggerFactory, handler)

			mockServer := &mocks.MockServer{}
			defer mockServer.AssertExpectations(t)

			// Act
			err := r.AddToServer(mockServer)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrorMsg)
		})
	}
}

func TestResource_AddToServer_HappyPath(t *testing.T) {
	// Arrange
	const (
		name        = "test_resource"
		title       = "Test Resource"
		description = "A test resource"
		mimeType    = "text/plain"
		size        = 100
		uri         = "test://resource"
	)

	mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	handler := func(ctx context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		return &baseresource.ReadResourceResult{}, nil
	}

	r := baseresource.New(name, title, description, mimeType, size, uri, mockLoggerFactory, handler)

	mockServer := &mocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockServer.EXPECT().AddResource(
		&mcp.Resource{
			Name:        name,
			Title:       title,
			Description: description,
			MIMEType:    mimeType,
			Size:        size,
			URI:         uri,
		},
		mock.AnythingOfType("mcp.ResourceHandler"),
	).Return()

	// Act
	err := r.AddToServer(mockServer)

	// Assert
	require.NoError(t, err)
}

func TestResource_ResourceHandler_HappyPath(t *testing.T) {
	// Arrange
	const (
		name        = "test_resource"
		title       = "Test Resource"
		description = "A test resource"
		mimeType    = "text/plain"
		size        = 100
		uri         = "test://resource"
	)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(mock.Anything).
		Return(mockLogger, nil).
		Once()

	expectedContents := []baseresource.ResourceContents{
		{MIMEType: "text/plain", Text: "test content"},
	}

	handler := func(ctx context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		assert.NotNil(t, logger)
		return &baseresource.ReadResourceResult{
			Contents: expectedContents,
		}, nil
	}

	r := baseresource.New(name, title, description, mimeType, size, uri, mockLoggerFactory, handler)

	var capturedHandler mcp.ResourceHandler
	mockServer := &mocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockServer.EXPECT().AddResource(
		mock.Anything,
		mock.AnythingOfType("mcp.ResourceHandler"),
	).Run(func(resource *mcp.Resource, h mcp.ResourceHandler) {
		capturedHandler = h
	}).Return()

	err := r.AddToServer(mockServer)
	require.NoError(t, err)

	// Act
	result, handlerErr := capturedHandler(t.Context(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: uri,
		},
	})

	// Assert
	require.NoError(t, handlerErr)
	require.Len(t, result.Contents, 1)
	assert.Equal(t, "text/plain", result.Contents[0].MIMEType)
	assert.Equal(t, "test content", result.Contents[0].Text)
}

func TestResource_ResourceHandler_HandlerError(t *testing.T) {
	// Arrange
	const (
		name        = "test_resource"
		title       = "Test Resource"
		description = "A test resource"
		mimeType    = "text/plain"
		size        = 100
		uri         = "test://resource"
	)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(mock.Anything).
		Return(mockLogger, nil).
		Once()

	expectedError := assert.AnError

	handler := func(ctx context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		return nil, expectedError
	}

	r := baseresource.New(name, title, description, mimeType, size, uri, mockLoggerFactory, handler)

	var capturedHandler mcp.ResourceHandler
	mockServer := &mocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockServer.EXPECT().AddResource(
		mock.Anything,
		mock.AnythingOfType("mcp.ResourceHandler"),
	).Run(func(resource *mcp.Resource, h mcp.ResourceHandler) {
		capturedHandler = h
	}).Return()

	err := r.AddToServer(mockServer)
	require.NoError(t, err)

	// Act
	result, handlerErr := capturedHandler(t.Context(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: uri,
		},
	})

	// Assert
	require.ErrorIs(t, handlerErr, expectedError)
	assert.Nil(t, result)
}

func TestResource_ResourceHandler_NilHandler(t *testing.T) {
	// Arrange
	const (
		name        = "test_resource"
		title       = "Test Resource"
		description = "A test resource"
		mimeType    = "text/plain"
		size        = 100
		uri         = "test://resource"
	)

	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(mock.Anything).
		Return(mockLogger, nil).
		Once()

	r := baseresource.New(name, title, description, mimeType, size, uri, mockLoggerFactory, nil)

	var capturedHandler mcp.ResourceHandler
	mockServer := &mocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockServer.EXPECT().AddResource(
		mock.Anything,
		mock.AnythingOfType("mcp.ResourceHandler"),
	).Run(func(resource *mcp.Resource, h mcp.ResourceHandler) {
		capturedHandler = h
	}).Return()

	err := r.AddToServer(mockServer)
	require.NoError(t, err)

	// Act
	result, handlerErr := capturedHandler(t.Context(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: uri,
		},
	})

	// Assert
	require.Error(t, handlerErr)
	assert.Contains(t, handlerErr.Error(), baseresource.UnexpectedErrorPrefix)
	assert.Nil(t, result)
}

func TestResource_ResourceHandler_NewMCPSessionLoggerError(t *testing.T) {
	// Arrange
	const (
		name        = "test_resource"
		title       = "Test Resource"
		description = "A test resource"
		mimeType    = "text/plain"
		size        = 100
		uri         = "test://resource"
	)

	mockLoggerFactory := &baseresourcemocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedError := messages.AnError

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(mock.Anything).
		Return(nil, expectedError).
		Once()

	handler := func(ctx context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		return &baseresource.ReadResourceResult{}, nil
	}

	r := baseresource.New(name, title, description, mimeType, size, uri, mockLoggerFactory, handler)

	var capturedHandler mcp.ResourceHandler
	mockServer := &mocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockServer.EXPECT().AddResource(
		mock.Anything,
		mock.AnythingOfType("mcp.ResourceHandler"),
	).Run(func(resource *mcp.Resource, h mcp.ResourceHandler) {
		capturedHandler = h
	}).Return()

	err := r.AddToServer(mockServer)
	require.NoError(t, err)

	// Act
	result, handlerErr := capturedHandler(t.Context(), &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: uri,
		},
	})

	// Assert
	require.ErrorIs(t, handlerErr, expectedError)
	assert.Nil(t, result)
}
