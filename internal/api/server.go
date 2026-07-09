package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/auth"
	"vessel.dev/vessel/internal/backup"
	"vessel.dev/vessel/internal/canvas"
	"vessel.dev/vessel/internal/database"
	"vessel.dev/vessel/internal/deployment"
	"vessel.dev/vessel/internal/domain"
	"vessel.dev/vessel/internal/env"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/git"
	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/notification"
	"vessel.dev/vessel/internal/notifier"
	"vessel.dev/vessel/internal/oauth"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/project_settings"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/service_var"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/storage"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/team"
	"vessel.dev/vessel/internal/terminal"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/updater"
	"vessel.dev/vessel/internal/user"
	"vessel.dev/vessel/internal/workspace"
)

// Server encapsulates HTTP routing, API handler dependencies, and authentication guards for the Vessel control plane.
type Server struct {
	router                 *http.ServeMux
	store                  *store.Store
	deployer               *orchestrator.Deployer
	proxyManager           *proxy.ProxyManager
	dockerClient           *client.Client
	tokenService           *services.TokenService
	authGuard              *middleware.AuthGuard
	dbDeployer             *orchestrator.DatabaseDeployer
	storageDeployer        *orchestrator.StorageDeployer
	cronManager            *orchestrator.CronManager
	cronService            *services.CronService
	serviceLinker          *services.ServiceLinker
	gitService             *git.Service
	serviceHandler         *service.Handler
	dbHandler              *database.Handler
	storageHandler         *storage.Handler
	jobHandler             *job.Handler
	canvasHandler          *canvas.Handler
	terminalHandler        *terminal.Handler
	deploymentHandler      *deployment.Handler
	serviceVarHandler      *service_var.Handler
	projectSettingsHandler *project_settings.Handler
	backupHandler          *backup.Handler
	teamHandler            *team.Handler
	workspaceHandler       *workspace.Handler
	settingsHandler        *settings.Handler
	updaterHandler         *updater.Handler
	userHandler            *user.Handler
	authHandler            *auth.Handler
	oauthHandler           *oauth.Handler
	gitHandler             *git.Handler
	webhookHandler         *git.WebhookHandler
	projectHandler         *project.Handler
	environmentHandler     *environment.Handler
	domainHandler          *domain.Handler
	projectEnvHandler      *env.Handler
	notifierService        *notifier.NotifierService
	notificationHandler    *notification.Handler
	updaterService         *updater.UpdaterService
}

