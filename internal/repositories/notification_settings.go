package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"codedock.dev/codedock/internal/models"
)

type NotificationSettingsRepository interface {
	GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error
}

type NotificationSettingsRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewNotificationSettingsRepo(db *sql.DB) *NotificationSettingsRepo {
	return &NotificationSettingsRepo{db: sqlx.NewDb(db, "sqlite")}
}

const notificationSettingsColumns = `id, discord_webhook_url, discord_ping_enabled, discord_enabled, slack_webhook_url, slack_enabled, telegram_bot_token, telegram_chat_id, telegram_enabled, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, smtp_enabled, resend_api_key, resend_enabled, pushover_user_key, pushover_api_token, pushover_enabled, generic_webhook_url, generic_webhook_enabled, notification_alerts, created_at, updated_at`

func notificationSettingsPlaceholders() string {
	columns := strings.Split(notificationSettingsColumns, ",")
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func scanNotificationSettings(scanner interface{ Scan(dest ...any) error }, cfg *models.NotificationSettings) error {
	return scanner.Scan(
		&cfg.ID, &cfg.DiscordWebhookURL, &cfg.DiscordPingEnabled, &cfg.DiscordEnabled, &cfg.SlackWebhookURL, &cfg.SlackEnabled, &cfg.TelegramBotToken, &cfg.TelegramChatID, &cfg.TelegramEnabled,
		&cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.SMTPEnabled,
		&cfg.ResendAPIKey, &cfg.ResendEnabled, &cfg.PushoverUserKey, &cfg.PushoverAPIToken, &cfg.PushoverEnabled, &cfg.GenericWebhookURL, &cfg.GenericWebhookEnabled,
		&cfg.NotificationAlerts,
		&cfg.CreatedAt, &cfg.UpdatedAt,
	)
}

func notificationSettingsArgs(cfg *models.NotificationSettings) []any {
	return []any{
		cfg.ID, cfg.DiscordWebhookURL, cfg.DiscordPingEnabled, cfg.DiscordEnabled, cfg.SlackWebhookURL, cfg.SlackEnabled, cfg.TelegramBotToken, cfg.TelegramChatID, cfg.TelegramEnabled,
		cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.SMTPEnabled,
		cfg.ResendAPIKey, cfg.ResendEnabled, cfg.PushoverUserKey, cfg.PushoverAPIToken, cfg.PushoverEnabled, cfg.GenericWebhookURL, cfg.GenericWebhookEnabled,
		cfg.NotificationAlerts,
		cfg.CreatedAt, cfg.UpdatedAt,
	}
}

func (r *NotificationSettingsRepo) GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.NotificationSettings
	err := scanNotificationSettings(r.db.QueryRowContext(ctx, `SELECT `+notificationSettingsColumns+` FROM notification_settings WHERE id = 'global' LIMIT 1`), &cfg)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.NotificationSettings{
			ID:                 "global",
			NotificationAlerts: true,
			SMTPPort:           587,
			CreatedAt:          time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:          time.Now().UTC().Format(time.RFC3339),
		}
		query := fmt.Sprintf(`INSERT INTO notification_settings (%s) VALUES (%s)`, notificationSettingsColumns, notificationSettingsPlaceholders())
		_, _ = r.db.ExecContext(ctx, query, notificationSettingsArgs(defaultSettings)...)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification settings: %w", err)
	}
	return &cfg, nil
}

func (r *NotificationSettingsRepo) UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf(`INSERT INTO notification_settings (%s)
	          VALUES (%s)
	          ON CONFLICT(id) DO UPDATE SET
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
	          updated_at = excluded.updated_at`, notificationSettingsColumns, notificationSettingsPlaceholders())
	_, err := r.db.ExecContext(ctx, query, notificationSettingsArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("failed to update notification settings: %w", err)
	}
	return nil
}
