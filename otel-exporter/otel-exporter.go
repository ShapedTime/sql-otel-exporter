package otel_exporter

import (
	"go.opentelemetry.io/otel/attribute"
	"log"
	otel_helpers "sqlquery_otel_exporter/otel-helpers"
)

type OtelExporter struct {
	metricName       string
	valueColumn      string
	attributeColumns []string
	staticAttributes map[string]string
	description      string
	done             chan bool
	stopped          bool
	explicitBuckets  []float64
	metricExporter   MetricExporter
	ignoreError      bool
}

type Metric struct {
	name        string
	value       float64
	attributes  []attribute.KeyValue
	description string
}

type Transformer interface {
	Transform(data []map[string]string) ([]Metric, error)
}

func NewOtelExporter(metricName, valueColumn, description string, attributeColumns []string, staticAttributes map[string]string, metricType string, explicitBuckets []float64, ignoreError bool) *OtelExporter {
	var metricExporter MetricExporter

	switch {
	case metricType == "gauge":
		metricExporter = NewGaugeMetric(metricName, description)
	case metricType == "histogram":
		metricExporter = NewHistogramMetric(metricName, description, explicitBuckets)
	default:
		log.Fatal("Unknown metric type", metricType)
	}

	return &OtelExporter{
		metricName:       metricName,
		valueColumn:      valueColumn,
		attributeColumns: attributeColumns,
		staticAttributes: staticAttributes,
		description:      description,
		done:             make(chan bool),
		explicitBuckets:  explicitBuckets,
		metricExporter:   metricExporter,
		ignoreError:      ignoreError,
	}
}

func (o *OtelExporter) Start(data chan []map[string]string) {
	log.Println("Starting otel exporter")
	t := NewDefaultTransformer(o.metricName, o.valueColumn, o.attributeColumns, o.staticAttributes)

	go func() {
		defer func() { o.stopped = true }()
		for {
			select {
			case <-o.done:
				log.Println("Received done signal")
				return
			case d := <-data:
				metrics, err := t.Transform(d)
				if err != nil && !o.ignoreError {
					otel_helpers.LogFatalError(err)
				}

				o.metricExporter.Record(metrics)
			}
		}
	}()
}

func (o *OtelExporter) Stop() {
	if o.stopped {
		return
	}
	log.Println("Stopping otel exporter")
	o.done <- true
	log.Println("Stopped otel exporter")
}
