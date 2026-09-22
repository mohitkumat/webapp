package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// RequestIDHeader is the header a client can set to propagate its own request
// ID (e.g. from the frontend), or that the API will set itself if absent.
const RequestIDHeader = "X-Request-Id"

// RequestID reads the request ID from context, set by the RequestLogger
// middleware. Handlers use this to attach the same ID to any log lines they
// emit, so one request's full path is grep-able by a single ID.
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// responseRecorder captures the status code so we can log it after the
// handler runs, since http.ResponseWriter doesn't expose it.
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// RequestLogger assigns a request ID (reusing an inbound one if the client
// sent one), logs one structured line per request with method, path, status,
// and duration, and makes the request ID available to handlers via
// RequestID(ctx). This is the "trace" for local development: grep any log
// file for one request_id and you have the full path through the system.
func RequestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set(RequestIDHeader, id)

		ctx := context.WithValue(r.Context(), requestIDKey, id)
		r = r.WithContext(ctx)

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(rec, r)

		logger.Info("request",
			"request_id", id,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func newRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
