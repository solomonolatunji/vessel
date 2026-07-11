package repos

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"vessel.dev/vessel/internal/cloud/models"
)

type CloudRepo interface {
	GetTeamByID(id uint) (*models.CloudTeam, error)
	GetTeamByName(name string) (*models.CloudTeam, error)
	GetActiveServerCount(teamID uint) (int64, error)
	GetDeploymentsInLastHour(teamID uint) (int64, error)
	LogUsage(usage *models.CloudUsageLog) error
	// Telemetry
	LogTelemetry(logEntry *models.CloudTelemetryLog) error
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
		// Might be null if no deployments
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
