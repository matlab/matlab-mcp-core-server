// Copyright 2026 The MathWorks, Inc.

//go:build linux

package system_test

import (
	"os/exec"
	"slices"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/system/testdata"
	"github.com/matlab/matlab-mcp-server/tests/testutils/otel"
	"github.com/stretchr/testify/suite"
)

const (
	magicSquareToolName = "generate_magic_square"
	// magicSquareToolHash is sha256(magicSquareToolName)[:16], the value the
	// server must emit for the extension tool instead of the cleartext name.
	magicSquareToolHash = "5b82b619eb287ccc"
)

type TelemetryTestSuite struct {
	SystemTestSuite
}

func (s *TelemetryTestSuite) TestTelemetry_StartupTelemetry() {
	otelCollector := otel.StartCollector(s.T(), otel.DefaultConfig())
	defer otelCollector.Stop(s.T())

	// Step 1: run the server with --version, which starts and exits at once.
	cmd := exec.Command(s.mcpServerPath, //nolint:gosec // Trusted test path
		"--version",
		"--telemetry-collector-endpoint="+otelCollector.Endpoint(),
	)
	_, err := cmd.CombinedOutput()
	s.Require().NoError(err, "version flag should execute successfully")

	// Step 2: wait for the single startup flush the run emits.
	telemetryTimeout := 30 * time.Second
	exports := otelCollector.WaitForMetrics(
		s.T(),
		telemetryTimeout,
		func(t otel.Telemetry) bool {
			return len(t) >= 1
		},
		"expected startup telemetry within %s", telemetryTimeout,
	)
	otelCollector.Stop(s.T())

	// Step 3: assert the run recorded exactly one server.starts metric.
	s.Require().Len(exports, 1)
	metrics := exports[0]

	s.Require().Equal(1, metrics.MetricCount())

	resourceMetric := metrics.ResourceMetrics().At(0)

	// Resource level: who emitted the telemetry.
	resourceAttrs := resourceMetric.Resource().Attributes()

	serviceName, exists := resourceAttrs.Get("service.name")
	s.True(exists, "service.name resource attribute should exist")
	s.Equal("matlab-mcp-server", serviceName.Str())

	serviceVersion, exists := resourceAttrs.Get("service.version")
	s.True(exists, "service.version resource attribute should exist")
	s.NotEmpty(serviceVersion.Str(), "service.version resource attribute should not be empty")

	// Metric level: the counter's own identity.
	metric := resourceMetric.
		ScopeMetrics().
		At(0).
		Metrics().
		At(0)

	s.Equal("server.starts", metric.Name())
	s.Equal("{start}", metric.Unit())

	// Datapoint level: the attributes recorded against the single start event.
	attrs := metric.Sum().DataPoints().At(0).Attributes()

	serverName, exists := attrs.Get("server.name")
	s.True(exists, "server.name attribute should exist")
	s.Equal("matlab-mcp-server", serverName.Str())

	serverOS, exists := attrs.Get("server.os")
	s.True(exists, "server.os attribute should exist")
	s.Equal("linux", serverOS.Str())

	specifiedParameters, exists := attrs.Get("server.specified_parameters")
	s.True(exists, "server.specified_parameters attribute should exist")

	paramValues := specifiedParameters.Slice()
	params := make([]string, paramValues.Len())
	for i := 0; i < paramValues.Len(); i++ {
		params[i] = paramValues.At(i).Str()
	}
	s.ElementsMatch([]string{
		"TelemetryCollectorEndpoint",
		"VersionMode",
	}, params)
}

func (s *TelemetryTestSuite) TestTelemetry_ToolCallRequestTelemetry() {
	otelCollector := otel.StartCollector(s.T(), otel.DefaultConfig())
	defer otelCollector.Stop(s.T())

	// Step 1: start a server session with the magic-square extension loaded.
	ctx := s.T().Context()
	extensionFile := testdata.WriteMagicSquareExtension(s.T())
	session := s.CreateMCPSession(ctx, nil, nil,
		"--telemetry-collector-endpoint="+otelCollector.Endpoint(),
		"--extension-file="+extensionFile,
	)
	defer func() { _ = session.Close() }() // leak net if an assertion below fails early

	// Step 2: call a built-in tool and the extension tool.
	_, err := session.EvaluateCode(ctx, "2+2", s.testDataDir)
	s.Require().NoError(err, "should invoke the built-in tool")

	_, err = session.CallTool(ctx, magicSquareToolName, map[string]any{"n": float64(3)})
	s.Require().NoError(err, "should invoke the extension tool")

	// Step 3: close the session so the server flushes telemetry on shutdown.
	s.CleanupSession(session, true)

	// Step 4: wait for telemetry. The two tool calls may arrive in separate
	// flushes, so wait until both names are present, not just one flush.
	telemetryTimeout := 30 * time.Second
	exports := otelCollector.WaitForMetrics(
		s.T(),
		telemetryTimeout,
		func(t otel.Telemetry) bool {
			toolNames := t.AttributeValues("tool.calls_request", "tool.name")
			return slices.Contains(toolNames, "evaluate_matlab_code") &&
				slices.Contains(toolNames, magicSquareToolHash)
		},
		"expected both tool names in telemetry within %s", telemetryTimeout,
	)

	otelCollector.Stop(s.T())

	// Step 5: Step 4 already gated on both names arriving. What is left to prove
	// is the privacy contract — the extension tool's cleartext name must never
	// reach the collector, only its hash.
	containsCleartext, err := exports.ContainsString(magicSquareToolName)
	s.Require().NoError(err)
	s.False(containsCleartext, "cleartext extension tool name must never reach the collector")
}

func TestTelemetrySuite(t *testing.T) {
	suite.Run(t, new(TelemetryTestSuite))
}
