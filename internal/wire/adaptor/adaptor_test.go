// Copyright 2026 The MathWorks, Inc.

package adaptor_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/logger"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/messagecatalog"
	"github.com/matlab/matlab-mcp-core-server/internal/wire"
	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
	"github.com/stretchr/testify/assert"
)

func TestAdaptor_ModeSelector_HappyPath(t *testing.T) {
	// Arrange
	expectedModeSelector := &modeselector.ModeSelector{}

	app := adaptor.NewAdaptor(&wire.Application{
		ModeSelector: expectedModeSelector,
	})

	// Act
	result := app.ModeSelector()

	// Assert
	assert.Equal(t, expectedModeSelector, result)
}

func TestAdaptor_MessageCatalog_HappyPath(t *testing.T) {
	// Arrange
	expectedMessageCatalog := messagecatalog.New()

	app := adaptor.NewAdaptor(&wire.Application{
		MessageCatalog: expectedMessageCatalog,
	})

	// Act
	result := app.MessageCatalog()

	// Assert
	assert.Equal(t, expectedMessageCatalog, result)
}

func TestAdaptor_LoggerFactory_HappyPath(t *testing.T) {
	// Arrange
	expectedLoggerFactory := &logger.Factory{}

	app := adaptor.NewAdaptor(&wire.Application{
		LoggerFactory: expectedLoggerFactory,
	})

	// Act
	result := app.LoggerFactory()

	// Assert
	assert.Equal(t, expectedLoggerFactory, result)
}
