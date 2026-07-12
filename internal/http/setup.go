package http

import (
	"context"
	"database/sql"
	"log"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"vessl.dev/vessl/internal/core"
	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/handlers"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/mcp"
	"vessl.dev/vessl/internal/notifications"
	"vessl.dev/vessl/internal/proxy"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/vault"
)

func NewServer(db *sql.DB, v *vault.Vault, deployer *engine.Deployer, traefikManager *proxy.TraefikManager, dockerClient *client.Client) *Server {

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
	teamEmailSettingsSQLiteRepository := repositories.NewWorkspaceEmailSettingsSQLiteRepository(db, v)
	canvasSQLiteRepository := repositories.NewCanvasSQLiteRepository(db, environmentSQLiteRepository)
	deploymentSQLiteRepository := repositories.NewDeploymentSQLiteRepository(db)
	workspaceSQLiteRepository := repositories.NewWorkspaceSQLiteRepository(db)
	oAuthSQLiteRepository := repositories.NewOAuthSQLiteRepository(db)
	gitSQLiteRepository := repositories.NewGitSQLiteRepository(db, v)
	prPreviewRepository := repositories.NewPRPreviewRepository(db)
	domainSQLiteRepository := repositories.NewDomainSQLiteRepository(db)
	gitAppSQLiteRepository := repositories.NewGitAppSQLiteRepository(db, v)
	teamAISettingsSQLiteRepository := repositories.NewWorkspaceAISettingsSQLiteRepository(db, v)
	vercelRepository := repositories.NewVercelRepository(db, v)

	httpEngineAdapter := newEngineAdapter(settingsSQLiteRepository, appServiceSQLiteRepository, envSQLiteRepository, databaseSQLiteRepository, storageSQLiteRepository, projectSQLiteRepository, jobSQLiteRepository, backupSQLiteRepository, s3DestinationSQLiteRepository, serviceVarSQLiteRepository, serverlessRepository)
	databaseDeployer := engine.NewDatabaseDeployer(dockerClient, httpEngineAdapter)
	storageDeployer := engine.NewStorageDeployer(dockerClient, httpEngineAdapter)

	cronManager := engine.NewCronManager(dockerClient, httpEngineAdapter)
	_ = cronManager.Start()

	backupManager := engine.NewBackupManager(dockerClient, httpEngineAdapter, "")
	_ = backupManager.Start()

	projectService := services.NewProjectService(projectSQLiteRepository, environmentSQLiteRepository, appServiceSQLiteRepository, serviceVarSQLiteRepository)
	appService := services.NewAppService(appServiceSQLiteRepository, serviceVarSQLiteRepository)
	databaseService := services.NewDatabaseService(databaseSQLiteRepository, databaseDeployer)
	tokenService := services.NewTokenService()
	settingsService := services.NewSettingsService(settingsSQLiteRepository, notificationSQLiteRepository)
	projectSettingsService := services.NewProjectSettingsService(projectSettingsSQLiteRepository, userSQLiteRepository)
	serviceLinker := services.NewServiceLinker(databaseSQLiteRepository, storageSQLiteRepository)
	emailSettingsService := services.NewEmailSettingsService(teamEmailSettingsSQLiteRepository)
	mailerService := notifications.NewMailerService(emailSettingsService)
	dispatcherService := core.NewDispatcherService(notificationSQLiteRepository, settingsSQLiteRepository, userSQLiteRepository, mailerService)
	storageService := services.NewStorageService(storageSQLiteRepository, storageDeployer)
	jobService := services.NewJobService(jobSQLiteRepository, cronManager)
	canvasService := services.NewCanvasService(canvasSQLiteRepository)
	deploymentService := services.NewDeploymentService(deploymentSQLiteRepository, appServiceSQLiteRepository, projectSQLiteRepository, deployer)
	backupService := services.NewBackupService(backupSQLiteRepository, s3DestinationSQLiteRepository, backupManager)
	workspaceService := services.NewWorkspaceService(workspaceSQLiteRepository, userSQLiteRepository)
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

	bridge := mcp.NewBridge(projectService, appService, databaseService)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService)

	appHandler := handlers.NewAppHandler(appService, projectService)
	databaseHandler := handlers.NewDatabaseHandler(databaseService, projectService)
	storageHandler := handlers.NewStorageHandler(storageService)
	jobHandler := handlers.NewJobHandler(jobService)
	canvasHandler := handlers.NewCanvasHandler(canvasService)
	terminalHandler := handlers.NewTerminalHandler(dockerClient, tokenService, appService)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, appService)
	serviceVarHandler := handlers.NewServiceVarHandler(appService)
	projectSettingsHandler := handlers.NewProjectSettingsHandler(projectSettingsService)
	backupHandler := handlers.NewBackupHandler(backupService)
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
