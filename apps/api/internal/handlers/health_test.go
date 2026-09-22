package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/handlers"
)

func TestHealth_WhenDBReachable_ReturnsOK(t *testing.T) {
	conn := setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handlers.Health(conn)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body handlers.HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Status != "ok" || body.DB != "ok" {
		t.Fatalf("expected ok/ok, got %+v", body)
	}
}

func TestHealth_WhenDBClosed_ReturnsDegraded(t *testing.T) {
	conn := setupTestDB(t)
	conn.Close() // simulate DB unreachable

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handlers.Health(conn)(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var body handlers.HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Status != "degraded" || body.DB != "unreachable" {
		t.Fatalf("expected degraded/unreachable, got %+v", body)
	}
}
