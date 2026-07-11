package repos

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func RunCloudMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS cloud_users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS cloud_teams (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS cloud_team_members (
			user_id VARCHAR(255) REFERENCES cloud_users(id) ON DELETE CASCADE,
			team_id VARCHAR(255) REFERENCES cloud_teams(id) ON DELETE CASCADE,
			role VARCHAR(50) DEFAULT 'member',
			PRIMARY KEY (user_id, team_id)
		);`,
		`CREATE TABLE IF NOT EXISTS cloud_servers (
			id VARCHAR(255) PRIMARY KEY,
			team_id VARCHAR(255) REFERENCES cloud_teams(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			agent_token VARCHAR(255) UNIQUE NOT NULL,
			is_connected BOOLEAN DEFAULT FALSE,
			last_ip VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS cloud_subscriptions (
			id VARCHAR(255) PRIMARY KEY,
			team_id VARCHAR(255) UNIQUE REFERENCES cloud_teams(id) ON DELETE CASCADE,
			stripe_customer_id VARCHAR(255),
			stripe_subscription_id VARCHAR(255),
			plan_id VARCHAR(50),
			status VARCHAR(50),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w\nQuery: %s", err, query)
		}
	}

	log.Println("Cloud PostgreSQL migrations applied successfully.")
	return nil
}
