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
	emailSettingsHandler   *handlers.EmailSettingsHandler
	aiDiagnosticsHandler   *handlers.AIDiagnosticsHandler
	vercelHandler          *handlers.VercelHandler
	serverlessHandler      *handlers.ServerlessHandler
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
