package handlers_test

import (
	"database/sql"
	"os"
	"testing"

	"api/internal/db"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable"
	}

	conn, err := db.Connect(url)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}

	if _, err := conn.Exec(`DROP TABLE IF EXISTS items, schema_migrations`); err != nil {
		t.Fatalf("reset test db: %v", err)
	}

	if err := db.RunMigrations(conn, "../../migrations"); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	t.Cleanup(func() { conn.Close() })

	return conn
}
