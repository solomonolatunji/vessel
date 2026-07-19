package repositories

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

//go:embed schema/*.sql
var schemaFS embed.FS

func RunMigrations(db *sql.DB) error {
	files, err := listSchemaFiles()
	if err != nil {
		return err
	}

	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	applied, err := loadAppliedMigrations(db)
	if err != nil {
		return err
	}

	for _, file := range files {
		if applied[file] {
			continue
		}

		content, err := schemaFS.ReadFile("schema/" + file)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration failed for %s: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (filename) VALUES (?)", file); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		slog.Info("applied migration", "file", file)
	}

	slog.Info("schema migrations up to date")
	return nil
}

func listSchemaFiles() ([]string, error) {
	entries, err := schemaFS.ReadDir("schema")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded schema directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		return fmt.Errorf("failed to count applied migrations: %w", err)
	}

	if count > 0 {
		return nil
	}

	var serverSettingsExists bool
	if err := db.QueryRow("SELECT COUNT(*) > 0 FROM sqlite_master WHERE type = 'table' AND name = 'server_settings'").Scan(&serverSettingsExists); err != nil {
		return fmt.Errorf("failed to check existing schema: %w", err)
	}

	if !serverSettingsExists {
		return nil
	}

	files, err := listSchemaFiles()
	if err != nil {
		return err
	}

	var baselineFiles []string
	for _, f := range files {
		if strings.HasPrefix(f, "001_") || strings.HasPrefix(f, "002_") {
			baselineFiles = append(baselineFiles, f)
		}
	}

	if len(baselineFiles) > 0 {
		var placeholders []string
		var args []any
		for _, file := range baselineFiles {
			placeholders = append(placeholders, "(?)")
			args = append(args, file)
		}
		query := fmt.Sprintf("INSERT INTO schema_migrations (filename) VALUES %s", strings.Join(placeholders, ", "))
		if _, err := db.Exec(query, args...); err != nil {
			return fmt.Errorf("failed to seed migrations: %w", err)
		}
	}

	slog.Info("seeded schema_migrations for existing database")
	return nil
}

func loadAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT filename FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to load applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		applied[filename] = true
	}

	return applied, rows.Err()
}
