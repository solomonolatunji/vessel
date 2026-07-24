package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type GitAppRepository interface {
	ListGithubApps(ctx context.Context) ([]models.GithubApp, error)
	GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error)
	SaveGithubApp(ctx context.Context, app *models.GithubApp) error
	DeleteGithubApp(ctx context.Context, id string) error
}

type GitAppRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewGitAppRepo(db *sql.DB, vault Vault) *GitAppRepo {
	return &GitAppRepo{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}

func saveApp(ctx context.Context, db *sqlx.DB, tableName string, columns []string, values []any) error {
	if tableName != "github_apps" {
		return errors.New("invalid table name")
	}
	placeholders := make([]string, len(columns))
	updates := make([]string, len(columns))
	for i, col := range columns {
		placeholders[i] = "?"
		if col != "id" && col != "created_at" {
			updates[i] = fmt.Sprintf("%s=excluded.%s", col, col)
		}
	}

	updateCols := []string{}
	for _, up := range updates {
		if up != "" {
			updateCols = append(updateCols, up)
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (%s)
		ON CONFLICT(id) DO UPDATE SET
			%s,
			updated_at=CURRENT_TIMESTAMP
	`, tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "), strings.Join(updateCols, ",\n\t\t\t"))

	_, err := db.ExecContext(ctx, query, values...)
	return err
}

func deleteApp(ctx context.Context, db *sqlx.DB, tableName, id string) error {
	if tableName != "github_apps" {
		return errors.New("invalid table name")
	}
	_, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), id)
	return err
}

func (r *GitAppRepo) ListGithubApps(ctx context.Context) ([]models.GithubApp, error) {
	query := `SELECT id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at FROM github_apps`
	var apps []models.GithubApp
	if err := r.db.SelectContext(ctx, &apps, query); err != nil {
		return nil, err
	}
	if apps == nil {
		apps = make([]models.GithubApp, 0)
	}
	for i := range apps {
		if apps[i].ClientSecret != "" {
			apps[i].ClientSecret = "********"
		}
		if apps[i].WebhookSecret != "" {
			apps[i].WebhookSecret = "********"
		}
		if apps[i].PrivateKey != "" {
			apps[i].PrivateKey = "********"
		}
	}
	return apps, nil
}

func (r *GitAppRepo) GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error) {
	query := `SELECT id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at FROM github_apps WHERE id = ?`
	var a models.GithubApp
	if err := r.db.GetContext(ctx, &a, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("GithubApp", id)
		}
		return nil, err
	}
	if cs, err := r.vault.Decrypt(a.ClientSecret); err == nil && cs != "" {
		a.ClientSecret = cs
	}
	if ws, err := r.vault.Decrypt(a.WebhookSecret); err == nil && ws != "" {
		a.WebhookSecret = ws
	}
	if pk, err := r.vault.Decrypt(a.PrivateKey); err == nil && pk != "" {
		a.PrivateKey = pk
	}
	return &a, nil
}

func (r *GitAppRepo) SaveGithubApp(ctx context.Context, app *models.GithubApp) error {
	if app.ID != "" {
		existing, err := r.GetGithubApp(ctx, app.ID)
		if err == nil && existing != nil {
			if app.ClientSecret == "" || app.ClientSecret == "********" {
				app.ClientSecret = existing.ClientSecret
			}
			if app.WebhookSecret == "" || app.WebhookSecret == "********" {
				app.WebhookSecret = existing.WebhookSecret
			}
			if app.PrivateKey == "" || app.PrivateKey == "********" {
				app.PrivateKey = existing.PrivateKey
			}
		}
	}

	cs, _ := r.vault.Encrypt(app.ClientSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)
	pk, _ := r.vault.Encrypt(app.PrivateKey)
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	cols := []string{"id", "name", "app_id", "installation_id", "client_id", "client_secret", "webhook_secret", "private_key", "is_public", "created_at", "updated_at"}
	vals := []any{app.ID, app.Name, app.AppID, app.InstallationID, app.ClientID, cs, ws, pk, app.IsPublic, app.CreatedAt, app.UpdatedAt}
	return saveApp(ctx, r.db, "github_apps", cols, vals)
}

func (r *GitAppRepo) DeleteGithubApp(ctx context.Context, id string) error {
	return deleteApp(ctx, r.db, "github_apps", id)
}
