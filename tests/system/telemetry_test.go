// Copyright 2026 The MathWorks, Inc.

//go:build linux

package system_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/tests/testutils/otel"
	"github.com/stretchr/testify/suite"
)

// TelemetryTestSuite tests that the telemetry system is working correctly.
type TelemetryTestSuite struct {
	SystemTestSuite
}

func (s *TelemetryTestSuite) TestTelemetry() {
	// Start an OTEL Collector fixture
	otelCollector := otel.StartCollector(s.T(), otel.DefaultConfig())
	defer otelCollector.Stop(s.T())

	// Start the server, with a basic --version, pointing to the above collector
	cmd := exec.Command(s.mcpServerPath, //nolint:gosec // Trusted test path
		"--version",
		"--telemetry-collector-endpoint="+otelCollector.Endpoint(),
	)
	_, err := cmd.CombinedOutput()
	s.Require().NoError(err, "version flag should execute successfully")

	// Wait for the collector to receive and flush telemetry to disk
	telemetryTimeout := 30 * time.Second
	otelCollector.WaitForTelemetry(s.T(), telemetryTimeout)

	// Stop the collector
	otelCollector.Stop(s.T())

	// Check the telemetry was indeed written
	metrics, err := otelCollector.ReadTelemetry(s.T())
	s.Require().NoError(err)

	s.Require().Equal(1, metrics.MetricCount())

	// Check some additional details, to make sure we sending the correct telemetry
	resourceMetric := metrics.ResourceMetrics().At(0)

	resourceAttrs := resourceMetric.Resource().Attributes()

	serviceName, exists := resourceAttrs.Get("service.name")
	s.True(exists, "service.name resource attribute should exist")
	s.Equal("matlab-mcp-server", serviceName.Str())

	serviceVersion, exists := resourceAttrs.Get("service.version")
	s.True(exists, "service.version resource attribute should exist")
	s.NotEmpty(serviceVersion.Str(), "service.version resource attribute should not be empty")

	metric := resourceMetric.
		ScopeMetrics().
		At(0).
		Metrics().
		At(0)

	s.Equal("server.starts", metric.Name())
	s.Equal("{start}", metric.Unit())

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

func TestTelemetrySuite(t *testing.T) {
	suite.Run(t, new(TelemetryTestSuite))
}
