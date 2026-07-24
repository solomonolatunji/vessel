package http

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	"codedock.run/codedock/internal/core"
	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/handlers"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
)

type Server struct {
	router                 *echo.Echo
	mcpBridge              *Bridge
	authRateLimiter        *middleware.RateLimiter
	otpRateLimiter         *middleware.RateLimiter
	deployer               *engine.Deployer
	traefikManager         *engine.TraefikManager
	dockerClient           *client.Client
	tokenService           *services.TokenService
	authGuard              *middleware.AuthGuard
	cronManager            *engine.CronManager
	serviceLinker          *services.ServiceLinker
	dispatcherService      *core.DispatcherService
	projectService         *services.ProjectService
	appService             *services.AppService
	appServiceHandler      *handlers.AppHandler
	dbHandler              *handlers.DatabaseHandler
	scheduledTaskHandler   *handlers.ScheduledTaskHandler
	canvasHandler          *handlers.CanvasHandler
	terminalHandler        *handlers.TerminalHandler
	deploymentHandler      *handlers.DeploymentHandler
	serviceVarHandler      *handlers.ServiceVarHandler
	projectSettingsHandler *handlers.ProjectSettingsHandler
	backupHandler          *handlers.BackupHandler
	settingsHandler        *handlers.SettingsHandler
	notifSettingsHandler   *handlers.NotificationSettingsHandler
	aiSettingsHandler      *handlers.AISettingsHandler
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
	serverlessHandler      *handlers.ServerlessHandler
	systemHandler          *handlers.SystemHandler
	composeHandler         *handlers.ComposeHandler
	oneClickHandler        *handlers.OneClickHandler
	archiveHandler         *handlers.ArchiveHandler
	migrationHandler       *handlers.MigrationHandler
	onboardingHandler      *handlers.OnboardingHandler
	dnsHandler             *handlers.DNSHandler
	metricsHandler         *handlers.MetricsHandler
	logHandler             *handlers.LogHandler
	auditLogHandler        *handlers.AuditLogHandler
	exampleHandler         *handlers.ExampleHandler
	serverHandler          *handlers.ServerHandler
	workerWSHandler        *handlers.WorkerWSHandler
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
