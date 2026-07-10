package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"vessel.dev/vessel/internal/models"
)

type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*models.ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
}

type NotificationRepository interface {
	ListChannelsByTeam(ctx context.Context, teamID string) ([]models.TeamNotificationChannel, error)
	GetChannel(ctx context.Context, id string) (*models.TeamNotificationChannel, error)
	SaveChannel(ctx context.Context, c *models.TeamNotificationChannel) error
	DeleteChannel(ctx context.Context, id string) error
}

type SettingsSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSettingsSQLiteRepository(db *sql.DB) *SettingsSQLiteRepository {
	return &SettingsSQLiteRepository{db: db}
}

const serverSettingsColumns = `id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts, registration_enabled, registration_domain_allowlist, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, default_wildcard_domain, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at`

func scanServerSettings(scanner interface{ Scan(dest ...any) error }, cfg *models.ServerSettings) error {
	return scanner.Scan(
		&cfg.ID, &cfg.CaddyWildcardIP, &cfg.DiscordWebhookURL, &cfg.SlackWebhookURL, &cfg.TelegramBotToken, &cfg.TelegramChatID,
		&cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.NotificationAlerts,
		&cfg.RegistrationEnabled, &cfg.RegistrationDomainAllowlist, &cfg.CustomDNSResolvers, &cfg.DNSValidationEnabled, &cfg.IPAllowlist, &cfg.MCPServerEnabled, &cfg.DefaultWildcardDomain,
		&cfg.UpdateCheckCron, &cfg.AutoUpdateEnabled, &cfg.CurrentVersion, &cfg.LatestVersion, &cfg.LastUpdateCheck, &cfg.UpdatedAt,
	)
}

func serverSettingsArgs(cfg *models.ServerSettings) []any {
	return []any{
		cfg.ID, cfg.CaddyWildcardIP, cfg.DiscordWebhookURL, cfg.SlackWebhookURL, cfg.TelegramBotToken, cfg.TelegramChatID,
		cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.NotificationAlerts,
		cfg.RegistrationEnabled, cfg.RegistrationDomainAllowlist, cfg.CustomDNSResolvers, cfg.DNSValidationEnabled, cfg.IPAllowlist, cfg.MCPServerEnabled, cfg.DefaultWildcardDomain,
		cfg.UpdateCheckCron, cfg.AutoUpdateEnabled, cfg.CurrentVersion, cfg.LatestVersion, cfg.LastUpdateCheck, cfg.UpdatedAt,
	}
}

func (r *SettingsSQLiteRepository) GetServerSettings(ctx context.Context) (*models.ServerSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.ServerSettings
	row := r.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT %s FROM server_settings LIMIT 1`, serverSettingsColumns))
	err := scanServerSettings(row, &cfg)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.ServerSettings{
			ID:                   "global",
			CaddyWildcardIP:      "127.0.0.1",
			NotificationAlerts:   true,
			RegistrationEnabled:  true,
			DNSValidationEnabled: true,
			MCPServerEnabled:     true,
			UpdateCheckCron:      "0 * * * *",
			AutoUpdateEnabled:    false,
			CurrentVersion:       "0.1.0",
			LatestVersion:        "0.1.0",
			UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
		}
		query := fmt.Sprintf(`INSERT INTO server_settings (%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, serverSettingsColumns)
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
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(id) DO UPDATE SET
	          caddy_wildcard_ip = excluded.caddy_wildcard_ip,
	          discord_webhook_url = excluded.discord_webhook_url,
	          slack_webhook_url = excluded.slack_webhook_url,
	          telegram_bot_token = excluded.telegram_bot_token,
	          telegram_chat_id = excluded.telegram_chat_id,
	          smtp_host = excluded.smtp_host,
	          smtp_port = excluded.smtp_port,
	          smtp_user = excluded.smtp_user,
	          smtp_password = excluded.smtp_password,
	          smtp_from_name = excluded.smtp_from_name,
	          smtp_from_address = excluded.smtp_from_address,
	          notification_alerts = excluded.notification_alerts,
	          registration_enabled = excluded.registration_enabled,
	          registration_domain_allowlist = excluded.registration_domain_allowlist,
	          custom_dns_resolvers = excluded.custom_dns_resolvers,
	          dns_validation_enabled = excluded.dns_validation_enabled,
	          ip_allowlist = excluded.ip_allowlist,
	          mcp_server_enabled = excluded.mcp_server_enabled,
	          default_wildcard_domain = excluded.default_wildcard_domain,
	          update_check_cron = excluded.update_check_cron,
	          auto_update_enabled = excluded.auto_update_enabled,
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

type NotificationSQLiteRepository struct {
	db *sql.DB
}

func NewNotificationSQLiteRepository(db *sql.DB) *NotificationSQLiteRepository {
	return &NotificationSQLiteRepository{db: db}
}

func (r *NotificationSQLiteRepository) ListChannelsByTeam(ctx context.Context, teamID string) ([]models.TeamNotificationChannel, error) {
	query := `SELECT id, team_id, provider, config, events, is_enabled, created_at, updated_at FROM team_notification_channels WHERE team_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}
	defer rows.Close()

	var channels []models.TeamNotificationChannel
	for rows.Next() {
		var c models.TeamNotificationChannel
		var configStr, eventsStr string
		if err := rows.Scan(&c.ID, &c.TeamID, &c.Provider, &configStr, &eventsStr, &c.IsEnabled, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		c.Config = []byte(configStr)
		c.Events = []byte(eventsStr)
		channels = append(channels, c)
	}
	return channels, nil
}

func (r *NotificationSQLiteRepository) GetChannel(ctx context.Context, id string) (*models.TeamNotificationChannel, error) {
	query := `SELECT id, team_id, provider, config, events, is_enabled, created_at, updated_at FROM team_notification_channels WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	var c models.TeamNotificationChannel
	var configStr, eventsStr string
	if err := row.Scan(&c.ID, &c.TeamID, &c.Provider, &configStr, &eventsStr, &c.IsEnabled, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	c.Config = []byte(configStr)
	c.Events = []byte(eventsStr)
	return &c, nil
}

func (r *NotificationSQLiteRepository) SaveChannel(ctx context.Context, c *models.TeamNotificationChannel) error {
	now := time.Now().UTC()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	query := `INSERT INTO team_notification_channels (
		id, team_id, provider, config, events, is_enabled, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		provider = excluded.provider,
		config = excluded.config,
		events = excluded.events,
		is_enabled = excluded.is_enabled,
		updated_at = excluded.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.TeamID, c.Provider, string(c.Config), string(c.Events), c.IsEnabled, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save channel: %w", err)
	}
	return nil
}

func (r *NotificationSQLiteRepository) DeleteChannel(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM team_notification_channels WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}
	return nil
}
