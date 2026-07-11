package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/utils"
)

type GitAppRepository interface {
	// Github App
	ListGithubApps(ctx context.Context, teamID string) ([]models.GithubApp, error)
	GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error)
	SaveGithubApp(ctx context.Context, app *models.GithubApp) error
	DeleteGithubApp(ctx context.Context, id string) error

	// Gitlab App
	ListGitlabApps(ctx context.Context, teamID string) ([]models.GitlabApp, error)
	GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error)
	SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error
	DeleteGitlabApp(ctx context.Context, id string) error

	// Bitbucket App
	ListBitbucketApps(ctx context.Context, teamID string) ([]models.BitbucketApp, error)
	GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error)
	SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error
	DeleteBitbucketApp(ctx context.Context, id string) error
}

type GitAppSQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewGitAppSQLiteRepository(db *sql.DB, vault Vault) *GitAppSQLiteRepository {
	return &GitAppSQLiteRepository{db: db, vault: vault}
}

type scanner interface {
	Scan(dest ...any) error
}

func listApps[T any](ctx context.Context, db *sql.DB, query string, teamID string, scanFn func(scanner, *T) error) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []T
	for rows.Next() {
		var a T
		if err := scanFn(rows, &a); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func getApp[T any](ctx context.Context, db *sql.DB, query string, id string, modelName string, scanFn func(scanner, *T) error) (*T, error) {
	row := db.QueryRowContext(ctx, query, id)
	var a T
	if err := scanFn(row, &a); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError(modelName, id)
		}
		return nil, err
	}
	return &a, nil
}

// ---- GitHub Apps ----

func (r *GitAppSQLiteRepository) ListGithubApps(ctx context.Context, teamID string) ([]models.GithubApp, error) {
	query := `SELECT id, team_id, name, app_id, installation_id, client_id, is_public, created_at, updated_at FROM github_apps WHERE team_id = ?`
	return listApps(ctx, r.db, query, teamID, func(s scanner, a *models.GithubApp) error {
		return s.Scan(&a.ID, &a.TeamID, &a.Name, &a.AppID, &a.InstallationID, &a.ClientID, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt)
	})
}

func (r *GitAppSQLiteRepository) GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error) {
	query := `SELECT id, team_id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at FROM github_apps WHERE id = ?`
	return getApp(ctx, r.db, query, id, "GithubApp", func(s scanner, a *models.GithubApp) error {
		var cs, ws, pk string
		if err := s.Scan(&a.ID, &a.TeamID, &a.Name, &a.AppID, &a.InstallationID, &a.ClientID, &cs, &ws, &pk, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return err
		}
		a.ClientSecret, _ = r.vault.Decrypt(cs)
		if a.ClientSecret == "" {
			a.ClientSecret = cs
		}
		a.WebhookSecret, _ = r.vault.Decrypt(ws)
		if a.WebhookSecret == "" {
			a.WebhookSecret = ws
		}
		a.PrivateKey, _ = r.vault.Decrypt(pk)
		if a.PrivateKey == "" {
			a.PrivateKey = pk
		}
		return nil
	})
}

func (r *GitAppSQLiteRepository) SaveGithubApp(ctx context.Context, app *models.GithubApp) error {
	query := `
		INSERT INTO github_apps (id, team_id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name,
			app_id=excluded.app_id,
			installation_id=excluded.installation_id,
			client_id=excluded.client_id,
			client_secret=excluded.client_secret,
			webhook_secret=excluded.webhook_secret,
			private_key=excluded.private_key,
			is_public=excluded.is_public,
			updated_at=CURRENT_TIMESTAMP
	`
	cs, _ := r.vault.Encrypt(app.ClientSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)
	pk, _ := r.vault.Encrypt(app.PrivateKey)

	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, app.ID, app.TeamID, app.Name, app.AppID, app.InstallationID, app.ClientID, cs, ws, pk, app.IsPublic, app.CreatedAt, app.UpdatedAt)
	return err
}

