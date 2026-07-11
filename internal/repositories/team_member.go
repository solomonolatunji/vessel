package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vessl.dev/vessl/internal/utils"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
)

func (r *TeamSQLiteRepository) AddMember(ctx context.Context, member *models.TeamMember) error {
	if member.ID == "" {
		member.ID = uuid.NewString()
	}
	if member.JoinedAt.IsZero() {
		member.JoinedAt = time.Now().UTC()
	}
	if member.Role == "" {
		member.Role = "Member"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO team_members (id, team_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		member.ID, member.TeamID, member.UserID, member.UserEmail, member.Role, member.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("add team member: %w", err)
	}
	return nil
}

func (r *TeamSQLiteRepository) RemoveMember(ctx context.Context, teamID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `DELETE FROM team_members WHERE team_id = ? AND user_id = ? AND role != 'Owner'`, teamID, userID)
	if err != nil {
		return fmt.Errorf("remove team member: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("member not found or cannot remove team Owner")
	}
	return nil
}

func (r *TeamSQLiteRepository) GetMember(ctx context.Context, teamID, userID string) (*models.TeamMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var m models.TeamMember
	var joinedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, team_id, user_id, user_email, role, joined_at FROM team_members WHERE team_id = ? AND user_id = ?`, teamID, userID).
		Scan(&m.ID, &m.TeamID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get team member: %w", err)
	}
	m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
	return &m, nil
}

func (r *TeamSQLiteRepository) ListMembers(ctx context.Context, teamID string) ([]*models.TeamMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx, `SELECT id, team_id, user_id, user_email, role, joined_at FROM team_members WHERE team_id = ? ORDER BY joined_at ASC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("list team members: %w", err)
	}
	defer rows.Close()
	var list []*models.TeamMember
	for rows.Next() {
		var m models.TeamMember
		var joinedStr string
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr); err != nil {
			return nil, err
		}
		m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
		list = append(list, &m)
	}
	return list, nil
}

func (r *TeamSQLiteRepository) CreateInvite(ctx context.Context, invite *models.TeamInvite) error {
	if invite.ID == "" {
		invite.ID = uuid.NewString()
	}
	if invite.Token == "" {
		invite.Token = uuid.NewString()
	}
	now := time.Now().UTC()
	if invite.CreatedAt.IsZero() {
		invite.CreatedAt = now
	}
	if invite.ExpiresAt.IsZero() {
		invite.ExpiresAt = now.Add(7 * 24 * time.Hour)
	}
	if invite.Role == "" {
		invite.Role = "Member"
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO team_invites (id, team_id, email, role, token, invited_by, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		invite.ID, invite.TeamID, invite.Email, invite.Role, invite.Token, invite.InvitedBy, invite.ExpiresAt.Format(time.RFC3339), invite.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create team invite: %w", err)
	}
	return nil
}

func (r *TeamSQLiteRepository) GetInviteByToken(ctx context.Context, token string) (*models.TeamInvite, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var inv models.TeamInvite
	var expStr, createdStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, team_id, email, role, token, invited_by, expires_at, created_at FROM team_invites WHERE token = ?`, token).
		Scan(&inv.ID, &inv.TeamID, &inv.Email, &inv.Role, &inv.Token, &inv.InvitedBy, &expStr, &createdStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Invite", token)
	}
	if err != nil {
		return nil, fmt.Errorf("get team invite: %w", err)
	}
	inv.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
	inv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	return &inv, nil
}

func (r *TeamSQLiteRepository) DeleteInvite(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `DELETE FROM team_invites WHERE id = ?`, id)
	return err
}
