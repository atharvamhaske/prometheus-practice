package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpReqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "https_request_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http-request_duration_seconds",
			Help: "HTTP request latency",
		},
		[]string{"method", "path"},
	)
)

func main() {
	//register metrics
	prometheus.MustRegister(httpReqTotal)
	prometheus.MustRegister(httpReqDuration)

	//sample APi endpoint
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from golang and prometheus"))

		duration := time.Since(start).Seconds()
		httpReqTotal.WithLabelValues(r.Method, "/", "200").Inc()
		httpReqDuration.WithLabelValues(r.Method, "/").Observe(duration)

	})
	http.Handle("/", h)
	
	http.Handle("/metrics", promhttp.Handler())
	
	log.Panicln("Go app is running on port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
