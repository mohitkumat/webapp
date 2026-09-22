package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/middleware"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestCORS_AllowedOrigin_SetsHeaders(t *testing.T) {
	h := middleware.CORS([]string{"http://127.0.0.1:3000"}, okHandler())

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("Origin", "http://127.0.0.1:3000")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://127.0.0.1:3000" {
		t.Fatalf("expected allow-origin header, got %q", got)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected request to reach the handler, got %d", rec.Code)
	}
}

func TestCORS_DisallowedOrigin_OmitsHeaders(t *testing.T) {
	h := middleware.CORS([]string{"http://127.0.0.1:3000"}, okHandler())

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("Origin", "http://evil.example.com")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allow-origin header for disallowed origin, got %q", got)
	}
}

func TestCORS_Preflight_ReturnsNoContentWithoutReachingHandler(t *testing.T) {
	reached := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached = true })
	h := middleware.CORS([]string{"http://127.0.0.1:3000"}, next)

	req := httptest.NewRequest(http.MethodOptions, "/items", nil)
	req.Header.Set("Origin", "http://127.0.0.1:3000")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if reached {
		t.Fatal("preflight request should not reach the wrapped handler")
	}
}
