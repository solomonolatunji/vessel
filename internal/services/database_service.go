package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/utils"
)

type DatabaseService struct {
	repo     repositories.DatabaseRepository
	deployer *engine.DatabaseDeployer
}

func NewDatabaseService(r repositories.DatabaseRepository, d *engine.DatabaseDeployer) *DatabaseService {
	return &DatabaseService{
		repo:     r,
		deployer: d,
	}
}

func (s *DatabaseService) CreateDatabase(ctx context.Context, db *models.Database) (*models.Database, error) {
	if db == nil || db.Name == "" || db.Engine == "" {
		return nil, errors.New("database name and engine required")
	}
	if db.ID == "" {
		db.ID = uuid.New().String()
	}
	if db.Status == "" {
		db.Status = "stopped"
	}
	now := time.Now()
	if db.CreatedAt.IsZero() {
		db.CreatedAt = now
	}
	db.UpdatedAt = now
	if err := s.repo.Create(ctx, db); err != nil {
		return nil, err
	}
	if s.deployer != nil {
		containerID, err := s.deployer.SpinUp(ctx, db)
		if err == nil && containerID != "" {
			db.ContainerID = containerID
			db.Status = "running"
			_ = s.repo.Update(ctx, db)
		} else if err != nil {
			db.Status = "error"
			_ = s.repo.Update(ctx, db)
		}
	}
	return db, nil
}

func (s *DatabaseService) CreateDatabaseFromRequest(ctx context.Context, req *models.CreateDatabaseRequest) (*models.Database, error) {
	if req.Name == "" || req.Engine == "" {
		return nil, errors.New("name and engine fields are required")
	}
	if req.Port <= 0 {
		switch strings.ToLower(req.Engine) {
		case "postgres", "postgresql":
			req.Port = 5432
		case "mysql":
			req.Port = 3306
		case "redis":
			req.Port = 6379
		case "mongodb", "mongo":
			req.Port = 27017
		default:
			req.Port = 5432
		}
	}
	if req.Username == "" && strings.ToLower(req.Engine) != "redis" {
		req.Username = "vessladmin"
	}
	if req.DatabaseName == "" {
		req.DatabaseName = "appdb"
	}
	db := &models.Database{
		ProjectID:     req.ProjectID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Engine:        req.Engine,
		Version:       req.Version,
		Port:          req.Port,
		Username:      req.Username,
		Password:      req.Password,
		DatabaseName:  req.DatabaseName,
		VolumePath:    req.VolumePath,
		CustomArgs:    req.CustomArgs,
		Status:        "stopped",
	}
	return s.CreateDatabase(ctx, db)
}

func (s *DatabaseService) GetDatabase(ctx context.Context, id string) (*models.Database, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DatabaseService) ListDatabases(ctx context.Context) ([]*models.Database, error) {
	return s.repo.List(ctx)
}

func (s *DatabaseService) ListDatabasesByProject(ctx context.Context, projectID string) ([]*models.Database, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.repo.ListByProject(ctx, projectID)
}

func (s *DatabaseService) UpdateDatabase(ctx context.Context, db *models.Database) error {
	if db == nil || db.ID == "" {
		return errors.New("valid database required for update")
	}
	db.UpdatedAt = time.Now()
	return s.repo.Update(ctx, db)
}

func (s *DatabaseService) DeleteDatabase(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if s.deployer != nil {
		_ = s.deployer.Stop(ctx, id)
	}
	return s.repo.Delete(ctx, id)
}

func (s *DatabaseService) StartDatabase(ctx context.Context, id string) (*models.Database, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, utils.NewNotFoundError("Database", id)
	}
	if s.deployer == nil {
		return nil, errors.New("database deployer unavailable")
	}
	containerID, err := s.deployer.SpinUp(ctx, db)
	if err != nil {
		db.Status = "error"
		_ = s.repo.Update(ctx, db)
		return nil, err
	}
	if containerID != "" {
		db.ContainerID = containerID
	}
	db.Status = "running"
	db.UpdatedAt = time.Now()
	_ = s.repo.Update(ctx, db)
	return db, nil
}

func (s *DatabaseService) StopDatabase(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return errors.New("database not found")
	}
	if s.deployer != nil {
		_ = s.deployer.Stop(ctx, id)
	}
	db.Status = "stopped"
	db.UpdatedAt = time.Now()
	return s.repo.Update(ctx, db)
}
