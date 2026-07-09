package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateTeam inserts a new workspace team and automatically assigns the creator as Owner.
func (s *Store) CreateTeam(team *types.Team) error {
	if team.ID == "" {
		team.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	if team.CreatedAt.IsZero() {
		team.CreatedAt = now
	}
	team.UpdatedAt = team.CreatedAt

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO teams (id, name, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		team.ID, team.Name, team.OwnerID, team.CreatedAt.Format(time.RFC3339), team.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to insert team: %w", err)
	}

	ownerMember := &types.TeamMember{
		ID:        uuid.New().String(),
		TeamID:    team.ID,
		UserID:    team.OwnerID,
		UserEmail: "",
		Role:      "Owner",
		JoinedAt:  now,
	}
	_, err = tx.Exec(`INSERT INTO team_members (id, team_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		ownerMember.ID, ownerMember.TeamID, ownerMember.UserID, ownerMember.UserEmail, ownerMember.Role, ownerMember.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to insert team owner membership: %w", err)
	}

	return tx.Commit()
}

// GetTeam retrieves a team by ID.
func (s *Store) GetTeam(id string) (*types.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var t types.Team
	var createdStr, updatedStr string
	err := s.db.QueryRow(`SELECT id, name, owner_id, created_at, updated_at FROM teams WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.OwnerID, &createdStr, &updatedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team %s: %w", id, err)
	}
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &t, nil
}

// ListUserTeams returns all teams where the user is a member or owner.
func (s *Store) ListUserTeams(userID string) ([]*types.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT t.id, t.name, t.owner_id, t.created_at, t.updated_at
	          FROM teams t
	          JOIN team_members m ON t.id = m.team_id
	          WHERE m.user_id = ? ORDER BY t.created_at DESC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user teams: %w", err)
	}
	defer rows.Close()

	var list []*types.Team
	for rows.Next() {
		var t types.Team
		var createdStr, updatedStr string
		if err := rows.Scan(&t.ID, &t.Name, &t.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &t)
	}
	return list, nil
}

// DeleteTeam removes a team and all its memberships if requested by the owner.
func (s *Store) DeleteTeam(id, ownerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM teams WHERE id = ? AND owner_id = ?`, id, ownerID)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("team not found or unauthorized (must be owner)")
	}
	_, _ = s.db.Exec(`DELETE FROM team_members WHERE team_id = ?`, id)
	_, _ = s.db.Exec(`DELETE FROM team_invites WHERE team_id = ?`, id)
	return nil
}

// AddTeamMember adds a user to a team.
func (s *Store) AddTeamMember(member *types.TeamMember) error {
	if member.ID == "" {
		member.ID = uuid.New().String()
	}
	if member.JoinedAt.IsZero() {
		member.JoinedAt = time.Now().UTC()
	}
	if member.Role == "" {
		member.Role = "Member"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`INSERT INTO team_members (id, team_id, user_id, user_email, role, joined_at) VALUES (?, ?, ?, ?, ?, ?)`,
		member.ID, member.TeamID, member.UserID, member.UserEmail, member.Role, member.JoinedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}
	return nil
}

// ListTeamMembers returns all members of a team.
func (s *Store) ListTeamMembers(teamID string) ([]*types.TeamMember, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`SELECT id, team_id, user_id, user_email, role, joined_at FROM team_members WHERE team_id = ? ORDER BY joined_at ASC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	defer rows.Close()

	var list []*types.TeamMember
	for rows.Next() {
		var m types.TeamMember
		var joinedStr string
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr); err != nil {
			return nil, err
		}
		m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
		list = append(list, &m)
	}
	return list, nil
}

// GetTeamMember retrieves a user's membership in a team.
func (s *Store) GetTeamMember(teamID, userID string) (*types.TeamMember, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var m types.TeamMember
	var joinedStr string
	err := s.db.QueryRow(`SELECT id, team_id, user_id, user_email, role, joined_at FROM team_members WHERE team_id = ? AND user_id = ?`, teamID, userID).
		Scan(&m.ID, &m.TeamID, &m.UserID, &m.UserEmail, &m.Role, &joinedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team member: %w", err)
	}
	m.JoinedAt, _ = time.Parse(time.RFC3339, joinedStr)
	return &m, nil
}

// RemoveTeamMember removes a user from a team.
func (s *Store) RemoveTeamMember(teamID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM team_members WHERE team_id = ? AND user_id = ? AND role != 'Owner'`, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("member not found or cannot remove team Owner")
	}
	return nil
}

// CreateTeamInvite stores a new team invitation.
func (s *Store) CreateTeamInvite(invite *types.TeamInvite) error {
	if invite.ID == "" {
		invite.ID = uuid.New().String()
	}
	if invite.Token == "" {
		invite.Token = uuid.New().String()
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

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`INSERT INTO team_invites (id, team_id, email, role, token, invited_by, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		invite.ID, invite.TeamID, invite.Email, invite.Role, invite.Token, invite.InvitedBy, invite.ExpiresAt.Format(time.RFC3339), invite.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create team invite: %w", err)
	}
	return nil
}

// GetTeamInviteByToken retrieves a team invite by its token string.
func (s *Store) GetTeamInviteByToken(token string) (*types.TeamInvite, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var inv types.TeamInvite
	var expStr, createdStr string
	err := s.db.QueryRow(`SELECT id, team_id, email, role, token, invited_by, expires_at, created_at FROM team_invites WHERE token = ?`, token).
		Scan(&inv.ID, &inv.TeamID, &inv.Email, &inv.Role, &inv.Token, &inv.InvitedBy, &expStr, &createdStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team invite: %w", err)
	}
	inv.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
	inv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	return &inv, nil
}

// ListTeamInvites returns all pending invitations for a team.
func (s *Store) ListTeamInvites(teamID string) ([]*types.TeamInvite, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`SELECT id, team_id, email, role, token, invited_by, expires_at, created_at FROM team_invites WHERE team_id = ? ORDER BY created_at DESC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team invites: %w", err)
	}
	defer rows.Close()

	var list []*types.TeamInvite
	for rows.Next() {
		var inv types.TeamInvite
		var expStr, createdStr string
		if err := rows.Scan(&inv.ID, &inv.TeamID, &inv.Email, &inv.Role, &inv.Token, &inv.InvitedBy, &expStr, &createdStr); err != nil {
			return nil, err
		}
		inv.ExpiresAt, _ = time.Parse(time.RFC3339, expStr)
		inv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		list = append(list, &inv)
	}
	return list, nil
}

// DeleteTeamInvite removes an invite.
func (s *Store) DeleteTeamInvite(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`DELETE FROM team_invites WHERE id = ?`, id)
	return err
}
