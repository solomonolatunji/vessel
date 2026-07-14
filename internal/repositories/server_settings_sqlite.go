package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
)

type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*models.ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
}

type SettingsSQLiteRepository struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewSettingsSQLiteRepository(db *sql.DB) *SettingsSQLiteRepository {
	return &SettingsSQLiteRepository{db: sqlx.NewDb(db, "sqlite")}
}

const serverSettingsColumns = `id, traefik_wildcard_ip, discord_webhook_url, discord_ping_enabled, discord_enabled, slack_webhook_url, slack_enabled, telegram_bot_token, telegram_chat_id, telegram_enabled, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, smtp_enabled, resend_api_key, resend_enabled, pushover_user_key, pushover_api_token, pushover_enabled, generic_webhook_url, generic_webhook_enabled, notification_alerts, registration_enabled, registration_domain_allowlist, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, default_wildcard_domain, site_name, public_ipv4, public_ipv6, show_sponsorship_popup, disable_two_step_confirmation, default_openai_key, default_anthropic_key, update_check_cron, auto_update_enabled, concurrent_builds, deployment_timeout, server_timezone, docker_cleanup_cron, disk_usage_threshold, disk_usage_cron, current_version, latest_version, last_update_check, updated_at`

func scanServerSettings(scanner interface{ Scan(dest ...any) error }, cfg *models.ServerSettings) error {
	return scanner.Scan(
		&cfg.ID, &cfg.TraefikWildcardIP, &cfg.DiscordWebhookURL, &cfg.DiscordPingEnabled, &cfg.DiscordEnabled, &cfg.SlackWebhookURL, &cfg.SlackEnabled, &cfg.TelegramBotToken, &cfg.TelegramChatID, &cfg.TelegramEnabled,
		&cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.SMTPEnabled,
		&cfg.ResendAPIKey, &cfg.ResendEnabled, &cfg.PushoverUserKey, &cfg.PushoverAPIToken, &cfg.PushoverEnabled, &cfg.GenericWebhookURL, &cfg.GenericWebhookEnabled,
		&cfg.NotificationAlerts,
		&cfg.RegistrationEnabled, &cfg.RegistrationDomainAllowlist, &cfg.CustomDNSResolvers, &cfg.DNSValidationEnabled, &cfg.IPAllowlist, &cfg.MCPServerEnabled, &cfg.DefaultWildcardDomain,
		&cfg.SiteName, &cfg.PublicIPv4, &cfg.PublicIPv6, &cfg.ShowSponsorshipPopup, &cfg.DisableTwoStepConfirmation,
		&cfg.DefaultOpenAIKey, &cfg.DefaultAnthropicKey,
		&cfg.UpdateCheckCron, &cfg.AutoUpdateEnabled,
		&cfg.ConcurrentBuilds, &cfg.DeploymentTimeout, &cfg.ServerTimezone, &cfg.DockerCleanupCron, &cfg.DiskUsageThreshold, &cfg.DiskUsageCron,
		&cfg.CurrentVersion, &cfg.LatestVersion, &cfg.LastUpdateCheck, &cfg.UpdatedAt,
	)
}

func serverSettingsArgs(cfg *models.ServerSettings) []any {
	return []any{
		cfg.ID, cfg.TraefikWildcardIP, cfg.DiscordWebhookURL, cfg.DiscordPingEnabled, cfg.DiscordEnabled, cfg.SlackWebhookURL, cfg.SlackEnabled, cfg.TelegramBotToken, cfg.TelegramChatID, cfg.TelegramEnabled,
		cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.SMTPEnabled,
		cfg.ResendAPIKey, cfg.ResendEnabled, cfg.PushoverUserKey, cfg.PushoverAPIToken, cfg.PushoverEnabled, cfg.GenericWebhookURL, cfg.GenericWebhookEnabled,
		cfg.NotificationAlerts,
		cfg.RegistrationEnabled, cfg.RegistrationDomainAllowlist, cfg.CustomDNSResolvers, cfg.DNSValidationEnabled, cfg.IPAllowlist, cfg.MCPServerEnabled, cfg.DefaultWildcardDomain,
		cfg.SiteName, cfg.PublicIPv4, cfg.PublicIPv6, cfg.ShowSponsorshipPopup, cfg.DisableTwoStepConfirmation,
		cfg.DefaultOpenAIKey, cfg.DefaultAnthropicKey,
		cfg.UpdateCheckCron, cfg.AutoUpdateEnabled, cfg.CurrentVersion, cfg.LatestVersion, cfg.LastUpdateCheck, cfg.UpdatedAt,
	}
}

