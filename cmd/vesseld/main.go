package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"vessel.dev/vessel/internal/agent"
	vesseldb "vessel.dev/vessel/internal/db"
	"vessel.dev/vessel/internal/engine"
	vesselhttp "vessel.dev/vessel/internal/http"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/repositories"
	"vessel.dev/vessel/internal/vault"
)

const vesselVersion = "0.1.0-alpha"

// Unused proxy listers removed

type dbDeployerStore struct {
	db    *sql.DB
	vault *vault.Vault
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
	isAgent := flag.Bool("agent", false, "Run in agent mode")
	agentToken := flag.String("token", "", "Agent auth token")
	serverURL := flag.String("server", "", "Controller server WSS URL")
	flag.Parse()
	log.Printf(" Booting Vessel Daemon (`vesseld`) v%s [%s/%s]...", vesselVersion, runtime.GOOS, runtime.GOARCH)
	if *isAgent {
		if *serverURL == "" {
			log.Fatal(" Error: --server is required in agent mode (e.g. wss://vessel.domain.com/api/agent)")
		}
		if *agentToken == "" {
			log.Fatal(" Error: --token is required in agent mode")
		}
		if err := agent.Run(context.Background(), *serverURL, *agentToken); err != nil {
			log.Fatalf(" Agent mode exited: %v", err)
		}
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dataDir := os.Getenv("VESSEL_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		log.Fatalf(" Failed to create data directory: %v", err)
	}
	vlt, err := vault.NewVault(dataDir)
	if err != nil {
		log.Fatalf(" Failed to initialize secrets vault: %v", err)
	}
	dbPath := filepath.Join(dataDir, "vessel.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		log.Fatalf(" Failed to open SQLite database: %v", err)
	}
	defer db.Close()
	if err := vesseldb.RunMigrations(db); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" Docker daemon connection warning: %v (container deployment features disabled)", err)
	}
	traefikMgr := proxy.NewTraefikManager(dockerClient, os.Getenv("VESSEL_TLS_EMAIL"))
	if err := traefikMgr.EnsureTraefikRunning(context.Background()); err != nil {
		log.Printf(" Warning: Failed to start Traefik proxy: %v", err)
	}

	host := os.Getenv("HOST")

	addr := host + ":" + port

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})
	apiServer := vesselhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient)
	log.Printf(" Vessel control plane listening on %s", addr)
	if err := http.ListenAndServe(addr, apiServer.Handler()); err != nil {
		log.Fatalf(" Server crashed: %v", err)
	}
}
