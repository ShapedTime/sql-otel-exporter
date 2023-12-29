package otel_helpers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"log"
)

var (
	receivedRowsCounter metric.Int64Counter
	dbPollerRunCounter  metric.Int64Counter
	dbReadError         metric.Int64Counter
	fatalError          metric.Int64Counter
)

func initMetrics(metricProvider *sdkMetric.MeterProvider) {
	var err error
	meter := metricProvider.Meter("sqlquery-otel-exporter-meter")

	receivedRowsCounter, err = meter.Int64Counter("received_rows", metric.WithDescription("Number of received rows"))
	if err != nil {
		log.Fatal(err)
	}

	dbPollerRunCounter, err = meter.Int64Counter("dbpoller_runs", metric.WithDescription("Number of times dbpoller was run"))
	if err != nil {
		log.Fatal(err)
	}

	dbReadError, err = meter.Int64Counter("db_read_error", metric.WithDescription("Number of times db read failed"))
	if err != nil {
		log.Fatal(err)
	}

	fatalError, err = meter.Int64Counter("fatal_error", metric.WithDescription("Number of times fatal error occurred"))
	if err != nil {
		log.Fatal(err)
	}
}

func IncReceivedRowsCounter(value int64, name string) {
	receivedRowsCounter.Add(context.Background(), value, metric.WithAttributes(attribute.String("name", name)))
}

func IncDbPollerRunCounter(name string) {
	dbPollerRunCounter.Add(context.Background(), 1, metric.WithAttributes(attribute.String("name", name)))
}

func IncDbReadErrorCounter(name string) {
	dbReadError.Add(context.Background(), 1, metric.WithAttributes(attribute.String("name", name)))
}

func IncFatalErrorCounter() {
	fatalError.Add(context.Background(), 1)
}

func LogFatalError(err error) {
	IncFatalErrorCounter()
	log.Fatal(err)
}
