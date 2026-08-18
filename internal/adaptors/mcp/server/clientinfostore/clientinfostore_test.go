// Copyright 2026 The MathWorks, Inc.

package clientinfostore_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/mcp/server/clientinfostore"
	"github.com/matlab/matlab-mcp-server/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestClientInfoStore_New_HappyPath(t *testing.T) {
	// Arrange
	// (no setup needed)

	// Act
	store := clientinfostore.New()

	// Assert
	assert.NotNil(t, store, "ClientInfoStore instance should not be nil")
}

func TestClientInfoStore_GetClientInfo_ReturnsZeroValueWhenNoInfoSet(t *testing.T) {
	// Arrange
	store := clientinfostore.New()

	// Act
	info := store.GetClientInfo()

	// Assert
	assert.Empty(t, info.Name, "GetClientInfo should return an empty name when no info has been set")
	assert.Empty(t, info.Title, "GetClientInfo should return an empty title when no info has been set")
}

func TestClientInfoStore_SetClientInfo_ThenGetClientInfoRoundTripsBothFields(t *testing.T) {
	// Arrange
	store := clientinfostore.New()

	expectedName := "vscode"
	expectedTitle := "Visual Studio Code"

	// Act
	store.SetClientInfo(entities.MCPClientInfo{Name: expectedName, Title: expectedTitle})
	clientInfo := store.GetClientInfo()

	// Assert
	assert.Equal(t, expectedName, clientInfo.Name, "GetClientInfo should return the name that was set")
	assert.Equal(t, expectedTitle, clientInfo.Title, "GetClientInfo should return the title that was set")
}

func TestClientInfoStore_SetClientInfo_ReplacesExistingInfo(t *testing.T) {
	// Arrange
	store := clientinfostore.New()

	store.SetClientInfo(entities.MCPClientInfo{Name: "cursor", Title: "Cursor"})

	expectedName := "vscode"
	expectedTitle := "Visual Studio Code"

	// Act
	store.SetClientInfo(entities.MCPClientInfo{Name: expectedName, Title: expectedTitle})
	info := store.GetClientInfo()

	// Assert
	assert.Equal(t, expectedName, info.Name, "GetClientInfo should return the name from the second Set call")
	assert.Equal(t, expectedTitle, info.Title, "GetClientInfo should return the title from the second Set call")
}

func TestClientInfoStore_SetClientInfo_PreservesEmptyTitle(t *testing.T) {
	// Arrange
	store := clientinfostore.New()

	expectedName := "cursor"

	// Act
	store.SetClientInfo(entities.MCPClientInfo{Name: expectedName, Title: ""})
	info := store.GetClientInfo()

	// Assert
	assert.Equal(t, expectedName, info.Name, "GetClientInfo should return the name that was set")
	assert.Empty(t, info.Title, "GetClientInfo should preserve an empty title without defaulting it")
}
