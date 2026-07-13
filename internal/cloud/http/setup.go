package http

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/cloud/handlers"
	"vessl.dev/vessl/internal/cloud/notifications"
	repos "vessl.dev/vessl/internal/cloud/repositories"
	"vessl.dev/vessl/internal/cloud/services"
)

// MountCloudRoutes mounts all cloud-specific routes onto the main echo instance and connects to PostgreSQL.
func MountCloudRoutes(e *echo.Echo) {
	db := InitDatabase()

	dsn := os.Getenv("CLOUD_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=vessl password=vessl dbname=vesslcloud port=5432 sslmode=disable"
	}
	cloudDB, dbErr := repos.NewCloudDB(dsn)
	if dbErr != nil {
		log.Printf("Warning: failed to init CloudDB for auth repo: %v", dbErr)
	}

	cloudRepo := repos.NewCloudRepo(db)
	var authRepo repos.AuthRepo
	if cloudDB != nil {
		authRepo = repos.NewAuthRepo(cloudDB.DB())
	}

	ctx := context.Background()
	mailerService, err := notifications.NewMailerService(ctx)
	if err != nil {
		log.Fatalf("failed to init mailer service: %v", err)
	}
	authService := services.NewAuthService(authRepo, mailerService)
	meteringService := services.NewMeteringService(cloudRepo, mailerService)

	ssoHandler, err := handlers.NewSSOHandler()
	if err != nil {
		log.Printf("SSO disabled: %v", err)
	}

	s := &Server{
		router:           e,
		db:               db,
		repo:             cloudRepo,
		authRepo:         authRepo,
		authService:      authService,
		mailerSvc:        mailerService,
		agentHandler:     handlers.NewAgentHandler(cloudRepo, meteringService),
		wizardHandler:    handlers.NewWizardHandler(cloudRepo),
		billingHandler:   handlers.NewBillingHandler(cloudRepo),
		authHandler:      handlers.NewAuthHandler(authService),
		userHandler:      handlers.NewUserHandler(cloudRepo, authRepo),
		adminHandler:     handlers.NewAdminHandler(cloudRepo, authRepo),
		meteringHandler:  handlers.NewMeteringHandler(cloudRepo, meteringService),
		telemetryHandler: handlers.NewTelemetryHandler(cloudRepo),
		teamHandler:      handlers.NewTeamHandler(cloudRepo),
		ssoHandler:       ssoHandler,
	}

	s.registerRoutes()
}
