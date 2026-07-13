package http

import (
	"log"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"vessl.dev/vessl/internal/cloud/handlers"
	"vessl.dev/vessl/internal/cloud/notifications"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/cloud/services"
	vesslMiddleware "vessl.dev/vessl/internal/http/middleware"
)

// Server is the Vessl Cloud API server.
type Server struct {
	router           *echo.Echo
	db               *gorm.DB
	repo             repos.CloudRepo
	authRepo         repos.AuthRepo
	authService      *services.AuthService
	mailerSvc        *notifications.MailerService
	agentHandler     *handlers.AgentHandler
	wizardHandler    *handlers.WizardHandler
	billingHandler   *handlers.BillingHandler
	authHandler      *handlers.AuthHandler
	userHandler      *handlers.UserHandler
	adminHandler     *handlers.AdminHandler
	meteringHandler  *handlers.MeteringHandler
	telemetryHandler *handlers.TelemetryHandler
	teamHandler      *handlers.TeamHandler
	ssoHandler       *handlers.SSOHandler
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api")

	if s.ssoHandler != nil {
		s.ssoHandler.RegisterRoutes(s.router.Group("/api/sso"))
		api.GET("/sso/session", func(c echo.Context) error {
			return c.JSON(200, map[string]interface{}{
				"status":  "success",
				"session": c.Get("saml_session"),
			})
		}, s.ssoHandler.RequireSAML())
	}

	api.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			log.Printf("REQUEST: %s %s | status: %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	api.Use(echoMiddleware.Recover())

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "service": "vessl-cloud"})
	})

	api.GET("/system/public", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"data": map[string]interface{}{
				"registrationEnabled": true,
				"siteName":            "Vessl Cloud",
				"emailEnabled":        true,
				"isCloudMode":         true,
			},
		})
	})

	api.GET("/agent/connect", s.agentHandler.AcceptConnection)
	api.Any("/servers/:serverId/proxy/*", s.agentHandler.ProxyToServer, vesslMiddleware.RequireCloudAuth())
	api.POST("/wizard/token", s.wizardHandler.GenerateAgentToken, vesslMiddleware.SeatLimitGuard(s.repo))

	api.POST("/billing/stripe/webhook", s.billingHandler.HandleStripeWebhook)
	api.POST("/billing/stripe/checkout", s.billingHandler.CreateStripeCheckout)
	api.POST("/billing/stripe/portal", s.billingHandler.CreateStripePortal)
	api.POST("/billing/paddle/webhook", s.billingHandler.HandlePaddleWebhook)
	api.POST("/billing/paddle/checkout", s.billingHandler.CreatePaddleCheckout)

	api.POST("/billing/usage/report", s.meteringHandler.ReportUsage)

	api.POST("/auth/signup", s.authHandler.Register)
	api.POST("/auth/signin", s.authHandler.Login)
	api.POST("/auth/forgot-password", s.authHandler.ForgotPassword)
	api.POST("/auth/reset-password", s.authHandler.ResetPassword)
	api.GET("/auth/verify-email", s.authHandler.VerifyEmail)

	api.GET("/users/me", s.userHandler.GetProfile)

	api.GET("/teams/:id/servers", s.teamHandler.ListServers, vesslMiddleware.RequireCloudAuth(), vesslMiddleware.RequireTeamRole(s.repo, "owner", "admin", "member"))
	api.PATCH("/teams/:id/branding", s.teamHandler.UpdateBranding, vesslMiddleware.RequireCloudAuth(), vesslMiddleware.RequireTeamRole(s.repo, "owner", "admin"))

	api.GET("/admin/stats", s.adminHandler.GetSystemStats, vesslMiddleware.RequireAdmin())
	api.GET("/admin/audit-logs", s.adminHandler.GetAuditLogs, vesslMiddleware.RequireAdmin())
	api.POST("/admin/licenses", s.adminHandler.GenerateOfflineLicense, vesslMiddleware.RequireAdmin())

	api.POST("/fleet/deploy", s.agentHandler.DeployToFleet, vesslMiddleware.DeploymentRateLimiter(s.repo))
	api.POST("/telemetry/ping", s.telemetryHandler.ReceivePing)
}

// Start starts the HTTP server on the given address.
func (s *Server) Start(address string) error {
	return s.router.Start(address)
}
