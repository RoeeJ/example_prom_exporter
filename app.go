package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	ps "github.com/mitchellh/go-ps"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	mc_status     prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		mc_status: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "process_mc_status",
			Help: "Status of the mc process.",
		}),
	}
	reg.MustRegister(m.mc_status)
	return m
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	go func() {
		for {
      updateMetrics(m)
      <-time.After(5 * time.Second)
		}
	}()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateMetrics(metrics *metrics) {
  //Reset metrics
  metrics.mc_status.Set(0)
  proc, _ := findProcess("mc")
  if proc != nil { //Process found
    metrics.mc_status.Set(1)
  } 
}

func findProcess(name string) (*ps.Process, error) {
  processes, err := ps.Processes()
  if err != nil {
    return nil, err
  }
  for _, proc := range processes {
    if proc.Executable() == name {
      return &proc, nil
    }
  }
  return nil, errors.New("Process not found")
}
