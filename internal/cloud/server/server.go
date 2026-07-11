package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"vessel.dev/vessel/internal/cloud/handlers"
)

type Server struct {
	router        *echo.Echo
	agentHandler  *handlers.AgentHandler
	wizardHandler *handlers.WizardHandler
}

func NewServer() *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	s := &Server{
		router:        e,
		agentHandler:  handlers.NewAgentHandler(),
		wizardHandler: handlers.NewWizardHandler(),
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

	// TODO: Mount handlers for billing, audit, etc.
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}
