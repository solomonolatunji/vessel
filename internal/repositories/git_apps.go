package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type GitAppRepository interface {
	ListGithubApps(ctx context.Context, workspaceID string) ([]models.GithubApp, error)
	GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error)
	SaveGithubApp(ctx context.Context, app *models.GithubApp) error
	DeleteGithubApp(ctx context.Context, id string) error

	ListGitlabApps(ctx context.Context, workspaceID string) ([]models.GitlabApp, error)
	GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error)
	SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error
	DeleteGitlabApp(ctx context.Context, id string) error

	ListBitbucketApps(ctx context.Context, workspaceID string) ([]models.BitbucketApp, error)
	GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error)
	SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error
	DeleteBitbucketApp(ctx context.Context, id string) error
}

type GitAppSQLiteRepository struct {
	db    *sqlx.DB
	vault Vault
}

func NewGitAppSQLiteRepository(db *sql.DB, vault Vault) *GitAppSQLiteRepository {
	return &GitAppSQLiteRepository{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}

func saveApp(ctx context.Context, db *sqlx.DB, tableName string, columns []string, values []any) error {
	placeholders := make([]string, len(columns))
	updates := make([]string, len(columns))
	for i, col := range columns {
		placeholders[i] = "?"
		if col != "id" && col != "workspace_id" && col != "created_at" {
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
	_, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), id)
	return err
}

func (r *GitAppSQLiteRepository) ListGithubApps(ctx context.Context, workspaceID string) ([]models.GithubApp, error) {
	query := `SELECT id, workspace_id, name, app_id, installation_id, client_id, is_public, created_at, updated_at FROM github_apps WHERE workspace_id = ?`
	var apps []models.GithubApp
	if err := r.db.SelectContext(ctx, &apps, query, workspaceID); err != nil {
		return nil, err
	}
	if apps == nil {
		apps = make([]models.GithubApp, 0)
	}
	return apps, nil
}

func (r *GitAppSQLiteRepository) GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error) {
	query := `SELECT id, workspace_id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at FROM github_apps WHERE id = ?`
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

func (r *GitAppSQLiteRepository) SaveGithubApp(ctx context.Context, app *models.GithubApp) error {
	cs, _ := r.vault.Encrypt(app.ClientSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)
	pk, _ := r.vault.Encrypt(app.PrivateKey)
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	cols := []string{"id", "workspace_id", "name", "app_id", "installation_id", "client_id", "client_secret", "webhook_secret", "private_key", "is_public", "created_at", "updated_at"}
	vals := []any{app.ID, app.WorkspaceID, app.Name, app.AppID, app.InstallationID, app.ClientID, cs, ws, pk, app.IsPublic, app.CreatedAt, app.UpdatedAt}
	return saveApp(ctx, r.db, "github_apps", cols, vals)
}

func (r *GitAppSQLiteRepository) DeleteGithubApp(ctx context.Context, id string) error {
	return deleteApp(ctx, r.db, "github_apps", id)
}

func (r *GitAppSQLiteRepository) ListGitlabApps(ctx context.Context, workspaceID string) ([]models.GitlabApp, error) {
	query := `SELECT id, workspace_id, name, app_id, api_url, is_public, created_at, updated_at FROM gitlab_apps WHERE workspace_id = ?`
	var apps []models.GitlabApp
	if err := r.db.SelectContext(ctx, &apps, query, workspaceID); err != nil {
		return nil, err
	}
	if apps == nil {
		apps = make([]models.GitlabApp, 0)
	}
	return apps, nil
}

func (r *GitAppSQLiteRepository) GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error) {
	query := `SELECT id, workspace_id, name, app_id, app_secret, webhook_secret, api_url, is_public, created_at, updated_at FROM gitlab_apps WHERE id = ?`
	var a models.GitlabApp
	if err := r.db.GetContext(ctx, &a, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("GitlabApp", id)
		}
		return nil, err
	}
	if as, err := r.vault.Decrypt(a.AppSecret); err == nil && as != "" {
		a.AppSecret = as
	}
	if ws, err := r.vault.Decrypt(a.WebhookSecret); err == nil && ws != "" {
		a.WebhookSecret = ws
	}
	return &a, nil
}

func (r *GitAppSQLiteRepository) SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error {
	as, _ := r.vault.Encrypt(app.AppSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	cols := []string{"id", "workspace_id", "name", "app_id", "app_secret", "webhook_secret", "api_url", "is_public", "created_at", "updated_at"}
	vals := []any{app.ID, app.WorkspaceID, app.Name, app.AppID, as, ws, app.APIURL, app.IsPublic, app.CreatedAt, app.UpdatedAt}
	return saveApp(ctx, r.db, "gitlab_apps", cols, vals)
}

func (r *GitAppSQLiteRepository) DeleteGitlabApp(ctx context.Context, id string) error {
	return deleteApp(ctx, r.db, "gitlab_apps", id)
}

func (r *GitAppSQLiteRepository) ListBitbucketApps(ctx context.Context, workspaceID string) ([]models.BitbucketApp, error) {
	query := `SELECT id, workspace_id, name, workspace, client_id, is_public, created_at, updated_at FROM bitbucket_apps WHERE workspace_id = ?`
	var apps []models.BitbucketApp
	if err := r.db.SelectContext(ctx, &apps, query, workspaceID); err != nil {
		return nil, err
	}
	if apps == nil {
		apps = make([]models.BitbucketApp, 0)
	}
	return apps, nil
}

func (r *GitAppSQLiteRepository) GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error) {
	query := `SELECT id, workspace_id, name, workspace, client_id, client_secret, webhook_secret, is_public, created_at, updated_at FROM bitbucket_apps WHERE id = ?`
	var a models.BitbucketApp
	if err := r.db.GetContext(ctx, &a, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("BitbucketApp", id)
		}
		return nil, err
	}
	if cs, err := r.vault.Decrypt(a.ClientSecret); err == nil && cs != "" {
		a.ClientSecret = cs
	}
	if ws, err := r.vault.Decrypt(a.WebhookSecret); err == nil && ws != "" {
		a.WebhookSecret = ws
	}
	return &a, nil
}

func (r *GitAppSQLiteRepository) SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error {
	cs, _ := r.vault.Encrypt(app.ClientSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)
	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	cols := []string{"id", "workspace_id", "name", "workspace", "client_id", "client_secret", "webhook_secret", "is_public", "created_at", "updated_at"}
	vals := []any{app.ID, app.WorkspaceID, app.Name, app.Workspace, app.ClientID, cs, ws, app.IsPublic, app.CreatedAt, app.UpdatedAt}
	return saveApp(ctx, r.db, "bitbucket_apps", cols, vals)
}

func (r *GitAppSQLiteRepository) DeleteBitbucketApp(ctx context.Context, id string) error {
	return deleteApp(ctx, r.db, "bitbucket_apps", id)
}
