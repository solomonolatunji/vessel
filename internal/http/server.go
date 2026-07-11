package http

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mark3labs/mcp-go/server"

	"vessel.dev/vessel/internal/core"
	"vessel.dev/vessel/internal/engine"
	"vessel.dev/vessel/internal/handlers"
	"vessel.dev/vessel/internal/http/middleware"
	"vessel.dev/vessel/internal/mcp"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/repositories"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/vault"
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
	aiDiagnosticsHandler   *handlers.AIDiagnosticsHandler
	vercelHandler          *handlers.VercelHandler
	serverlessHandler      *handlers.ServerlessHandler
}

type appRepositories struct {
	settings     repositories.SettingsRepository
	user         repositories.UserRepository
	oauth        repositories.OAuthRepository
	notification repositories.NotificationRepository
	service      repositories.AppServiceRepository
	git          repositories.GitRepository
	env          repositories.EnvRepository
	environment  repositories.EnvironmentRepository
	domain       repositories.DomainRepository
	project      repositories.ProjectRepository
	database     repositories.DatabaseRepository
	storage      repositories.StorageRepository
	job          repositories.JobRepository
	canvas       repositories.CanvasRepository
	deployment   repositories.DeploymentRepository
	backup       repositories.BackupRepository
	team         repositories.TeamRepository
	ws           repositories.WorkspaceRepository
	ps           repositories.ProjectSettingsRepository
	svVar        repositories.ServiceVarRepository
	s3           repositories.S3DestinationRepository
	prPreview    repositories.PRPreviewRepository
	serverless   repositories.ServerlessRepository
	gitApps      repositories.GitAppRepository
	aiSettings   repositories.TeamAISettingsRepository
	vercel       *repositories.VercelRepository
}

type appServices struct {
	settings     *services.SettingsService
	user         *services.UserService
	token        *services.TokenService
	auth         *services.AuthService
	oauth        *services.OAuthService
	updater      *services.UpdaterService
	git          *services.GitService
	environment  *services.EnvironmentService
	app          *services.AppService
	project      *services.ProjectService
	db           *services.DatabaseService
	storage      *services.StorageService
	job          *services.JobService
	canvas       *services.CanvasService
	backup       *services.BackupService
	team         *services.TeamService
	ws           *services.WorkspaceService
	ps           *services.ProjectSettingsService
	notification *services.NotificationService
	deployment   *services.DeploymentService
	prPreview    *services.PRPreviewService
	gitApps      *services.GitAppsService
	aiSettings   *services.AISettingsService
	vercel       *services.VercelService
	serverless   services.ServerlessService
	svcLinker    *services.ServiceLinker
	dispatcher   *core.DispatcherService
	cronMgr      *engine.CronManager
	backupMgr    *engine.BackupManager
}

func initRepositories(db *sql.DB, v *vault.Vault) *appRepositories {
	envRepo := repositories.NewEnvironmentSQLiteRepository(db)
	return &appRepositories{
		settings:     repositories.NewSettingsSQLiteRepository(db),
		user:         repositories.NewUserSQLiteRepository(db),
		oauth:        repositories.NewOAuthSQLiteRepository(db),
		notification: repositories.NewNotificationSQLiteRepository(db),
		service:      repositories.NewAppServiceSQLiteRepository(db),
		git:          repositories.NewGitSQLiteRepository(db, v),
		env:          repositories.NewEnvSQLiteRepository(db, v),
		environment:  envRepo,
		domain:       repositories.NewDomainSQLiteRepository(db),
		project:      repositories.NewProjectSQLiteRepository(db, envRepo),
		database:     repositories.NewDatabaseSQLiteRepository(db, v),
		storage:      repositories.NewStorageSQLiteRepository(db, v),
		job:          repositories.NewJobSQLiteRepository(db),
		canvas:       repositories.NewCanvasSQLiteRepository(db, envRepo),
		deployment:   repositories.NewDeploymentSQLiteRepository(db),
		backup:       repositories.NewBackupSQLiteRepository(db),
		team:         repositories.NewTeamSQLiteRepository(db),
		ws:           repositories.NewWorkspaceSQLiteRepository(db),
		ps:           repositories.NewProjectSettingsSQLiteRepository(db),
		svVar:        repositories.NewServiceVarSQLiteRepository(db),
		s3:           repositories.NewS3DestinationSQLiteRepository(db),
		prPreview:    repositories.NewPRPreviewRepository(db),
		serverless:   repositories.NewServerlessRepository(db),
		gitApps:      repositories.NewGitAppSQLiteRepository(db, v),
		aiSettings:   repositories.NewTeamAISettingsSQLiteRepository(db, v),
		vercel:       repositories.NewVercelRepository(db, v),
	}
}

