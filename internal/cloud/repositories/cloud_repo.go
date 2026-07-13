package repos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"vessl.dev/vessl/internal/models"
)

type CloudRepo interface {
	GetTeamByID(id uint) (*models.CloudTeam, error)
	GetTeamByName(name string) (*models.CloudTeam, error)
	GetTeamMember(teamID uint, userID string) (*models.WorkspaceMember, error)
	GetTeamOwners(teamID uint) ([]models.WorkspaceMember, error)
	UpdateTeam(team *models.CloudTeam) error
	GetActiveServerCount(teamID uint) (int64, error)
	GetServerByToken(token string) (*models.CloudServer, error)
	GetDeploymentsInLastHour(teamID uint) (int64, error)
	RegisterServer(ctx context.Context, teamID uint, token string, name string, ip string) error
	GetTeamByStripeCustomerID(customerID string) (*models.CloudTeam, error)
	GetTeamByPaddleCustomerID(customerID string) (*models.CloudTeam, error)
	LogUsage(usage *models.CloudUsageLog) error
	LogTelemetry(logEntry *models.CloudTelemetryLog) error
	InsertAuditLog(ctx context.Context, entry *models.AuditLog) error
	GetTotalServers(ctx context.Context) (int64, error)
	GetActiveSubscriptions(ctx context.Context) (int64, error)
	ListAuditLogs(ctx context.Context, limit int, offset int) ([]models.AuditLog, error)
	GetCurrentMonthUsage(teamID uint) (int, int, error)
	GetUserTeams(userID string) ([]models.CloudTeam, error)
	GetTeamServers(teamID uint) ([]models.CloudServer, error)
}

type cloudRepo struct {
	db *gorm.DB
}

func NewCloudRepo(db *gorm.DB) CloudRepo {
	return &cloudRepo{db: db}
}

func (r *cloudRepo) GetTeamByID(id uint) (*models.CloudTeam, error) {
	var team models.CloudTeam
	if err := r.db.First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *cloudRepo) GetTeamByName(name string) (*models.CloudTeam, error) {
	var team models.CloudTeam
	if err := r.db.Where("name = ?", name).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *cloudRepo) GetTeamMember(teamID uint, userID string) (*models.WorkspaceMember, error) {
	var member models.WorkspaceMember
	if err := r.db.Where("team_id = ? AND user_id = ?", fmt.Sprintf("%d", teamID), userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *cloudRepo) GetTeamOwners(teamID uint) ([]models.WorkspaceMember, error) {
	var members []models.WorkspaceMember
	if err := r.db.Where("team_id = ? AND role = ?", fmt.Sprintf("%d", teamID), "owner").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *cloudRepo) GetTeamByStripeCustomerID(customerID string) (*models.CloudTeam, error) {
	var team models.CloudTeam
	if err := r.db.Where("stripe_customer_id = ?", customerID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *cloudRepo) GetTeamByPaddleCustomerID(customerID string) (*models.CloudTeam, error) {
	var team models.CloudTeam
	if err := r.db.Where("paddle_customer_id = ?", customerID).First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *cloudRepo) UpdateTeam(team *models.CloudTeam) error {
	return r.db.Save(team).Error
}

func (r *cloudRepo) GetActiveServerCount(teamID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.CloudServer{}).Where("team_id = ? AND is_active = ?", teamID, true).Count(&count).Error
	return count, err
}

func (r *cloudRepo) GetDeploymentsInLastHour(teamID uint) (int64, error) {
	var count int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err := r.db.Model(&models.CloudUsageLog{}).
		Where("team_id = ? AND reported_at >= ?", teamID, oneHourAgo).
		Select("SUM(deployments)").
		Row().
		Scan(&count)

	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (r *cloudRepo) LogUsage(usage *models.CloudUsageLog) error {
	return r.db.Create(usage).Error
}

func (r *cloudRepo) LogTelemetry(logEntry *models.CloudTelemetryLog) error {
	return r.db.Create(logEntry).Error
}

func (r *cloudRepo) InsertAuditLog(ctx context.Context, entry *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *cloudRepo) GetTotalServers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CloudServer{}).Count(&count).Error
	return count, err
}

func (r *cloudRepo) GetActiveSubscriptions(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CloudTeam{}).Where("plan != ?", "hobby").Count(&count).Error
	return count, err
}

func (r *cloudRepo) RegisterServer(ctx context.Context, teamID uint, token string, name string, ip string) error {
	server := &models.CloudServer{
		WorkspaceID: teamID,
		ServerID:    token,
		Name:        name,
		IPAddress:   ip,
		IsActive:    true,
		LastPing:    time.Now(),
	}
	return r.db.WithContext(ctx).Create(server).Error
}

func (r *cloudRepo) GetServerByToken(token string) (*models.CloudServer, error) {
	var server models.CloudServer
	if err := r.db.Where("server_id = ?", token).First(&server).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &server, nil
}

func (r *cloudRepo) ListAuditLogs(ctx context.Context, limit int, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.WithContext(ctx).Order("timestamp desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}

func (r *cloudRepo) GetCurrentMonthUsage(teamID uint) (int, int, error) {
	var result struct {
		TotalHours     int
		TotalBandwidth int
	}
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	err := r.db.Model(&models.CloudUsageLog{}).
		Select("COALESCE(SUM(container_hours), 0) as total_hours, COALESCE(SUM(bandwidth_gb), 0) as total_bandwidth").
		Where("team_id = ? AND reported_at >= ?", teamID, startOfMonth).
		Scan(&result).Error
	return result.TotalHours, result.TotalBandwidth, err
}

func (r *cloudRepo) GetUserTeams(userID string) ([]models.CloudTeam, error) {
	// Find all workspace memberships for this user
	var memberships []models.WorkspaceMember
	if err := r.db.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return nil, err
	}

	if len(memberships) == 0 {
		return []models.CloudTeam{}, nil
	}

	// Extract team IDs
	var teamIDs []string
	for _, m := range memberships {
		teamIDs = append(teamIDs, m.WorkspaceID)
	}

	// Find the actual CloudTeams
	var teams []models.CloudTeam
	if err := r.db.Where("id IN ?", teamIDs).Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *cloudRepo) GetTeamServers(teamID uint) ([]models.CloudServer, error) {
	var servers []models.CloudServer
	if err := r.db.Where("workspace_id = ?", teamID).Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}
