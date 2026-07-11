package server

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"vessel.dev/vessel/internal/cloud/handlers"
	vesselMiddleware "vessel.dev/vessel/internal/cloud/middleware"
	"vessel.dev/vessel/internal/cloud/services"
)

type Server struct {
	router          *echo.Echo
	agentHandler    *handlers.AgentHandler
	wizardHandler   *handlers.WizardHandler
	billingHandler  *handlers.BillingHandler
	authHandler     *handlers.AuthHandler
	userHandler     *handlers.UserHandler
	adminHandler    *handlers.AdminHandler
	meteringHandler *handlers.MeteringHandler
}

func NewServer() *Server {
	e := echo.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	s := &Server{
		router:          e,
		agentHandler:    handlers.NewAgentHandler(),
		wizardHandler:   handlers.NewWizardHandler(),
		billingHandler:  handlers.NewBillingHandler(),
		authHandler:     handlers.NewAuthHandler(),
		userHandler:     handlers.NewUserHandler(),
		adminHandler:    handlers.NewAdminHandler(),
		meteringHandler: handlers.NewMeteringHandler(services.NewMeteringService()),
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api/cloud")

	// Global middleware
	api.Use(echoMiddleware.Logger())
	api.Use(echoMiddleware.Recover())

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "service": "vessel-cloud"})
	})

	api.GET("/agent/connect", s.agentHandler.AcceptConnection)

	// Agent & Wizard routes
	api.POST("/wizard/token", s.wizardHandler.GenerateAgentToken, vesselMiddleware.SeatLimitGuard())

	api.POST("/billing/stripe/webhook", s.billingHandler.HandleStripeWebhook)
	api.POST("/billing/stripe/checkout", s.billingHandler.CreateStripeCheckout)

	api.POST("/billing/paddle/webhook", s.billingHandler.HandlePaddleWebhook)
	api.POST("/billing/paddle/checkout", s.billingHandler.CreatePaddleCheckout)

	api.POST("/billing/usage/report", s.meteringHandler.ReportUsage)

	api.POST("/auth/register", s.authHandler.Register)
	api.POST("/auth/login", s.authHandler.Login)

	api.GET("/users/me", s.userHandler.GetProfile)

	api.GET("/admin/stats", s.adminHandler.GetSystemStats)
	api.GET("/admin/audit-logs", s.adminHandler.GetAuditLogs)

	api.POST("/fleet/deploy", s.agentHandler.DeployToFleet, vesselMiddleware.DeploymentRateLimiter())
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}
