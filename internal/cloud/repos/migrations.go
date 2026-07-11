package repos

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RunCloudMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS cloud_users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			full_name VARCHAR(255) DEFAULT '',
			password_hash TEXT NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'user',
			email_verified BOOLEAN DEFAULT FALSE,
			verified_at TIMESTAMP WITH TIME ZONE,
			verify_token TEXT,
			verify_token_expires_at TIMESTAMP WITH TIME ZONE,
			otp_code VARCHAR(6),
			otp_expires_at TIMESTAMP WITH TIME ZONE,
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
		`CREATE INDEX IF NOT EXISTS idx_cloud_users_email ON cloud_users(email);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w\nQuery: %s", err, query)
		}
	}

	log.Println("Cloud PostgreSQL migrations applied successfully.")
	return nil
}

// SeedAdminUser inserts a default admin user if CLOUD_ADMIN_EMAIL and CLOUD_ADMIN_PASSWORD are set.
func SeedAdminUser(db *sql.DB) error {
	adminEmail := os.Getenv("CLOUD_ADMIN_EMAIL")
	adminPassword := os.Getenv("CLOUD_ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		return nil
	}

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM cloud_users WHERE email = $1`, adminEmail).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check admin user existence: %w", err)
	}
	if count > 0 {
		log.Printf("Admin user %s already exists, skipping seed.", adminEmail)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	id := "admin-" + adminEmail
	_, err = db.Exec(
		`INSERT INTO cloud_users (id, email, full_name, password_hash, role, email_verified) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, adminEmail, "Admin", string(hash), "admin", true,
	)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	log.Printf("Seeded admin user: %s", adminEmail)
	return nil
}
