package http

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"vessl.dev/vessl/dashboard"
	_ "vessl.dev/vessl/docs"
)

func (s *Server) registerRoutes() {
	apiGroup := s.router.Group("/api")

	s.router.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
	s.router.GET("/docs/*", echoSwagger.WrapHandler)

	authGroup := apiGroup.Group("")
	authGroup.Use(s.authGuard.RequireAuth())
	apiGroup.POST("/auth/signup", s.authHandler.Register)
	apiGroup.POST("/auth/signin", s.authHandler.Login)
	authGroup.GET("/auth/me", s.userHandler.GetProfile)
	apiGroup.POST("/auth/logout", s.authHandler.Logout)
	authGroup.GET("/projects", s.projectHandler.ListProjects)
	authGroup.POST("/projects", s.projectHandler.CreateProject)
	authGroup.GET("/projects/:id", s.projectHandler.GetProject)
	authGroup.DELETE("/projects/:id", s.projectHandler.DeleteProject)
	authGroup.POST("/projects/:id/deploy", s.deploymentHandler.DeployProject, s.authGuard.RequireScope("deploy:write"))
	authGroup.GET("/projects/:id/domains", s.domainHandler.ListByProject)
	authGroup.POST("/projects/:id/domains", s.domainHandler.Create)
	authGroup.DELETE("/domains/:id", s.domainHandler.Delete)
	authGroup.GET("/projects/:id/env", s.projectEnvHandler.GetVars, s.authGuard.RequireScope("env:read"))
	authGroup.PUT("/projects/:id/env", s.projectEnvHandler.SetVars, s.authGuard.RequireScope("env:write"))
	authGroup.GET("/databases", s.dbHandler.ListDatabases, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases", s.dbHandler.CreateDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.GET("/databases/:id", s.dbHandler.GetDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.DELETE("/databases/:id", s.dbHandler.DeleteDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/start", s.dbHandler.StartDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/stop", s.dbHandler.StopDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.POST("/databases/:id/query", s.dbHandler.QueryDatabase, s.authGuard.RequireScope("database:manage"))
	authGroup.GET("/storage", s.storageHandler.ListStorage)
	authGroup.POST("/storage", s.storageHandler.CreateStorage)
	authGroup.GET("/storage/:id", s.storageHandler.GetStorage)
	authGroup.DELETE("/storage/:id", s.storageHandler.DeleteStorage)
	authGroup.POST("/storage/:id/start", s.storageHandler.StartStorage)
	authGroup.POST("/storage/:id/stop", s.storageHandler.StopStorage)
	authGroup.GET("/jobs", s.jobHandler.ListProjectJobs)
	authGroup.POST("/jobs", s.jobHandler.Create)
	authGroup.GET("/jobs/:id", s.jobHandler.Get)
	authGroup.DELETE("/jobs/:id", s.jobHandler.Delete)
	authGroup.POST("/jobs/:id/trigger", s.jobHandler.Run)
	authGroup.POST("/git/connect", s.gitHandler.Connect)
	authGroup.GET("/git/status", s.gitHandler.Status)
	authGroup.DELETE("/git/connect/:provider", s.gitHandler.Disconnect)
	authGroup.GET("/git/repos", s.gitHandler.ListRepos)
	apiGroup.POST("/webhooks/git/:projectId", s.webhookHandler.HandleProjectWebhook)
	apiGroup.POST("/webhooks/git/services/:serviceId", s.webhookHandler.HandleServiceWebhook)
	apiGroup.POST("/webhooks/github/services/:serviceId", s.webhookHandler.HandleGitHubWebhook)
	authGroup.GET("/canvas/projects", s.canvasHandler.ListCanvasSummaries)
	authGroup.GET("/projects/:id/summary", s.canvasHandler.GetCanvasSummary)
	authGroup.GET("/environments/:id/canvas", s.canvasHandler.GetEnvironmentCanvas)
	authGroup.POST("/projects/:id/environments", s.environmentHandler.Create)
	authGroup.GET("/projects/:id/environments", s.environmentHandler.ListByProject)
	authGroup.DELETE("/environments/:id", s.environmentHandler.Delete)
	authGroup.POST("/environments/:id/apps", s.appServiceHandler.Create)
	authGroup.GET("/environments/:id/apps", s.appServiceHandler.ListByEnvironment)
	authGroup.GET("/projects/:id/apps", s.appServiceHandler.ListByProject)
	authGroup.GET("/apps/:id", s.appServiceHandler.Get)
	authGroup.PUT("/apps/:id", s.appServiceHandler.Update)
	authGroup.DELETE("/apps/:id", s.appServiceHandler.Delete)
	authGroup.GET("/services/:serviceId/deployments", s.deploymentHandler.ListServiceDeployments)
	authGroup.POST("/services/:serviceId/deploy", s.deploymentHandler.Trigger)
	authGroup.POST("/deployments/:id/rollback", s.deploymentHandler.Rollback)
	authGroup.GET("/deployments/:id/logs", s.deploymentHandler.GetLogs, s.authGuard.RequireScope("logs:read"))
	authGroup.GET("/services/:serviceId/metrics", s.deploymentHandler.GetMetrics)
	authGroup.GET("/services/:serviceId/variables", s.serviceVarHandler.List)
	authGroup.POST("/services/:serviceId/variables", s.serviceVarHandler.Create)
	authGroup.PUT("/services/:serviceId/variables/:id", s.serviceVarHandler.Update)
	authGroup.DELETE("/services/:serviceId/variables/:id", s.serviceVarHandler.Delete)
	authGroup.GET("/services/:serviceId/serverless/code", s.serverlessHandler.GetCode)
	authGroup.POST("/services/:serviceId/serverless/code", s.serverlessHandler.SaveCode)
	authGroup.GET("/projects/:projectId/webhooks", s.projectSettingsHandler.ListWebhooks)
	authGroup.POST("/projects/:projectId/webhooks", s.projectSettingsHandler.CreateWebhook)
	authGroup.DELETE("/projects/:projectId/webhooks/:id", s.projectSettingsHandler.DeleteWebhook)
	authGroup.GET("/projects/:projectId/tokens", s.projectSettingsHandler.ListTokens, s.authGuard.RequireScope("env:read"))
	authGroup.POST("/projects/:projectId/tokens", s.projectSettingsHandler.CreateToken, s.authGuard.RequireScope("env:write"))
	authGroup.DELETE("/projects/:projectId/tokens/:id", s.projectSettingsHandler.DeleteToken, s.authGuard.RequireScope("env:write"))
	authGroup.GET("/projects/:projectId/members", s.projectSettingsHandler.ListMembers)
	authGroup.POST("/projects/:projectId/members", s.projectSettingsHandler.AddMember)
	authGroup.DELETE("/projects/:projectId/members/:id", s.projectSettingsHandler.RemoveMember)
	authGroup.GET("/backups", s.backupHandler.List)
	authGroup.POST("/backups", s.backupHandler.Create)
	authGroup.GET("/backups/:id", s.backupHandler.Get)
	authGroup.DELETE("/backups/:id", s.backupHandler.Delete)
	authGroup.POST("/backups/:id/trigger", s.backupHandler.Trigger)
	authGroup.GET("/backups/:id/records", s.backupHandler.ListRecords)
	authGroup.GET("/s3-destinations", s.backupHandler.ListS3Destinations)
	authGroup.POST("/s3-destinations", s.backupHandler.CreateS3Destination)
	authGroup.DELETE("/s3-destinations/:id", s.backupHandler.DeleteS3Destination)
	authGroup.GET("/teams", s.workspaceHandler.List)
	authGroup.POST("/teams", s.workspaceHandler.Create)
	authGroup.GET("/teams/:id", s.workspaceHandler.Get)
	authGroup.DELETE("/teams/:id", s.workspaceHandler.Delete)
	authGroup.GET("/teams/:id/members", s.workspaceHandler.ListMembers)
	authGroup.POST("/teams/:id/invite", s.workspaceHandler.InviteMember)
	authGroup.DELETE("/teams/:id/members/:userId", s.workspaceHandler.RemoveMember)
	apiGroup.GET("/team-invites/:token", s.workspaceHandler.GetInvite)
	authGroup.POST("/team-invites/:token/accept", s.workspaceHandler.AcceptInvite)
	authGroup.GET("/workspaces", s.workspaceHandler.List)
	authGroup.POST("/workspaces", s.workspaceHandler.Create)
	authGroup.GET("/workspaces/:id", s.workspaceHandler.Get)
	authGroup.PUT("/workspaces/:id", s.workspaceHandler.Update)
	authGroup.DELETE("/workspaces/:id", s.workspaceHandler.Delete)
	authGroup.GET("/teams/:teamId/trusted-domains", s.workspaceHandler.ListTrustedDomains)
	authGroup.POST("/teams/:teamId/trusted-domains", s.workspaceHandler.CreateTrustedDomain)
	authGroup.DELETE("/trusted-domains/:id", s.workspaceHandler.DeleteTrustedDomain)
	authGroup.GET("/teams/:teamId/ssh-keys", s.workspaceHandler.ListSSHKeys)
	authGroup.POST("/teams/:teamId/ssh-keys", s.workspaceHandler.CreateSSHKey)
	authGroup.DELETE("/ssh-keys/:id", s.workspaceHandler.DeleteSSHKey)
	authGroup.GET("/teams/:teamId/audit-logs", s.workspaceHandler.ListAuditLogs)
	authGroup.GET("/settings", s.settingsHandler.GetSettings)
	apiGroup.PUT("/settings", s.settingsHandler.UpdateSettings, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/license", s.settingsHandler.ActivateLicense, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/updates/status", s.updaterHandler.GetUpdateStatus)
	apiGroup.POST("/settings/updates/check", s.updaterHandler.CheckUpdate, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/updates/deploy", s.updaterHandler.DeployUpdate, s.authGuard.RequireRole("admin"))
	authGroup.GET("/mcp/sse", s.HandleMCPSSE)
	authGroup.POST("/mcp/messages", s.HandleMCPMessage)
	authGroup.GET("/profile", s.userHandler.GetProfile)
	authGroup.PUT("/profile", s.userHandler.UpdateProfile)
	authGroup.GET("/profile/tokens", s.userHandler.ListPATs)
	authGroup.POST("/profile/tokens", s.userHandler.CreatePAT)
	authGroup.DELETE("/profile/tokens/:id", s.userHandler.DeletePAT)
	authGroup.GET("/settings/notifications", s.notificationHandler.ListChannels)
	apiGroup.PUT("/settings/notifications", s.notificationHandler.SaveChannel, s.authGuard.RequireRole("admin"))
	apiGroup.POST("/settings/notifications/test", s.notificationHandler.TestNotification, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/notifications/:id", s.settingsHandler.GetTeamNotificationChannel)
	authGroup.DELETE("/settings/notifications/:id", s.notificationHandler.DeleteChannel)
	authGroup.GET("/settings/oauth/providers", s.oauthHandler.ListProviders)
	apiGroup.PUT("/settings/oauth/providers", s.oauthHandler.SaveProvider, s.authGuard.RequireRole("admin"))
	apiGroup.GET("/auth/oauth/:provider", s.oauthHandler.OAuthRedirect)
	apiGroup.GET("/auth/oauth/:provider/callback", s.oauthHandler.OAuthCallback)
	authGroup.POST("/auth/2fa/setup", s.oauthHandler.Setup2FA)
	authGroup.POST("/auth/2fa/verify", s.oauthHandler.Verify2FA)
	authGroup.POST("/auth/2fa/disable", s.oauthHandler.Disable2FA)
	apiGroup.GET("/ws/terminal/:id", s.terminalHandler.HandleWebSocket)

	apiGroup.POST("/settings/git_apps/github/manifest-callback", s.gitAppsHandler.ExchangeGithubManifestCode, s.authGuard.RequireRole("admin"))
	authGroup.GET("/settings/git_apps/github", s.gitAppsHandler.ListGithubApps)
	authGroup.GET("/settings/git_apps/github/:id", s.gitAppsHandler.GetGithubApp)
	apiGroup.PUT("/settings/git_apps/github", s.gitAppsHandler.SaveGithubApp, s.authGuard.RequireRole("admin"))
	apiGroup.DELETE("/settings/git_apps/github/:id", s.gitAppsHandler.DeleteGithubApp, s.authGuard.RequireRole("admin"))

	authGroup.GET("/settings/git_apps/gitlab", s.gitAppsHandler.ListGitlabApps)
	authGroup.GET("/settings/git_apps/gitlab/:id", s.gitAppsHandler.GetGitlabApp)
	apiGroup.PUT("/settings/git_apps/gitlab", s.gitAppsHandler.SaveGitlabApp, s.authGuard.RequireRole("admin"))
	apiGroup.DELETE("/settings/git_apps/gitlab/:id", s.gitAppsHandler.DeleteGitlabApp, s.authGuard.RequireRole("admin"))

	authGroup.GET("/settings/git_apps/bitbucket", s.gitAppsHandler.ListBitbucketApps)
	authGroup.GET("/settings/git_apps/bitbucket/:id", s.gitAppsHandler.GetBitbucketApp)
	apiGroup.PUT("/settings/git_apps/bitbucket", s.gitAppsHandler.SaveBitbucketApp, s.authGuard.RequireRole("admin"))
	apiGroup.DELETE("/settings/git_apps/bitbucket/:id", s.gitAppsHandler.DeleteBitbucketApp, s.authGuard.RequireRole("admin"))

	authGroup.GET("/teams/:teamId/ai_settings", s.aiSettingsHandler.Get)
	apiGroup.PUT("/teams/:teamId/ai_settings", s.aiSettingsHandler.Save, s.authGuard.RequireAuth())

	authGroup.GET("/teams/:teamId/email_settings", s.emailSettingsHandler.GetTeamEmailSettings)
	apiGroup.PUT("/teams/:teamId/email_settings", s.emailSettingsHandler.SaveTeamEmailSettings, s.authGuard.RequireAuth())
	authGroup.POST("/deployments/:id/diagnostics", s.aiDiagnosticsHandler.Analyze)

	authGroup.GET("/oauth/vercel/callback", s.vercelHandler.Callback)
	authGroup.GET("/vercel/projects", s.vercelHandler.ListProjects)
	authGroup.GET("/vercel/projects/:id/env", s.vercelHandler.GetProjectEnv)

	apiGroup.GET("/ws/services/:id/terminal", s.terminalHandler.HandleWebSocket)
	s.setupSPAFallback()
}

func (s *Server) setupSPAFallback() {
	staticDir := os.Getenv("VESSL_STATIC_DIR")

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

	s.router.GET("/*", func(c echo.Context) error {
		reqPath := filepath.Clean(c.Request().URL.Path)
		if reqPath == "/" || reqPath == "." {
			reqPath = "index.html"
		}

		content, err := dashboard.DistFS.ReadFile("dist/" + reqPath)
		if err != nil {
			indexContent, err := dashboard.DistFS.ReadFile("dist/index.html")
			if err != nil {
				return c.String(http.StatusNotFound, "Dashboard not built. Please run 'npm run build' in the dashboard directory.")
			}
			return c.HTMLBlob(http.StatusOK, indexContent)
		}

		contentType := http.DetectContentType(content)
		if filepath.Ext(reqPath) == ".css" {
			contentType = "text/css"
		} else if filepath.Ext(reqPath) == ".js" {
			contentType = "application/javascript"
		} else if filepath.Ext(reqPath) == ".svg" {
			contentType = "image/svg+xml"
		}
		return c.Blob(http.StatusOK, contentType, content)
	})
}
