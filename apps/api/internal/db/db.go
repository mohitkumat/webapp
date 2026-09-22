package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	pingMaxAttempts = 5
	pingBaseDelay   = 500 * time.Millisecond
	pingMaxDelay    = 5 * time.Second
)

// Connect opens a pooled connection to Postgres and verifies it is reachable,
// retrying with backoff to ride out transient startup DNS/connection blips.
func Connect(databaseURL string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err := pingWithRetry(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return conn, nil
}

func pingWithRetry(conn *sql.DB) error {
	var err error
	delay := pingBaseDelay
	for attempt := 1; attempt <= pingMaxAttempts; attempt++ {
		if err = conn.Ping(); err == nil {
			return nil
		}
		if attempt == pingMaxAttempts {
			break
		}
		time.Sleep(delay)
		if delay *= 2; delay > pingMaxDelay {
			delay = pingMaxDelay
		}
	}
	return err
}
