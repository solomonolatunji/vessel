package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"vessel.dev/vessel/internal/cloud/handlers"
	vesselMiddleware "vessel.dev/vessel/internal/cloud/middleware"
	"vessel.dev/vessel/internal/cloud/repos"
	"vessel.dev/vessel/internal/cloud/services"
)

type Server struct {
	router           *echo.Echo
	db               *gorm.DB
	repo             repos.CloudRepo
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

func NewServer(db *gorm.DB) *Server {
	e := echo.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	repo := repos.NewCloudRepo(db)

	ssoHandler, err := newSSOHandler()
	if err != nil {
		log.Printf("SSO disabled: %v", err)
	}

	s := &Server{
		router:           e,
		db:               db,
		repo:             repo,
		agentHandler:     handlers.NewAgentHandler(),
		wizardHandler:    handlers.NewWizardHandler(repo),
		billingHandler:   handlers.NewBillingHandler(),
		authHandler:      handlers.NewAuthHandler(),
		userHandler:      handlers.NewUserHandler(),
		adminHandler:     handlers.NewAdminHandler(),
		meteringHandler:  handlers.NewMeteringHandler(services.NewMeteringService(repo)),
		telemetryHandler: handlers.NewTelemetryHandler(repo),
		ssoHandler:       ssoHandler,
		teamHandler:      handlers.NewTeamHandler(repo),
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api/cloud")

	if s.ssoHandler != nil {
		s.ssoHandler.RegisterRoutes(s.router.Group("/api/cloud/sso"))
		api.GET("/sso/session", func(c echo.Context) error {
			return c.JSON(200, map[string]interface{}{
				"status":  "success",
				"session": c.Get("saml_session"),
			})
		}, s.ssoHandler.RequireSAML())
	}

	api.Use(echoMiddleware.Logger())
	api.Use(echoMiddleware.Recover())

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "service": "vessel-cloud"})
	})

	api.GET("/agent/connect", s.agentHandler.AcceptConnection)
	api.POST("/wizard/token", s.wizardHandler.GenerateAgentToken, vesselMiddleware.SeatLimitGuard(s.repo))

	api.POST("/billing/stripe/webhook", s.billingHandler.HandleStripeWebhook)
	api.POST("/billing/stripe/checkout", s.billingHandler.CreateStripeCheckout)
	api.POST("/billing/paddle/webhook", s.billingHandler.HandlePaddleWebhook)
	api.POST("/billing/paddle/checkout", s.billingHandler.CreatePaddleCheckout)

	api.POST("/billing/usage/report", s.meteringHandler.ReportUsage)

	api.POST("/auth/register", s.authHandler.Register)
	api.POST("/auth/login", s.authHandler.Login)

	api.GET("/users/me", s.userHandler.GetProfile)

	api.PATCH("/teams/:id/branding", s.teamHandler.UpdateBranding)

	api.GET("/admin/stats", s.adminHandler.GetSystemStats)
	api.GET("/admin/audit-logs", s.adminHandler.GetAuditLogs)
	api.POST("/admin/licenses", s.adminHandler.GenerateOfflineLicense)

	api.POST("/fleet/deploy", s.agentHandler.DeployToFleet, vesselMiddleware.DeploymentRateLimiter(s.repo))
	api.POST("/telemetry/ping", s.telemetryHandler.ReceivePing)
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}

func newSSOHandler() (*handlers.SSOHandler, error) {
	metadataURL := os.Getenv("SAML_IDP_METADATA_URL")
	if metadataURL == "" {
		return nil, fmt.Errorf("SAML_IDP_METADATA_URL not set")
	}

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Vessel Cloud"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	cert, _ := x509.ParseCertificate(certBytes)

	baseURL := os.Getenv("VESSEL_CLOUD_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return handlers.NewSSOHandler(baseURL, metadataURL, key, cert)
}