// NewServer initializes a Server wired to the database store, container orchestrator, reverse proxy, and Docker client.
func NewServer(s *store.Store, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	cronMgr := orchestrator.NewCronManager(dockerClient, s)
	_ = cronMgr.Start()

	backupMgr := orchestrator.NewBackupManager(dockerClient, s, "")
	_ = backupMgr.Start()

	tokenService := services.NewTokenService()

	// Domain repositories
	settingsRepo := settings.NewSQLiteRepository(s.DB())
	userRepo := user.NewSQLiteRepository(s.DB())
	oauthRepo := oauth.NewSQLiteRepository(s.DB())
	notifRepo := notification.NewSQLiteRepository(s.DB())

	notifierService := notifier.NewNotifierService(notifRepo)

	// Domain services
	settingsService := settings.NewService(settingsRepo)
	userService := user.NewService(userRepo)
	authService := auth.NewService(userRepo, settingsRepo, tokenService)
	oauthService := oauth.NewService(oauthRepo, userRepo, tokenService)

	updaterService := updater.NewUpdaterService(settingsRepo)
	updaterService.Start(context.Background())

	// Claims extractor helpers
	extractUserID := func(r *http.Request) string {
		if c := GetUserClaimsFromContext(r.Context()); c != nil {
			return c.UserID
		}
		return ""
	}
	extractClaims := func(r *http.Request) (userID, email string) {
		if c := GetUserClaimsFromContext(r.Context()); c != nil {
			return c.UserID, c.Email
		}
		return "", ""
	}
	extractClaims3 := func(r *http.Request) (userID, email, role string) {
		if c := GetUserClaimsFromContext(r.Context()); c != nil {
			return c.UserID, c.Email, c.Role
		}
		return "", "", ""
	}

	// Git domain
	gitRepo := git.NewSQLiteRepository(s.DB(), s.Vault())
	gitService := git.NewService(gitRepo, nil)
	gitService.WithProjectService(&gitProjectAdapter{store: s})
	gitHandler := git.NewHandler(gitService, extractUserID)

	// Project, environment, domain, and project-env domains
	envRepo := environment.NewSQLiteRepository(s.DB())
	envService := environment.NewService(envRepo)
	envHandler := environment.NewHandler(envService)

	domainRepo := domain.NewSQLiteRepository(s.DB())
	domainService := domain.NewService(domainRepo)
	domainHandler := domain.NewHandler(domainService, proxyManager)

	projectEnvRepo := env.NewSQLiteRepository(s.DB(), s.Vault())
	projectEnvService := env.NewService(projectEnvRepo)
	projectEnvHandler := env.NewHandler(projectEnvService)

	projectRepo := project.NewSQLiteRepository(s.DB(), envRepo)
	projectService := project.NewService(projectRepo, &appServiceRepoAdapter{store: s})
	projectHandler := project.NewHandler(projectService, proxyManager, extractUserID)

	// ── New domain packages ──────────────────────────────────────────

	// Service (app services)
	serviceRepo := service.NewSQLiteRepository(s.DB())
	serviceHandler := service.NewHandler(serviceRepo)

	// Database
	dbRepo := database.NewSQLiteRepository(s.DB(), s.Vault())
	dbHandler := database.NewHandler(dbRepo, &dbDeployerAdapter{inner: orchestrator.NewDatabaseDeployer(dockerClient, s)})

	// Storage
	storageRepo := storage.NewSQLiteRepository(s.DB(), s.Vault())
	storageHandler := storage.NewHandler(storageRepo, &storageDeployerAdapter{inner: orchestrator.NewStorageDeployer(dockerClient, s)})

	// Jobs
	jobRepo := job.NewSQLiteRepository(s.DB())
	jobHandler := job.NewHandler(jobRepo)

	// Canvas
	canvasRepo := canvas.NewSQLiteRepository(s.DB(), envRepo)
	canvasHandler := canvas.NewHandler(canvasRepo)

	// Deployment
	deploymentRepo := deployment.NewSQLiteRepository(s.DB())
	deploymentHandler := deployment.NewHandler(deploymentRepo, &deploymentSvcAdapter{svcRepo: serviceRepo}, &deploymentProjectStoreAdapter{store: s}, &deploymentProjectDeployerAdapter{gitService: gitService, deployer: deployer, proxyManager: proxyManager})

	// Service Variables
	svVarRepo := service_var.NewSQLiteRepository(s.DB(), &serviceVarSvcAdapter{svcRepo: serviceRepo})
	svVarHandler := service_var.NewHandler(svVarRepo, &serviceVarSvcAdapter{svcRepo: serviceRepo})

	// Terminal
	terminalHandler := terminal.NewHandler(dockerClient, &tokenValidatorAdapter{inner: tokenService}, &terminalSvcAdapter{svcRepo: serviceRepo})

	// Backup
	backupRepo := backup.NewSQLiteRepository(s.DB())
	backupHandler := backup.NewHandler(backupRepo, &backupManagerAdapter{inner: backupMgr})

	// Team
	teamRepo := team.NewSQLiteRepository(s.DB())
	teamHandler := team.NewHandler(teamRepo, &teamUserProviderAdapter{userRepo: userRepo}, extractClaims3)

	// Workspace
	wsRepo := workspace.NewSQLiteRepository(s.DB())
	workspaceHandler := workspace.NewHandler(wsRepo, extractClaims3)

	// Project Settings
	psRepo := project_settings.NewSQLiteRepository(s.DB())
	projectSettingsHandler := project_settings.NewHandler(psRepo, &projectSettingsUserProviderAdapter{userRepo: userRepo}, extractUserID)

	srv := &Server{
		router:                 http.NewServeMux(),
		store:                  s,
		deployer:               deployer,
		proxyManager:           proxyManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              middleware.NewAuthGuard(tokenService, s),
		dbDeployer:             orchestrator.NewDatabaseDeployer(dockerClient, s),
		storageDeployer:        orchestrator.NewStorageDeployer(dockerClient, s),
		cronManager:            cronMgr,
		cronService:            services.NewCronService(s, cronMgr),
		serviceLinker:          services.NewServiceLinker(s),
		serviceHandler:         serviceHandler,
		dbHandler:              dbHandler,
		storageHandler:         storageHandler,
		jobHandler:             jobHandler,
		canvasHandler:          canvasHandler,
		terminalHandler:        terminalHandler,
		deploymentHandler:      deploymentHandler,
		serviceVarHandler:      svVarHandler,
		projectSettingsHandler: projectSettingsHandler,
		backupHandler:          backupHandler,
		teamHandler:            teamHandler,
		workspaceHandler:       workspaceHandler,
		settingsHandler:        settings.NewHandler(settingsService, dockerClient),
		updaterHandler:         updater.NewHandler(updaterService),
		userHandler:            user.NewHandler(userService, extractUserID),
		authHandler:            auth.NewHandler(authService, extractUserID),
		oauthHandler:           oauth.NewHandler(oauthService, extractClaims),
		gitHandler:             gitHandler,
		webhookHandler:         git.NewWebhookHandler(s, gitService, deployer, proxyManager),
		projectHandler:         projectHandler,
		environmentHandler:     envHandler,
		domainHandler:          domainHandler,
		projectEnvHandler:      projectEnvHandler,
		notifierService:        notifierService,
		notificationHandler: func() *notification.Handler {
			notifService := notification.NewService(notifRepo, notifierService)
			return notification.NewHandler(notifService)
		}(),
		updaterService: updaterService,
	}
	if srv.deployer != nil {
		srv.deployer.EnvProvider = srv.serviceLinker.GetLinkedEnvironmentVariables
	}
	srv.registerRoutes()
	return srv
}

