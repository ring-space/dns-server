package metrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mu            sync.Mutex
	counts        = make(map[string]int)
	gaugeRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dns_requests_by_country",
			Help: "Number of DNS A-requests by client country in the current interval",
		},
		[]string{"country"},
	)
	flushInterval time.Duration
)

// Init registers metrics and starts HTTP server and flusher
func Init(addr string, interval time.Duration) {
	flushInterval = interval
	prometheus.MustRegister(gaugeRequests)
	go startHTTPServer(addr)
	go startFlusher()
}

// Record increments the counter for a given country
func Record(country string) {
	mu.Lock()
	counts[country]++
	mu.Unlock()
}

// startHTTPServer serves /metrics on the given address
func startHTTPServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, nil)
}

// startFlusher periodically flushes counts into Prometheus and resets them
func startFlusher() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()
	for range ticker.C {
		flush()
	}
}

// flush updates Prometheus gauges and resets internal counters
func flush() {
	mu.Lock()
	defer mu.Unlock()
	for country, count := range counts {
		gaugeRequests.WithLabelValues(country).Set(float64(count))
		counts[country] = 0
	}
}
