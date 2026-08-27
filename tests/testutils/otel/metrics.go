// Copyright 2026 The MathWorks, Inc.

package otel

import (
	"bytes"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Telemetry []pmetric.Metrics

func (t Telemetry) AttributeValues(metricName, attributeKey string) []string {
	var values []string
	for _, metrics := range t {
		for _, resourceMetric := range metrics.ResourceMetrics().All() {
			for _, scopeMetric := range resourceMetric.ScopeMetrics().All() {
				for _, metric := range scopeMetric.Metrics().All() {
					// Sum().DataPoints() panics on a non-Sum metric, so skip
					// any metric that is not the named counter.
					if metric.Name() != metricName || metric.Type() != pmetric.MetricTypeSum {
						continue
					}
					for _, dataPoint := range metric.Sum().DataPoints().All() {
						if value, ok := dataPoint.Attributes().Get(attributeKey); ok {
							values = append(values, value.Str())
						}
					}
				}
			}
		}
	}
	return values
}

func (t Telemetry) ContainsString(substr string) (bool, error) {
	marshaler := &pmetric.JSONMarshaler{}
	needle := []byte(substr)
	for i, metrics := range t {
		data, err := marshaler.MarshalMetrics(metrics)
		if err != nil {
			return false, fmt.Errorf("failed to marshal telemetry export %d: %w", i, err)
		}
		if bytes.Contains(data, needle) {
			return true, nil
		}
	}
	return false, nil
}
