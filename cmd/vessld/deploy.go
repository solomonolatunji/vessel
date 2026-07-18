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
	composePath := ""
	archivePath := ""
	projectID := ""
	branch := "main"
	rootDir := ""
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
		case "--compose", "-c":
			if i+1 < len(args) {
				composePath = args[i+1]
				i++
			}
		case "--archive", "-a":
			if i+1 < len(args) {
				archivePath = args[i+1]
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
		case "--dir", "-d":
			if i+1 < len(args) {
				rootDir = args[i+1]
				i++
			}
		case "--template", "-t":
			if i+1 < len(args) {
				templateName := args[i+1]
				gitURL = "https://github.com/vesslhq/vessl-examples.git"
				branch = "main"
				rootDir = templateName
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

	if gitURL == "" && imageRef == "" && composePath == "" && archivePath == "" {
		exitError("Usage: vessld deploy <git-url> | --template <t> | --image <img> | --compose <file> | --archive <file>")
	}

	count := 0
	for _, v := range []bool{gitURL != "", imageRef != "", composePath != "", archivePath != ""} {
		if v {
			count++
		}
	}
	if count > 1 {
		exitError("Specify only one: Git URL, --image, --compose, or --archive")
	}

	dataDir, db, vlt := initDataDir()
	defer db.Close()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		exitError("Failed to connect to Docker: %v", err)
	}

	envRepo := repositories.NewEnvironmentRepo(db)
	appRepo := repositories.NewAppServiceRepo(db)
	projectRepo := repositories.NewProjectRepo(db, envRepo)
	settingsRepo := repositories.NewSettingsRepo(db)

	appName := extractRepoName(gitURL)
	if appName == "app" && imageRef != "" {
		appName = imageRef
		if idx := strings.LastIndex(appName, "/"); idx >= 0 {
			appName = appName[idx+1:]
		}
		appName = strings.Split(appName, ":")[0]
	}

	if projectID == "" {
		projects, _, _ := projectRepo.List(context.Background(), 1000, 0)
		if len(projects) > 0 {
			fmt.Println("📁 Select a project for this deployment:")
			fmt.Println("  [0] Create a new project")
			for i, p := range projects {
				fmt.Printf("  [%d] %s (%s)\n", i+1, p.Name, p.ID[:8])
			}
			for {
				fmt.Printf("Enter choice [0-%d]: ", len(projects))
				var choice int
				_, err := fmt.Scanln(&choice)
				if err != nil || choice < 0 || choice > len(projects) {
					fmt.Println("Invalid choice.")
					continue
				}
				if choice > 0 {
					projectID = projects[choice-1].ID
					fmt.Printf("📁 Using project: %s (%s)\n", projects[choice-1].Name, projectID[:8])
					break
				} else {
					break // choice == 0
				}
			}
		}

		if projectID == "" {
			defaultName := utils.GenerateRandomName()
			fmt.Printf("📁 New project name (press Enter to use '%s'): ", defaultName)
			var newName string
			fmt.Scanln(&newName)
			newName = strings.TrimSpace(newName)
			if newName == "" {
				newName = defaultName
			}
			p := &models.ProjectConfig{
				ID:        uuid.New().String(),
				Name:      newName,
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
			RepositoryURL: gitURL,
			Branch:        branch,
			RootDirectory: rootDir,
			InternalPort:  port,
			Status:        "created",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if imageRef != "" {
			svc.ImageRef = imageRef
		}
		if err := appRepo.Create(context.Background(), svc); err != nil {
			exitError("Failed to create app: %v", err)
		}
		fmt.Printf("📦 Created app: %s (%s)\n", appName, svc.ID[:8])
	} else {
		fmt.Printf("📦 Using existing app: %s (%s)\n", appName, svc.ID[:8])
	}

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})

	switch {
	case gitURL != "":
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
		fmt.Println("🔨 Building and deploying...")
		srcDir := cloneDir
		if svc.RootDirectory != "" {
			srcDir = filepath.Join(cloneDir, svc.RootDirectory)
		}
		containerID, err := deployer.DeployAppService(context.Background(), svc, srcDir, os.Stdout)
		if err != nil {
			exitError("Deployment failed: %v", err)
		}
		fmt.Printf("\n✅ Deployed! Container: %s\n", containerID)

	case imageRef != "":
		fmt.Printf("🐳 Deploying image %s...\n", imageRef)
		containerID, err := deployer.DeployImage(context.Background(), svc, os.Stdout)
		if err != nil {
			exitError("Image deploy failed: %v", err)
		}
		slog.Info("container started from image", "image", imageRef, "containerID", containerID)
		fmt.Printf("\n✅ Deployed! Container: %s\n", containerID)

	case composePath != "":
		fmt.Printf("📦 Deploying compose file %s...\n", composePath)
		composeDeployer := engine.NewComposeDeployer(dockerClient)
		services, err := composeDeployer.Deploy(context.Background(), composePath, projectID)
		if err != nil {
			exitError("Compose deploy failed: %v", err)
		}
		fmt.Printf("\n✅ Deployed %d services from compose file\n", len(services))
		for _, s := range services {
			fmt.Printf("   - %s (%s)\n", s.Name, s.ContainerID[:12])
		}

	case archivePath != "":
		fmt.Printf("📦 Deploying archive %s...\n", archivePath)
		archiveDir := filepath.Join(dataDir, "builds", svc.ID, "archive")
		_ = os.RemoveAll(archiveDir)
		_ = os.MkdirAll(archiveDir, 0o755)
		f, err := os.Open(archivePath)
		if err != nil {
			exitError("Failed to open archive: %v", err)
		}
		if err := extractArchiveTo(archiveDir, f); err != nil {
			f.Close()
			exitError("Failed to extract archive: %v", err)
		}
		f.Close()

		srcDir := findSourceDir(archiveDir)
		fmt.Println("🔨 Building and deploying...")
		containerID, err := deployer.DeployAppService(context.Background(), svc, srcDir, os.Stdout)
		if err != nil {
			exitError("Deployment failed: %v", err)
		}
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
		magicDomain := os.Getenv("VESSL_MAGIC_DOMAIN")
		if magicDomain == "" {
			magicDomain = "sslip.io"
		}
		fmt.Printf("   URL: http://%s.%s.%s\n", cleanName, cleanIP, magicDomain)
	}
}
