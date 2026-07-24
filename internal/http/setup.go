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

	"codedock.run/codedock/internal/core"
	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/handlers"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/notifications"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
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
	if dashboardURL := os.Getenv("CODEDOCK_DASHBOARD_URL"); dashboardURL != "" {
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
	notifRepo := repositories.NewNotificationSettingsRepo(db)
	aiRepo := repositories.NewAISettingsRepo(db)
	envVarRepo := repositories.NewEnvRepo(db, v)
	scheduledTaskRepo := repositories.NewScheduledTaskRepo(db)
	backupRepo := repositories.NewBackupRepo(db, v)
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
	volumeRepo := repositories.NewServiceVolumeRepo(db)

	httpEngineAdapter := newEngineAdapter(settingsRepo, appRepo, envVarRepo, dbRepo, projectRepo, scheduledTaskRepo, backupRepo, s3DestinationRepo, serviceVarRepo, serverlessRepository)
	databaseDeployer := engine.NewDatabaseDeployer(dockerClient, httpEngineAdapter)

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

	projectService := services.NewProjectService(projectRepo, environmentRepo, appRepo, serviceVarRepo, settingsRepo, projectSettingsRepo)
	appService := services.NewAppService(appRepo, serviceVarRepo, volumeRepo)
	databaseService := services.NewDatabaseService(dbRepo, databaseDeployer)
	tokenService, err := services.NewTokenService()
	if err != nil {
		return nil, fmt.Errorf("token service: %w", err)
	}
	settingsService := services.NewSettingsService(settingsRepo)
	notifSettingsService := services.NewNotificationSettingsService(notifRepo)
	aiSettingsService := services.NewAISettingsService(aiRepo)
	serviceLinker := services.NewServiceLinker(dbRepo)
	mailerService, err := notifications.NewMailerService(notifSettingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	authService := services.NewAuthService(userRepo, settingsRepo, notifRepo, projectSettingsRepo, tokenService, mailerService)
	projectSettingsService := services.NewProjectSettingsService(projectSettingsRepo, userRepo, authService)
	dispatcherService := core.NewDispatcherService(settingsRepo, notifRepo, userRepo, mailerService)

	deploymentListeners := core.NewDeploymentListeners(dispatcherService, appRepo)
	deploymentListeners.Register()

	serverRepo := repositories.NewServerRepository(db)
	workerHub := engine.NewWorkerHub(serverRepo)

	scheduledTaskService := services.NewScheduledTaskService(scheduledTaskRepo, cronManager)
	canvasService := services.NewCanvasService(canvasRepo)
	gitService := services.NewGitService(gitRepo)
	statsMonitor := engine.NewStatsMonitor(dockerClient)
	deploymentService := services.NewDeploymentService(deployRepo, appRepo, projectRepo, deployer, gitService, statsMonitor, volumeRepo, workerHub)
	aiAnalysisService := services.NewAIAnalysisService(deployRepo, appRepo, aiRepo)

	autoscaler := engine.NewAutoscalerWorker(appRepo, statsMonitor, deploymentService)
	autoscaler.Start()

	backupService := services.NewBackupService(backupRepo, s3DestinationRepo, backupManager)
	userService := services.NewUserService(userRepo)
	oAuthService := services.NewOAuthService(oauthRepo, userRepo, tokenService)
	prPreviewService := services.NewPRPreviewService(prPreviewRepository, appService, gitService, deployer, workerHub, projectRepo)
	dnsProviderService := services.NewDNSProviderService(settingsRepo)
	environmentService := services.NewEnvironmentService(environmentRepo, domainRepo, envVarRepo, dnsProviderService)
	notificationService := services.NewNotificationService(dispatcherService)
	gitAppsService := services.NewGitAppsService(gitAppRepo)
	serverlessService := services.NewServerlessService(serverlessRepository)
	dnsService := services.NewDNSService(dnsRepo, dnsProviderService)
	envSuggestionService := services.NewEnvSuggestionService(gitService)
	metricsService := services.NewMetricsService()
	logService := services.NewLogService()
	auditService := services.NewAuditService(auditRepository)

	updaterService := services.NewUpdaterService(settingsRepo)
	updaterService.Start(context.Background())

	bridge := NewBridge(projectService, appService, databaseService)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService, projectSettingsRepo)

	appHandler := handlers.NewAppHandler(appService, projectService, deployer, deploymentService, environmentService)
	databaseHandler := handlers.NewDatabaseHandler(databaseService, projectService)
	scheduledTaskHandler := handlers.NewScheduledTaskHandler(scheduledTaskService)
	canvasHandler := handlers.NewCanvasHandler(canvasService)
	terminalHandler := handlers.NewTerminalHandler(dockerClient, tokenService, appService)
	projectHandler := handlers.NewProjectHandler(projectService, projectSettingsService)
	environmentHandler := handlers.NewEnvironmentHandler(environmentService)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, appService, auditService, aiAnalysisService, prPreviewService, projectService)
	serviceVarHandler := handlers.NewServiceVarHandler(appService, auditService, envSuggestionService)
	projectSettingsHandler := handlers.NewProjectSettingsHandler(projectSettingsService)
	backupHandler := handlers.NewBackupHandler(backupService)
	settingsHandler := handlers.NewSettingsHandler(settingsService, notifSettingsService)
	notifSettingsHandler := handlers.NewNotificationSettingsHandler(notifSettingsService)
	aiSettingsHandler := handlers.NewAISettingsHandler(aiSettingsService)
	updaterHandler := handlers.NewUpdaterHandler(updaterService)
	userHandler := handlers.NewUserHandler(userService, mailerService)
	authHandler := handlers.NewAuthHandler(authService)
	oAuthHandler := handlers.NewOAuthHandler(oAuthService)
	gitHandler := handlers.NewGitHandler(gitService)
	webhookHandler := handlers.NewWebhookHandler(gitService, projectService, appService, deploymentService, prPreviewService, gitAppsService)

	domainHandler := handlers.NewDomainHandler(environmentService)
	projectEnvHandler := handlers.NewProjectEnvHandler(environmentService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	gitAppsHandler := handlers.NewGitAppsHandler(gitAppsService)
	tmplMgr, _ := engine.NewTemplateManager()
	composeParserService := services.NewComposeParserService()
	composeHandler := handlers.NewComposeHandler(projectService, appService, databaseService, environmentRepo, appRepo, composeParserService)
	oneClickService := services.NewOneClickService(tmplMgr, databaseDeployer, environmentRepo, dbRepo)
	oneClickHandler := handlers.NewOneClickHandler(oneClickService)
	archiveService := services.NewArchiveService(appService, deploymentService)
	archiveHandler := handlers.NewArchiveHandler(archiveService)
	serverlessHandler := handlers.NewServerlessHandler(serverlessService)
	systemService := services.NewSystemService()
	systemHandler := handlers.NewSystemHandler(systemService)
	migrationService := services.NewMigrationService(dbRepo, dataDir)
	migrationHandler := handlers.NewMigrationHandler(migrationService)
	onboardingService := services.NewOnboardingService(userService, authService, settingsService, gitAppsService, backupService)
	onboardingHandler := handlers.NewOnboardingHandler(userService, onboardingService)
	dnsHandler := handlers.NewDNSHandler(dnsService)
	metricsHandler := handlers.NewMetricsHandler(metricsService)
	logHandler := handlers.NewLogHandler(logService)
	auditLogHandler := handlers.NewAuditLogHandler(auditService)
	exampleService := services.NewExampleService()
	exampleHandler := handlers.NewExampleHandler(exampleService)


	workerWSHandler := handlers.NewWorkerWSHandler(workerHub, serverRepo)

	authLimiter := middleware.NewRateLimiter(10, time.Minute)
	otpLimiter := middleware.NewRateLimiter(5, time.Minute)

	srv := &Server{
		router:                 e,
		mcpBridge:              bridge,
		authRateLimiter:        authLimiter,
		otpRateLimiter:         otpLimiter,
		deployer:               deployer,
		traefikManager:         traefikManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              authGuard,
		cronManager:            cronManager,
		serviceLinker:          serviceLinker,
		dispatcherService:      dispatcherService,
		projectService:         projectService,
		appService:             appService,
		appServiceHandler:      appHandler,
		dbHandler:              databaseHandler,
		scheduledTaskHandler:   scheduledTaskHandler,
		canvasHandler:          canvasHandler,
		terminalHandler:        terminalHandler,
		deploymentHandler:      deploymentHandler,
		serviceVarHandler:      serviceVarHandler,
		projectSettingsHandler: projectSettingsHandler,
		backupHandler:          backupHandler,
		settingsHandler:        settingsHandler,
		notifSettingsHandler:   notifSettingsHandler,
		aiSettingsHandler:      aiSettingsHandler,
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
		serverlessHandler:      serverlessHandler,
		systemHandler:          systemHandler,
		composeHandler:         composeHandler,
		oneClickHandler:        oneClickHandler,
		archiveHandler:         archiveHandler,
		migrationHandler:       migrationHandler,
		onboardingHandler:      onboardingHandler,
		dnsHandler:             dnsHandler,
		metricsHandler:         metricsHandler,
		logHandler:             logHandler,
		auditLogHandler:        auditLogHandler,
		exampleHandler:         exampleHandler,
		workerWSHandler:        workerWSHandler,
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
