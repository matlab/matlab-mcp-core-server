// Copyright 2025-2026 The MathWorks, Inc.

package logger_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/logger"
	loggermocks "github.com/matlab/matlab-mcp-server/mocks/adaptors/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMultiHandler_HappyPath(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	// Act
	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Assert
	assert.NotNil(t, multiHandler, "MultiHandler instance should not be nil")
}

func TestNewMultiHandler_NoHandlers(t *testing.T) {
	// Arrange & Act
	multiHandler := logger.NewMultiHandler()

	// Assert
	assert.NotNil(t, multiHandler, "MultiHandler instance should not be nil even with no handlers")
}

func TestSlogMultiHandler_Enabled_AnyHandlerEnabled_ReturnsTrue(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	ctx := t.Context()
	expectedLevel := slog.LevelInfo

	mockHandler1.EXPECT().
		Enabled(ctx, expectedLevel).
		Return(false).
		Once()

	mockHandler2.EXPECT().
		Enabled(ctx, expectedLevel).
		Return(true).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	result := multiHandler.Enabled(ctx, expectedLevel)

	// Assert
	assert.True(t, result, "Should return true when any handler is enabled")
}

func TestSlogMultiHandler_Enabled_NoHandlersEnabled_ReturnsFalse(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	ctx := t.Context()
	expectedLevel := slog.LevelInfo

	mockHandler1.EXPECT().
		Enabled(ctx, expectedLevel).
		Return(false).
		Once()

	mockHandler2.EXPECT().
		Enabled(ctx, expectedLevel).
		Return(false).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	result := multiHandler.Enabled(ctx, expectedLevel)

	// Assert
	assert.False(t, result, "Should return false when no handlers are enabled")
}

func TestSlogMultiHandler_Enabled_NoHandlers_ReturnsFalse(t *testing.T) {
	// Arrange
	ctx := t.Context()
	level := slog.LevelInfo

	multiHandler := logger.NewMultiHandler()

	// Act
	result := multiHandler.Enabled(ctx, level)

	// Assert
	assert.False(t, result, "Should return false when no handlers exist")
}

func TestSlogMultiHandler_Handle_AllHandlersEnabled_CallsAll(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	ctx := t.Context()
	expectedRecord := slog.Record{}
	expectedRecord.Level = slog.LevelInfo

	mockHandler1.EXPECT().
		Enabled(ctx, expectedRecord.Level).
		Return(true).
		Once()

	mockHandler1.EXPECT().
		Handle(ctx, expectedRecord).
		Return(nil).
		Once()

	mockHandler2.EXPECT().
		Enabled(ctx, expectedRecord.Level).
		Return(true).
		Once()

	mockHandler2.EXPECT().
		Handle(ctx, expectedRecord).
		Return(nil).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	err := multiHandler.Handle(ctx, expectedRecord)

	// Assert
	require.NoError(t, err, "Handle should not return an error")
}

func TestSlogMultiHandler_Handle_SomeHandlersDisabled_OnlyCallsEnabled(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	ctx := t.Context()
	expectedRecord := slog.Record{}
	expectedRecord.Level = slog.LevelInfo

	mockHandler1.EXPECT().
		Enabled(ctx, expectedRecord.Level).
		Return(false).
		Once()

	mockHandler2.EXPECT().
		Enabled(ctx, expectedRecord.Level).
		Return(true).
		Once()

	mockHandler2.EXPECT().
		Handle(ctx, expectedRecord).
		Return(nil).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	err := multiHandler.Handle(ctx, expectedRecord)

	// Assert
	require.NoError(t, err, "Handle should not return an error")
}

func TestSlogMultiHandler_Handle_FirstHandlerError_ReturnsFirstError(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	ctx := t.Context()
	expectedRecord := slog.Record{}
	expectedRecord.Level = slog.LevelInfo

	expectedError := assert.AnError

	mockHandler1.EXPECT().
		Enabled(ctx, expectedRecord.Level).
		Return(true).
		Once()

	mockHandler1.EXPECT().
		Handle(ctx, expectedRecord).
		Return(expectedError).
		Once()

	mockHandler2.EXPECT().
		Enabled(ctx, expectedRecord.Level).
		Return(true).
		Once()

	mockHandler2.EXPECT().
		Handle(ctx, expectedRecord).
		Return(nil).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	err := multiHandler.Handle(ctx, expectedRecord)

	// Assert
	assert.ErrorIs(t, err, expectedError, "Should return the first error encountered")
}

func TestSlogMultiHandler_Handle_NoHandlers_ReturnsNil(t *testing.T) {
	// Arrange
	ctx := t.Context()
	record := slog.Record{}
	record.Level = slog.LevelInfo

	multiHandler := logger.NewMultiHandler()

	// Act
	err := multiHandler.Handle(ctx, record)

	// Assert
	require.NoError(t, err, "Handle should not return an error when no handlers exist")
}

func TestSlogMultiHandler_WithAttrs_CallsAllHandlers(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	expectedAttrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.String("key2", "value2"),
	}

	mockHandler1.EXPECT().
		WithAttrs(expectedAttrs).
		Return(mockHandler1).
		Once()

	mockHandler2.EXPECT().
		WithAttrs(expectedAttrs).
		Return(mockHandler2).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	result := multiHandler.WithAttrs(expectedAttrs)

	// Assert
	assert.NotSame(t, multiHandler, result, "WithAttrs should not return the same multiHandler instance")
}

func TestSlogMultiHandler_WithAttrs_NoHandlers_DoesNotReturnsItself(t *testing.T) {
	// Arrange
	attrs := []slog.Attr{
		slog.String("key1", "value1"),
	}

	multiHandler := logger.NewMultiHandler()

	// Act
	result := multiHandler.WithAttrs(attrs)

	// Assert
	assert.NotEqual(t, multiHandler, result, "WithAttrs should return a new handler instance")
}

func TestSlogMultiHandler_WithAttrs_DoesNotMutateParent(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	child := slog.NewJSONHandler(&buf, nil)
	multiHandler := logger.NewMultiHandler(child)
	parent := slog.New(multiHandler)

	// Act - derive a child logger, then log through the parent
	_ = parent.With("child-attr", "child-value")
	parent.Info("parent-message")

	// Assert - the parent's output must not carry the derived logger's attribute
	assert.NotContains(t, buf.String(), "child-attr")
}

func TestSlogMultiHandler_WithGroup_CallsAllHandlers(t *testing.T) {
	// Arrange
	mockHandler1 := &loggermocks.MockHandler{}
	defer mockHandler1.AssertExpectations(t)

	mockHandler2 := &loggermocks.MockHandler{}
	defer mockHandler2.AssertExpectations(t)

	expectedGroupName := "test-group"

	mockHandler1.EXPECT().
		WithGroup(expectedGroupName).
		Return(mockHandler1).
		Once()

	mockHandler2.EXPECT().
		WithGroup(expectedGroupName).
		Return(mockHandler2).
		Once()

	multiHandler := logger.NewMultiHandler(mockHandler1, mockHandler2)

	// Act
	result := multiHandler.WithGroup(expectedGroupName)

	// Assert
	assert.NotSame(t, multiHandler, result, "WithGroup should return a new handler instance")
}

func TestSlogMultiHandler_WithGroup_NoHandlers_DoesNotReturnsItself(t *testing.T) {
	// Arrange
	groupName := "test-group"

	multiHandler := logger.NewMultiHandler()

	// Act
	result := multiHandler.WithGroup(groupName)

	// Assert
	assert.NotSame(t, multiHandler, result, "WithGroup should return a new handler instance")
}
