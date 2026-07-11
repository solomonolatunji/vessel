package server

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vessel.dev/vessel/internal/models"
)

func InitDatabase() *gorm.DB {
	dsn := os.Getenv("CLOUD_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=vessel password=vessel dbname=vesselcloud port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to cloud database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.CloudTeam{},
		&models.CloudServer{},
		&models.CloudUsageLog{},
		&models.CloudTelemetryLog{},
	); err != nil {
		log.Fatalf("Failed to run cloud database migrations: %v", err)
	}

	return db
}
