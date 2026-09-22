package db_test

import (
	"database/sql"
	"os"
	"testing"

	"api/internal/db"

	_ "github.com/lib/pq"
)

func testConn(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable"
	}
	conn, err := db.Connect(url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestRunMigrations_IsIdempotent(t *testing.T) {
	conn := testConn(t)

	if _, err := conn.Exec(`DROP TABLE IF EXISTS items, schema_migrations`); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if err := db.RunMigrations(conn, "../../migrations"); err != nil {
		t.Fatalf("first run: %v", err)
	}
	if err := db.RunMigrations(conn, "../../migrations"); err != nil {
		t.Fatalf("second run should be a no-op, got error: %v", err)
	}

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("count migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 recorded migration, got %d", count)
	}
}
