package http

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"vessel.dev/vessel/internal/dispatch"
	"vessel.dev/vessel/internal/engine"
	"vessel.dev/vessel/internal/handlers"
	"vessel.dev/vessel/internal/listeners"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/repositories"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/vault"
)

type Server struct {
	router                 *echo.Echo
	deployer               *engine.Deployer
	traefikManager         *proxy.TraefikManager
	dockerClient           *client.Client
	tokenService           *services.TokenService
	authGuard              *middleware.AuthGuard
	cronManager            *engine.CronManager
	serviceLinker          *services.ServiceLinker
	dispatcherService      *dispatch.DispatcherService
	appServiceHandler      *handlers.AppHandler
	dbHandler              *handlers.DatabaseHandler
	storageHandler         *handlers.StorageHandler
	jobHandler             *handlers.JobHandler
	canvasHandler          *handlers.CanvasHandler
	terminalHandler        *handlers.TerminalHandler
	deploymentHandler      *handlers.DeploymentHandler
	serviceVarHandler      *handlers.ServiceVarHandler
	projectSettingsHandler *handlers.ProjectSettingsHandler
	backupHandler          *handlers.BackupHandler
	teamHandler            *handlers.TeamHandler
	workspaceHandler       *handlers.WorkspaceHandler
	settingsHandler        *handlers.SettingsHandler
	updaterHandler         *handlers.UpdaterHandler
	userHandler            *handlers.UserHandler
	authHandler            *handlers.AuthHandler
	oauthHandler           *handlers.OAuthHandler
	gitHandler             *handlers.GitHandler
	webhookHandler         *handlers.WebhookHandler
	projectHandler         *handlers.ProjectHandler
	environmentHandler     *handlers.EnvironmentHandler
	domainHandler          *handlers.DomainHandler
	projectEnvHandler      *handlers.ProjectEnvHandler
	notificationHandler    *handlers.NotificationHandler
	gitAppsHandler         *handlers.GitAppsHandler
	aiSettingsHandler      *handlers.AISettingsHandler
	aiDiagnosticsHandler   *handlers.AIDiagnosticsHandler
	vercelHandler          *handlers.VercelHandler
}

