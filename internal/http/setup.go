package http

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"vessl.dev/vessl/internal/core"
	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/handlers"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/notifications"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

func NewServer(db *sql.DB, v *utils.Vault, deployer *engine.Deployer, traefikManager *engine.TraefikManager, dockerClient *client.Client, dataDir string) (*Server, error) {

	e := echo.New()
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			slog.Info("request", "method", v.Method, "uri", v.URI, "status", v.Status)
			return nil
		},
	}))
	e.Use(echomiddleware.Recover())

	allowOrigins := []string{"http://localhost:3000", "http://localhost:8080"}
	if dashboardURL := os.Getenv("VESSL_DASHBOARD_URL"); dashboardURL != "" {
		allowOrigins = append(allowOrigins, dashboardURL)
	}

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	e.Use(echomiddleware.CSRFWithConfig(echomiddleware.CSRFConfig{
		TokenLength:  32,
		TokenLookup:  "header:X-CSRF-Token",
		CookieName:   "csrf_token",
		CookieMaxAge: 86400,
	}))

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
	projectSettingsSQLiteRepository := repositories.NewProjectSettingsSQLiteRepository(db)
	userSQLiteRepository := repositories.NewUserSQLiteRepository(db)
	canvasSQLiteRepository := repositories.NewCanvasSQLiteRepository(db, environmentSQLiteRepository)
	deploymentSQLiteRepository := repositories.NewDeploymentSQLiteRepository(db)
	oAuthSQLiteRepository := repositories.NewOAuthSQLiteRepository(db)
	gitSQLiteRepository := repositories.NewGitSQLiteRepository(db, v)
	prPreviewRepository := repositories.NewPRPreviewRepository(db)
	domainSQLiteRepository := repositories.NewDomainSQLiteRepository(db)
	gitAppSQLiteRepository := repositories.NewGitAppSQLiteRepository(db, v)
	dnsSQLiteRepository := repositories.NewDNSSQLiteRepository(db)
	auditLogRepository := repositories.NewAuditLogSQLiteRepository(db)
	vercelRepository := repositories.NewVercelRepository(db, v)

	httpEngineAdapter := newEngineAdapter(settingsSQLiteRepository, appServiceSQLiteRepository, envSQLiteRepository, databaseSQLiteRepository, storageSQLiteRepository, projectSQLiteRepository, jobSQLiteRepository, backupSQLiteRepository, s3DestinationSQLiteRepository, serviceVarSQLiteRepository, serverlessRepository)
	databaseDeployer := engine.NewDatabaseDeployer(dockerClient, httpEngineAdapter)
	storageDeployer := engine.NewStorageDeployer(dockerClient, httpEngineAdapter)

	cronManager := engine.NewCronManager(dockerClient, httpEngineAdapter)

	settings, _ := settingsSQLiteRepository.GetServerSettings(context.Background())
	if settings != nil && settings.DockerCleanupCron != "" {
		_ = cronManager.ScheduleDockerCleanup(settings.DockerCleanupCron)
	}
	if settings != nil && settings.DiskUsageCron != "" {
		_ = cronManager.ScheduleDiskUsageCheck(settings.DiskUsageCron, settings.DiskUsageThreshold)
	}

	_ = cronManager.Start()

	backupManager := engine.NewBackupManager(dockerClient, httpEngineAdapter, "")
	_ = backupManager.Start()

	projectService := services.NewProjectService(projectSQLiteRepository, environmentSQLiteRepository, appServiceSQLiteRepository, serviceVarSQLiteRepository, settingsSQLiteRepository)
	appService := services.NewAppService(appServiceSQLiteRepository, serviceVarSQLiteRepository)
	databaseService := services.NewDatabaseService(databaseSQLiteRepository, databaseDeployer)
	tokenService, err := services.NewTokenService()
	if err != nil {
		return nil, fmt.Errorf("token service: %w", err)
	}
	settingsService := services.NewSettingsService(settingsSQLiteRepository)
	serviceLinker := services.NewServiceLinker(databaseSQLiteRepository, storageSQLiteRepository)
	mailerService, err := notifications.NewMailerService(settingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	authService := services.NewAuthService(userSQLiteRepository, settingsSQLiteRepository, projectSettingsSQLiteRepository, tokenService, mailerService)
	projectSettingsService := services.NewProjectSettingsService(projectSettingsSQLiteRepository, userSQLiteRepository, authService)
	dispatcherService := core.NewDispatcherService(settingsSQLiteRepository, userSQLiteRepository, mailerService)
	storageService := services.NewStorageService(storageSQLiteRepository, storageDeployer)
	jobService := services.NewJobService(jobSQLiteRepository, cronManager)
	canvasService := services.NewCanvasService(canvasSQLiteRepository)
	gitService := services.NewGitService(gitSQLiteRepository)
	statsMonitor := engine.NewStatsMonitor(dockerClient)
	deploymentService := services.NewDeploymentService(deploymentSQLiteRepository, appServiceSQLiteRepository, projectSQLiteRepository, deployer, gitService, statsMonitor)
	backupService := services.NewBackupService(backupSQLiteRepository, s3DestinationSQLiteRepository, backupManager)
	userService := services.NewUserService(userSQLiteRepository)
	oAuthService := services.NewOAuthService(oAuthSQLiteRepository, userSQLiteRepository, tokenService)
	prPreviewService := services.NewPRPreviewService(prPreviewRepository, appService, gitService, deployer)
	dnsProviderService := services.NewDNSProviderService(settingsSQLiteRepository)
	environmentService := services.NewEnvironmentService(environmentSQLiteRepository, domainSQLiteRepository, envSQLiteRepository, dnsProviderService)
	notificationService := services.NewNotificationService(dispatcherService)
	gitAppsService := services.NewGitAppsService(gitAppSQLiteRepository)
	vercelService := services.NewVercelService(vercelRepository)
	serverlessService := services.NewServerlessService(serverlessRepository)
	dnsService := services.NewDNSService(dnsSQLiteRepository, dnsProviderService)
	metricsService := services.NewMetricsService()
	logService := services.NewLogService()
	auditService := services.NewAuditService(auditLogRepository)

	updaterService := services.NewUpdaterService(settingsSQLiteRepository)
	updaterService.Start(context.Background())

	bridge := NewBridge(projectService, appService, databaseService)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService)

	appHandler := handlers.NewAppHandler(appService, projectService)
	databaseHandler := handlers.NewDatabaseHandler(databaseService, projectService)
	storageHandler := handlers.NewStorageHandler(storageService)
	jobHandler := handlers.NewJobHandler(jobService)
	canvasHandler := handlers.NewCanvasHandler(canvasService)
	terminalHandler := handlers.NewTerminalHandler(dockerClient, tokenService, appService)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, appService, auditService)
	serviceVarHandler := handlers.NewServiceVarHandler(appService, auditService)
	projectSettingsHandler := handlers.NewProjectSettingsHandler(projectSettingsService)
	backupHandler := handlers.NewBackupHandler(backupService)
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
	vercelHandler := handlers.NewVercelHandler(vercelService)
	tmplMgr, _ := engine.NewTemplateManager()
	composeDeployer := engine.NewComposeDeployer(dockerClient)
	composeHandler := handlers.NewComposeHandler(composeDeployer, projectService, appService, environmentSQLiteRepository, appServiceSQLiteRepository)
	oneClickService := services.NewOneClickService(tmplMgr, databaseDeployer, environmentSQLiteRepository, databaseSQLiteRepository)
	oneClickHandler := handlers.NewOneClickHandler(oneClickService)
	archiveService := services.NewArchiveService(appService, deploymentService)
	archiveHandler := handlers.NewArchiveHandler(archiveService)
	serverlessHandler := handlers.NewServerlessHandler(serverlessService)
	systemService := services.NewSystemService()
	systemHandler := handlers.NewSystemHandler(systemService)
	migrationService := services.NewMigrationService(databaseSQLiteRepository, dataDir)
	migrationHandler := handlers.NewMigrationHandler(migrationService)
	onboardingHandler := handlers.NewOnboardingHandler(userService, authService, settingsService)
	railwayService := services.NewRailwayService(projectService, environmentService, appService, databaseService)
	railwayHandler := handlers.NewRailwayHandler(railwayService)
	dnsHandler := handlers.NewDNSHandler(dnsService)
	metricsHandler := handlers.NewMetricsHandler(metricsService)
	logHandler := handlers.NewLogHandler(logService)
	auditLogHandler := handlers.NewAuditLogHandler(auditService)
	authLimiter := middleware.NewRateLimiter(10, time.Minute)

	srv := &Server{
		router:                 e,
		mcpBridge:              bridge,
		authRateLimiter:        authLimiter,
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
		vercelHandler:          vercelHandler,
		serverlessHandler:      serverlessHandler,
		systemHandler:          systemHandler,
		composeHandler:         composeHandler,
		oneClickHandler:        oneClickHandler,
		archiveHandler:         archiveHandler,
		migrationHandler:       migrationHandler,
		onboardingHandler:      onboardingHandler,
		railwayHandler:         railwayHandler,
		dnsHandler:             dnsHandler,
		metricsHandler:         metricsHandler,
		logHandler:             logHandler,
		auditLogHandler:        auditLogHandler,
	}

	if srv.deployer != nil {
		srv.deployer.EnvProvider = func(projectID string) (map[string]string, error) {
			return srv.serviceLinker.GetLinkedEnvironmentVariables(context.Background(), projectID)
		}
		srv.deployer.EnvInterpolator = func(projectID string) (map[string]map[string]string, error) {
			return srv.serviceLinker.GetNamespacedVariables(context.Background(), projectID)
		}
	}

	srv.registerRoutes()
	return srv, nil
}