func (r *SettingsSQLiteRepository) GetServerSettings(ctx context.Context) (*models.ServerSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.ServerSettings
	err := scanServerSettings(r.db.QueryRowContext(ctx, `SELECT `+serverSettingsColumns+` FROM server_settings WHERE id = 'global' LIMIT 1`), &cfg)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.ServerSettings{
			ID:                   "global",
			RegistrationEnabled:  true,
			NotificationAlerts:   true,
			DNSValidationEnabled: true,
			CustomDNSResolvers:   "1.1.1.1",
			MCPServerEnabled:     true,
			UpdateCheckCron:      "0 * * * *",
			AutoUpdateEnabled:    false,
			CurrentVersion:       "0.1.0",
			LatestVersion:        "0.1.0",
			UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
			ShowSponsorshipPopup: true,
		}
		query := fmt.Sprintf(`INSERT INTO server_settings (%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, serverSettingsColumns)
		_, _ = r.db.ExecContext(ctx, query, serverSettingsArgs(defaultSettings)...)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server settings: %w", err)
	}
	return &cfg, nil
}

func (r *SettingsSQLiteRepository) UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf(`INSERT INTO server_settings (%s)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(id) DO UPDATE SET
	          traefik_wildcard_ip = excluded.traefik_wildcard_ip,
	          discord_webhook_url = excluded.discord_webhook_url,
	          discord_ping_enabled = excluded.discord_ping_enabled,
	          discord_enabled = excluded.discord_enabled,
	          slack_webhook_url = excluded.slack_webhook_url,
	          slack_enabled = excluded.slack_enabled,
	          telegram_bot_token = excluded.telegram_bot_token,
	          telegram_chat_id = excluded.telegram_chat_id,
	          telegram_enabled = excluded.telegram_enabled,
	          smtp_host = excluded.smtp_host,
	          smtp_port = excluded.smtp_port,
	          smtp_user = excluded.smtp_user,
	          smtp_password = excluded.smtp_password,
	          smtp_from_name = excluded.smtp_from_name,
	          smtp_from_address = excluded.smtp_from_address,
	          smtp_enabled = excluded.smtp_enabled,
	          resend_api_key = excluded.resend_api_key,
	          resend_enabled = excluded.resend_enabled,
	          pushover_user_key = excluded.pushover_user_key,
	          pushover_api_token = excluded.pushover_api_token,
	          pushover_enabled = excluded.pushover_enabled,
	          generic_webhook_url = excluded.generic_webhook_url,
	          generic_webhook_enabled = excluded.generic_webhook_enabled,
	          notification_alerts = excluded.notification_alerts,
	          registration_enabled = excluded.registration_enabled,
	          registration_domain_allowlist = excluded.registration_domain_allowlist,
	          custom_dns_resolvers = excluded.custom_dns_resolvers,
	          dns_validation_enabled = excluded.dns_validation_enabled,
	          ip_allowlist = excluded.ip_allowlist,
	          mcp_server_enabled = excluded.mcp_server_enabled,
	          default_wildcard_domain = excluded.default_wildcard_domain,
	          site_name = excluded.site_name,
	          public_ipv4 = excluded.public_ipv4,
	          public_ipv6 = excluded.public_ipv6,
	          show_sponsorship_popup = excluded.show_sponsorship_popup,
	          disable_two_step_confirmation = excluded.disable_two_step_confirmation,
	          default_openai_key = excluded.default_openai_key,
	          default_anthropic_key = excluded.default_anthropic_key,
	          update_check_cron = excluded.update_check_cron,
	          auto_update_enabled = excluded.auto_update_enabled,
	          concurrent_builds = excluded.concurrent_builds,
	          deployment_timeout = excluded.deployment_timeout,
	          server_timezone = excluded.server_timezone,
	          docker_cleanup_cron = excluded.docker_cleanup_cron,
	          disk_usage_threshold = excluded.disk_usage_threshold,
	          disk_usage_cron = excluded.disk_usage_cron,
	          current_version = excluded.current_version,
	          latest_version = excluded.latest_version,
	          last_update_check = excluded.last_update_check,
	          updated_at = excluded.updated_at`, serverSettingsColumns)
	_, err := r.db.ExecContext(ctx, query, serverSettingsArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("failed to update server settings: %w", err)
	}
	return nil
}

func (r *SettingsSQLiteRepository) ListProjects(ctx context.Context) ([]map[string]any, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, repo_url, branch, status, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()
	var projects []map[string]any
	for rows.Next() {
		var id, name, repoURL, branch, status, updatedAt string
		if err := rows.Scan(&id, &name, &repoURL, &branch, &status, &updatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, map[string]any{
			"id":        id,
			"name":      name,
			"repoUrl":   repoURL,
			"branch":    branch,
			"status":    status,
			"updatedAt": updatedAt,
		})
	}
	return projects, nil
}
