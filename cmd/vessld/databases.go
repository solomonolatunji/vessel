package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

func runDatabases(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vessld db:<command> [args]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  list                  List all databases")
		fmt.Println("  show <id>             Show database details")
		fmt.Println("  create <name> <engine>  Create a database (--project required)")
		fmt.Println("  destroy <id>          Delete a database")
		return
	}

	_, db, vlt := initDataDir()
	defer db.Close()

	dbRepo := repositories.NewDatabaseSQLiteRepository(db, vlt)

	cmd := args[0]

	switch cmd {
	case "list":
		dbs, err := dbRepo.List(context.Background())
		if err != nil {
			exitError("Failed to list databases: %v", err)
		}
		if len(dbs) == 0 {
			fmt.Println("  No databases found.")
			return
		}
		for _, d := range dbs {
			status := d.Status
			if status == "" {
				status = "stopped"
			}
			fmt.Printf("  %s  %s  %s  port=%d  status=%s\n",
				d.ID[:8], d.Name, d.Engine, d.Port, status)
		}

	case "show":
		if len(args) < 2 {
			exitError("Usage: vessld db:show <id>")
		}
		d, err := dbRepo.GetByID(context.Background(), args[1])
		if err != nil {
			exitError("Database not found: %v", err)
		}
		fmt.Printf("  ID:       %s\n", d.ID)
		fmt.Printf("  Name:     %s\n", d.Name)
		fmt.Printf("  Engine:   %s\n", d.Engine)
		fmt.Printf("  Version:  %s\n", d.Version)
		fmt.Printf("  Port:     %d\n", d.Port)
		fmt.Printf("  Status:   %s\n", d.Status)
		fmt.Printf("  Project:  %s\n", d.ProjectID)
		fmt.Printf("  DNS:      %s\n", d.InternalDNS)
		fmt.Printf("  Username: %s\n", d.Username)
		fmt.Printf("  DB Name:  %s\n", d.DatabaseName)
		fmt.Printf("  Volume:   %s\n", d.VolumePath)

	case "create":
		if len(args) < 3 {
			exitError("Usage: vessld db:create <name> <engine> --project <id>")
		}
		name := args[1]
		engine := args[2]
		projectID := ""
		version := ""
		for i := 3; i < len(args); i++ {
			switch args[i] {
			case "--project":
				if i+1 < len(args) {
					projectID = args[i+1]
					i++
				}
			case "--version":
				if i+1 < len(args) {
					version = args[i+1]
					i++
				}
			}
		}
		if projectID == "" {
			exitError("--project <id> is required")
		}

		engine = mapEngine(engine)
		port := defaultPortForEngine(engine)

		database := &models.Database{
			ID:           uuid.New().String(),
			ProjectID:    projectID,
			Name:         name,
			Engine:       models.DatabaseEngine(engine),
			Version:      version,
			Port:         port,
			Status:       models.DatabaseStatusCreated,
			Username:     "vessl",
			DatabaseName: "vessl",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		database.Password = uuid.New().String()[:16]

		if err := dbRepo.Create(context.Background(), database); err != nil {
			exitError("Failed to create database: %v", err)
		}
		fmt.Printf("✅ Database created: %s (%s %s)\n", name, engine, version)
		fmt.Printf("   Connection: %s://%s:%s@%s:%d/%s\n",
			engine, database.Username, database.Password, name, port, database.DatabaseName)

	case "destroy":
		if len(args) < 2 {
			exitError("Usage: vessld db:destroy <id>")
		}
		d, err := dbRepo.GetByID(context.Background(), args[1])
		if err != nil {
			exitError("Database not found: %v", err)
		}
		fmt.Printf("Are you sure you want to delete '%s' (%s %s)? (y/N): ", d.Name, d.Engine, d.ID[:8])
		var confirm string
		fmt.Scanln(&confirm)
		if !isYes(confirm) {
			fmt.Println("Cancelled.")
			return
		}
		if err := dbRepo.Delete(context.Background(), args[1]); err != nil {
			exitError("Failed to delete database: %v", err)
		}
		fmt.Printf("✅ Database deleted: %s\n", d.Name)

	default:
		fmt.Printf("Unknown db command: %s\n", cmd)
		fmt.Println("Try: list, show <id>, create <name> <engine>, destroy <id>")
	}
}

func mapEngine(e string) string {
	switch e {
	case "postgres", "postgresql":
		return "postgres"
	case "mysql":
		return "mysql"
	case "mariadb":
		return "mariadb"
	case "redis":
		return "redis"
	case "mongo", "mongodb":
		return "mongodb"
	case "clickhouse":
		return "clickhouse"
	case "kafka":
		return "kafka"
	case "rabbitmq":
		return "rabbitmq"
	case "nats":
		return "nats"
	default:
		return e
	}
}

func defaultPortForEngine(e string) int {
	switch e {
	case "postgres":
		return 5432
	case "mysql", "mariadb":
		return 3306
	case "redis":
		return 6379
	case "mongodb":
		return 27017
	case "clickhouse":
		return 9000
	case "kafka":
		return 9092
	case "rabbitmq":
		return 5672
	case "nats":
		return 4222
	default:
		return 0
	}
}
