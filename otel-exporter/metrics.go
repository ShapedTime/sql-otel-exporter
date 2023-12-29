package otel_exporter

type MetricExporter interface {
	Record(metrics []Metric)
}
