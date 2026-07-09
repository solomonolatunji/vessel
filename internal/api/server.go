package api

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/notifier"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/services/oauth"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/updater"
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
	gitService             *services.GitService
	deploymentHandler      *DeploymentHandler
	serviceVarHandler      *ServiceVarHandler
	projectSettingsHandler *ProjectSettingsHandler
	backupManager          *orchestrator.BackupManager
	backupHandler          *BackupHandler
	teamHandler            *TeamHandler
	workspaceHandler       *WorkspaceHandler
	settingsHandler        *SettingsHandler
	notifierService        *notifier.NotifierService
	notificationHandler    *NotificationHandler
	oauthService           *oauth.OAuthService
	oauthHandler           *OAuthHandler
	updaterService         *updater.UpdaterService
}

// NewServer initializes a Server wired to the database store, container orchestrator, reverse proxy, and Docker client.
func NewServer(s *store.Store, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	cronMgr := orchestrator.NewCronManager(dockerClient, s)
	_ = cronMgr.Start()

	backupMgr := orchestrator.NewBackupManager(dockerClient, s, "")
	_ = backupMgr.Start()

	tokenService := services.NewTokenService()
	notifierService := notifier.NewNotifierService(s)
	oauthService := oauth.NewOAuthService()

	updaterService := updater.NewUpdaterService(s)
	updaterService.Start(context.Background())

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
		gitService:             services.NewGitService(s),
		deploymentHandler:      NewDeploymentHandler(s),
		serviceVarHandler:      NewServiceVarHandler(s),
		projectSettingsHandler: NewProjectSettingsHandler(s),
		backupManager:          backupMgr,
		backupHandler:          NewBackupHandler(s, backupMgr),
		teamHandler:            NewTeamHandler(s),
		workspaceHandler:       NewWorkspaceHandler(s),
		settingsHandler:        NewSettingsHandler(s, dockerClient, updaterService),
		notifierService:        notifierService,
		notificationHandler:    NewNotificationHandler(s, notifierService),
		oauthService:           oauthService,
		oauthHandler:           NewOAuthHandler(s, oauthService, tokenService),
		updaterService:         updaterService,
	}
	if srv.deployer != nil {
		srv.deployer.EnvProvider = srv.serviceLinker.GetLinkedEnvironmentVariables
	}
	srv.registerRoutes()
	return srv
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
func GetUserClaimsFromContext(ctx context.Context) *types.UserClaims {
	return middleware.GetUserClaimsFromContext(ctx)
}
