package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sort"
)

//go:embed schema/*.sql
var schemaFS embed.FS

func RunMigrations(db *sql.DB) error {
	entries, err := schemaFS.ReadDir("schema")
	if err != nil {
		return fmt.Errorf("failed to read embedded schema directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := schemaFS.ReadFile("schema/" + file)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("migration failed for %s: %w", file, err)
		}
	}

	log.Println(" Vessel SQLite schema initialized successfully from embedded files")
	return nil
}
