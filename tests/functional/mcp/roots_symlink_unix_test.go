// Copyright 2026 The MathWorks, Inc.

//go:build !windows

package mcp_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-server/tests/testutils/mcpclient"
	"github.com/stretchr/testify/suite"
)

// RootsSymlinkTestSuite tests that an MCP root passed as a symlink is matched
// against MATLAB's resolved working directory. Building the symlink
// in-test reproduces the macOS /var vs /private/var mismatch on any Unix host.
//
// Unix-only: the divergence has no Windows equivalent and os.Symlink needs
// elevated privilege there, so the file is excluded rather than skipped.
type RootsSymlinkTestSuite struct {
	MCPTestSuite
}

func TestRootsSymlinkTestSuite(t *testing.T) {
	suite.Run(t, new(RootsSymlinkTestSuite))
}

func (s *RootsSymlinkTestSuite) TestMATLABStartsInSymlinkedRoot_ResolvesToRealPath() {
	// Arrange
	base := s.T().TempDir()
	realDir := filepath.Join(base, "real-dir")
	s.Require().NoError(os.Mkdir(realDir, 0o750))

	linkDir := filepath.Join(base, "link-dir")
	s.Require().NoError(os.Mkdir(linkDir, 0o750))
	linkPath := filepath.Join(linkDir, "root-link")
	s.Require().NoError(os.Symlink(realDir, linkPath))

	// Act
	_ = s.CreateSession(
		mcpclient.WithRoots(NewRootFromDir(linkPath, "workspace")),
	)

	// Assert
	startupInfo := s.WaitForStartupInfo()
	s.Equal(
		s.normalizedPath(realDir),
		s.normalizedPath(startupInfo.WorkingDir),
		"MATLAB should start in the symlink-resolved MCP root directory",
	)
}
