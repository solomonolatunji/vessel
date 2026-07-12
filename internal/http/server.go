package http

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mark3labs/mcp-go/server"

	"vessl.dev/vessl/internal/core"
	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/handlers"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/mcp"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/notifications"
	"vessl.dev/vessl/internal/proxy"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/vault"
)

type Server struct {
	router                 *echo.Echo
	mcpBridge              *mcp.Bridge
	deployer               *engine.Deployer
	traefikManager         *proxy.TraefikManager
	dockerClient           *client.Client
	tokenService           *services.TokenService
	authGuard              *middleware.AuthGuard
	cronManager            *engine.CronManager
	serviceLinker          *services.ServiceLinker
	dispatcherService      *core.DispatcherService
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
	emailSettingsHandler   *handlers.EmailSettingsHandler
	aiDiagnosticsHandler   *handlers.AIDiagnosticsHandler
	vercelHandler          *handlers.VercelHandler
	serverlessHandler      *handlers.ServerlessHandler
}

func NewServer(db *sql.DB, v *vault.Vault, deployer *engine.Deployer, traefikManager *proxy.TraefikManager, dockerClient *client.Client) *Server {
	// Base router setup
	e := echo.New()
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			log.Printf("REQUEST: %s %s | status: %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Repositories
	environmentSQLiteRepository := repositories.NewEnvironmentSQLiteRepository(db)
	projectSQLiteRepository := repositories.NewProjectSQLiteRepository(db, environmentSQLiteRepository)
	appServiceSQLiteRepository := repositories.NewAppServiceSQLiteRepository(db)
	serviceVarSQLiteRepository := repositories.NewServiceVarSQLiteRepository(db)
	databaseSQLiteRepository := repositories.NewDatabaseSQLiteRepository(db, v)
	settingsSQLiteRepository := repositories.NewSettingsSQLiteRepository(db)
	envSQLiteRepository := repositories.NewEnvSQLiteRepository(db, v)
	storageSQLiteRepository := repositories.NewStorageSQLiteRepository(db, v)
	jobSQLiteRepository := repositories.NewJobSQLiteRepository(db)
	backupSQLiteRepository := repositories.NewBackupSQLiteRepository(db)
	s3DestinationSQLiteRepository := repositories.NewS3DestinationSQLiteRepository(db)
	serverlessRepository := repositories.NewServerlessRepository(db)
	notificationSQLiteRepository := repositories.NewNotificationSQLiteRepository(db)
	projectSettingsSQLiteRepository := repositories.NewProjectSettingsSQLiteRepository(db)
	userSQLiteRepository := repositories.NewUserSQLiteRepository(db)
	teamEmailSettingsSQLiteRepository := repositories.NewTeamEmailSettingsSQLiteRepository(db, v)
	canvasSQLiteRepository := repositories.NewCanvasSQLiteRepository(db, environmentSQLiteRepository)
	deploymentSQLiteRepository := repositories.NewDeploymentSQLiteRepository(db)
	teamSQLiteRepository := repositories.NewTeamSQLiteRepository(db)
	workspaceSQLiteRepository := repositories.NewWorkspaceSQLiteRepository(db)
	oAuthSQLiteRepository := repositories.NewOAuthSQLiteRepository(db)
	gitSQLiteRepository := repositories.NewGitSQLiteRepository(db, v)
	prPreviewRepository := repositories.NewPRPreviewRepository(db)
	domainSQLiteRepository := repositories.NewDomainSQLiteRepository(db)
	gitAppSQLiteRepository := repositories.NewGitAppSQLiteRepository(db, v)
	teamAISettingsSQLiteRepository := repositories.NewTeamAISettingsSQLiteRepository(db, v)
	vercelRepository := repositories.NewVercelRepository(db, v)

	// Engine & Deployers
	httpEngineAdapter := newEngineAdapter(settingsSQLiteRepository, appServiceSQLiteRepository, envSQLiteRepository, databaseSQLiteRepository, storageSQLiteRepository, projectSQLiteRepository, jobSQLiteRepository, backupSQLiteRepository, s3DestinationSQLiteRepository, serviceVarSQLiteRepository, serverlessRepository)
	databaseDeployer := engine.NewDatabaseDeployer(dockerClient, httpEngineAdapter)
	storageDeployer := engine.NewStorageDeployer(dockerClient, httpEngineAdapter)

	cronManager := engine.NewCronManager(dockerClient, httpEngineAdapter)
	_ = cronManager.Start()

	backupManager := engine.NewBackupManager(dockerClient, httpEngineAdapter, "")
	_ = backupManager.Start()

	// Services
	projectService := services.NewProjectService(projectSQLiteRepository, environmentSQLiteRepository, appServiceSQLiteRepository, serviceVarSQLiteRepository)
	appService := services.NewAppService(appServiceSQLiteRepository, serviceVarSQLiteRepository)
	databaseService := services.NewDatabaseService(databaseSQLiteRepository, databaseDeployer)
	tokenService := services.NewTokenService()
	settingsService := services.NewSettingsService(settingsSQLiteRepository, notificationSQLiteRepository)
	projectSettingsService := services.NewProjectSettingsService(projectSettingsSQLiteRepository, userSQLiteRepository)
	serviceLinker := services.NewServiceLinker(databaseSQLiteRepository, storageSQLiteRepository)
	emailSettingsService := services.NewEmailSettingsService(teamEmailSettingsSQLiteRepository)
	mailerService := notifications.NewMailerService(emailSettingsService)
	dispatcherService := core.NewDispatcherService(notificationSQLiteRepository, settingsSQLiteRepository, mailerService)
	storageService := services.NewStorageService(storageSQLiteRepository, storageDeployer)
	jobService := services.NewJobService(jobSQLiteRepository, cronManager)
	canvasService := services.NewCanvasService(canvasSQLiteRepository)
	deploymentService := services.NewDeploymentService(deploymentSQLiteRepository, appServiceSQLiteRepository, projectSQLiteRepository, deployer)
	backupService := services.NewBackupService(backupSQLiteRepository, s3DestinationSQLiteRepository, backupManager)
	teamService := services.NewTeamService(teamSQLiteRepository, userSQLiteRepository)
	workspaceService := services.NewWorkspaceService(workspaceSQLiteRepository)
	userService := services.NewUserService(userSQLiteRepository)
	authService := services.NewAuthService(userSQLiteRepository, settingsSQLiteRepository, tokenService)
	oAuthService := services.NewOAuthService(oAuthSQLiteRepository, userSQLiteRepository, tokenService)
	gitService := services.NewGitService(gitSQLiteRepository)
	prPreviewService := services.NewPRPreviewService(prPreviewRepository, appService, gitService, deployer)
	environmentService := services.NewEnvironmentService(environmentSQLiteRepository, domainSQLiteRepository, envSQLiteRepository)
	notificationService := services.NewNotificationService(notificationSQLiteRepository, dispatcherService)
	gitAppsService := services.NewGitAppsService(gitAppSQLiteRepository)
	aiSettingsService := services.NewAISettingsService(teamAISettingsSQLiteRepository)
	vercelService := services.NewVercelService(vercelRepository)
	serverlessService := services.NewServerlessService(serverlessRepository)

	updaterService := services.NewUpdaterService(settingsSQLiteRepository)
	updaterService.Start(context.Background())

	// MCP Bridge
	bridge := mcp.NewBridge(projectService, appService, databaseService)

	// Auth Guard
	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService)

	// Handlers
	appHandler := handlers.NewAppHandler(appService)
	databaseHandler := handlers.NewDatabaseHandler(databaseService)
	storageHandler := handlers.NewStorageHandler(storageService)
	jobHandler := handlers.NewJobHandler(jobService)
	canvasHandler := handlers.NewCanvasHandler(canvasService)
	terminalHandler := handlers.NewTerminalHandler(dockerClient, tokenService, appService)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, appService)
	serviceVarHandler := handlers.NewServiceVarHandler(appService)
	projectSettingsHandler := handlers.NewProjectSettingsHandler(projectSettingsService)
	backupHandler := handlers.NewBackupHandler(backupService)
	teamHandler := handlers.NewTeamHandler(teamService)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	updaterHandler := handlers.NewUpdaterHandler(updaterService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	oAuthHandler := handlers.NewOAuthHandler(oAuthService)
	gitHandler := handlers.NewGitHandler(gitService)
	webhookHandler := handlers.NewWebhookHandler(gitService, projectService, appService, deploymentService, prPreviewService)
	projectHandler := handlers.NewProjectHandler(projectService)
	environmentHandler := handlers.NewEnvironmentHandler(environmentService)
	domainHandler := handlers.NewDomainHandler(environmentService)
	projectEnvHandler := handlers.NewProjectEnvHandler(environmentService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	gitAppsHandler := handlers.NewGitAppsHandler(gitAppsService)
	aiSettingsHandler := handlers.NewAISettingsHandler(aiSettingsService)
	emailSettingsHandler := handlers.NewEmailSettingsHandler(emailSettingsService)
	aiDiagnosticsHandler := handlers.NewAIDiagnosticsHandler(aiSettingsService, deploymentService, projectService)
	vercelHandler := handlers.NewVercelHandler(vercelService)
	serverlessHandler := handlers.NewServerlessHandler(serverlessService)

	srv := &Server{
		router:                 e,
		mcpBridge:              bridge,
		deployer:               deployer,
		traefikManager:         traefikManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              authGuard,
		cronManager:            cronManager,
		serviceLinker:          serviceLinker,
		dispatcherService:      dispatcherService,
		appServiceHandler:      appHandler,
		dbHandler:              databaseHandler,
		storageHandler:         storageHandler,
		jobHandler:             jobHandler,
		canvasHandler:          canvasHandler,
		terminalHandler:        terminalHandler,
		deploymentHandler:      deploymentHandler,
		serviceVarHandler:      serviceVarHandler,
		projectSettingsHandler: projectSettingsHandler,
		backupHandler:          backupHandler,
		teamHandler:            teamHandler,
		workspaceHandler:       workspaceHandler,
		settingsHandler:        settingsHandler,
		updaterHandler:         updaterHandler,
		userHandler:            userHandler,
		authHandler:            authHandler,
		oauthHandler:           oAuthHandler,
		gitHandler:             gitHandler,
		webhookHandler:         webhookHandler,
		projectHandler:         projectHandler,
		environmentHandler:     environmentHandler,
		domainHandler:          domainHandler,
		projectEnvHandler:      projectEnvHandler,
		notificationHandler:    notificationHandler,
		gitAppsHandler:         gitAppsHandler,
		aiSettingsHandler:      aiSettingsHandler,
		emailSettingsHandler:   emailSettingsHandler,
		aiDiagnosticsHandler:   aiDiagnosticsHandler,
		vercelHandler:          vercelHandler,
		serverlessHandler:      serverlessHandler,
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

func (s *Server) StartMCPStdio() error {
	mcpServer := s.mcpBridge.MCPServer()
	return server.ServeStdio(mcpServer)
}

func (s *Server) HandleMCPSSE(c echo.Context) error {
	mcpServer := s.mcpBridge.MCPServer()
	sseServer := server.NewSSEServer(mcpServer)
	sseServer.SSEHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) HandleMCPMessage(c echo.Context) error {
	mcpServer := s.mcpBridge.MCPServer()
	sseServer := server.NewSSEServer(mcpServer)
	sseServer.MessageHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