// ServeHTTP satisfies the http.Handler interface, routing through the registered mux with CORS middleware.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.CORSMiddleware(s.router).ServeHTTP(w, r)
}

// Handler returns the root HTTP handler wrapped with global CORS and authentication middleware.
func (s *Server) Handler() http.Handler {
	return middleware.CORSMiddleware(s.router)
}

// RequireAuth validates Bearer tokens or query parameters via middleware before invoking the handler.
func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return s.authGuard.RequireAuth(next)
}

// RequireRole enforces that the authenticated user possesses the specified role via middleware.
func (s *Server) RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return s.authGuard.RequireRole(requiredRole, next)
}

// GetUserClaimsFromContext retrieves the authenticated user's claims from request context via middleware.
func GetUserClaimsFromContext(ctx context.Context) *user.UserClaims {
	return middleware.GetUserClaimsFromContext(ctx)
}

// ── Legacy adapters ────────────────────────────────────────────────────

// gitProjectAdapter bridges the legacy store app-service query to the git.ProjectService interface.
type gitProjectAdapter struct {
	store *store.Store
}

func (a *gitProjectAdapter) ListAppServicesByProject(projectID string) ([]*git.AppService, error) {
	apps, err := a.store.ListAppServicesByProject(projectID)
	if err != nil {
		return nil, err
	}
	var result []*git.AppService
	for _, app := range apps {
		result = append(result, &git.AppService{
			ID:            app.ID,
			ProjectID:     app.ProjectID,
			EnvironmentID: app.EnvironmentID,
			Name:          app.Name,
			RepositoryURL: app.RepositoryURL,
			Branch:        app.Branch,
			ContainerID:   app.ContainerID,
		})
	}
	return result, nil
}

// appServiceRepoAdapter bridges the legacy store app-service creation to the project.AppServiceRepository interface.
type appServiceRepoAdapter struct {
	store *store.Store
}

func (a *appServiceRepoAdapter) CreateAppService(_ context.Context, app *types.AppServiceConfig) error {
	return a.store.CreateAppService(app)
}

