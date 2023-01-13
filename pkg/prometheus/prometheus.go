package prometheus

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServePrometheusHandler(prometheusConfig *Config) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	s := &http.Server{
		Addr:        prometheusConfig.ListenAddress,
		Handler:     mux,
		ReadTimeout: time.Duration(prometheusConfig.ReadTimeout) * time.Second,
	}

	return s.ListenAndServe()
}
