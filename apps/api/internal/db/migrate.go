package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations applies every .sql file in dir, in filename order, that has not
// already been recorded in schema_migrations. Each file runs in its own transaction.
// Filenames should sort in the order they must run, e.g. 0001_init.sql, 0002_add_x.sql.
func RunMigrations(conn *sql.DB, dir string) error {
	if _, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename    TEXT PRIMARY KEY,
			applied_at  TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	files, err := migrationFiles(dir)
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}

	for _, file := range files {
		var alreadyApplied bool
		err := conn.QueryRow(
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE filename = $1)`,
			filepath.Base(file),
		).Scan(&alreadyApplied)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", file, err)
		}
		if alreadyApplied {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		tx, err := conn.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", file, err)
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", file, err)
		}

		if _, err := tx.Exec(
			`INSERT INTO schema_migrations (filename) VALUES ($1)`,
			filepath.Base(file),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", file, err)
		}
	}

	return nil
}

func migrationFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}
