package otel_exporter

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	otel_helpers "sqlquery_otel_exporter/otel-helpers"
)

type HistogramMetric struct {
	histogram metric.Float64Histogram
}

func NewHistogramMetric(metricName string, description string, explicitBuckets []float64) *HistogramMetric {
	meter := otel.GetMeterProvider().Meter("sqlquery-otel-exporter-meter")

	var histogram metric.Float64Histogram
	var err error

	if explicitBuckets != nil && len(explicitBuckets) > 0 {
		histogram, err = meter.Float64Histogram(
			metricName,
			metric.WithDescription(description),
			metric.WithExplicitBucketBoundaries(explicitBuckets...),
		)
	} else {
		histogram, err = meter.Float64Histogram(
			metricName,
			metric.WithDescription(description),
		)
	}

	if err != nil {
		otel_helpers.LogFatalError(err)
	}

	return &HistogramMetric{
		histogram: histogram,
	}
}

func (h *HistogramMetric) Record(metrics []Metric) {
	for _, m := range metrics {
		h.histogram.Record(context.Background(), m.value, metric.WithAttributes(m.attributes...))
	}
}
