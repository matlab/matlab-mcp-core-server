// Copyright 2026 The MathWorks, Inc.

package functional

import (
	"github.com/stretchr/testify/suite"
)

// MockMATLABTestSuite provides common setup for functional tests
type MockMATLABTestSuite struct {
	suite.Suite
}

// SetupSuite runs once before all tests in a suite
func (s *MockMATLABTestSuite) SetupSuite() {
	// TODO: Put a mock MATLAB installation on path, and be that the only on path
}
