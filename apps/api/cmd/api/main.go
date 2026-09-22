package main

import (
	"net/http"
	"time"

	"api/internal/config"
	"api/internal/db"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/observability"
)

func main() {
	cfg := config.Load()
	logger := observability.NewLogger()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("db connect failed", "error", err)
		panic(err)
	}
	defer conn.Close()

	if err := db.RunMigrations(conn, "migrations"); err != nil {
		logger.Error("migrations failed", "error", err)
		panic(err)
	}

	metrics := observability.NewMetrics()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health(conn))
	mux.HandleFunc("GET /items", handlers.ListItems(conn))
	mux.HandleFunc("POST /items", handlers.CreateItem(conn))
	mux.HandleFunc("GET /metrics", metrics.Handler())

	// Order matters: CORS first (short-circuits preflight before anything
	// else runs), then metrics (counts every request including errors), then
	// request logging (innermost, so duration_ms reflects the actual handler
	// time, not middleware overhead).
	var handler http.Handler = mux
	handler = observability.RequestLogger(logger, handler)
	handler = metrics.Middleware(handler)
	handler = middleware.CORS(cfg.AllowedOrigin, handler)

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("api listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server error", "error", err)
		panic(err)
	}
}
