package observability_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/observability"
)

func TestMetrics_HandlerOutputsPrometheusFormat(t *testing.T) {
	m := observability.NewMetrics()
	m.IncCounter("http_requests_total")
	m.IncCounter("http_requests_total")
	m.SetGauge("in_flight_requests", 3)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	m.Handler()(rec, req)

	body := rec.Body.String()

	if !strings.Contains(body, "http_requests_total 2") {
		t.Fatalf("expected counter value 2, got body:\n%s", body)
	}
	if !strings.Contains(body, "in_flight_requests 3") {
		t.Fatalf("expected gauge value 3, got body:\n%s", body)
	}
	if !strings.Contains(body, "# TYPE http_requests_total counter") {
		t.Fatalf("expected TYPE comment for counter, got body:\n%s", body)
	}
}

func TestMetrics_MiddlewareCountsRequestsByStatusClass(t *testing.T) {
	m := observability.NewMetrics()

	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	serverErr := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) })

	for _, h := range []http.Handler{
		m.Middleware(ok), m.Middleware(ok), m.Middleware(notFound), m.Middleware(serverErr),
	} {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}

	rec := httptest.NewRecorder()
	m.Handler()(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()

	if !strings.Contains(body, "http_requests_total 4") {
		t.Fatalf("expected 4 total requests, got:\n%s", body)
	}
	if !strings.Contains(body, "http_requests_4xx_total 1") {
		t.Fatalf("expected 1 4xx request, got:\n%s", body)
	}
	if !strings.Contains(body, "http_requests_5xx_total 1") {
		t.Fatalf("expected 1 5xx request, got:\n%s", body)
	}
}
