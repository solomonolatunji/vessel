package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/utils"
)

type NotificationRepository interface {
	ListChannelsByTeam(ctx context.Context, teamID string) ([]models.TeamNotificationChannel, error)
	GetChannel(ctx context.Context, id string) (*models.TeamNotificationChannel, error)
	SaveChannel(ctx context.Context, c *models.TeamNotificationChannel) error
	DeleteChannel(ctx context.Context, id string) error
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Channel", id)
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
