package observability

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
)

// Metrics is a minimal, thread-safe counter/gauge registry. It exists because
// this sandbox can't reach prometheus/client_golang's dependencies (see
// docs/ADR-002-observability.md). Output is real Prometheus text exposition
// format, so it's still scrapeable by an actual Prometheus server.
type Metrics struct {
	mu       sync.Mutex
	counters map[string]float64
	gauges   map[string]float64
}

func NewMetrics() *Metrics {
	return &Metrics{
		counters: make(map[string]float64),
		gauges:   make(map[string]float64),
	}
}

// IncCounter increments a named counter by 1. Counters only ever go up
// (request counts, error counts) — that's the Prometheus convention.
func (m *Metrics) IncCounter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name]++
}

// SetGauge sets a named gauge to an absolute value (in-flight requests,
// last-request duration) — gauges can go up or down.
func (m *Metrics) SetGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
}

// Handler serves current values in Prometheus text exposition format at
// whatever path it's mounted on (conventionally /metrics).
func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		defer m.mu.Unlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		for _, name := range sortedKeys(m.counters) {
			fmt.Fprintf(w, "# TYPE %s counter\n%s %v\n", name, name, m.counters[name])
		}
		for _, name := range sortedKeys(m.gauges) {
			fmt.Fprintf(w, "# TYPE %s gauge\n%s %v\n", name, name, m.gauges[name])
		}
	}
}

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Middleware wraps a handler to count total requests, count requests by
// status class (2xx/4xx/5xx), and track the duration of the most recent
// request. Mount Metrics.Handler() separately at /metrics to read these.
func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		m.IncCounter("http_requests_total")
		switch {
		case rec.status >= 500:
			m.IncCounter("http_requests_5xx_total")
		case rec.status >= 400:
			m.IncCounter("http_requests_4xx_total")
		}
	})
}
