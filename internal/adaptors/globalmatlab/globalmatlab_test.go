// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab"
	"github.com/stretchr/testify/assert"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	// Act
	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
		mockConfigFactory,
	)

	// Assert
	assert.NotNil(t, globalMATLABSession)
}
