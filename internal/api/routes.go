package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Server) registerRoutes() {
	s.router.HandleFunc("POST /api/auth/signup", s.authHandler.Register)
	s.router.HandleFunc("POST /api/auth/signin", s.authHandler.Login)
	s.router.HandleFunc("GET /api/auth/me", s.RequireAuth(s.userHandler.GetProfile))
	s.router.HandleFunc("POST /api/auth/logout", s.authHandler.Logout)

	s.router.HandleFunc("GET /api/projects", s.RequireAuth(s.projectHandler.ListProjects))
	s.router.HandleFunc("POST /api/projects", s.RequireAuth(s.projectHandler.CreateProject))
	s.router.HandleFunc("GET /api/projects/{id}", s.RequireAuth(s.projectHandler.GetProject))
	s.router.HandleFunc("DELETE /api/projects/{id}", s.RequireAuth(s.projectHandler.DeleteProject))
	s.router.HandleFunc("POST /api/projects/{id}/deploy", s.RequireAuth(s.deploymentHandler.DeployProject))

	s.router.HandleFunc("GET /api/projects/{id}/domains", s.RequireAuth(s.domainHandler.ListByProject))
	s.router.HandleFunc("POST /api/projects/{id}/domains", s.RequireAuth(s.domainHandler.Create))
	s.router.HandleFunc("DELETE /api/domains/{id}", s.RequireAuth(s.domainHandler.Delete))

	s.router.HandleFunc("GET /api/projects/{id}/env", s.RequireAuth(s.projectEnvHandler.GetVars))
	s.router.HandleFunc("PUT /api/projects/{id}/env", s.RequireAuth(s.projectEnvHandler.SetVars))

	s.router.HandleFunc("GET /api/databases", s.RequireAuth(s.dbHandler.ListDatabases))
	s.router.HandleFunc("POST /api/databases", s.RequireAuth(s.dbHandler.CreateDatabase))
	s.router.HandleFunc("GET /api/databases/{id}", s.RequireAuth(s.dbHandler.GetDatabase))
	s.router.HandleFunc("DELETE /api/databases/{id}", s.RequireAuth(s.dbHandler.DeleteDatabase))
	s.router.HandleFunc("POST /api/databases/{id}/start", s.RequireAuth(s.dbHandler.StartDatabase))
	s.router.HandleFunc("POST /api/databases/{id}/stop", s.RequireAuth(s.dbHandler.StopDatabase))

	s.router.HandleFunc("GET /api/storage", s.RequireAuth(s.storageHandler.ListStorage))
	s.router.HandleFunc("POST /api/storage", s.RequireAuth(s.storageHandler.CreateStorage))
	s.router.HandleFunc("GET /api/storage/{id}", s.RequireAuth(s.storageHandler.GetStorage))
	s.router.HandleFunc("DELETE /api/storage/{id}", s.RequireAuth(s.storageHandler.DeleteStorage))
	s.router.HandleFunc("POST /api/storage/{id}/start", s.RequireAuth(s.storageHandler.StartStorage))
	s.router.HandleFunc("POST /api/storage/{id}/stop", s.RequireAuth(s.storageHandler.StopStorage))

	s.router.HandleFunc("GET /api/jobs", s.RequireAuth(s.jobHandler.ListProjectJobs))
	s.router.HandleFunc("POST /api/jobs", s.RequireAuth(s.jobHandler.Create))
	s.router.HandleFunc("GET /api/jobs/{id}", s.RequireAuth(s.jobHandler.Get))
	s.router.HandleFunc("DELETE /api/jobs/{id}", s.RequireAuth(s.jobHandler.Delete))
	s.router.HandleFunc("POST /api/jobs/{id}/trigger", s.RequireAuth(s.jobHandler.Run))

	s.router.HandleFunc("POST /api/git/connect", s.RequireAuth(s.gitHandler.Connect))
	s.router.HandleFunc("GET /api/git/status", s.RequireAuth(s.gitHandler.Status))
	s.router.HandleFunc("DELETE /api/git/connect/{provider}", s.RequireAuth(s.gitHandler.Disconnect))
	s.router.HandleFunc("GET /api/git/repos", s.RequireAuth(s.gitHandler.ListRepos))
	s.router.HandleFunc("POST /api/webhooks/git/{projectId}", s.webhookHandler.HandleProjectWebhook)
	s.router.HandleFunc("POST /api/webhooks/git/services/{serviceId}", s.webhookHandler.HandleServiceWebhook)

	s.router.HandleFunc("GET /api/canvas/projects", s.RequireAuth(s.canvasHandler.ListCanvasSummaries))
	s.router.HandleFunc("GET /api/projects/{id}/summary", s.RequireAuth(s.canvasHandler.GetCanvasSummary))
	s.router.HandleFunc("GET /api/environments/{id}/canvas", s.RequireAuth(s.canvasHandler.GetEnvironmentCanvas))

	s.router.HandleFunc("POST /api/projects/{id}/environments", s.RequireAuth(s.environmentHandler.Create))
	s.router.HandleFunc("GET /api/projects/{id}/environments", s.RequireAuth(s.environmentHandler.ListByProject))
	s.router.HandleFunc("DELETE /api/environments/{id}", s.RequireAuth(s.environmentHandler.Delete))

	s.router.HandleFunc("POST /api/environments/{id}/apps", s.RequireAuth(s.appServiceHandler.Create))
	s.router.HandleFunc("GET /api/environments/{id}/apps", s.RequireAuth(s.appServiceHandler.ListByEnvironment))
	s.router.HandleFunc("GET /api/projects/{id}/apps", s.RequireAuth(s.appServiceHandler.ListByProject))
	s.router.HandleFunc("GET /api/apps/{id}", s.RequireAuth(s.appServiceHandler.Get))
	s.router.HandleFunc("PUT /api/apps/{id}", s.RequireAuth(s.appServiceHandler.Update))
	s.router.HandleFunc("DELETE /api/apps/{id}", s.RequireAuth(s.appServiceHandler.Delete))

	s.router.HandleFunc("GET /api/services/{serviceId}/deployments", s.RequireAuth(s.deploymentHandler.ListServiceDeployments))
	s.router.HandleFunc("POST /api/services/{serviceId}/deploy", s.RequireAuth(s.deploymentHandler.Trigger))
	s.router.HandleFunc("POST /api/deployments/{id}/rollback", s.RequireAuth(s.deploymentHandler.Rollback))
	s.router.HandleFunc("GET /api/deployments/{id}/logs", s.RequireAuth(s.deploymentHandler.GetLogs))
	s.router.HandleFunc("GET /api/services/{serviceId}/metrics", s.RequireAuth(s.deploymentHandler.GetMetrics))

	s.router.HandleFunc("GET /api/services/{serviceId}/variables", s.RequireAuth(s.serviceVarHandler.List))
	s.router.HandleFunc("POST /api/services/{serviceId}/variables", s.RequireAuth(s.serviceVarHandler.Create))
	s.router.HandleFunc("PUT /api/services/{serviceId}/variables/{id}", s.RequireAuth(s.serviceVarHandler.Update))
	s.router.HandleFunc("DELETE /api/services/{serviceId}/variables/{id}", s.RequireAuth(s.serviceVarHandler.Delete))

	s.router.HandleFunc("GET /api/projects/{projectId}/webhooks", s.RequireAuth(s.projectSettingsHandler.ListWebhooks))
	s.router.HandleFunc("POST /api/projects/{projectId}/webhooks", s.RequireAuth(s.projectSettingsHandler.CreateWebhook))
	s.router.HandleFunc("DELETE /api/projects/{projectId}/webhooks/{id}", s.RequireAuth(s.projectSettingsHandler.DeleteWebhook))
	s.router.HandleFunc("GET /api/projects/{projectId}/tokens", s.RequireAuth(s.projectSettingsHandler.ListTokens))
	s.router.HandleFunc("POST /api/projects/{projectId}/tokens", s.RequireAuth(s.projectSettingsHandler.CreateToken))
	s.router.HandleFunc("DELETE /api/projects/{projectId}/tokens/{id}", s.RequireAuth(s.projectSettingsHandler.DeleteToken))
	s.router.HandleFunc("GET /api/projects/{projectId}/members", s.RequireAuth(s.projectSettingsHandler.ListMembers))
	s.router.HandleFunc("POST /api/projects/{projectId}/members", s.RequireAuth(s.projectSettingsHandler.AddMember))
	s.router.HandleFunc("DELETE /api/projects/{projectId}/members/{id}", s.RequireAuth(s.projectSettingsHandler.RemoveMember))

	s.router.HandleFunc("GET /api/backups", s.RequireAuth(s.backupHandler.List))
	s.router.HandleFunc("POST /api/backups", s.RequireAuth(s.backupHandler.Create))
	s.router.HandleFunc("GET /api/backups/{id}", s.RequireAuth(s.backupHandler.Get))
	s.router.HandleFunc("DELETE /api/backups/{id}", s.RequireAuth(s.backupHandler.Delete))
	s.router.HandleFunc("POST /api/backups/{id}/trigger", s.RequireAuth(s.backupHandler.Trigger))
	s.router.HandleFunc("GET /api/backups/{id}/records", s.RequireAuth(s.backupHandler.ListRecords))
	s.router.HandleFunc("GET /api/s3-destinations", s.RequireAuth(s.backupHandler.ListS3Destinations))
	s.router.HandleFunc("POST /api/s3-destinations", s.RequireAuth(s.backupHandler.CreateS3Destination))
	s.router.HandleFunc("DELETE /api/s3-destinations/{id}", s.RequireAuth(s.backupHandler.DeleteS3Destination))

	s.router.HandleFunc("GET /api/teams", s.RequireAuth(s.teamHandler.List))
	s.router.HandleFunc("POST /api/teams", s.RequireAuth(s.teamHandler.Create))
	s.router.HandleFunc("GET /api/teams/{id}", s.RequireAuth(s.teamHandler.Get))
	s.router.HandleFunc("DELETE /api/teams/{id}", s.RequireAuth(s.teamHandler.Delete))
	s.router.HandleFunc("GET /api/teams/{id}/members", s.RequireAuth(s.teamHandler.ListMembers))
	s.router.HandleFunc("POST /api/teams/{id}/invite", s.RequireAuth(s.teamHandler.InviteMember))
	s.router.HandleFunc("DELETE /api/teams/{id}/members/{userId}", s.RequireAuth(s.teamHandler.RemoveMember))
	s.router.HandleFunc("GET /api/team-invites/{token}", s.teamHandler.GetInvite)
	s.router.HandleFunc("POST /api/team-invites/{token}/accept", s.RequireAuth(s.teamHandler.AcceptInvite))

	s.router.HandleFunc("GET /api/workspaces", s.RequireAuth(s.workspaceHandler.List))
	s.router.HandleFunc("POST /api/workspaces", s.RequireAuth(s.workspaceHandler.Create))
	s.router.HandleFunc("GET /api/workspaces/{id}", s.RequireAuth(s.workspaceHandler.Get))
	s.router.HandleFunc("PUT /api/workspaces/{id}", s.RequireAuth(s.workspaceHandler.Update))
	s.router.HandleFunc("DELETE /api/workspaces/{id}", s.RequireAuth(s.workspaceHandler.Delete))
	s.router.HandleFunc("GET /api/teams/{teamId}/trusted-domains", s.RequireAuth(s.workspaceHandler.ListTrustedDomains))
	s.router.HandleFunc("POST /api/teams/{teamId}/trusted-domains", s.RequireAuth(s.workspaceHandler.CreateTrustedDomain))
	s.router.HandleFunc("DELETE /api/trusted-domains/{id}", s.RequireAuth(s.workspaceHandler.DeleteTrustedDomain))
	s.router.HandleFunc("GET /api/teams/{teamId}/ssh-keys", s.RequireAuth(s.workspaceHandler.ListSSHKeys))
	s.router.HandleFunc("POST /api/teams/{teamId}/ssh-keys", s.RequireAuth(s.workspaceHandler.CreateSSHKey))
	s.router.HandleFunc("DELETE /api/ssh-keys/{id}", s.RequireAuth(s.workspaceHandler.DeleteSSHKey))
	s.router.HandleFunc("GET /api/teams/{teamId}/audit-logs", s.RequireAuth(s.workspaceHandler.ListAuditLogs))

	s.router.HandleFunc("GET /api/settings", s.RequireAuth(s.settingsHandler.GetSettings))
	s.router.HandleFunc("PUT /api/settings", s.RequireRole("admin", s.settingsHandler.UpdateSettings))
	s.router.HandleFunc("GET /api/settings/updates/status", s.RequireAuth(s.updaterHandler.GetUpdateStatus))
	s.router.HandleFunc("POST /api/settings/updates/check", s.RequireRole("admin", s.updaterHandler.CheckUpdate))
	s.router.HandleFunc("POST /api/settings/updates/deploy", s.RequireRole("admin", s.updaterHandler.DeployUpdate))
	s.router.HandleFunc("GET /api/mcp", s.RequireAuth(s.settingsHandler.HandleMCPRequest))
	s.router.HandleFunc("POST /api/mcp", s.RequireAuth(s.settingsHandler.HandleMCPRequest))

	s.router.HandleFunc("GET /api/profile", s.RequireAuth(s.userHandler.GetProfile))
	s.router.HandleFunc("PUT /api/profile", s.RequireAuth(s.userHandler.UpdateProfile))
	s.router.HandleFunc("GET /api/profile/tokens", s.RequireAuth(s.userHandler.ListPATs))
	s.router.HandleFunc("POST /api/profile/tokens", s.RequireAuth(s.userHandler.CreatePAT))
	s.router.HandleFunc("DELETE /api/profile/tokens/{id}", s.RequireAuth(s.userHandler.DeletePAT))

	s.router.HandleFunc("GET /api/settings/notifications", s.RequireAuth(s.notificationHandler.GetIntegrations))
	s.router.HandleFunc("PUT /api/settings/notifications", s.RequireRole("admin", s.notificationHandler.SaveIntegrations))
	s.router.HandleFunc("POST /api/settings/notifications/test", s.RequireRole("admin", s.notificationHandler.TestNotification))
	s.router.HandleFunc("GET /api/projects/{id}/notifications", s.RequireAuth(s.notificationHandler.GetProjectPreferences))
	s.router.HandleFunc("PUT /api/projects/{id}/notifications", s.RequireAuth(s.notificationHandler.SaveProjectPreferences))

	s.router.HandleFunc("GET /api/settings/oauth/providers", s.RequireAuth(s.oauthHandler.ListProviders))
	s.router.HandleFunc("PUT /api/settings/oauth/providers", s.RequireRole("admin", s.oauthHandler.SaveProvider))
	s.router.HandleFunc("GET /api/auth/oauth/{provider}", s.oauthHandler.OAuthRedirect)
	s.router.HandleFunc("GET /api/auth/oauth/{provider}/callback", s.oauthHandler.OAuthCallback)
	s.router.HandleFunc("POST /api/auth/2fa/setup", s.RequireAuth(s.oauthHandler.Setup2FA))
	s.router.HandleFunc("POST /api/auth/2fa/verify", s.RequireAuth(s.oauthHandler.Verify2FA))
	s.router.HandleFunc("POST /api/auth/2fa/disable", s.RequireAuth(s.oauthHandler.Disable2FA))

	s.router.HandleFunc("GET /ws/terminal/{id}", s.terminalHandler.HandleWebSocket)
	s.router.HandleFunc("GET /ws/services/{id}/terminal", s.terminalHandler.HandleWebSocket)

	s.setupSPAFallback()
}

func (s *Server) setupSPAFallback() {
	staticDir := os.Getenv("VESSEL_STATIC_DIR")
	if staticDir == "" {
		staticDir = "dashboard/dist"
	}
	if stat, err := os.Stat(staticDir); err == nil && stat.IsDir() {
		fileServer := http.FileServer(http.Dir(staticDir))
		s.router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws/") {
				http.NotFound(w, r)
				return
			}
			path := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
			if _, err := os.Stat(path); os.IsNotExist(err) {
				http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	}
}
