package otel_exporter

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	otel_helpers "sqlquery_otel_exporter/otel-helpers"
	"sync"
)

type GaugeMetric struct {
	gauge   metric.Float64ObservableGauge
	metrics []Metric
	mx      sync.RWMutex
}

func NewGaugeMetric(metricName string, description string) *GaugeMetric {
	gaugeMetric := GaugeMetric{
		metrics: []Metric{},
	}
	meter := otel.GetMeterProvider().Meter("sqlquery-otel-exporter-meter")

	gauge, err := meter.Float64ObservableGauge(
		metricName,
		metric.WithDescription(description),
		metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
			metrics := gaugeMetric.getMetrics()
			for _, m := range metrics {
				o.Observe(m.value, metric.WithAttributes(m.attributes...))
			}

			return nil
		}),
	)

	if err != nil {
		otel_helpers.LogFatalError(err)
	}

	gaugeMetric.gauge = gauge

	return &gaugeMetric
}

func (h *GaugeMetric) Record(metrics []Metric) {
	h.mx.Lock()
	defer h.mx.Unlock()

	h.metrics = metrics
}

func (h *GaugeMetric) getMetrics() []Metric {
	h.mx.RLock()
	defer h.mx.RUnlock()

	return h.metrics
}
