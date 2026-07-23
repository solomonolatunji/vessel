package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"codedock.dev/codedock/internal/engine"
	codedockhttp "codedock.dev/codedock/internal/http"
	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/repositories"
	"codedock.dev/codedock/internal/services"
	"codedock.dev/codedock/internal/telemetry"
	"codedock.dev/codedock/internal/utils"
)

var codedockVersion = "dev"

type dbDeployerStore struct {
	db    *sql.DB
	vault *utils.Vault
}

func (a *dbDeployerStore) GetServerSettings() (*models.ServerSettings, error) {
	return repositories.NewSettingsRepo(a.db).GetServerSettings(context.Background())
}

func (a *dbDeployerStore) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return repositories.NewAppServiceRepo(a.db).ListByProject(context.Background(), projectID)
}

func (a *dbDeployerStore) GetEnvVars(projectID string) (map[string]string, error) {
	return repositories.NewEnvRepo(a.db, a.vault).GetVars(context.Background(), projectID)
}

func (a *dbDeployerStore) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	svVarRepo := repositories.NewServiceVarRepo(a.db)
	return svVarRepo.ListByService(context.Background(), serviceID)
}

func (a *dbDeployerStore) ListLogDrainsByService(serviceID string) ([]*models.LogDrain, error) {
	return repositories.NewAppServiceRepo(a.db).ListLogDrainsByService(context.Background(), serviceID)
}

func (a *dbDeployerStore) GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error) {
	svlsRepo := repositories.NewServerlessRepository(a.db)
	return svlsRepo.GetCodeByServiceID(context.Background(), serviceID)
}

func (a *dbDeployerStore) UpdateAppService(app *models.AppService) error {
	repo := repositories.NewAppServiceRepo(a.db)
	return repo.Update(context.Background(), app)
}

func main() {
	_ = godotenv.Load(".env")
	mainCLI()
}

func initDataDir() (string, *sql.DB, *utils.Vault) {
	dataDir := os.Getenv("CODEDOCK_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		slog.Error("failed to create data directory", "err", err)
		os.Exit(1)
	}
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		slog.Error("failed to initialize secrets vault", "err", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(dataDir, "codedock.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		slog.Error("failed to open SQLite database", "err", err)
		os.Exit(1)
	}
	if err := repositories.RunMigrations(db); err != nil {
		slog.Error("failed to run database migrations", "err", err)
		os.Exit(1)
	}
	return dataDir, db, vlt
}

func startServer() {
	slog.Info("booting daemon", "version", codedockVersion, "os", runtime.GOOS, "arch", runtime.GOARCH)
	dataDir, db, vlt := initDataDir()
	defer db.Close()

	telemetry.Init()
	defer telemetry.Close()
	telemetry.Track("system", "daemon_start", map[string]interface{}{
		"version": codedockVersion,
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	})

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Warn("Docker daemon connection warning", "err", err, "detail", "container deployment features disabled")
	}

	traefikMgr := engine.NewTraefikManager(dockerClient, os.Getenv("CODEDOCK_TLS_EMAIL"))
	if err := traefikMgr.EnsureTraefikRunning(context.Background()); err != nil {
		slog.Warn("failed to start Traefik proxy", "err", err)
	}

	tsdbMgr := engine.NewTSDBManager(dockerClient)
	if err := tsdbMgr.EnsureTSDBRunning(context.Background()); err != nil {
		slog.Warn("failed to start TSDB", "err", err)
	}

	lokiMgr := engine.NewLokiManager(dockerClient)
	if err := lokiMgr.EnsureLokiRunning(context.Background()); err != nil {
		slog.Warn("failed to start Loki", "err", err)
	}

	metricsWorker := engine.NewMetricsWorker(dockerClient)
	metricsWorker.Start()

	logWorker := engine.NewLogWorker(dockerClient)
	logWorker.Start(context.Background())

	services.StartTelemetryReporter(db, codedockVersion)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})
	apiServer, err := codedockhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient, dataDir)
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: apiServer.Handler(),
	}

	go func() {
		slog.Info("control plane listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server crashed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}
	slog.Info("server exited")
}

func runMCP() {
	slog.Info("starting MCP stdio server")
	_, db, vlt := initDataDir()
	defer db.Close()

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	deployer := engine.NewDeployer(dockerClient, &dbDeployerStore{db: db, vault: vlt})
	traefikMgr := engine.NewTraefikManager(dockerClient, os.Getenv("CODEDOCK_TLS_EMAIL"))
	apiServer, err := codedockhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient, "")
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		os.Exit(1)
	}

	if err := apiServer.StartMCPStdio(); err != nil {
		slog.Error("MCP server exited", "err", err)
		os.Exit(1)
	}
}
