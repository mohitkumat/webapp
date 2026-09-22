package observability

import (
	"log/slog"
	"os"
)

// NewLogger returns a JSON structured logger. Every log line includes a
// timestamp, level, message, and whatever attributes (request_id, path,
// duration_ms, etc.) the caller attaches.
func NewLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}
