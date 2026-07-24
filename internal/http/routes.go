package http

import (
	"os"
	"path/filepath"

	"codedock.run/codedock/apps/dashboard"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerRoutes() {
	apiGroup := s.router.Group("/api")

	authGroup := apiGroup.Group("")
	authGroup.Use(s.authGuard.RequireAuth())

	s.registerAuthRoutes(apiGroup, authGroup)
	s.registerSystemRoutes(apiGroup, authGroup)
	s.registerUserRoutes(apiGroup, authGroup)
	s.registerProjectRoutes(apiGroup, authGroup)
	s.registerServerRoutes(authGroup)
	s.registerDatabaseRoutes(authGroup)
	s.registerAppRoutes(apiGroup, authGroup)
	s.registerDeploymentRoutes(authGroup)
	s.registerBackupRoutes(authGroup)
	s.registerSettingsRoutes(apiGroup, authGroup)
	s.registerMiscRoutes(apiGroup, authGroup)

	s.setupSPAFallback()
}

func (s *Server) RequireServiceRole(minPermission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			serviceID := c.Param("serviceId")
			if serviceID == "" {
				serviceID = c.Param("id")
			}
			if serviceID == "" {
				return next(c)
			}

			userClaims := GetUserClaimsFromContext(c.Request().Context())
			if userClaims == nil {
				return utils.Error(c, 401, "unauthorized")
			}
			if userClaims.Role == "admin" {
				return next(c)
			}

			svc, err := s.appService.GetAppService(c.Request().Context(), serviceID)
			if err != nil || svc == nil {
				return utils.Error(c, 404, "service not found")
			}

			if userClaims.Role == "api" {
				if c.Get("project_id") != svc.ProjectID {
					return utils.Error(c, 403, "api token not authorized for this project")
				}
				if minPermission != "" && minPermission != models.MemberPermissionMember {
					return utils.Error(c, 403, "api tokens cannot perform admin/owner actions")
				}
				return next(c)
			}

			if !s.projectService.HasPermission(c.Request().Context(), svc.ProjectID, userClaims.UserID, userClaims.Role, minPermission) {
				return utils.Error(c, 403, "insufficient project permissions")
			}
			return next(c)
		}
	}
}

