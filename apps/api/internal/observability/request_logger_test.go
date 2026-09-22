package observability_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/observability"
)

func TestRequestLogger_GeneratesIDWhenAbsent(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	var seenID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenID = observability.RequestID(r.Context())
	})

	h := observability.RequestLogger(logger, next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if seenID == "" {
		t.Fatal("expected a generated request ID to be available in context")
	}
	if got := rec.Header().Get(observability.RequestIDHeader); got != seenID {
		t.Fatalf("expected response header to match context ID, got %q vs %q", got, seenID)
	}
}

func TestRequestLogger_ReusesInboundID(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h := observability.RequestLogger(logger, next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(observability.RequestIDHeader, "client-supplied-id")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get(observability.RequestIDHeader); got != "client-supplied-id" {
		t.Fatalf("expected inbound request ID to be reused, got %q", got)
	}
}

func TestRequestLogger_LogsStructuredLineWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	h := observability.RequestLogger(logger, next)

	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var logLine map[string]any
	line := strings.TrimSpace(buf.String())
	if err := json.Unmarshal([]byte(line), &logLine); err != nil {
		t.Fatalf("expected valid JSON log line, got %q: %v", line, err)
	}

	for _, field := range []string{"request_id", "method", "path", "status", "duration_ms"} {
		if _, ok := logLine[field]; !ok {
			t.Fatalf("expected log line to contain %q, got %+v", field, logLine)
		}
	}
	if logLine["status"] != float64(http.StatusCreated) {
		t.Fatalf("expected status 201 in log line, got %v", logLine["status"])
	}
}
