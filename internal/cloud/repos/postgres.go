package repos

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type CloudDB struct {
	db *sql.DB
}

// DB returns the underlying *sql.DB for repos that need raw SQL access.
func (c *CloudDB) DB() *sql.DB {
	return c.db
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

	if err := SeedAdminUser(db); err != nil {
		log.Printf("Warning: failed to seed admin user: %v", err)
	}

	return &CloudDB{db: db}, nil
}
