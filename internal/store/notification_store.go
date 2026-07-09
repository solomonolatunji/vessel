package store

import (
	"database/sql"
	"fmt"
	"time"

	"vessel.dev/vessel/internal/types"
)

func (s *Store) GetNotificationIntegration() (*types.NotificationIntegration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, smtp_enabled, COALESCE(smtp_host, ''), COALESCE(smtp_port, 587), COALESCE(smtp_user, ''), COALESCE(smtp_password, ''), COALESCE(smtp_from_name, ''), COALESCE(smtp_from_address, ''), resend_enabled, COALESCE(resend_api_key, ''), slack_enabled, COALESCE(slack_webhook_url, ''), discord_enabled, COALESCE(discord_webhook_url, ''), discord_ping_enabled, telegram_enabled, COALESCE(telegram_bot_token, ''), COALESCE(telegram_chat_id, ''), pushover_enabled, COALESCE(pushover_user_key, ''), COALESCE(pushover_api_token, ''), webhook_enabled, COALESCE(webhook_url, ''), COALESCE(updated_at, '') FROM notification_integrations WHERE id = 'global'`

	row := s.db.QueryRow(query)
	var n types.NotificationIntegration
	var smtpHost, smtpUser, smtpPassword, smtpFromName, smtpFromAddress, resendKey, slackUrl, discordUrl, telegramBot, telegramChat, pushoverUser, pushoverToken, webhookUrl, updatedAt string
	var smtpPort int

	err := row.Scan(
		&n.ID, &n.SMTPEnabled, &smtpHost, &smtpPort, &smtpUser, &smtpPassword, &smtpFromName, &smtpFromAddress,
		&n.ResendEnabled, &resendKey, &n.SlackEnabled, &slackUrl,
		&n.DiscordEnabled, &discordUrl, &n.DiscordPingEnabled,
		&n.TelegramEnabled, &telegramBot, &telegramChat,
		&n.PushoverEnabled, &pushoverUser, &pushoverToken,
		&n.WebhookEnabled, &webhookUrl, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return &types.NotificationIntegration{ID: "global"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan notification integration: %w", err)
	}

	n.SMTPHost = smtpHost
	n.SMTPPort = smtpPort
	n.SMTPUser = smtpUser
	n.SMTPPassword = smtpPassword
	n.SMTPFromName = smtpFromName
	n.SMTPFromAddress = smtpFromAddress
	n.ResendAPIKey = resendKey
	n.SlackWebhookURL = slackUrl
	n.DiscordWebhookURL = discordUrl
	n.TelegramBotToken = telegramBot
	n.TelegramChatID = telegramChat
	n.PushoverUserKey = pushoverUser
	n.PushoverAPIToken = pushoverToken
	n.WebhookURL = webhookUrl
	n.UpdatedAt = updatedAt

	return &n, nil
}

func (s *Store) SaveNotificationIntegration(n *types.NotificationIntegration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	n.ID = "global"
	n.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	query := `INSERT INTO notification_integrations (
		id, smtp_enabled, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address,
		resend_enabled, resend_api_key, slack_enabled, slack_webhook_url,
		discord_enabled, discord_webhook_url, discord_ping_enabled,
		telegram_enabled, telegram_bot_token, telegram_chat_id,
		pushover_enabled, pushover_user_key, pushover_api_token,
		webhook_enabled, webhook_url, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		smtp_enabled = excluded.smtp_enabled,
		smtp_host = excluded.smtp_host,
		smtp_port = excluded.smtp_port,
		smtp_user = excluded.smtp_user,
		smtp_password = excluded.smtp_password,
		smtp_from_name = excluded.smtp_from_name,
		smtp_from_address = excluded.smtp_from_address,
		resend_enabled = excluded.resend_enabled,
		resend_api_key = excluded.resend_api_key,
		slack_enabled = excluded.slack_enabled,
		slack_webhook_url = excluded.slack_webhook_url,
		discord_enabled = excluded.discord_enabled,
		discord_webhook_url = excluded.discord_webhook_url,
		discord_ping_enabled = excluded.discord_ping_enabled,
		telegram_enabled = excluded.telegram_enabled,
		telegram_bot_token = excluded.telegram_bot_token,
		telegram_chat_id = excluded.telegram_chat_id,
		pushover_enabled = excluded.pushover_enabled,
		pushover_user_key = excluded.pushover_user_key,
		pushover_api_token = excluded.pushover_api_token,
		webhook_enabled = excluded.webhook_enabled,
		webhook_url = excluded.webhook_url,
		updated_at = excluded.updated_at`

	_, err := s.db.Exec(query,
		n.ID, n.SMTPEnabled, n.SMTPHost, n.SMTPPort, n.SMTPUser, n.SMTPPassword, n.SMTPFromName, n.SMTPFromAddress,
		n.ResendEnabled, n.ResendAPIKey, n.SlackEnabled, n.SlackWebhookURL,
		n.DiscordEnabled, n.DiscordWebhookURL, n.DiscordPingEnabled,
		n.TelegramEnabled, n.TelegramBotToken, n.TelegramChatID,
		n.PushoverEnabled, n.PushoverUserKey, n.PushoverAPIToken,
		n.WebhookEnabled, n.WebhookURL, n.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save notification integration: %w", err)
	}

	return nil
}

func (s *Store) GetProjectNotificationPref(projectID string) (*types.ProjectNotificationPref, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT project_id, email_enabled, slack_enabled, discord_enabled, telegram_enabled, pushover_enabled, webhook_enabled, COALESCE(events, 'deploy.success,deploy.failure,invite'), updated_at FROM project_notification_prefs WHERE project_id = ?`

	row := s.db.QueryRow(query, projectID)
	var pref types.ProjectNotificationPref
	err := row.Scan(&pref.ProjectID, &pref.EmailEnabled, &pref.SlackEnabled, &pref.DiscordEnabled, &pref.TelegramEnabled, &pref.PushoverEnabled, &pref.WebhookEnabled, &pref.Events, &pref.UpdatedAt)
	if err == sql.ErrNoRows {
		return &types.ProjectNotificationPref{
			ProjectID:       projectID,
			EmailEnabled:    true,
			SlackEnabled:    true,
			DiscordEnabled:  true,
			TelegramEnabled: true,
			PushoverEnabled: true,
			WebhookEnabled:  true,
			Events:          "deploy.success,deploy.failure,invite",
			UpdatedAt:       time.Now().UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan project notification preferences: %w", err)
	}

	return &pref, nil
}

func (s *Store) SaveProjectNotificationPref(pref *types.ProjectNotificationPref) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pref.UpdatedAt = time.Now().UTC()

	query := `INSERT INTO project_notification_prefs (
		project_id, email_enabled, slack_enabled, discord_enabled, telegram_enabled, pushover_enabled, webhook_enabled, events, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(project_id) DO UPDATE SET
		email_enabled = excluded.email_enabled,
		slack_enabled = excluded.slack_enabled,
		discord_enabled = excluded.discord_enabled,
		telegram_enabled = excluded.telegram_enabled,
		pushover_enabled = excluded.pushover_enabled,
		webhook_enabled = excluded.webhook_enabled,
		events = excluded.events,
		updated_at = excluded.updated_at`

	_, err := s.db.Exec(query,
		pref.ProjectID, pref.EmailEnabled, pref.SlackEnabled, pref.DiscordEnabled,
		pref.TelegramEnabled, pref.PushoverEnabled, pref.WebhookEnabled, pref.Events, pref.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save project notification preferences: %w", err)
	}

	return nil
}