func (r *GitAppSQLiteRepository) DeleteGithubApp(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM github_apps WHERE id = ?`, id)
	return err
}

// ---- GitLab Apps ----

func (r *GitAppSQLiteRepository) ListGitlabApps(ctx context.Context, teamID string) ([]models.GitlabApp, error) {
	query := `SELECT id, team_id, name, app_id, api_url, is_public, created_at, updated_at FROM gitlab_apps WHERE team_id = ?`
	return listApps(ctx, r.db, query, teamID, func(s scanner, a *models.GitlabApp) error {
		return s.Scan(&a.ID, &a.TeamID, &a.Name, &a.AppID, &a.APIURL, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt)
	})
}

func (r *GitAppSQLiteRepository) GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error) {
	query := `SELECT id, team_id, name, app_id, app_secret, webhook_secret, api_url, is_public, created_at, updated_at FROM gitlab_apps WHERE id = ?`
	return getApp(ctx, r.db, query, id, "GitlabApp", func(s scanner, a *models.GitlabApp) error {
		var as, ws string
		if err := s.Scan(&a.ID, &a.TeamID, &a.Name, &a.AppID, &as, &ws, &a.APIURL, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return err
		}
		a.AppSecret, _ = r.vault.Decrypt(as)
		if a.AppSecret == "" {
			a.AppSecret = as
		}
		a.WebhookSecret, _ = r.vault.Decrypt(ws)
		if a.WebhookSecret == "" {
			a.WebhookSecret = ws
		}
		return nil
	})
}

func (r *GitAppSQLiteRepository) SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error {
	query := `
		INSERT INTO gitlab_apps (id, team_id, name, app_id, app_secret, webhook_secret, api_url, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name,
			app_id=excluded.app_id,
			app_secret=excluded.app_secret,
			webhook_secret=excluded.webhook_secret,
			api_url=excluded.api_url,
			is_public=excluded.is_public,
			updated_at=CURRENT_TIMESTAMP
	`
	as, _ := r.vault.Encrypt(app.AppSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)

	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, app.ID, app.TeamID, app.Name, app.AppID, as, ws, app.APIURL, app.IsPublic, app.CreatedAt, app.UpdatedAt)
	return err
}

func (r *GitAppSQLiteRepository) DeleteGitlabApp(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM gitlab_apps WHERE id = ?`, id)
	return err
}

// ---- Bitbucket Apps ----

func (r *GitAppSQLiteRepository) ListBitbucketApps(ctx context.Context, teamID string) ([]models.BitbucketApp, error) {
	query := `SELECT id, team_id, name, workspace, client_id, is_public, created_at, updated_at FROM bitbucket_apps WHERE team_id = ?`
	return listApps(ctx, r.db, query, teamID, func(s scanner, a *models.BitbucketApp) error {
		return s.Scan(&a.ID, &a.TeamID, &a.Name, &a.Workspace, &a.ClientID, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt)
	})
}

func (r *GitAppSQLiteRepository) GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error) {
	query := `SELECT id, team_id, name, workspace, client_id, client_secret, webhook_secret, is_public, created_at, updated_at FROM bitbucket_apps WHERE id = ?`
	return getApp(ctx, r.db, query, id, "BitbucketApp", func(s scanner, a *models.BitbucketApp) error {
		var cs, ws string
		if err := s.Scan(&a.ID, &a.TeamID, &a.Name, &a.Workspace, &a.ClientID, &cs, &ws, &a.IsPublic, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return err
		}
		a.ClientSecret, _ = r.vault.Decrypt(cs)
		if a.ClientSecret == "" {
			a.ClientSecret = cs
		}
		a.WebhookSecret, _ = r.vault.Decrypt(ws)
		if a.WebhookSecret == "" {
			a.WebhookSecret = ws
		}
		return nil
	})
}

func (r *GitAppSQLiteRepository) SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error {
	query := `
		INSERT INTO bitbucket_apps (id, team_id, name, workspace, client_id, client_secret, webhook_secret, is_public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name,
			workspace=excluded.workspace,
			client_id=excluded.client_id,
			client_secret=excluded.client_secret,
			webhook_secret=excluded.webhook_secret,
			is_public=excluded.is_public,
			updated_at=CURRENT_TIMESTAMP
	`
	cs, _ := r.vault.Encrypt(app.ClientSecret)
	ws, _ := r.vault.Encrypt(app.WebhookSecret)

	if app.CreatedAt.IsZero() {
		app.CreatedAt = time.Now()
	}
	app.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, app.ID, app.TeamID, app.Name, app.Workspace, app.ClientID, cs, ws, app.IsPublic, app.CreatedAt, app.UpdatedAt)
	return err
}

func (r *GitAppSQLiteRepository) DeleteBitbucketApp(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bitbucket_apps WHERE id = ?`, id)
	return err
}
