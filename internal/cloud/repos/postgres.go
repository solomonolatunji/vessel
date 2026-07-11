package repos

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type CloudDB struct {
	db *sql.DB
}

func NewCloudDB(dsn string) (*CloudDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("Warning: Failed to ping Postgres (if DSN is empty, this is expected)")
	}

	if err := RunCloudMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run cloud migrations: %w", err)
	}

	return &CloudDB{db: db}, nil
}
