package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByID(ctx context.Context, id string) (*models.Team, error)
	ListTeamsByUser(ctx context.Context, userID string) ([]*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id, ownerID string) error
	AddMember(ctx context.Context, member *models.TeamMember) error
	RemoveMember(ctx context.Context, teamID, userID string) error
	GetMember(ctx context.Context, teamID, userID string) (*models.TeamMember, error)
	ListMembers(ctx context.Context, teamID string) ([]*models.TeamMember, error)
	CreateInvite(ctx context.Context, invite *models.TeamInvite) error
	GetInviteByToken(ctx context.Context, token string) (*models.TeamInvite, error)
	DeleteInvite(ctx context.Context, id string) error
}

type TeamSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewTeamSQLiteRepository(db *sql.DB) *TeamSQLiteRepository {
	return &TeamSQLiteRepository{db: db}
}

func (r *TeamSQLiteRepository) Migrate(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar_url TEXT NOT NULL DEFAULT '',
			preferred_region TEXT NOT NULL DEFAULT 'local',
			owner_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			user_id TEXT NOT NULL,
			user_email TEXT NOT NULL DEFAULT '',
			role TEXT NOT NULL DEFAULT 'Member',
			joined_at TEXT NOT NULL,
			UNIQUE(team_id, user_id)
		);
		CREATE TABLE IF NOT EXISTS team_invites (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
			email TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'Member',
			token TEXT NOT NULL UNIQUE,
			invited_by TEXT NOT NULL,
			expires_at TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
	`)
	return err
}

func (r *TeamSQLiteRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	if team.ID == "" {
		team.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if team.CreatedAt.IsZero() {
		team.CreatedAt = now
	}
	team.UpdatedAt = team.CreatedAt
	r.mu.Lock()
	defer r.mu.Unlock()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO teams (id, name, avatar_url, preferred_region, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		team.ID, team.Name, team.AvatarURL, team.PreferredRegion, team.OwnerID, team.CreatedAt.Format(time.RFC3339), team.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}
	ownerMember := &models.TeamMember{
		ID:        uuid.NewString(),
		TeamID:    team.ID,
		UserID:    team.OwnerID,
		UserEmail: "",
		Role:      "Owner",
		JoinedAt:  now,
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO team_members (id, team_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		ownerMember.ID, ownerMember.TeamID, ownerMember.UserID, ownerMember.UserEmail, ownerMember.Role, ownerMember.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("insert team owner: %w", err)
	}
	return tx.Commit()
}

func (r *TeamSQLiteRepository) GetTeamByID(ctx context.Context, id string) (*models.Team, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var t models.Team
	var createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM teams WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.AvatarURL, &t.PreferredRegion, &t.OwnerID, &createdStr, &updatedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Team", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get team %s: %w", id, err)
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &t, nil
}

func (r *TeamSQLiteRepository) ListTeamsByUser(ctx context.Context, userID string) ([]*models.Team, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := `SELECT t.id, t.name, t.avatar_url, t.preferred_region, t.owner_id, t.created_at, t.updated_at
	          FROM teams t
	          JOIN team_members m ON t.id = m.team_id
	          WHERE m.user_id = ? ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list teams by user: %w", err)
	}
	defer rows.Close()
	var list []*models.Team
	for rows.Next() {
		var t models.Team
		var createdStr, updatedStr string
		if err := rows.Scan(&t.ID, &t.Name, &t.AvatarURL, &t.PreferredRegion, &t.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &t)
	}
	return list, nil
}

func (r *TeamSQLiteRepository) UpdateTeam(ctx context.Context, team *models.Team) error {
	team.UpdatedAt = time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE teams SET name = ?, avatar_url = ?, preferred_region = ?, updated_at = ? WHERE id = ?`,
		team.Name, team.AvatarURL, team.PreferredRegion, team.UpdatedAt.Format(time.RFC3339), team.ID)
	return err
}

func (r *TeamSQLiteRepository) DeleteTeam(ctx context.Context, id, ownerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM teams WHERE id = ? AND owner_id = ?`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete team: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("team not found or unauthorized (must be owner)")
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM team_members WHERE team_id = ?`, id)
	_, _ = r.db.ExecContext(ctx, `DELETE FROM team_invites WHERE team_id = ?`, id)
	return nil
}
