package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"vessel.dev/vessel/internal/cloud/handlers"
)

type Server struct {
	router         *echo.Echo
	agentHandler   *handlers.AgentHandler
	wizardHandler  *handlers.WizardHandler
	billingHandler *handlers.BillingHandler
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	adminHandler   *handlers.AdminHandler
}

func NewServer() *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	s := &Server{
		router:         e,
		agentHandler:   handlers.NewAgentHandler(),
		wizardHandler:  handlers.NewWizardHandler(),
		billingHandler: handlers.NewBillingHandler(),
		authHandler:    handlers.NewAuthHandler(),
		userHandler:    handlers.NewUserHandler(),
		adminHandler:   handlers.NewAdminHandler(),
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api/cloud")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "service": "vessel-cloud"})
	})

	api.GET("/agent/connect", s.agentHandler.AcceptConnection)
	api.POST("/wizard/token", s.wizardHandler.GenerateAgentToken)

	api.POST("/billing/stripe/webhook", s.billingHandler.HandleStripeWebhook)
	api.POST("/billing/paddle/webhook", s.billingHandler.HandlePaddleWebhook)

	api.POST("/auth/register", s.authHandler.Register)
	api.POST("/auth/login", s.authHandler.Login)

	api.GET("/users/me", s.userHandler.GetProfile)

	api.GET("/admin/stats", s.adminHandler.GetSystemStats)
	api.GET("/admin/audit-logs", s.adminHandler.GetAuditLogs)
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}