func (s *Server) registerAuthRoutes(apiGroup, authGroup *echo.Group) {
	apiGroup.POST("/auth/signup", s.authHandler.Register, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/signin", s.authHandler.Login, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/refresh", s.authHandler.Refresh)
	apiGroup.POST("/auth/forgot-password", s.authHandler.ForgotPassword, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/reset-password", s.authHandler.ResetPassword, s.authRateLimiter.Middleware)
	apiGroup.POST("/auth/logout", s.authHandler.Logout)
	authGroup.GET("/auth/me", s.userHandler.GetProfile)

	apiGroup.GET("/auth/oauth/providers/enabled", s.oauthHandler.ListEnabledProviders)
	apiGroup.GET("/auth/oauth/:provider", s.oauthHandler.OAuthRedirect)
	apiGroup.GET("/auth/oauth/:provider/callback", s.oauthHandler.OAuthCallback)
	authGroup.POST("/auth/2fa/setup", s.oauthHandler.Setup2FA)
	authGroup.POST("/auth/2fa/verify", s.oauthHandler.Verify2FA)
	authGroup.POST("/auth/2fa/disable", s.oauthHandler.Disable2FA, s.otpRateLimiter.Middleware)
}

func (s *Server) registerSystemRoutes(apiGroup, authGroup *echo.Group) {
	apiGroup.GET("/system/public", s.settingsHandler.GetPublicSettings)
	apiGroup.GET("/system/setup-status", s.onboardingHandler.SetupStatus)
	apiGroup.POST("/system/setup", s.onboardingHandler.Setup)
	authGroup.GET("/system/stats", s.systemHandler.GetStats)
	apiGroup.POST("/system/restart", s.systemHandler.Restart, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/maintenance/cleanup", s.systemHandler.Cleanup, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/export", s.migrationHandler.Export, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/system/import", s.migrationHandler.Import, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerUserRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/users", s.userHandler.ListUsers, s.authGuard.RequireRole("admin"))
	authGroup.POST("/users/invite", s.authHandler.AdminInviteUser, s.authGuard.RequireRole("admin"))
	authGroup.DELETE("/users/:id", s.userHandler.DeleteUser, s.authGuard.RequireRole("admin"))
	authGroup.GET("/profile", s.userHandler.GetProfile)
	authGroup.PUT("/profile", s.userHandler.UpdateProfile)
	authGroup.POST("/profile/email/request", s.userHandler.RequestEmailChange)
	authGroup.POST("/profile/email/verify", s.userHandler.VerifyEmailChange, s.otpRateLimiter.Middleware)
	authGroup.PUT("/profile/password", s.userHandler.ChangePassword)
	authGroup.GET("/profile/tokens", s.userHandler.ListPATs)
	authGroup.POST("/profile/tokens", s.userHandler.CreatePAT)
	authGroup.DELETE("/profile/tokens/:id", s.userHandler.DeletePAT)
}

func (s *Server) registerProjectRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/projects", s.projectHandler.ListProjects)
	authGroup.POST("/projects", s.projectHandler.CreateProject)

	projectAuth := s.authGuard.RequireProjectRole("")
	projectAuthAdmin := s.authGuard.RequireProjectRole(models.MemberPermissionAdmin)
	projectAuthOwner := s.authGuard.RequireProjectRole(models.MemberPermissionOwner)

	authGroup.GET("/projects/:id", s.projectHandler.GetProject, projectAuth)
	authGroup.DELETE("/projects/:id", s.projectHandler.DeleteProject, projectAuthOwner)

	authGroup.GET("/services/:id/domains", s.domainHandler.ListByService)
	authGroup.POST("/services/:id/domains", s.domainHandler.Create)
	authGroup.DELETE("/domains/:id", s.domainHandler.Delete)

	authGroup.GET("/projects/:id/env", s.projectEnvHandler.GetVars, projectAuth, s.authGuard.RequireScope("env:read"))
	authGroup.PUT("/projects/:id/env", s.projectEnvHandler.SetVars, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.POST("/projects/:id/environments", s.environmentHandler.Create, projectAuthAdmin)
	authGroup.GET("/projects/:id/environments", s.environmentHandler.ListByProject, projectAuth)
	authGroup.GET("/projects/:id/apps", s.appServiceHandler.ListByProject, projectAuth)

	authGroup.GET("/projects/:projectId/tokens", s.projectSettingsHandler.ListTokens, projectAuthAdmin, s.authGuard.RequireScope("env:read"))
	authGroup.POST("/projects/:projectId/tokens", s.projectSettingsHandler.CreateToken, projectAuthAdmin, s.authGuard.RequireScope("env:write"))
	authGroup.DELETE("/projects/:projectId/tokens/:id", s.projectSettingsHandler.DeleteToken, projectAuthAdmin, s.authGuard.RequireScope("env:write"))

	authGroup.GET("/projects/:projectId/members", s.projectSettingsHandler.ListMembers, projectAuth)
	authGroup.POST("/projects/:projectId/members", s.projectSettingsHandler.AddMember, projectAuthAdmin)
	authGroup.DELETE("/projects/:projectId/members/:id", s.projectSettingsHandler.RemoveMember, projectAuthAdmin)
}

func (s *Server) registerServerRoutes(authGroup *echo.Group) {
	authGroup.GET("/servers", s.serverHandler.List, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/servers", s.serverHandler.Create, s.authGuard.RequireScope("server:write"))
}

func (s *Server) registerDatabaseRoutes(authGroup *echo.Group) {
	authGroup.GET("/databases", s.dbHandler.ListDatabases, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases", s.dbHandler.CreateDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.GET("/databases/:id", s.dbHandler.GetDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.PUT("/databases/:id", s.dbHandler.UpdateDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.DELETE("/databases/:id", s.dbHandler.DeleteDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/start", s.dbHandler.StartDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/stop", s.dbHandler.StopDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/restart", s.dbHandler.RestartDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/query", s.dbHandler.QueryDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/import", s.dbHandler.ImportData, s.authGuard.RequireScope("database:manage"))

	authGroup.GET("/databases/:id/schemas", s.dbHandler.GetSchemas, s.authGuard.RequireScope("database:manage"))
	authGroup.GET("/databases/:id/data/:table", s.dbHandler.GetTableData, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/data/:table", s.dbHandler.InsertTableRow, s.authGuard.RequireScope("database:manage"))
	authGroup.PUT("/databases/:id/data/:table", s.dbHandler.UpdateTableRow, s.authGuard.RequireScope("database:manage"))
	authGroup.DELETE("/databases/:id/data/:table", s.dbHandler.DeleteTableRow, s.authGuard.RequireScope("database:manage"))
}

func (s *Server) registerAppRoutes(apiGroup, authGroup *echo.Group) {

	serviceAuthAdmin := s.RequireServiceRole(models.MemberPermissionAdmin)
	serviceAuthOwner := s.RequireServiceRole(models.MemberPermissionOwner)
	serviceAuth := s.RequireServiceRole("")

	authGroup.GET("/environments/:id/apps", s.appServiceHandler.ListByEnvironment)
	authGroup.POST("/environments/:id/apps", s.appServiceHandler.Create)
	authGroup.DELETE("/environments/:id", s.environmentHandler.Delete)
	authGroup.GET("/apps/:id", s.appServiceHandler.Get, serviceAuth)
	authGroup.PUT("/apps/:id", s.appServiceHandler.Update, serviceAuthAdmin)
	authGroup.DELETE("/apps/:id", s.appServiceHandler.Delete, serviceAuthOwner)
	authGroup.POST("/apps/:id/stop", s.appServiceHandler.StopService, serviceAuthAdmin)
	authGroup.POST("/apps/:id/redeploy", s.appServiceHandler.RedeployService, serviceAuthAdmin)
	authGroup.POST("/apps/:id/restart", s.appServiceHandler.RestartService, serviceAuthAdmin)
	appsGroup := authGroup.Group("/apps")

	appsGroup.GET("/:id/webhooks", s.appServiceHandler.ListWebhooks, serviceAuth)
	appsGroup.POST("/:id/webhooks", s.appServiceHandler.CreateWebhook, serviceAuthAdmin)
	appsGroup.DELETE("/:id/webhooks/:webhookId", s.appServiceHandler.DeleteWebhook, serviceAuthAdmin)
	appsGroup.GET("/:id/volumes", s.appServiceHandler.ListVolumes, serviceAuth)
	appsGroup.POST("/:id/volumes", s.appServiceHandler.CreateVolume, serviceAuthAdmin)
	appsGroup.DELETE("/:id/volumes/:volumeId", s.appServiceHandler.DeleteVolume, serviceAuthAdmin)
	appsGroup.GET("/:id/log-drains", s.appServiceHandler.ListLogDrains, serviceAuth)
	appsGroup.POST("/:id/log-drains", s.appServiceHandler.CreateLogDrain, serviceAuthAdmin)
	appsGroup.DELETE("/:id/log-drains/:drainId", s.appServiceHandler.DeleteLogDrain, serviceAuthAdmin)

	authGroup.GET("/services/:serviceId/variables", s.serviceVarHandler.List, serviceAuth)
	authGroup.GET("/services/:serviceId/env-suggestions", s.serviceVarHandler.Suggest, serviceAuth)
	authGroup.POST("/services/:serviceId/variables", s.serviceVarHandler.Create, serviceAuthAdmin)
	authGroup.PUT("/services/:serviceId/variables/:id", s.serviceVarHandler.Update, serviceAuthAdmin)
	authGroup.DELETE("/services/:serviceId/variables/:id", s.serviceVarHandler.Delete, serviceAuthAdmin)

	authGroup.GET("/services/:serviceId/serverless/code", s.serverlessHandler.GetCode, serviceAuth)
	authGroup.POST("/services/:serviceId/serverless/code", s.serverlessHandler.SaveCode, serviceAuthAdmin)
}

func (s *Server) registerDeploymentRoutes(authGroup *echo.Group) {
	serviceAuthAdmin := s.RequireServiceRole(models.MemberPermissionAdmin)
	serviceAuth := s.RequireServiceRole("")

	authGroup.GET("/services/:serviceId/deployments", s.deploymentHandler.ListServiceDeployments, serviceAuth)
	authGroup.GET("/services/:serviceId/previews", s.deploymentHandler.ListPRPreviews, serviceAuth)
	authGroup.POST("/services/:serviceId/deploy", s.deploymentHandler.Trigger, serviceAuthAdmin)
	authGroup.POST("/deployments/:id/rollback", s.deploymentHandler.Rollback)
	authGroup.GET("/deployments/:id/logs", s.deploymentHandler.GetLogs, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/deployments/:id/explain", s.deploymentHandler.ExplainFailure)
	authGroup.GET("/services/:serviceId/metrics", s.deploymentHandler.GetMetrics, serviceAuth)
	authGroup.GET("/services/:serviceId/metrics/historical", s.metricsHandler.GetHistoricalMetrics, serviceAuth)
	authGroup.GET("/services/:serviceId/logs/historical", s.logHandler.GetHistoricalLogs, serviceAuth)
}

func (s *Server) registerBackupRoutes(authGroup *echo.Group) {
	authGroup.GET("/backups", s.backupHandler.List)
	authGroup.POST("/backups", s.backupHandler.Create)
	authGroup.GET("/backups/:id", s.backupHandler.Get)
	authGroup.PUT("/backups/:id", s.backupHandler.Update, s.authGuard.RequireScope("backup:write"))
	authGroup.DELETE("/backups/:id", s.backupHandler.Delete)
	authGroup.POST("/backups/:id/trigger", s.backupHandler.Trigger)
	authGroup.POST("/backups/:id/restore", s.backupHandler.Restore)
	authGroup.GET("/backups/:id/records", s.backupHandler.ListRecords)
	authGroup.GET("/backups/:id/records/:recordId/download", s.backupHandler.DownloadRecord)
	authGroup.DELETE("/backups/:id/records/:recordId", s.backupHandler.DeleteRecord)
	authGroup.GET("/s3-destinations", s.backupHandler.ListS3Destinations)
	authGroup.POST("/s3-destinations", s.backupHandler.CreateS3Destination)
	authGroup.DELETE("/s3-destinations/:id", s.backupHandler.DeleteS3Destination)
}

func (s *Server) registerSettingsRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.GET("/settings", s.settingsHandler.GetSettings)
	apiGroup.PUT("/settings", s.settingsHandler.UpdateSettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/ai", s.aiSettingsHandler.GetAISettings)
	authGroup.POST("/ai/diagnose", s.aiSettingsHandler.DiagnoseLogs)
	apiGroup.PUT("/ai", s.aiSettingsHandler.UpdateAISettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/notifications", s.notifSettingsHandler.GetNotificationSettings)
	apiGroup.PUT("/notifications", s.notifSettingsHandler.UpdateNotificationSettings, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/updates/status", s.updaterHandler.GetUpdateStatus)
	apiGroup.POST("/settings/updates/check", s.updaterHandler.CheckUpdate, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/updates/deploy", s.updaterHandler.DeployUpdate, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/oauth/providers", s.oauthHandler.ListProviders)
	apiGroup.PUT("/settings/oauth/providers", s.oauthHandler.SaveProvider, s.authGuard.RequireRole("admin"))

	apiGroup.POST("/settings/git_apps/github/manifest-callback", s.gitAppsHandler.ExchangeGithubManifestCode, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/git_apps/github", s.gitAppsHandler.ListGithubApps)
	authGroup.GET("/settings/git_apps/github/:id", s.gitAppsHandler.GetGithubApp)
	apiGroup.PUT("/settings/git_apps/github", s.gitAppsHandler.SaveGithubApp, s.authGuard.RequireRole("admin"))
	apiGroup.DELETE("/settings/git_apps/github/:id", s.gitAppsHandler.DeleteGithubApp, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/notifications/test", s.notificationHandler.TestNotification, s.authGuard.RequireRole("admin"))
}

func (s *Server) registerMiscRoutes(apiGroup, authGroup *echo.Group) {
	authGroup.POST("/compose/deploy", s.composeHandler.Deploy)
	authGroup.POST("/compose/analyze", s.composeHandler.Analyze)
	authGroup.POST("/deploy/archive", s.archiveHandler.DeployArchive)
	authGroup.GET("/examples", s.exampleHandler.List)
	authGroup.GET("/one-click", s.oneClickHandler.List)
	authGroup.POST("/one-click/deploy", s.oneClickHandler.Deploy)
	authGroup.POST("/dns", s.dnsHandler.Create)
	authGroup.GET("/dns", s.dnsHandler.List)
	authGroup.PUT("/dns/:id", s.dnsHandler.Update)
	authGroup.DELETE("/dns/:id", s.dnsHandler.Delete)
	authGroup.GET("/scheduled-tasks", s.scheduledTaskHandler.ListProjectScheduledTasks)
	authGroup.POST("/scheduled-tasks", s.scheduledTaskHandler.Create)
	authGroup.GET("/scheduled-tasks/:id", s.scheduledTaskHandler.Get)
	authGroup.DELETE("/scheduled-tasks/:id", s.scheduledTaskHandler.Delete)
	authGroup.POST("/scheduled-tasks/:id/trigger", s.scheduledTaskHandler.Run)
	authGroup.POST("/git/connect", s.gitHandler.Connect)
	authGroup.GET("/git/status", s.gitHandler.Status)
	authGroup.DELETE("/git/connect/:provider", s.gitHandler.Disconnect)
	authGroup.GET("/git/repos", s.gitHandler.ListRepos)
	apiGroup.POST("/webhooks/git/services/:serviceId", s.webhookHandler.HandleServiceWebhook)
	apiGroup.POST("/webhooks/github/services/:serviceId", s.webhookHandler.HandleGitHubWebhook)
	authGroup.GET("/canvas/projects", s.canvasHandler.ListCanvasSummaries)
	authGroup.GET("/projects/:id/summary", s.canvasHandler.GetCanvasSummary)
	authGroup.GET("/environments/:id/canvas", s.canvasHandler.GetEnvironmentCanvas)
	authGroup.GET("/audit-logs", s.auditLogHandler.List)
	authGroup.GET("/mcp/sse", s.HandleMCPSSE)
	authGroup.POST("/mcp/messages", s.HandleMCPMessage)
	apiGroup.GET("/ws/terminal/:id", s.terminalHandler.HandleWebSocket)
	apiGroup.GET("/ws/services/:id/terminal", s.terminalHandler.HandleWebSocket)
	apiGroup.GET("/ws/worker", s.workerWSHandler.Connect)
}

func (s *Server) setupSPAFallback() {
	staticDir := os.Getenv("CODEDOCK_STATIC_DIR")

	if staticDir != "" {
		if stat, err := os.Stat(staticDir); err == nil && stat.IsDir() {
			s.router.GET("/*", func(c echo.Context) error {
				path := filepath.Join(staticDir, filepath.Clean(c.Request().URL.Path))
				if _, err := os.Stat(path); os.IsNotExist(err) {
					return c.File(filepath.Join(staticDir, "index.html"))
				}
				return c.File(path)
			})
			return
		}
	}

	dashboard.RegisterHandlers(s.router)
}
