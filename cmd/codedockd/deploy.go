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

	"codedock.dev/codedock/internal/engine"
	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/repositories"
	"codedock.dev/codedock/internal/utils"
)

type deployArgs struct {
	gitURL      string
	imageRef    string
	archivePath string
	projectID   string
	branch      string
	rootDir     string
	port        int
}

func runDeploy(args []string) {
	dArgs := parseDeployArgs(args)

	if dArgs.gitURL == "" && dArgs.imageRef == "" && dArgs.archivePath == "" {
		exitError("Usage: codedockd deploy <git-url> | --template <t> | --image <img> | --archive <file>")
	}

	count := 0
	for _, v := range []bool{dArgs.gitURL != "", dArgs.imageRef != "", dArgs.archivePath != ""} {
		if v {
			count++
		}
	}
	if count > 1 {
		exitError("Specify only one: Git URL, --image, or --archive")
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

	appName := resolveAppName(dArgs.gitURL, dArgs.imageRef)

	dArgs.projectID = selectOrCreateProject(projectRepo, dArgs.projectID)
	envID := setupEnvironment(envRepo, dArgs.projectID)
	svc := setupAppService(appRepo, dArgs.projectID, envID, appName, dArgs)

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})

	performDeployment(deployer, dockerClient, svc, dArgs, dataDir)
	printDeploymentURL(settingsRepo, appName)
}

func parseDeployArgs(args []string) deployArgs {
	d := deployArgs{
		branch: "main",
		port:   3000,
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--git", "-g":
			if i+1 < len(args) {
				d.gitURL = args[i+1]
				i++
			}
		case "--image", "-i":
			if i+1 < len(args) {
				d.imageRef = args[i+1]
				i++
			}
		case "--archive", "-a":
			if i+1 < len(args) {
				d.archivePath = args[i+1]
				i++
			}
		case "--project", "-p":
			if i+1 < len(args) {
				d.projectID = args[i+1]
				i++
			}
		case "--branch", "-b":
			if i+1 < len(args) {
				d.branch = args[i+1]
				i++
			}
		case "--dir", "-d":
			if i+1 < len(args) {
				d.rootDir = args[i+1]
				i++
			}
		case "--template", "-t":
			if i+1 < len(args) {
				templateName := args[i+1]
				d.gitURL = "https://github.com/buildwithtechx/codedock-examples.git"
				d.branch = "main"
				d.rootDir = templateName
				i++
			}
		case "--port":
			if i+1 < len(args) {
				if p, err := parseUint(args[i+1]); err == nil {
					d.port = p
				}
				i++
			}
		default:
			if strings.HasPrefix(args[i], "http") || strings.HasPrefix(args[i], "git@") {
				d.gitURL = args[i]
			} else if strings.Contains(args[i], ":") || strings.Contains(args[i], "/") {
				d.imageRef = args[i]
			}
		}
	}
	return d
}

func resolveAppName(gitURL, imageRef string) string {
	appName := extractRepoName(gitURL)
	if appName == "app" && imageRef != "" {
		appName = imageRef
		if idx := strings.LastIndex(appName, "/"); idx >= 0 {
			appName = appName[idx+1:]
		}
		appName = strings.Split(appName, ":")[0]
	}
	return appName
}

func selectOrCreateProject(projectRepo *repositories.ProjectRepo, projectID string) string {
	if projectID != "" {
		return projectID
	}
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
				fmt.Printf("📁 Using project: %s (%s)\n", projects[choice-1].Name, projects[choice-1].ID[:8])
				return projects[choice-1].ID
			} else {
				break
			}
		}
	}

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
	fmt.Printf("📁 Created project: %s (%s)\n", p.Name, p.ID[:8])
	return p.ID
}

func setupEnvironment(envRepo *repositories.EnvironmentRepo, projectID string) string {
	envs, _ := envRepo.ListByProject(context.Background(), projectID)
	if len(envs) > 0 {
		return envs[0].ID
	}
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
	fmt.Println("  Created environment: production")
	return env.ID
}

func setupAppService(appRepo *repositories.AppServiceRepo, projectID, envID, appName string, dArgs deployArgs) *models.AppService {
	apps, _ := appRepo.ListByProject(context.Background(), projectID)
	for _, a := range apps {
		if a.Name == appName {
			fmt.Printf("📦 Using existing app: %s (%s)\n", appName, a.ID[:8])
			return a
		}
	}

	svc := &models.AppService{
		ID:            uuid.New().String(),
		ProjectID:     projectID,
		EnvironmentID: envID,
		Name:          appName,
		RepositoryURL: dArgs.gitURL,
		Branch:        dArgs.branch,
		RootDirectory: dArgs.rootDir,
		InternalPort:  dArgs.port,
		Status:        "created",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if dArgs.imageRef != "" {
		svc.ImageRef = dArgs.imageRef
	}
	if err := appRepo.Create(context.Background(), svc); err != nil {
		exitError("Failed to create app: %v", err)
	}
	fmt.Printf("📦 Created app: %s (%s)\n", appName, svc.ID[:8])
	return svc
}

func performDeployment(deployer *engine.Deployer, dockerClient *client.Client, svc *models.AppService, dArgs deployArgs, dataDir string) {
	switch {
	case dArgs.gitURL != "":
		cloneDir := filepath.Join(dataDir, "builds", svc.ID)
		_ = os.RemoveAll(cloneDir)
		_ = os.MkdirAll(cloneDir, 0o755)
		fmt.Printf("📥 Cloning %s (branch: %s)...\n", dArgs.gitURL, dArgs.branch)
		cloneCmd := exec.Command("git", "clone", "--depth", "1", "--branch", dArgs.branch, dArgs.gitURL, cloneDir)
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

	case dArgs.imageRef != "":
		fmt.Printf("🐳 Deploying image %s...\n", dArgs.imageRef)
		containerID, err := deployer.DeployImage(context.Background(), svc, os.Stdout)
		if err != nil {
			exitError("Image deploy failed: %v", err)
		}
		slog.Info("container started from image", "image", dArgs.imageRef, "containerID", containerID)
		fmt.Printf("\n✅ Deployed! Container: %s\n", containerID)

	case dArgs.archivePath != "":
		fmt.Printf("📦 Deploying archive %s...\n", dArgs.archivePath)
		archiveDir := filepath.Join(dataDir, "builds", svc.ID, "archive")
		_ = os.RemoveAll(archiveDir)
		_ = os.MkdirAll(archiveDir, 0o755)
		f, err := os.Open(dArgs.archivePath)
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

	fmt.Printf("   App: %s (%s)\n", svc.Name, svc.ID[:8])
}

func printDeploymentURL(settingsRepo *repositories.SettingsRepo, appName string) {
	wildcard := ""
	settings, err := settingsRepo.GetServerSettings(context.Background())
	if err == nil {
		wildcard = settings.DefaultWildcardDomain
	}
	if wildcard == "" {
		wildcard = os.Getenv("CODEDOCK_DOMAIN")
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
		hostIP := os.Getenv("CODEDOCK_HOST_IP")
		if hostIP == "" {
			hostIP = "127.0.0.1"
		}
		cleanName := utils.SanitizeDomainName(appName)
		cleanIP := strings.ReplaceAll(hostIP, ".", "-")
		magicDomain := os.Getenv("CODEDOCK_MAGIC_DOMAIN")
		if magicDomain == "" {
			magicDomain = "sslip.io"
		}
		fmt.Printf("   URL: http://%s.%s.%s\n", cleanName, cleanIP, magicDomain)
	}
}