func NewServer(db *sql.DB, vault *vault.Vault, deployer *engine.Deployer, traefikManager *proxy.TraefikManager, dockerClient *client.Client) *Server {
	settingsRepo := repositories.NewSettingsSQLiteRepository(db)
	userRepo := repositories.NewUserSQLiteRepository(db)
	oauthRepo := repositories.NewOAuthSQLiteRepository(db)
	notifRepo := repositories.NewNotificationSQLiteRepository(db)
	serviceRepo := repositories.NewAppServiceSQLiteRepository(db)
	gitRepo := repositories.NewGitSQLiteRepository(db, vault)
	envRepo := repositories.NewEnvSQLiteRepository(db, vault)
	environmentRepo := repositories.NewEnvironmentSQLiteRepository(db)
	domainRepo := repositories.NewDomainSQLiteRepository(db)
	projectRepo := repositories.NewProjectSQLiteRepository(db, environmentRepo)
	databaseRepo := repositories.NewDatabaseSQLiteRepository(db, vault)
	storageRepo := repositories.NewStorageSQLiteRepository(db, vault)
	jobRepo := repositories.NewJobSQLiteRepository(db)
	canvasRepo := repositories.NewCanvasSQLiteRepository(db, environmentRepo)
	deploymentRepo := repositories.NewDeploymentSQLiteRepository(db)
	backupRepo := repositories.NewBackupSQLiteRepository(db)
	teamRepo := repositories.NewTeamSQLiteRepository(db)
	wsRepo := repositories.NewWorkspaceSQLiteRepository(db)
	psRepo := repositories.NewProjectSettingsSQLiteRepository(db)
	svVarRepo := repositories.NewServiceVarSQLiteRepository(db)
	s3Repo := repositories.NewS3DestinationSQLiteRepository(db)
	prPreviewRepo := repositories.NewPRPreviewRepository(db)
	ea := newEngineAdapter(settingsRepo, serviceRepo, envRepo, databaseRepo, storageRepo, projectRepo, jobRepo, backupRepo, s3Repo, svVarRepo)
	cronMgr := engine.NewCronManager(dockerClient, ea)
	_ = cronMgr.Start()
	backupMgr := engine.NewBackupManager(dockerClient, ea, "")
	_ = backupMgr.Start()
	dbDeployer := engine.NewDatabaseDeployer(dockerClient, ea)
	storageDeployer := engine.NewStorageDeployer(dockerClient, ea)
	svcLinker := services.NewServiceLinker(databaseRepo, storageRepo)
	dispatcherSvc := dispatch.NewDispatcherService(notifRepo)
	deploymentListeners := listeners.NewDeploymentListeners(dispatcherSvc)
	deploymentListeners.Register()
	settingsService := services.NewSettingsService(settingsRepo, notifRepo)
	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService()
	authService := services.NewAuthService(userRepo, settingsRepo, tokenService)
	oauthService := services.NewOAuthService(oauthRepo, userRepo, tokenService)
	updaterService := services.NewUpdaterService(settingsRepo)
	updaterService.Start(context.Background())
	gitService := services.NewGitService(gitRepo)
	environmentService := services.NewEnvironmentService(environmentRepo, domainRepo, envRepo)
	appService := services.NewAppService(serviceRepo, svVarRepo)
	projectService := services.NewProjectService(projectRepo, environmentRepo, serviceRepo, svVarRepo)
	dbService := services.NewDatabaseService(databaseRepo, dbDeployer)
	storageService := services.NewStorageService(storageRepo, storageDeployer)
	jobService := services.NewJobService(jobRepo, cronMgr)
	canvasService := services.NewCanvasService(canvasRepo)
	backupService := services.NewBackupService(backupRepo, s3Repo, backupMgr)
	teamService := services.NewTeamService(teamRepo, userRepo)
	wsService := services.NewWorkspaceService(wsRepo)
	psService := services.NewProjectSettingsService(psRepo, userRepo)
	notificationService := services.NewNotificationService(notifRepo, dispatcherSvc)
	deploymentService := services.NewDeploymentService(deploymentRepo, serviceRepo, projectRepo, deployer)
	prPreviewService := services.NewPRPreviewService(prPreviewRepo, appService, gitService, deployer)

	gitAppsRepo := repositories.NewGitAppSQLiteRepository(db, vault)
	gitAppsService := services.NewGitAppsService(gitAppsRepo)

	aiSettingsRepo := repositories.NewTeamAISettingsSQLiteRepository(db, vault)
	aiSettingsService := services.NewAISettingsService(aiSettingsRepo)

	vercelRepo := repositories.NewVercelRepository(db, vault)
	vercelService := services.NewVercelService(vercelRepo)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, psService)
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	srv := &Server{
		router:                 e,
		deployer:               deployer,
		traefikManager:         traefikManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              authGuard,
		cronManager:            cronMgr,
		serviceLinker:          svcLinker,
		dispatcherService:      dispatcherSvc,
		appServiceHandler:      handlers.NewAppHandler(appService),
		dbHandler:              handlers.NewDatabaseHandler(dbService),
		storageHandler:         handlers.NewStorageHandler(storageService),
		jobHandler:             handlers.NewJobHandler(jobService),
		canvasHandler:          handlers.NewCanvasHandler(canvasService),
		terminalHandler:        handlers.NewTerminalHandler(dockerClient, tokenService, appService),
		deploymentHandler:      handlers.NewDeploymentHandler(deploymentService, appService),
		serviceVarHandler:      handlers.NewServiceVarHandler(appService),
		projectSettingsHandler: handlers.NewProjectSettingsHandler(psService),
		backupHandler:          handlers.NewBackupHandler(backupService),
		teamHandler:            handlers.NewTeamHandler(teamService),
		workspaceHandler:       handlers.NewWorkspaceHandler(wsService),
		settingsHandler:        handlers.NewSettingsHandler(settingsService),
		updaterHandler:         handlers.NewUpdaterHandler(updaterService),
		userHandler:            handlers.NewUserHandler(userService),
		authHandler:            handlers.NewAuthHandler(authService),
		oauthHandler:           handlers.NewOAuthHandler(oauthService),
		gitHandler:             handlers.NewGitHandler(gitService),
		webhookHandler:         handlers.NewWebhookHandler(gitService, projectService, appService, deploymentService, prPreviewService),
		projectHandler:         handlers.NewProjectHandler(projectService),
		environmentHandler:     handlers.NewEnvironmentHandler(environmentService),
		domainHandler:          handlers.NewDomainHandler(environmentService),
		projectEnvHandler:      handlers.NewProjectEnvHandler(environmentService),
		notificationHandler:    handlers.NewNotificationHandler(notificationService),
		gitAppsHandler:         handlers.NewGitAppsHandler(gitAppsService),
		aiSettingsHandler:      handlers.NewAISettingsHandler(aiSettingsService),
		aiDiagnosticsHandler:   handlers.NewAIDiagnosticsHandler(aiSettingsService, deploymentService, projectService),
		vercelHandler:          handlers.NewVercelHandler(vercelService),
	}
	if srv.deployer != nil {
		srv.deployer.EnvProvider = func(projectID string) (map[string]string, error) {
			return srv.serviceLinker.GetLinkedEnvironmentVariables(context.Background(), projectID)
		}
	}
	srv.registerRoutes()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func GetUserClaimsFromContext(ctx context.Context) *models.UserClaims {
	return middleware.GetUserClaimsFromContext(ctx)
}
