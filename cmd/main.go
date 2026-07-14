// @title Vessl API
// @version 1.0
// @description Vessl API Documentation
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"vessl.dev/vessl/internal/engine"
	vesslhttp "vessl.dev/vessl/internal/http"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

var vesslVersion = "dev"

type dbDeployerStore struct {
	db    *sql.DB
	vault *utils.Vault
}

func (a *dbDeployerStore) GetServerSettings() (*models.ServerSettings, error) {
	return repositories.NewSettingsSQLiteRepository(a.db).GetServerSettings(context.Background())
}

func (a *dbDeployerStore) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return repositories.NewAppServiceSQLiteRepository(a.db).ListByProject(context.Background(), projectID)
}

func (a *dbDeployerStore) GetEnvVars(projectID string) (map[string]string, error) {
	return repositories.NewEnvSQLiteRepository(a.db, a.vault).GetVars(context.Background(), projectID)
}

func (a *dbDeployerStore) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	svVarRepo := repositories.NewServiceVarSQLiteRepository(a.db)
	return svVarRepo.ListByService(context.Background(), serviceID)
}

func (a *dbDeployerStore) GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error) {
	svlsRepo := repositories.NewServerlessRepository(a.db)
	return svlsRepo.GetCodeByServiceID(context.Background(), serviceID)
}

func main() {
	_ = godotenv.Load()
	mainCLI()
}

func initDataDir() (string, *sql.DB, *utils.Vault) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dataDir := os.Getenv("VESSL_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatalf(" Failed to create data directory: %v", err)
	}
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		log.Fatalf(" Failed to initialize secrets vault: %v", err)
	}
	dbPath := filepath.Join(dataDir, "vessl.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		log.Fatalf(" Failed to open SQLite database: %v", err)
	}
	if err := repositories.RunMigrations(db); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
	return dataDir, db, vlt
}

func startServer() {
	log.Printf(" Booting Vessl Daemon (`vessld`) v%s [%s/%s]...", vesslVersion, runtime.GOOS, runtime.GOARCH)
	_, db, vlt := initDataDir()
	defer db.Close()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" Docker daemon connection warning: %v (container deployment features disabled)", err)
	}

	traefikMgr := engine.NewTraefikManager(dockerClient, os.Getenv("VESSL_TLS_EMAIL"))
	if err := traefikMgr.EnsureTraefikRunning(context.Background()); err != nil {
		log.Printf(" Warning: Failed to start Traefik proxy: %v", err)
	}

	services.StartTelemetryReporter(db, vesslVersion)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})
	apiServer := vesslhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient)

	log.Printf(" Vessl control plane listening on %s", addr)
	if err := http.ListenAndServe(addr, apiServer.Handler()); err != nil {
		log.Fatalf(" Server crashed: %v", err)
	}
}

func runMCP() {
	log.Printf("Starting MCP stdio server...")
	_, db, vlt := initDataDir()
	defer db.Close()

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})
	traefikMgr := engine.NewTraefikManager(dockerClient, os.Getenv("VESSL_TLS_EMAIL"))
	apiServer := vesslhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient)

	if err := apiServer.StartMCPStdio(); err != nil {
		log.Fatalf("MCP Server exited: %v", err)
	}
}
