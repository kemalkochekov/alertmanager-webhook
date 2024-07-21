package prometheuserver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	OpsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_operations_total",
		Help: "The total number of processed operations",
	})
)

func Server() error {
	reg := prometheus.NewRegistry()

	// Register the metric with the registry
	reg.MustRegister(OpsProcessed)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    "localhost:9090",
		Handler: mux,
	}

	return server.ListenAndServe()
}
