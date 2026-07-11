package http

import (
	"context"
	"database/sql"

	"github.com/docker/docker/client"

	"vessl.dev/vessl/internal/core"
	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/handlers"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/vault"
)

type appRepositories struct {
	settings      repositories.SettingsRepository
	user          repositories.UserRepository
	oauth         repositories.OAuthRepository
	notification  repositories.NotificationRepository
	service       repositories.AppServiceRepository
	git           repositories.GitRepository
	env           repositories.EnvRepository
	environment   repositories.EnvironmentRepository
	domain        repositories.DomainRepository
	project       repositories.ProjectRepository
	database      repositories.DatabaseRepository
	storage       repositories.StorageRepository
	job           repositories.JobRepository
	canvas        repositories.CanvasRepository
	deployment    repositories.DeploymentRepository
	backup        repositories.BackupRepository
	team          repositories.TeamRepository
	ws            repositories.WorkspaceRepository
	ps            repositories.ProjectSettingsRepository
	svVar         repositories.ServiceVarRepository
	s3            repositories.S3DestinationRepository
	prPreview     repositories.PRPreviewRepository
	serverless    repositories.ServerlessRepository
	gitApps       repositories.GitAppRepository
	aiSettings    repositories.TeamAISettingsRepository
	emailSettings repositories.TeamEmailSettingsRepository
	vercel        *repositories.VercelRepository
}

type appServices struct {
	settings      *services.SettingsService
	user          *services.UserService
	token         *services.TokenService
	auth          *services.AuthService
	oauth         *services.OAuthService
	updater       *services.UpdaterService
	git           *services.GitService
	environment   *services.EnvironmentService
	app           *services.AppService
	project       *services.ProjectService
	db            *services.DatabaseService
	storage       *services.StorageService
	job           *services.JobService
	canvas        *services.CanvasService
	backup        *services.BackupService
	team          *services.TeamService
	ws            *services.WorkspaceService
	ps            *services.ProjectSettingsService
	notification  *services.NotificationService
	deployment    *services.DeploymentService
	prPreview     *services.PRPreviewService
	gitApps       *services.GitAppsService
	aiSettings    *services.AISettingsService
	emailSettings *services.EmailSettingsService
	vercel        *services.VercelService
	serverless    services.ServerlessService
	svcLinker     *services.ServiceLinker
	dispatcher    *core.DispatcherService
	cronMgr       *engine.CronManager
	backupMgr     *engine.BackupManager
}

func initRepositories(db *sql.DB, v *vault.Vault) *appRepositories {
	envRepo := repositories.NewEnvironmentSQLiteRepository(db)
	return &appRepositories{
		settings:      repositories.NewSettingsSQLiteRepository(db),
		user:          repositories.NewUserSQLiteRepository(db),
		oauth:         repositories.NewOAuthSQLiteRepository(db),
		notification:  repositories.NewNotificationSQLiteRepository(db),
		service:       repositories.NewAppServiceSQLiteRepository(db),
		git:           repositories.NewGitSQLiteRepository(db, v),
		env:           repositories.NewEnvSQLiteRepository(db, v),
		environment:   envRepo,
		domain:        repositories.NewDomainSQLiteRepository(db),
		project:       repositories.NewProjectSQLiteRepository(db, envRepo),
		database:      repositories.NewDatabaseSQLiteRepository(db, v),
		storage:       repositories.NewStorageSQLiteRepository(db, v),
		job:           repositories.NewJobSQLiteRepository(db),
		canvas:        repositories.NewCanvasSQLiteRepository(db, envRepo),
		deployment:    repositories.NewDeploymentSQLiteRepository(db),
		backup:        repositories.NewBackupSQLiteRepository(db),
		team:          repositories.NewTeamSQLiteRepository(db),
		ws:            repositories.NewWorkspaceSQLiteRepository(db),
		ps:            repositories.NewProjectSettingsSQLiteRepository(db),
		svVar:         repositories.NewServiceVarSQLiteRepository(db),
		s3:            repositories.NewS3DestinationSQLiteRepository(db),
		prPreview:     repositories.NewPRPreviewRepository(db),
		serverless:    repositories.NewServerlessRepository(db),
		gitApps:       repositories.NewGitAppSQLiteRepository(db, v),
		aiSettings:    repositories.NewTeamAISettingsSQLiteRepository(db, v),
		emailSettings: repositories.NewTeamEmailSettingsSQLiteRepository(db, v),
		vercel:        repositories.NewVercelRepository(db, v),
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
		settings:      services.NewSettingsService(repos.settings, repos.notification),
		user:          services.NewUserService(repos.user),
		token:         tokenService,
		auth:          services.NewAuthService(repos.user, repos.settings, tokenService),
		oauth:         services.NewOAuthService(repos.oauth, repos.user, tokenService),
		updater:       updaterService,
		git:           gitService,
		environment:   services.NewEnvironmentService(repos.environment, repos.domain, repos.env),
		app:           appService,
		project:       services.NewProjectService(repos.project, repos.environment, repos.service, repos.svVar),
		db:            services.NewDatabaseService(repos.database, dbDeployer),
		storage:       services.NewStorageService(repos.storage, storageDeployer),
		job:           services.NewJobService(repos.job, cronMgr),
		canvas:        services.NewCanvasService(repos.canvas),
		backup:        services.NewBackupService(repos.backup, repos.s3, backupMgr),
		team:          services.NewTeamService(repos.team, repos.user),
		ws:            services.NewWorkspaceService(repos.ws),
		ps:            services.NewProjectSettingsService(repos.ps, repos.user),
		notification:  services.NewNotificationService(repos.notification, dispatcherSvc),
		deployment:    services.NewDeploymentService(repos.deployment, repos.service, repos.project, deployer),
		prPreview:     services.NewPRPreviewService(repos.prPreview, appService, gitService, deployer),
		gitApps:       services.NewGitAppsService(repos.gitApps),
		aiSettings:    services.NewAISettingsService(repos.aiSettings),
		emailSettings: services.NewEmailSettingsService(repos.emailSettings),
		vercel:        services.NewVercelService(repos.vercel),
		serverless:    services.NewServerlessService(repos.serverless),
		svcLinker:     svcLinker,
		dispatcher:    dispatcherSvc,
		cronMgr:       cronMgr,
		backupMgr:     backupMgr,
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
	srv.emailSettingsHandler = handlers.NewEmailSettingsHandler(svcs.emailSettings)
	srv.aiDiagnosticsHandler = handlers.NewAIDiagnosticsHandler(svcs.aiSettings, svcs.deployment, svcs.project)
	srv.vercelHandler = handlers.NewVercelHandler(svcs.vercel)
	srv.serverlessHandler = handlers.NewServerlessHandler(svcs.serverless)
}
