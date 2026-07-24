package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

func runProjects(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codedockd project:<command> [args]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  list                 List all projects")
		fmt.Println("  show <id>            Show project details")
		fmt.Println("  create <name>        Create a project")
		fmt.Println("  destroy <id>         Delete a project")
		return
	}

	_, db, _ := initDataDir()
	defer db.Close()

	envRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, envRepo)

	cmd := args[0]

	switch cmd {
	case "list":
		projects, total, err := projectRepo.List(context.Background(), 100, 0)
		if err != nil {
			exitError("Failed to list projects: %v", err)
		}
		fmt.Printf("  %d projects found:\n", total)
		for _, p := range projects {
			fmt.Printf("  %-8s  %-20s  %s\n", p.ID[:8], p.Name, p.Description)
		}

	case "show":
		if len(args) < 2 {
			exitError("Usage: codedockd project:show <id>")
		}
		p, err := projectRepo.Get(context.Background(), args[1])
		if err != nil {
			exitError("Project not found: %v", err)
		}
		fmt.Printf("  ID:          %s\n", p.ID)
		fmt.Printf("  Name:        %s\n", p.Name)
		fmt.Printf("  Description: %s\n", p.Description)
		fmt.Printf("  Created:     %s\n", p.CreatedAt.Format(time.RFC3339))

		envs, _ := envRepo.ListByProject(context.Background(), p.ID)
		if len(envs) > 0 {
			fmt.Println("  Environments:")
			for _, e := range envs {
				fmt.Printf("    ├─ %s  %s\n", e.ID[:8], e.Name)
			}
		}

	case "create":
		if len(args) < 2 {
			exitError("Usage: codedockd project:create <name>")
		}
		name := args[1]
		description := ""
		for i := 2; i < len(args); i++ {
			if args[i] == "--description" && i+1 < len(args) {
				description = args[i+1]
				i++
			}
		}

		proj := &models.ProjectConfig{
			ID:          uuid.New().String(),
			Name:        name,
			Description: description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := projectRepo.Create(context.Background(), proj); err != nil {
			exitError("Failed to create project: %v", err)
		}
		fmt.Printf("✅ Project created: %s (%s)\n", name, proj.ID[:8])

		env := &models.EnvironmentConfig{
			ID:        uuid.New().String(),
			ProjectID: proj.ID,
			Name:      "production",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := envRepo.Create(context.Background(), env); err != nil {
			fmt.Printf("  ⚠️  Failed to create default environment: %v\n", err)
		} else {
			fmt.Printf("  └─ Default environment 'production' created (%s)\n", env.ID[:8])
		}

	case "destroy":
		if len(args) < 2 {
			exitError("Usage: codedockd project:destroy <id>")
		}
		p, err := projectRepo.Get(context.Background(), args[1])
		if err != nil {
			exitError("Project not found: %v", err)
		}
		fmt.Printf("Are you sure you want to delete project '%s' (%s)? (y/N): ", p.Name, p.ID[:8])
		var confirm string
		fmt.Scanln(&confirm)
		if !isYes(confirm) {
			fmt.Println("Cancelled.")
			return
		}
		if err := projectRepo.Delete(context.Background(), args[1]); err != nil {
			exitError("Failed to delete project: %v", err)
		}
		fmt.Printf("✅ Project deleted: %s\n", p.Name)

	default:
		fmt.Printf("Unknown project command: %s\n", cmd)
		fmt.Println("Try: list, show <id>, create <name>, destroy <id>")
	}
}

func runEnvVars(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codedockd env:<command> [args]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  list --project <id>       List all env vars for a project")
		fmt.Println("  set KEY=VALUE --project <id>  Set one or more env vars")
		fmt.Println("  unset KEY --project <id>  Remove an env var")
		return
	}

	_, db, vlt := initDataDir()
	defer db.Close()

	envRepo := repositories.NewEnvRepo(db, vlt)

	projectID := ""
	for i := 1; i < len(args); i++ {
		if args[i] == "--project" && i+1 < len(args) {
			projectID = args[i+1]
		}
	}

	cmd := args[0]

	switch cmd {
	case "list":
		if projectID == "" {
			exitError("--project <id> is required")
		}
		vars, err := envRepo.GetVars(context.Background(), projectID)
		if err != nil {
			exitError("Failed to list env vars: %v", err)
		}
		if len(vars) == 0 {
			fmt.Println("  No environment variables set.")
			return
		}
		for k, v := range vars {
			fmt.Printf("  %s=%s\n", k, v)
		}

	case "set":
		if projectID == "" {
			exitError("--project <id> is required")
		}
		newVars := map[string]string{}
		for _, arg := range args[1:] {
			if arg == "--project" {
				break
			}
			if len(arg) > 0 && arg[0] != '-' {
				for j := 0; j < len(arg); j++ {
					if arg[j] == '=' {
						newVars[arg[:j]] = arg[j+1:]
						break
					}
				}
			}
		}
		if len(newVars) == 0 {
			exitError("No KEY=VALUE pairs provided")
		}
		for k, v := range newVars {
			if err := envRepo.SetVar(context.Background(), projectID, k, v); err != nil {
				exitError("Failed to set %s: %v", k, err)
			}
			fmt.Printf("  ✅ Set %s\n", k)
		}

	case "unset":
		if projectID == "" {
			exitError("--project <id> is required")
		}
		existing, _ := envRepo.GetVars(context.Background(), projectID)
		if existing == nil {
			fmt.Println("  No vars to unset.")
			return
		}
		for k := range existing {
			for _, arg := range args[1:] {
				if arg == k {
					delete(existing, k)
					fmt.Printf("  ✅ Unset %s\n", k)
				}
			}
		}
		for k, v := range existing {
			_ = envRepo.SetVar(context.Background(), projectID, k, v)
		}

	default:
		fmt.Printf("Unknown env command: %s\n", cmd)
		fmt.Println("Try: list, set KEY=VALUE, unset KEY")
	}
}
