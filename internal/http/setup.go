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
	e.Use(echomiddleware.GzipWithConfig(echomiddleware.GzipConfig{
		Level: 5,
	}))

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

	environmentRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, environmentRepo)
	appRepo := repositories.NewAppServiceRepo(db)
	serviceVarRepo := repositories.NewServiceVarRepo(db)
	dbRepo := repositories.NewDatabaseRepo(db, v)
	settingsRepo := repositories.NewSettingsRepo(db)
	envVarRepo := repositories.NewEnvRepo(db, v)
	storageRepo := repositories.NewStorageRepo(db, v)
	jobRepo := repositories.NewJobRepo(db)
	backupRepo := repositories.NewBackupRepo(db)
	s3DestinationRepo := repositories.NewS3DestinationRepo(db)
	serverlessRepository := repositories.NewServerlessRepository(db)
	projectSettingsRepo := repositories.NewProjectSettingsRepo(db)
	userRepo := repositories.NewUserRepo(db)
	canvasRepo := repositories.NewCanvasRepo(db, environmentRepo)
	deployRepo := repositories.NewDeploymentRepo(db)
	oauthRepo := repositories.NewOAuthRepo(db)
	gitRepo := repositories.NewGitRepo(db, v)
	prPreviewRepository := repositories.NewPRPreviewRepository(db)
	domainRepo := repositories.NewDomainRepo(db)
	gitAppRepo := repositories.NewGitAppRepo(db, v)
	dnsRepo := repositories.NewDNSRepo(db)
	auditRepository := repositories.NewAuditLogRepo(db)
	vercelRepository := repositories.NewVercelRepository(db, v)

	httpEngineAdapter := newEngineAdapter(settingsRepo, appRepo, envVarRepo, dbRepo, storageRepo, projectRepo, jobRepo, backupRepo, s3DestinationRepo, serviceVarRepo, serverlessRepository)
	databaseDeployer := engine.NewDatabaseDeployer(dockerClient, httpEngineAdapter)
	storageDeployer := engine.NewStorageDeployer(dockerClient, httpEngineAdapter)

	cronManager := engine.NewCronManager(dockerClient, httpEngineAdapter)

	settings, _ := settingsRepo.GetServerSettings(context.Background())
	if settings != nil && settings.DockerCleanupCron != "" {
		_ = cronManager.ScheduleDockerCleanup(settings.DockerCleanupCron)
	}
	if settings != nil && settings.DiskUsageCron != "" {
		_ = cronManager.ScheduleDiskUsageCheck(settings.DiskUsageCron, settings.DiskUsageThreshold)
	}

	_ = cronManager.Start()

	backupManager := engine.NewBackupManager(dockerClient, httpEngineAdapter, "")
	_ = backupManager.Start()

	projectService := services.NewProjectService(projectRepo, environmentRepo, appRepo, serviceVarRepo, settingsRepo)
	appService := services.NewAppService(appRepo, serviceVarRepo)
	databaseService := services.NewDatabaseService(dbRepo, databaseDeployer)
	tokenService, err := services.NewTokenService()
	if err != nil {
		return nil, fmt.Errorf("token service: %w", err)
	}
	settingsService := services.NewSettingsService(settingsRepo)
	serviceLinker := services.NewServiceLinker(dbRepo, storageRepo)
	mailerService, err := notifications.NewMailerService(settingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	authService := services.NewAuthService(userRepo, settingsRepo, projectSettingsRepo, tokenService, mailerService)
	projectSettingsService := services.NewProjectSettingsService(projectSettingsRepo, userRepo, authService)
	dispatcherService := core.NewDispatcherService(settingsRepo, userRepo, mailerService)
	storageService := services.NewStorageService(storageRepo, storageDeployer)
	jobService := services.NewJobService(jobRepo, cronManager)
	canvasService := services.NewCanvasService(canvasRepo)
	gitService := services.NewGitService(gitRepo)
	statsMonitor := engine.NewStatsMonitor(dockerClient)
	deploymentService := services.NewDeploymentService(deployRepo, appRepo, projectRepo, deployer, gitService, statsMonitor)

	autoscaler := engine.NewAutoscalerWorker(appRepo, statsMonitor, deploymentService)
	autoscaler.Start()

	backupService := services.NewBackupService(backupRepo, s3DestinationRepo, backupManager)
	userService := services.NewUserService(userRepo)
	oAuthService := services.NewOAuthService(oauthRepo, userRepo, tokenService)
	prPreviewService := services.NewPRPreviewService(prPreviewRepository, appService, gitService, deployer)
	dnsProviderService := services.NewDNSProviderService(settingsRepo)
	environmentService := services.NewEnvironmentService(environmentRepo, domainRepo, envVarRepo, dnsProviderService)
	notificationService := services.NewNotificationService(dispatcherService)
	gitAppsService := services.NewGitAppsService(gitAppRepo)
	vercelService := services.NewVercelService(vercelRepository)
	serverlessService := services.NewServerlessService(serverlessRepository)
	dnsService := services.NewDNSService(dnsRepo, dnsProviderService)
	metricsService := services.NewMetricsService()
	logService := services.NewLogService()
	auditService := services.NewAuditService(auditRepository)

	updaterService := services.NewUpdaterService(settingsRepo)
	updaterService.Start(context.Background())

	bridge := NewBridge(projectService, appService, databaseService)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService)

	appHandler := handlers.NewAppHandler(appService, projectService, deployer, deploymentService)
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
	webhookHandler := handlers.NewWebhookHandler(gitService, projectService, appService, deploymentService, prPreviewService, gitAppsService)
	projectHandler := handlers.NewProjectHandler(projectService)
	environmentHandler := handlers.NewEnvironmentHandler(environmentService)
	domainHandler := handlers.NewDomainHandler(environmentService)
	projectEnvHandler := handlers.NewProjectEnvHandler(environmentService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	gitAppsHandler := handlers.NewGitAppsHandler(gitAppsService)
	vercelHandler := handlers.NewVercelHandler(vercelService)
	tmplMgr, _ := engine.NewTemplateManager()
	composeDeployer := engine.NewComposeDeployer(dockerClient)
	composeHandler := handlers.NewComposeHandler(composeDeployer, projectService, appService, environmentRepo, appRepo)
	oneClickService := services.NewOneClickService(tmplMgr, databaseDeployer, environmentRepo, dbRepo)
	oneClickHandler := handlers.NewOneClickHandler(oneClickService)
	archiveService := services.NewArchiveService(appService, deploymentService)
	archiveHandler := handlers.NewArchiveHandler(archiveService)
	serverlessHandler := handlers.NewServerlessHandler(serverlessService)
	systemService := services.NewSystemService()
	systemHandler := handlers.NewSystemHandler(systemService)
	migrationService := services.NewMigrationService(dbRepo, dataDir)
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
