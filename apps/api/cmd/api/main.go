package main

import (
	"log"
	"net/http"
	"time"

	"api/internal/config"
	"api/internal/db"
	"api/internal/handlers"
	"api/internal/middleware"
)

func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer conn.Close()

	if err := db.RunMigrations(conn, "migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health(conn))
	mux.HandleFunc("GET /items", handlers.ListItems(conn))
	mux.HandleFunc("POST /items", handlers.CreateItem(conn))

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      middleware.CORS(cfg.AllowedOrigins, mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("api listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
