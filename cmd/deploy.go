package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

func runDeploy(args []string) {
	gitURL := ""
	projectID := ""
	branch := "main"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--git", "-g":
			if i+1 < len(args) {
				gitURL = args[i+1]
				i++
			}
		case "--project", "-p":
			if i+1 < len(args) {
				projectID = args[i+1]
				i++
			}
		case "--branch", "-b":
			if i+1 < len(args) {
				branch = args[i+1]
				i++
			}
		default:
			if strings.HasPrefix(args[i], "http") || strings.HasPrefix(args[i], "git@") {
				gitURL = args[i]
			}
		}
	}

	if gitURL == "" {
		exitError("Usage: vessld deploy <git-url> [--project <id>] [--branch <name>]")
	}

	dataDir, db, vlt := initDataDir()
	defer db.Close()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		exitError("Failed to connect to Docker: %v", err)
	}

	envRepo := repositories.NewEnvironmentSQLiteRepository(db)
	appRepo := repositories.NewAppServiceSQLiteRepository(db)
	projectRepo := repositories.NewProjectSQLiteRepository(db, envRepo)

	if projectID == "" {
		projects, _, _ := projectRepo.List(context.Background(), "", 100, 0)
		if len(projects) > 0 {
			projectID = projects[0].ID
			fmt.Printf("📁 Using project: %s (%s)\n", projects[0].Name, projectID[:8])
		} else {
			p := &models.ProjectConfig{
				ID:        uuid.New().String(),
				Name:      extractRepoName(gitURL),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := projectRepo.Create(context.Background(), p); err != nil {
				exitError("Failed to create project: %v", err)
			}
			projectID = p.ID
			fmt.Printf("📁 Created project: %s (%s)\n", p.Name, projectID[:8])
		}
	}

	envs, _ := envRepo.ListByProject(context.Background(), projectID)
	var envID string
	if len(envs) > 0 {
		envID = envs[0].ID
	} else {
		env := &models.EnvironmentConfig{
			ID:        uuid.New().String(),
			ProjectID: projectID,
			Name:      "production",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := envRepo.Create(context.Background(), env); err != nil {
			exitError("Failed to create environment: %v", err)
		}
		envID = env.ID
		fmt.Println("  Created environment: production")
	}

	appName := extractRepoName(gitURL)
	apps, _ := appRepo.ListByProject(context.Background(), projectID)
	var svc *models.AppService
	for _, a := range apps {
		if a.Name == appName {
			svc = a
			break
		}
	}
	if svc == nil {
		svc = &models.AppService{
			ID:            uuid.New().String(),
			ProjectID:     projectID,
			EnvironmentID: envID,
			Name:          appName,
			InternalPort:  3000,
			BuildEngine:   "railpack",
			Status:        "created",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := appRepo.Create(context.Background(), svc); err != nil {
			exitError("Failed to create app: %v", err)
		}
		fmt.Printf("📦 Created app: %s (%s)\n", appName, svc.ID[:8])
	} else {
		fmt.Printf("📦 Using existing app: %s (%s)\n", appName, svc.ID[:8])
	}

	cloneDir := filepath.Join(dataDir, "builds", svc.ID)
	_ = os.RemoveAll(cloneDir)
	_ = os.MkdirAll(cloneDir, 0o755)
	fmt.Printf("📥 Cloning %s (branch: %s)...\n", gitURL, branch)
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "--branch", branch, gitURL, cloneDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		exitError("Git clone failed: %v", err)
	}

	project, err := projectRepo.Get(context.Background(), projectID)
	if err != nil {
		exitError("Failed to load project: %v", err)
	}

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})
	fmt.Println("🔨 Building and deploying...")
	containerID, err := deployer.Deploy(context.Background(), project, cloneDir, os.Stdout)
	if err != nil {
		exitError("Deployment failed: %v", err)
	}

	fmt.Printf("\n✅ Deployed! Container: %s\n", containerID)
	fmt.Printf("   App: %s (%s)\n", appName, svc.ID[:8])
	if hostIP := os.Getenv("VESSL_HOST_IP"); hostIP != "" {
		cleanName := strings.ToLower(strings.ReplaceAll(appName, " ", "-"))
		cleanIP := strings.ReplaceAll(hostIP, ".", "-")
		fmt.Printf("   URL: http://%s.%s.sslip.io\n", cleanName, cleanIP)
	}
}

func extractRepoName(url string) string {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "app"
}
