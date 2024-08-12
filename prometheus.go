package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tagNames = []string{
		"server",
		"scheme",
		"method",
		"hostname",
		"status",
		"uri",
		"jrpc_method",
	}

	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_request_count",
		Help: "request count",
	}, tagNames)

	requestsSizeCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_request_size_bytes",
		Help: "request size in bytes",
	}, tagNames)

	responseSizeCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_response_size_bytes",
		Help: "response size in bytes",
	}, tagNames)

	requestDurationHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ngx_request_duration_seconds",
		Help: "request serving time in seconds",
	}, tagNames)
)

func init() {
	prometheus.MustRegister(
		requestCounter,
		requestsSizeCounter,
		responseSizeCounter,
		requestDurationHistogramVec,
	)
}

func prometheusMetricsRegister(l *logEntry) {
	tags := []string{
		l.server,
		l.scheme,
		l.method,
		l.hostname,
		l.status,
		l.uri,
		l.jrpc_method,
	}

	requestCounter.WithLabelValues(tags...).Inc()

	responseSizeCounter.WithLabelValues(tags...).Add(float64(l.bytesSent))
	requestsSizeCounter.WithLabelValues(tags...).Add(float64(l.bytesReceived))
	requestDurationHistogramVec.WithLabelValues(tags...).Observe(l.duration)
}

func prometheusListener(listen string) {
	r := http.NewServeMux()

	r.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    listen,
		Handler: r,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to listen on %s: %s", listen, err)
	}
}
