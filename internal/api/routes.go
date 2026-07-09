package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Server) registerRoutes() {
	s.router.HandleFunc("POST /api/auth/register", s.handleRegister)
	s.router.HandleFunc("POST /api/auth/login", s.handleLogin)
	s.router.HandleFunc("GET /api/auth/me", s.RequireAuth(s.handleGetCurrentUser))
	s.router.HandleFunc("POST /api/auth/logout", s.handleLogout)

	s.router.HandleFunc("GET /api/projects", s.RequireAuth(s.handleListProjects))
	s.router.HandleFunc("POST /api/projects", s.RequireAuth(s.handleCreateProject))
	s.router.HandleFunc("GET /api/projects/{id}", s.RequireAuth(s.handleGetProject))
	s.router.HandleFunc("DELETE /api/projects/{id}", s.RequireAuth(s.handleDeleteProject))
	s.router.HandleFunc("POST /api/projects/{id}/deploy", s.RequireAuth(s.handleDeployProject))

	s.router.HandleFunc("GET /api/projects/{id}/domains", s.RequireAuth(s.handleListDomains))
	s.router.HandleFunc("POST /api/projects/{id}/domains", s.RequireAuth(s.handleAddDomain))
	s.router.HandleFunc("DELETE /api/domains/{id}", s.RequireAuth(s.handleDeleteDomain))

	s.router.HandleFunc("GET /api/projects/{id}/env", s.RequireAuth(s.handleGetEnvVars))
	s.router.HandleFunc("PUT /api/projects/{id}/env", s.RequireAuth(s.handleSetEnvVars))

	s.router.HandleFunc("GET /api/databases", s.RequireAuth(s.handleListDatabases))
	s.router.HandleFunc("POST /api/databases", s.RequireAuth(s.handleCreateDatabase))
	s.router.HandleFunc("GET /api/databases/{id}", s.RequireAuth(s.handleGetDatabase))
	s.router.HandleFunc("DELETE /api/databases/{id}", s.RequireAuth(s.handleDeleteDatabase))
	s.router.HandleFunc("POST /api/databases/{id}/start", s.RequireAuth(s.handleStartDatabase))
	s.router.HandleFunc("POST /api/databases/{id}/stop", s.RequireAuth(s.handleStopDatabase))

	s.router.HandleFunc("GET /api/storage", s.RequireAuth(s.handleListStorage))
	s.router.HandleFunc("POST /api/storage", s.RequireAuth(s.handleCreateStorage))
	s.router.HandleFunc("GET /api/storage/{id}", s.RequireAuth(s.handleGetStorage))
	s.router.HandleFunc("DELETE /api/storage/{id}", s.RequireAuth(s.handleDeleteStorage))
	s.router.HandleFunc("POST /api/storage/{id}/start", s.RequireAuth(s.handleStartStorage))
	s.router.HandleFunc("POST /api/storage/{id}/stop", s.RequireAuth(s.handleStopStorage))

	s.router.HandleFunc("GET /api/jobs", s.RequireAuth(s.handleJobs))
	s.router.HandleFunc("POST /api/jobs", s.RequireAuth(s.handleJobs))
	s.router.HandleFunc("GET /api/jobs/{id}", s.RequireAuth(s.handleJobDetail))
	s.router.HandleFunc("DELETE /api/jobs/{id}", s.RequireAuth(s.handleJobDetail))
	s.router.HandleFunc("POST /api/jobs/{id}/trigger", s.RequireAuth(s.handleJobDetail))

	s.router.HandleFunc("POST /api/git/connect", s.RequireAuth(s.handleConnectGitProvider))
	s.router.HandleFunc("GET /api/git/status", s.RequireAuth(s.handleGetGitProvidersStatus))
	s.router.HandleFunc("DELETE /api/git/connect/{provider}", s.RequireAuth(s.handleDisconnectGitProvider))
	s.router.HandleFunc("GET /api/git/repos", s.RequireAuth(s.handleListGitRepositories))
	s.router.HandleFunc("POST /api/webhooks/git/{projectId}", s.handleGitWebhook)
	s.router.HandleFunc("POST /api/webhooks/git/services/{serviceId}", s.handleServiceGitWebhook)
	s.router.HandleFunc("GET /api/canvas/projects", s.RequireAuth(s.ListProjectCanvasSummaries))
	s.router.HandleFunc("GET /api/projects/{id}/summary", s.RequireAuth(s.GetProjectCanvasSummary))
	s.router.HandleFunc("GET /api/environments/{id}/canvas", s.RequireAuth(s.GetEnvironmentCanvas))

	s.router.HandleFunc("POST /api/projects/{id}/environments", s.RequireAuth(s.CreateEnvironment))
	s.router.HandleFunc("GET /api/projects/{id}/environments", s.RequireAuth(s.ListEnvironments))
	s.router.HandleFunc("DELETE /api/environments/{id}", s.RequireAuth(s.DeleteEnvironment))

	s.router.HandleFunc("POST /api/environments/{id}/apps", s.RequireAuth(s.CreateAppService))
	s.router.HandleFunc("GET /api/environments/{id}/apps", s.RequireAuth(s.ListAppServicesByEnvironment))
	s.router.HandleFunc("GET /api/projects/{id}/apps", s.RequireAuth(s.ListAppServicesByProject))
	s.router.HandleFunc("GET /api/apps/{id}", s.RequireAuth(s.GetAppService))
	s.router.HandleFunc("PUT /api/apps/{id}", s.RequireAuth(s.UpdateAppService))
	s.router.HandleFunc("DELETE /api/apps/{id}", s.RequireAuth(s.DeleteAppService))

	// --- Service Deployments & Metrics tabs ---
	s.router.HandleFunc("GET /api/services/{serviceId}/deployments", s.RequireAuth(s.deploymentHandler.ListServiceDeployments))
	s.router.HandleFunc("POST /api/services/{serviceId}/deploy", s.RequireAuth(s.deploymentHandler.TriggerServiceDeployment))
	s.router.HandleFunc("POST /api/deployments/{id}/rollback", s.RequireAuth(s.deploymentHandler.RollbackDeployment))
	s.router.HandleFunc("GET /api/deployments/{id}/logs", s.RequireAuth(s.deploymentHandler.GetDeploymentLogs))
	s.router.HandleFunc("GET /api/services/{serviceId}/metrics", s.RequireAuth(s.deploymentHandler.GetServiceMetrics))

	// --- Service Variables & Raw Editor tabs ---
	s.router.HandleFunc("GET /api/services/{serviceId}/variables", s.RequireAuth(s.serviceVarHandler.ListServiceVariables))
	s.router.HandleFunc("POST /api/services/{serviceId}/variables", s.RequireAuth(s.serviceVarHandler.SetServiceVariable))
	s.router.HandleFunc("PUT /api/services/{serviceId}/variables/bulk", s.RequireAuth(s.serviceVarHandler.BulkSetServiceVariables))
	s.router.HandleFunc("DELETE /api/services/{serviceId}/variables/{id}", s.RequireAuth(s.serviceVarHandler.DeleteServiceVariable))

	// --- Project Settings tabs (Billing, Webhooks, Tokens, Members) ---
	s.router.HandleFunc("GET /api/projects/{projectId}/billing", s.RequireAuth(s.projectSettingsHandler.GetProjectBilling))
	s.router.HandleFunc("GET /api/projects/{projectId}/webhooks", s.RequireAuth(s.projectSettingsHandler.ListWebhooks))
	s.router.HandleFunc("POST /api/projects/{projectId}/webhooks", s.RequireAuth(s.projectSettingsHandler.CreateWebhook))
	s.router.HandleFunc("DELETE /api/projects/{projectId}/webhooks/{id}", s.RequireAuth(s.projectSettingsHandler.DeleteWebhook))
	s.router.HandleFunc("GET /api/projects/{projectId}/tokens", s.RequireAuth(s.projectSettingsHandler.ListTokens))
	s.router.HandleFunc("POST /api/projects/{projectId}/tokens", s.RequireAuth(s.projectSettingsHandler.CreateToken))
	s.router.HandleFunc("DELETE /api/projects/{projectId}/tokens/{id}", s.RequireAuth(s.projectSettingsHandler.DeleteToken))
	s.router.HandleFunc("GET /api/projects/{projectId}/members", s.RequireAuth(s.projectSettingsHandler.ListMembers))
	s.router.HandleFunc("POST /api/projects/{projectId}/members", s.RequireAuth(s.projectSettingsHandler.InviteMember))
	s.router.HandleFunc("DELETE /api/projects/{projectId}/members/{id}", s.RequireAuth(s.projectSettingsHandler.RemoveMember))

	// --- Automated Database & Volume Backups and S3 Destinations ---
	s.router.HandleFunc("GET /api/backups", s.RequireAuth(s.backupHandler.ListBackups))
	s.router.HandleFunc("POST /api/backups", s.RequireAuth(s.backupHandler.CreateBackup))
	s.router.HandleFunc("GET /api/backups/{id}", s.RequireAuth(s.backupHandler.GetBackup))
	s.router.HandleFunc("DELETE /api/backups/{id}", s.RequireAuth(s.backupHandler.DeleteBackup))
	s.router.HandleFunc("POST /api/backups/{id}/trigger", s.RequireAuth(s.backupHandler.TriggerBackup))
	s.router.HandleFunc("GET /api/backups/{id}/records", s.RequireAuth(s.backupHandler.ListBackupRecords))
	s.router.HandleFunc("GET /api/s3-destinations", s.RequireAuth(s.backupHandler.ListS3Destinations))
	s.router.HandleFunc("POST /api/s3-destinations", s.RequireAuth(s.backupHandler.CreateS3Destination))
	s.router.HandleFunc("DELETE /api/s3-destinations/{id}", s.RequireAuth(s.backupHandler.DeleteS3Destination))

	// --- Teams & Organizations Collaboration ---
	s.router.HandleFunc("GET /api/teams", s.RequireAuth(s.teamHandler.ListTeams))
	s.router.HandleFunc("POST /api/teams", s.RequireAuth(s.teamHandler.CreateTeam))
	s.router.HandleFunc("GET /api/teams/{id}", s.RequireAuth(s.teamHandler.GetTeam))
	s.router.HandleFunc("DELETE /api/teams/{id}", s.RequireAuth(s.teamHandler.DeleteTeam))
	s.router.HandleFunc("GET /api/teams/{id}/members", s.RequireAuth(s.teamHandler.ListMembers))
	s.router.HandleFunc("POST /api/teams/{id}/invite", s.RequireAuth(s.teamHandler.InviteMember))
	s.router.HandleFunc("DELETE /api/teams/{id}/members/{userId}", s.RequireAuth(s.teamHandler.RemoveMember))
	s.router.HandleFunc("GET /api/team-invites/{token}", s.teamHandler.GetInvite)
	s.router.HandleFunc("POST /api/team-invites/{token}/accept", s.RequireAuth(s.teamHandler.AcceptInvite))

	// --- Workspace Level Settings (Trusted Domains, SSH Keys, Audit Logs) ---
	s.router.HandleFunc("GET /api/workspaces", s.RequireAuth(s.workspaceHandler.ListWorkspaces))
	s.router.HandleFunc("POST /api/workspaces", s.RequireAuth(s.workspaceHandler.CreateWorkspace))
	s.router.HandleFunc("GET /api/workspaces/{id}", s.RequireAuth(s.workspaceHandler.GetWorkspace))
	s.router.HandleFunc("PUT /api/workspaces/{id}", s.RequireAuth(s.workspaceHandler.UpdateWorkspace))
	s.router.HandleFunc("DELETE /api/workspaces/{id}", s.RequireAuth(s.workspaceHandler.DeleteWorkspace))
	s.router.HandleFunc("GET /api/workspaces/{id}/projects", s.RequireAuth(s.workspaceHandler.ListWorkspaceProjects))
	s.router.HandleFunc("GET /api/teams/{teamId}/trusted-domains", s.RequireAuth(s.workspaceHandler.ListTrustedDomains))
	s.router.HandleFunc("POST /api/teams/{teamId}/trusted-domains", s.RequireAuth(s.workspaceHandler.CreateTrustedDomain))
	s.router.HandleFunc("DELETE /api/trusted-domains/{id}", s.RequireAuth(s.workspaceHandler.DeleteTrustedDomain))
	s.router.HandleFunc("GET /api/teams/{teamId}/ssh-keys", s.RequireAuth(s.workspaceHandler.ListSSHKeys))
	s.router.HandleFunc("POST /api/teams/{teamId}/ssh-keys", s.RequireAuth(s.workspaceHandler.CreateSSHKey))
	s.router.HandleFunc("DELETE /api/ssh-keys/{id}", s.RequireAuth(s.workspaceHandler.DeleteSSHKey))
	s.router.HandleFunc("GET /api/teams/{teamId}/audit-logs", s.RequireAuth(s.workspaceHandler.ListAuditLogs))

	// --- Global Server Settings, System Prune, Profile & PATs ---
	s.router.HandleFunc("GET /api/settings", s.RequireAuth(s.settingsHandler.GetServerSettings))
	s.router.HandleFunc("PUT /api/settings", s.RequireRole("admin", s.settingsHandler.UpdateServerSettings))
	s.router.HandleFunc("POST /api/settings/prune", s.RequireRole("admin", s.settingsHandler.TriggerSystemPrune))
	s.router.HandleFunc("GET /api/settings/updates/status", s.RequireAuth(s.settingsHandler.GetUpdateStatus))
	s.router.HandleFunc("POST /api/settings/updates/check", s.RequireRole("admin", s.settingsHandler.CheckUpdate))
	s.router.HandleFunc("POST /api/settings/updates/deploy", s.RequireRole("admin", s.settingsHandler.DeployUpdate))
	s.router.HandleFunc("GET /api/mcp", s.RequireAuth(s.settingsHandler.HandleMCP))
	s.router.HandleFunc("POST /api/mcp", s.RequireAuth(s.settingsHandler.HandleMCP))
	s.router.HandleFunc("GET /api/profile", s.RequireAuth(s.settingsHandler.GetProfile))
	s.router.HandleFunc("PUT /api/profile", s.RequireAuth(s.settingsHandler.UpdateProfile))
	s.router.HandleFunc("GET /api/profile/tokens", s.RequireAuth(s.settingsHandler.ListPATs))
	s.router.HandleFunc("POST /api/profile/tokens", s.RequireAuth(s.settingsHandler.CreatePAT))
	s.router.HandleFunc("DELETE /api/profile/tokens/{id}", s.RequireAuth(s.settingsHandler.DeletePAT))

	// --- Multi-Channel Notification Integrations & Project Preferences ---
	s.router.HandleFunc("GET /api/settings/notifications", s.RequireAuth(s.notificationHandler.GetIntegrations))
	s.router.HandleFunc("PUT /api/settings/notifications", s.RequireRole("admin", s.notificationHandler.SaveIntegrations))
	s.router.HandleFunc("POST /api/settings/notifications/test", s.RequireRole("admin", s.notificationHandler.TestNotification))
	s.router.HandleFunc("GET /api/projects/{id}/notifications", s.RequireAuth(s.notificationHandler.GetProjectPreferences))
	s.router.HandleFunc("PUT /api/projects/{id}/notifications", s.RequireAuth(s.notificationHandler.SaveProjectPreferences))

	// --- OAuth 2.0 / OpenID Connect & TOTP 2FA Authentication ---
	s.router.HandleFunc("GET /api/settings/oauth/providers", s.RequireAuth(s.oauthHandler.ListProviders))
	s.router.HandleFunc("PUT /api/settings/oauth/providers", s.RequireRole("admin", s.oauthHandler.SaveProvider))
	s.router.HandleFunc("GET /api/auth/oauth/{provider}", s.oauthHandler.OAuthRedirect)
	s.router.HandleFunc("GET /api/auth/oauth/{provider}/callback", s.oauthHandler.OAuthCallback)
	s.router.HandleFunc("POST /api/auth/2fa/setup", s.RequireAuth(s.oauthHandler.Setup2FA))
	s.router.HandleFunc("POST /api/auth/2fa/verify", s.RequireAuth(s.oauthHandler.Verify2FA))
	s.router.HandleFunc("POST /api/auth/2fa/disable", s.RequireAuth(s.oauthHandler.Disable2FA))

	s.router.HandleFunc("GET /ws/terminal/{id}", s.handleTerminalWebSocket)
	s.router.HandleFunc("GET /ws/services/{id}/terminal", s.handleTerminalWebSocket)

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
