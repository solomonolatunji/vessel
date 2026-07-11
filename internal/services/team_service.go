package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
)

type TeamService struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func NewTeamService(tr repositories.TeamRepository, ur repositories.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: tr,
		userRepo: ur,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name, ownerID string) (*models.Team, error) {
	if name == "" || ownerID == "" {
		return nil, errors.New("team name and ownerId are required")
	}
	team := &models.Team{
		ID:        uuid.New().String(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.teamRepo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}
	member := &models.TeamMember{
		ID:       uuid.New().String(),
		TeamID:   team.ID,
		UserID:   ownerID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}
	_ = s.teamRepo.AddMember(ctx, member)
	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamID string) (*models.Team, error) {
	if teamID == "" {
		return nil, errors.New("teamId is required")
	}
	return s.teamRepo.GetTeamByID(ctx, teamID)
}

func (s *TeamService) ListTeamsByUser(ctx context.Context, userID string) ([]*models.Team, error) {
	if userID == "" {
		return nil, errors.New("userId is required")
	}
	return s.teamRepo.ListTeamsByUser(ctx, userID)
}

func (s *TeamService) UpdateTeam(ctx context.Context, team *models.Team) error {
	if team == nil || team.ID == "" {
		return errors.New("valid team is required for update")
	}
	team.UpdatedAt = time.Now()
	return s.teamRepo.UpdateTeam(ctx, team)
}

func (s *TeamService) DeleteTeam(ctx context.Context, teamID, ownerID string) error {
	if teamID == "" || ownerID == "" {
		return errors.New("teamId and ownerId are required")
	}
	return s.teamRepo.DeleteTeam(ctx, teamID, ownerID)
}

func (s *TeamService) AddMember(ctx context.Context, teamID, userID, role string) error {
	if teamID == "" || userID == "" {
		return errors.New("teamId and userId are required")
	}
	if role == "" {
		role = "member"
	}
	member := &models.TeamMember{
		ID:       uuid.New().String(),
		TeamID:   teamID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return s.teamRepo.AddMember(ctx, member)
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, userID string) error {
	if teamID == "" || userID == "" {
		return errors.New("teamId and userId are required")
	}
	return s.teamRepo.RemoveMember(ctx, teamID, userID)
}

func (s *TeamService) ListMembers(ctx context.Context, teamID string) ([]*models.TeamMember, error) {
	if teamID == "" {
		return nil, errors.New("teamId is required")
	}
	return s.teamRepo.ListMembers(ctx, teamID)
}

func (s *TeamService) InviteMember(ctx context.Context, teamID, email, role string) (*models.TeamInvite, error) {
	if teamID == "" || email == "" {
		return nil, errors.New("teamId and email are required")
	}
	if role == "" {
		role = "member"
	}
	invite := &models.TeamInvite{
		ID:        uuid.New().String(),
		TeamID:    teamID,
		Email:     email,
		Role:      role,
		Token:     uuid.New().String(),
		CreatedAt: time.Now(),
	}
	if err := s.teamRepo.CreateInvite(ctx, invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *TeamService) GetInvite(ctx context.Context, token string) (*models.TeamInvite, error) {
	if token == "" {
		return nil, errors.New("token required")
	}
	return s.teamRepo.GetInviteByToken(ctx, token)
}

func (s *TeamService) AcceptInvite(ctx context.Context, token, userID string) error {
	if token == "" || userID == "" {
		return errors.New("token and userId are required")
	}
	invite, err := s.teamRepo.GetInviteByToken(ctx, token)
	if err != nil || invite == nil {
		return errors.New("invalid or expired invite token")
	}
	if err := s.AddMember(ctx, invite.TeamID, userID, invite.Role); err != nil {
		return err
	}
	return s.teamRepo.DeleteInvite(ctx, invite.ID)
}