func initServices(repos *appRepositories, dockerClient *client.Client, deployer *engine.Deployer) *appServices {
	ea := newEngineAdapter(repos.settings, repos.service, repos.env, repos.database, repos.storage, repos.project, repos.job, repos.backup, repos.s3, repos.svVar, repos.serverless)
	cronMgr := engine.NewCronManager(dockerClient, ea)
	_ = cronMgr.Start()
	backupMgr := engine.NewBackupManager(dockerClient, ea, "")
	_ = backupMgr.Start()

	dbDeployer := engine.NewDatabaseDeployer(dockerClient, ea)
	storageDeployer := engine.NewStorageDeployer(dockerClient, ea)
	svcLinker := services.NewServiceLinker(repos.database, repos.storage)

	dispatcherSvc := core.NewDispatcherService(repos.notification, repos.settings)
	deploymentListeners := core.NewDeploymentListeners(dispatcherSvc)
	deploymentListeners.Register()

	tokenService := services.NewTokenService()
	updaterService := services.NewUpdaterService(repos.settings)
	updaterService.Start(context.Background())
	appService := services.NewAppService(repos.service, repos.svVar)
	gitService := services.NewGitService(repos.git)

	return &appServices{
		settings:     services.NewSettingsService(repos.settings, repos.notification),
		user:         services.NewUserService(repos.user),
		token:        tokenService,
		auth:         services.NewAuthService(repos.user, repos.settings, tokenService),
		oauth:        services.NewOAuthService(repos.oauth, repos.user, tokenService),
		updater:      updaterService,
		git:          gitService,
		environment:  services.NewEnvironmentService(repos.environment, repos.domain, repos.env),
		app:          appService,
		project:      services.NewProjectService(repos.project, repos.environment, repos.service, repos.svVar),
		db:           services.NewDatabaseService(repos.database, dbDeployer),
		storage:      services.NewStorageService(repos.storage, storageDeployer),
		job:          services.NewJobService(repos.job, cronMgr),
		canvas:       services.NewCanvasService(repos.canvas),
		backup:       services.NewBackupService(repos.backup, repos.s3, backupMgr),
		team:         services.NewTeamService(repos.team, repos.user),
		ws:           services.NewWorkspaceService(repos.ws),
		ps:           services.NewProjectSettingsService(repos.ps, repos.user),
		notification: services.NewNotificationService(repos.notification, dispatcherSvc),
		deployment:   services.NewDeploymentService(repos.deployment, repos.service, repos.project, deployer),
		prPreview:    services.NewPRPreviewService(repos.prPreview, appService, gitService, deployer),
		gitApps:      services.NewGitAppsService(repos.gitApps),
		aiSettings:   services.NewAISettingsService(repos.aiSettings),
		vercel:       services.NewVercelService(repos.vercel),
		serverless:   services.NewServerlessService(repos.serverless),
		svcLinker:    svcLinker,
		dispatcher:   dispatcherSvc,
		cronMgr:      cronMgr,
		backupMgr:    backupMgr,
	}
}

func initHandlers(srv *Server, svcs *appServices, dockerClient *client.Client) {
	srv.appServiceHandler = handlers.NewAppHandler(svcs.app)
	srv.dbHandler = handlers.NewDatabaseHandler(svcs.db)
	srv.storageHandler = handlers.NewStorageHandler(svcs.storage)
	srv.jobHandler = handlers.NewJobHandler(svcs.job)
	srv.canvasHandler = handlers.NewCanvasHandler(svcs.canvas)
	srv.terminalHandler = handlers.NewTerminalHandler(dockerClient, svcs.token, svcs.app)
	srv.deploymentHandler = handlers.NewDeploymentHandler(svcs.deployment, svcs.app)
	srv.serviceVarHandler = handlers.NewServiceVarHandler(svcs.app)
	srv.projectSettingsHandler = handlers.NewProjectSettingsHandler(svcs.ps)
	srv.backupHandler = handlers.NewBackupHandler(svcs.backup)
	srv.teamHandler = handlers.NewTeamHandler(svcs.team)
	srv.workspaceHandler = handlers.NewWorkspaceHandler(svcs.ws)
	srv.settingsHandler = handlers.NewSettingsHandler(svcs.settings)
	srv.updaterHandler = handlers.NewUpdaterHandler(svcs.updater)
	srv.userHandler = handlers.NewUserHandler(svcs.user)
	srv.authHandler = handlers.NewAuthHandler(svcs.auth)
	srv.oauthHandler = handlers.NewOAuthHandler(svcs.oauth)
	srv.gitHandler = handlers.NewGitHandler(svcs.git)
	srv.webhookHandler = handlers.NewWebhookHandler(svcs.git, svcs.project, svcs.app, svcs.deployment, svcs.prPreview)
	srv.projectHandler = handlers.NewProjectHandler(svcs.project)
	srv.environmentHandler = handlers.NewEnvironmentHandler(svcs.environment)
	srv.domainHandler = handlers.NewDomainHandler(svcs.environment)
	srv.projectEnvHandler = handlers.NewProjectEnvHandler(svcs.environment)
	srv.notificationHandler = handlers.NewNotificationHandler(svcs.notification)
	srv.gitAppsHandler = handlers.NewGitAppsHandler(svcs.gitApps)
	srv.aiSettingsHandler = handlers.NewAISettingsHandler(svcs.aiSettings)
	srv.aiDiagnosticsHandler = handlers.NewAIDiagnosticsHandler(svcs.aiSettings, svcs.deployment, svcs.project)
	srv.vercelHandler = handlers.NewVercelHandler(svcs.vercel)
	srv.serverlessHandler = handlers.NewServerlessHandler(svcs.serverless)
}

func NewServer(db *sql.DB, vault *vault.Vault, deployer *engine.Deployer, traefikManager *proxy.TraefikManager, dockerClient *client.Client) *Server {
	repos := initRepositories(db, vault)
	svcs := initServices(repos, dockerClient, deployer)

	mcpBridge := mcp.NewBridge(svcs.project, svcs.app, svcs.db)
	authGuard := middleware.NewAuthGuard(svcs.token, svcs.settings, svcs.ps)

	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	srv := &Server{
		router:            e,
		mcpBridge:         mcpBridge,
		deployer:          deployer,
		traefikManager:    traefikManager,
		dockerClient:      dockerClient,
		tokenService:      svcs.token,
		authGuard:         authGuard,
		cronManager:       svcs.cronMgr,
		serviceLinker:     svcs.svcLinker,
		dispatcherService: svcs.dispatcher,
	}

	initHandlers(srv, svcs, dockerClient)

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
