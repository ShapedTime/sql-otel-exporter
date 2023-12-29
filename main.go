package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sqlquery_otel_exporter/config"
	"sqlquery_otel_exporter/dbpoller"
	"sqlquery_otel_exporter/dbreader"
	otelexporter "sqlquery_otel_exporter/otel-exporter"
	"sqlquery_otel_exporter/otel-helpers"
	"syscall"
)

func main() {
	cfg, err := config.ReadConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	stoppers := []func(){}

	log.Println("Starting")
	for i, query := range cfg.Queries {
		name := query.Name
		if name == "" {
			name = fmt.Sprintf("query-%d", i)
		}

		columnTracker, err := dbreader.NewColumnTracker("./column_tracker_data/")
		if err != nil {
			log.Fatalln(err)
		}
		reader, err := dbreader.NewDbReader(query.SQL, query.TrackingColumn, query.DefaultTrackingValue, cfg.Driver, cfg.Datasource, columnTracker, name)
		if err != nil {
			log.Fatal(err)
		}

		stoppers = append(stoppers, reader.Close)

		dbResult := make(chan []map[string]string)
		poller := dbpoller.NewDbPoller(query.IntervalSeconds, name)
		poller.Start(reader, dbResult)
		stoppers = append(stoppers, poller.Stop)

		otelExporter := otelexporter.NewOtelExporter(query.MetricName, query.ValueColumn, query.Description, query.AttributeColumns, query.StaticAttributes, query.DataType, query.Buckets, query.IgnoreError)
		otelExporter.Start(dbResult)
		stoppers = append(stoppers, otelExporter.Stop)
	}

	stoppers = append(stoppers, otel_helpers.Shutdown)

	// Wait for Ctrl+C signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan // Block until a signal is received

	log.Println("Stopping")
	for _, stopper := range stoppers {
		stopper()
	}

	log.Println("Stopped")
}
