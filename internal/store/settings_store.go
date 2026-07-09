package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
)

// GetServerSettings retrieves the global server configurations or initializes default singleton settings.
func (s *Store) GetServerSettings() (*types.ServerSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var cfg types.ServerSettings
	err := s.db.QueryRow(`SELECT id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, COALESCE(smtp_from_name, ''), COALESCE(smtp_from_address, ''), notification_alerts,
	                             registration_enabled, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at
	                      FROM server_settings WHERE id = 'global'`).
		Scan(&cfg.ID, &cfg.CaddyWildcardIP, &cfg.DiscordWebhookURL, &cfg.SlackWebhookURL, &cfg.TelegramBotToken, &cfg.TelegramChatID, &cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.NotificationAlerts,
			&cfg.RegistrationEnabled, &cfg.CustomDNSResolvers, &cfg.DNSValidationEnabled, &cfg.IPAllowlist, &cfg.MCPServerEnabled, &cfg.UpdateCheckCron, &cfg.AutoUpdateEnabled, &cfg.CurrentVersion, &cfg.LatestVersion, &cfg.LastUpdateCheck, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &types.ServerSettings{
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
		query := `INSERT INTO server_settings (id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts,
		                                       registration_enabled, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, _ = s.db.Exec(query, defaultSettings.ID, defaultSettings.CaddyWildcardIP, defaultSettings.DiscordWebhookURL, defaultSettings.SlackWebhookURL, defaultSettings.TelegramBotToken, defaultSettings.TelegramChatID, defaultSettings.SMTPHost, defaultSettings.SMTPPort, defaultSettings.SMTPUser, defaultSettings.SMTPPassword, defaultSettings.SMTPFromName, defaultSettings.SMTPFromAddress, defaultSettings.NotificationAlerts,
			defaultSettings.RegistrationEnabled, defaultSettings.CustomDNSResolvers, defaultSettings.DNSValidationEnabled, defaultSettings.IPAllowlist, defaultSettings.MCPServerEnabled, defaultSettings.UpdateCheckCron, defaultSettings.AutoUpdateEnabled, defaultSettings.CurrentVersion, defaultSettings.LatestVersion, defaultSettings.LastUpdateCheck, defaultSettings.UpdatedAt)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server settings: %w", err)
	}
	return &cfg, nil
}

// UpdateServerSettings saves changes to the global daemon settings.
func (s *Store) UpdateServerSettings(cfg *types.ServerSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	s.mu.Lock()
	defer s.mu.Unlock()

	query := `INSERT INTO server_settings (id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts,
	                                       registration_enabled, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
	          custom_dns_resolvers = excluded.custom_dns_resolvers,
	          dns_validation_enabled = excluded.dns_validation_enabled,
	          ip_allowlist = excluded.ip_allowlist,
	          mcp_server_enabled = excluded.mcp_server_enabled,
	          update_check_cron = excluded.update_check_cron,
	          auto_update_enabled = excluded.auto_update_enabled,
	          current_version = excluded.current_version,
	          latest_version = excluded.latest_version,
	          last_update_check = excluded.last_update_check,
	          updated_at = excluded.updated_at`
	_, err := s.db.Exec(query, cfg.ID, cfg.CaddyWildcardIP, cfg.DiscordWebhookURL, cfg.SlackWebhookURL, cfg.TelegramBotToken, cfg.TelegramChatID, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.NotificationAlerts,
		cfg.RegistrationEnabled, cfg.CustomDNSResolvers, cfg.DNSValidationEnabled, cfg.IPAllowlist, cfg.MCPServerEnabled, cfg.UpdateCheckCron, cfg.AutoUpdateEnabled, cfg.CurrentVersion, cfg.LatestVersion, cfg.LastUpdateCheck, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update server settings: %w", err)
	}
	return nil
}

// CreatePersonalAccessToken inserts a new CLI/API access token for a user.
func (s *Store) CreatePersonalAccessToken(pat *types.PersonalAccessToken) error {
	if pat.ID == "" {
		pat.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	if pat.CreatedAt.IsZero() {
		pat.CreatedAt = now
	}
	if pat.ExpiresAt.IsZero() {
		pat.ExpiresAt = now.Add(365 * 24 * time.Hour)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`INSERT INTO personal_access_tokens (id, user_id, name, token_hash, prefix, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		pat.ID, pat.UserID, pat.Name, pat.TokenHash, pat.Prefix, pat.ExpiresAt.Format(time.RFC3339), pat.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create personal access token: %w", err)
	}
	return nil
}

// ListPersonalAccessTokens returns all PATs belonging to a user without exposing raw token hashes.
func (s *Store) ListPersonalAccessTokens(userID string) ([]*types.PersonalAccessToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`SELECT id, user_id, name, prefix, expires_at, created_at FROM personal_access_tokens WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list personal access tokens: %w", err)
	}
	defer rows.Close()

	var list []*types.PersonalAccessToken
	for rows.Next() {
		var pat types.PersonalAccessToken
		var expStr, createdStr string
		if err := rows.Scan(&pat.ID, &pat.UserID, &pat.Name, &pat.Prefix, &expStr, &createdStr); err != nil {
			return nil, err
		}
		pat.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
		pat.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		list = append(list, &pat)
	}
	return list, nil
}

// DeletePersonalAccessToken revokes a user's access token.
func (s *Store) DeletePersonalAccessToken(id, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM personal_access_tokens WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete personal access token: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("access token not found or unauthorized")
	}
	return nil
}
