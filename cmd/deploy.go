package main

import (
	"context"
	"fmt"
	"log/slog"
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
	"vessl.dev/vessl/internal/utils"
)

func runDeploy(args []string) {
	gitURL := ""
	imageRef := ""
	projectID := ""
	branch := "main"
	port := 3000

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--git", "-g":
			if i+1 < len(args) {
				gitURL = args[i+1]
				i++
			}
		case "--image", "-i":
			if i+1 < len(args) {
				imageRef = args[i+1]
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
		case "--port":
			if i+1 < len(args) {
				if p, err := parseUint(args[i+1]); err == nil {
					port = p
				}
				i++
			}
		default:
			if strings.HasPrefix(args[i], "http") || strings.HasPrefix(args[i], "git@") {
				gitURL = args[i]
			} else if strings.Contains(args[i], ":") || strings.Contains(args[i], "/") {
				imageRef = args[i]
			}
		}
	}

	if gitURL == "" && imageRef == "" {
		exitError("Usage: vessld deploy <git-url> or vessld deploy --image <image> [--port <n>]")
	}

	if gitURL != "" && imageRef != "" {
		exitError("Specify either a Git URL or --image, not both")
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
	settingsRepo := repositories.NewSettingsSQLiteRepository(db)

	appName := extractRepoName(gitURL)
	if appName == "app" && imageRef != "" {
		appName = imageRef
		if idx := strings.LastIndex(appName, "/"); idx >= 0 {
			appName = appName[idx+1:]
		}
		appName = strings.Split(appName, ":")[0]
	}

	if projectID == "" {
		projects, _, _ := projectRepo.List(context.Background(), "", 100, 0)
		if len(projects) > 0 {
			projectID = projects[0].ID
			fmt.Printf("📁 Using project: %s (%s)\n", projects[0].Name, projectID[:8])
		} else {
			p := &models.ProjectConfig{
				ID:        uuid.New().String(),
				Name:      appName,
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
			InternalPort:  port,
			Status:        "created",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if imageRef != "" {
			svc.RepositoryURL = imageRef
		}
		if err := appRepo.Create(context.Background(), svc); err != nil {
			exitError("Failed to create app: %v", err)
		}
		fmt.Printf("📦 Created app: %s (%s)\n", appName, svc.ID[:8])
	} else {
		fmt.Printf("📦 Using existing app: %s (%s)\n", appName, svc.ID[:8])
	}

	if gitURL != "" {
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
	} else {
		fmt.Printf("🐳 Pulling image %s...\n", imageRef)
		cm := engine.NewContainerManager(dockerClient, &dbDeployerStore{db: db, vault: vlt})
		containerName := fmt.Sprintf("vessl-app-%s", svc.ID[:8])
		containerID, err := cm.CreateAndStart(
			context.Background(),
			containerName,
			imageRef,
			svc.ID,
			"",
			port,
			[]string{},
			512,
			0.5,
			"",
		)
		if err != nil {
			exitError("Failed to start container: %v", err)
		}
		slog.Info("container started from image", "image", imageRef, "containerID", containerID)
		fmt.Printf("\n✅ Deployed! Container: %s\n", containerID)
	}

	fmt.Printf("   App: %s (%s)\n", appName, svc.ID[:8])

	wildcard := ""
	settings, err := settingsRepo.GetServerSettings(context.Background())
	if err == nil {
		wildcard = settings.DefaultWildcardDomain
	}
	if wildcard == "" {
		wildcard = os.Getenv("VESSL_DOMAIN")
	}

	if wildcard != "" {
		cleanName := utils.SanitizeDomainName(appName)
		base := strings.TrimPrefix(wildcard, "*.")
		if strings.HasPrefix(base, "http") {
			base = strings.TrimPrefix(base, "https://")
			base = strings.TrimPrefix(base, "http://")
		}
		fmt.Printf("   URL: https://%s.%s\n", cleanName, base)
	} else {
		hostIP := os.Getenv("VESSL_HOST_IP")
		if hostIP == "" {
			hostIP = "127.0.0.1"
		}
		cleanName := utils.SanitizeDomainName(appName)
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