// dbDeployerAdapter bridges the orchestrator.DatabaseDeployer to the database.Deployer interface.
type dbDeployerAdapter struct {
	inner *orchestrator.DatabaseDeployer
}

func (a *dbDeployerAdapter) SpinUp(ctx context.Context, db *database.Database) (string, error) {
	cfg := &types.DatabaseConfig{
		ID: db.ID, ProjectID: db.ProjectID, EnvironmentID: db.EnvironmentID,
		Name: db.Name, Engine: db.Engine, Version: db.Version,
		Port: db.Port, Username: db.Username, Password: db.Password,
		DatabaseName: db.DatabaseName, VolumePath: db.VolumePath,
		ContainerID: db.ContainerID, Status: db.Status,
		InternalDNS: db.InternalDNS, ExternalDNS: db.ExternalDNS,
		CreatedAt: db.CreatedAt, UpdatedAt: db.UpdatedAt,
	}
	return a.inner.SpinUp(ctx, cfg)
}

func (a *dbDeployerAdapter) Stop(ctx context.Context, id string) error {
	return a.inner.Stop(ctx, id)
}

// storageDeployerAdapter bridges the orchestrator.StorageDeployer to the storage.Deployer interface.
type storageDeployerAdapter struct {
	inner *orchestrator.StorageDeployer
}

func (a *storageDeployerAdapter) SpinUp(ctx context.Context, s *storage.Storage) (string, error) {
	cfg := &types.StorageConfig{
		ID: s.ID, ProjectID: s.ProjectID, EnvironmentID: s.EnvironmentID,
		Name: s.Name, Type: s.Type, APIPort: s.APIPort, ConsolePort: s.ConsolePort,
		AccessKey: s.AccessKey, SecretKey: s.SecretKey, BucketName: s.BucketName,
		VolumePath: s.VolumePath, ContainerID: s.ContainerID, Status: s.Status,
		InternalDNS: s.InternalDNS, ExternalDNS: s.ExternalDNS,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
	return a.inner.SpinUp(ctx, cfg)
}

func (a *storageDeployerAdapter) Stop(ctx context.Context, id string) error {
	return a.inner.Stop(ctx, id)
}

// backupManagerAdapter bridges the orchestrator.BackupManager to the backup.BackupManager interface.
type backupManagerAdapter struct {
	inner *orchestrator.BackupManager
}

func (a *backupManagerAdapter) RegisterBackup(cfg *backup.BackupConfig) error {
	return a.inner.RegisterBackup(toTypesBackupConfig(cfg))
}

func (a *backupManagerAdapter) UnregisterBackup(backupConfigID string) {
	a.inner.UnregisterBackup(backupConfigID)
}

func (a *backupManagerAdapter) TriggerBackup(ctx context.Context, backupConfigID string) (*backup.BackupRecord, error) {
	rec, err := a.inner.TriggerBackup(ctx, backupConfigID)
	if err != nil {
		return nil, err
	}
	return fromTypesBackupRecord(rec), nil
}

func toTypesBackupConfig(cfg *backup.BackupConfig) *types.BackupConfig {
	return &types.BackupConfig{
		ID: cfg.ID, ProjectID: cfg.ProjectID, DatabaseID: cfg.DatabaseID,
		StorageID: cfg.StorageID, S3DestinationID: cfg.S3DestinationID,
		Name: cfg.Name, Schedule: cfg.Schedule, RetentionDays: cfg.RetentionDays,
		Status: cfg.Status, CreatedAt: cfg.CreatedAt, UpdatedAt: cfg.UpdatedAt,
	}
}

func fromTypesBackupRecord(rec *types.BackupRecord) *backup.BackupRecord {
	return &backup.BackupRecord{
		ID: rec.ID, BackupConfigID: rec.BackupConfigID, ProjectID: rec.ProjectID,
		DatabaseID: rec.DatabaseID, Status: rec.Status, FilePath: rec.FilePath,
		FileSizeBytes: rec.FileSizeBytes, S3URL: rec.S3URL, Logs: rec.Logs,
		StartedAt: rec.StartedAt, CompletedAt: rec.CompletedAt,
	}
}

// serviceVarSvcAdapter bridges the service repository to the service_var.ServiceRepository interface.
type serviceVarSvcAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *serviceVarSvcAdapter) GetByID(ctx context.Context, id string) (*service_var.ServiceDTO, error) {
	app, err := a.svcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}
	return &service_var.ServiceDTO{
		ID:            app.ID,
		ProjectID:     app.ProjectID,
		EnvironmentID: app.EnvironmentID,
	}, nil
}

