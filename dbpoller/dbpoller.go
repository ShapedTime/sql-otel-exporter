package dbpoller

import (
	"log"
	"sqlquery_otel_exporter/dbreader"
	otel_helpers "sqlquery_otel_exporter/otel-helpers"
	"time"
)

type DbPoller struct {
	interval time.Duration
	done     chan bool
	stopped  bool
	name     string
}

func NewDbPoller(intervalSeconds int, name string) *DbPoller {
	return &DbPoller{
		interval: time.Duration(intervalSeconds) * time.Second,
		name:     name,
		done:     make(chan bool),
	}
}

func (p *DbPoller) Start(reader *dbreader.DbReader, result chan []map[string]string) {
	log.Println("Starting poller with interval", p.interval.Seconds(), "seconds")
	ticker := time.NewTicker(p.interval)

	go func() {
		defer func() { p.stopped = true }()
		for {
			select {
			case <-p.done:
				return
			case <-ticker.C:
				otel_helpers.IncDbPollerRunCounter(p.name)
				res, err := reader.Read()
				if err != nil {
					log.Println("cannot read db", err)
					otel_helpers.IncDbReadErrorCounter(p.name)
				}

				log.Println("Poller read", len(res), "rows")

				if len(res) > 0 {
					result <- res
				}
			}
		}
	}()
}

func (p *DbPoller) Stop() {
	if p.stopped {
		return
	}
	log.Println("Stopping poller")
	p.done <- true
}
