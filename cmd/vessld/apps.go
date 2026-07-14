package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

func runApps(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vessld apps:<command> [args]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  list                  List all apps")
		fmt.Println("  show <id>             Show app details")
		fmt.Println("  create <name>         Create an app (set --project and optional --env, --port)")
		fmt.Println("  destroy <id>          Delete an app")
		return
	}

	_, db, _ := initDataDir()
	defer db.Close()

	appRepo := repositories.NewAppServiceSQLiteRepository(db)
	svcVarRepo := repositories.NewServiceVarSQLiteRepository(db)
	envRepo := repositories.NewEnvironmentSQLiteRepository(db)
	projectRepo := repositories.NewProjectSQLiteRepository(db, envRepo)

	cmd := args[0]

	switch cmd {
	case "list":
		projects, _, _ := projectRepo.List(context.Background(), 1000, 0)
		for _, p := range projects {
			apps, _ := appRepo.ListByProject(context.Background(), p.ID)
			if len(apps) == 0 {
				fmt.Printf("  %s (%s): no apps\n", p.Name, p.ID[:8])
				continue
			}
			fmt.Printf("  %s (%s):\n", p.Name, p.ID[:8])
			for _, a := range apps {
				fmt.Printf("    ├─ %s  %s  port=%d  build=%s  status=%s\n",
					a.ID[:8], a.Name, a.InternalPort, a.BuildEngine, a.Status)
			}
		}

	case "show":
		if len(args) < 2 {
			exitError("Usage: vessld apps:show <id>")
		}
		app, err := appRepo.GetByID(context.Background(), args[1])
		if err != nil {
			exitError("App not found: %v", err)
		}
		fmt.Printf("  ID:        %s\n", app.ID)
		fmt.Printf("  Name:      %s\n", app.Name)
		fmt.Printf("  Project:   %s\n", app.ProjectID)
		fmt.Printf("  Env:       %s\n", app.EnvironmentID)
		fmt.Printf("  Port:      %d\n", app.InternalPort)
		fmt.Printf("  Status:    %s\n", app.Status)
		fmt.Printf("  Domain:    %s\n", app.Domain)
		fmt.Printf("  Build:     %s\n", app.BuildEngine)
		fmt.Printf("  Repo:      %s\n", app.RepositoryURL)
		fmt.Printf("  Branch:    %s\n", app.Branch)
		fmt.Printf("  Docker:    %s\n", app.DockerfilePath)
		vars, _ := svcVarRepo.ListByService(context.Background(), app.ID)
		if len(vars) > 0 {
			fmt.Printf("  Variables:\n")
			for _, v := range vars {
				fmt.Printf("    %s=%s\n", v.Key, v.Value)
			}
		}

	case "create":
		if len(args) < 2 {
			exitError("Usage: vessld apps:create <name> --project <id> [--env <id>] [--port <n>]")
		}
		name := args[1]
		projectID := ""
		envID := ""
		port := 3000
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--project":
				if i+1 < len(args) {
					projectID = args[i+1]
					i++
				}
			case "--env":
				if i+1 < len(args) {
					envID = args[i+1]
					i++
				}
			case "--port":
				if i+1 < len(args) {
					p, err := parseUint(args[i+1])
					if err == nil {
						port = p
					}
					i++
				}
			}
		}
		if projectID == "" {
			exitError("--project <id> is required")
		}

		if envID == "" {
			envs, _ := envRepo.ListByProject(context.Background(), projectID)
			if len(envs) > 0 {
				envID = envs[0].ID
			}
		}

		svc := &models.AppService{
			ID:            uuid.New().String(),
			ProjectID:     projectID,
			EnvironmentID: envID,
			Name:          name,
			InternalPort:  port,
			Status:        "created",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := appRepo.Create(context.Background(), svc); err != nil {
			exitError("Failed to create app: %v", err)
		}
		fmt.Printf("✅ App created: %s (%s)\n", name, svc.ID[:8])

	case "destroy":
		if len(args) < 2 {
			exitError("Usage: vessld apps:destroy <id>")
		}
		app, err := appRepo.GetByID(context.Background(), args[1])
		if err != nil {
			exitError("App not found: %v", err)
		}
		fmt.Printf("Are you sure you want to delete '%s' (%s)? (y/N): ", app.Name, app.ID[:8])
		var confirm string
		fmt.Scanln(&confirm)
		if !isYes(confirm) {
			fmt.Println("Cancelled.")
			return
		}
		if err := appRepo.Delete(context.Background(), args[1]); err != nil {
			exitError("Failed to delete app: %v", err)
		}
		fmt.Printf("✅ App deleted: %s\n", app.Name)

	default:
		fmt.Printf("Unknown apps command: %s\n", cmd)
		fmt.Println("Try: list, show <id>, create <name>, destroy <id>")
	}
}

func isYes(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "y" || s == "yes"
}