// terminalSvcAdapter bridges the service repository to the terminal.ServiceRepository interface.
type terminalSvcAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *terminalSvcAdapter) GetByID(ctx context.Context, id string) (*terminal.AppService, error) {
	app, err := a.svcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}
	return &terminal.AppService{ID: app.ID, ContainerID: app.ContainerID}, nil
}

// deploymentProjectStoreAdapter bridges the store to the deployment.ProjectStore interface.
type deploymentProjectStoreAdapter struct {
	store *store.Store
}

func (a *deploymentProjectStoreAdapter) GetByID(ctx context.Context, id string) (*deployment.ProjectConfig, error) {
	p, err := a.store.GetProject(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return &deployment.ProjectConfig{ID: p.ID, Name: p.Name, Description: p.Description, TeamID: p.TeamID}, nil
}

// deploymentProjectDeployerAdapter bridges gitService, deployer, and proxyManager into a single ProjectDeployer.
type deploymentProjectDeployerAdapter struct {
	gitService   *git.Service
	deployer     *orchestrator.Deployer
	proxyManager *proxy.ProxyManager
}

func (a *deploymentProjectDeployerAdapter) CloneOrPullRepository(ctx context.Context, projectID, sourceDir string) error {
	if a.gitService == nil {
		return nil
	}
	return a.gitService.CloneOrPullRepository(ctx, projectID, sourceDir, nil)
}

func (a *deploymentProjectDeployerAdapter) DeployProject(ctx context.Context, project *deployment.ProjectConfig, sourceDir string) (string, error) {
	if a.deployer == nil {
		return "", fmt.Errorf("deployer not available")
	}
	p := &types.ProjectConfig{ID: project.ID, Name: project.Name, Description: project.Description, TeamID: project.TeamID}
	return a.deployer.Deploy(ctx, p, sourceDir, nil)
}

func (a *deploymentProjectDeployerAdapter) ReloadProxy(ctx context.Context) error {
	if a.proxyManager == nil {
		return nil
	}
	return a.proxyManager.Reload(ctx)
}

// deploymentSvcAdapter bridges the service repository to the deployment.ServiceRepository interface.
type deploymentSvcAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *deploymentSvcAdapter) GetByID(ctx context.Context, id string) (any, error) {
	return a.svcRepo.GetByID(ctx, id)
}

// teamUserProviderAdapter bridges the user repository to the team.UserProvider interface.
type teamUserProviderAdapter struct {
	userRepo *user.SQLiteRepository
}

func (a *teamUserProviderAdapter) GetUserByEmail(email string) (*user.User, error) {
	return a.userRepo.GetUserByEmail(context.Background(), email)
}

// projectSettingsUserProviderAdapter bridges the user repository to the project_settings.UserProvider interface.
type projectSettingsUserProviderAdapter struct {
	userRepo *user.SQLiteRepository
}

func (a *projectSettingsUserProviderAdapter) GetUserByEmail(ctx context.Context, email string) (*project_settings.User, error) {
	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return &project_settings.User{ID: u.ID, Email: u.Email}, nil
}

// tokenValidatorAdapter bridges the services.TokenService to the terminal.TokenValidator interface.
type tokenValidatorAdapter struct {
	inner *services.TokenService
}

func (a *tokenValidatorAdapter) ValidateToken(tokenStr string) (*terminal.TokenClaim, error) {
	claims, err := a.inner.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	return &terminal.TokenClaim{UserID: sub, Email: email}, nil
}
