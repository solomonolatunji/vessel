package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type OneClickService struct {
	tmplManager *engine.TemplateManager
	dbDeployer  *engine.DatabaseDeployer
	envRepo     repositories.EnvironmentRepository
	dbRepo      repositories.DatabaseRepository
}

func NewOneClickService(
	tm *engine.TemplateManager,
	dd *engine.DatabaseDeployer,
	er repositories.EnvironmentRepository,
	dr repositories.DatabaseRepository,
) *OneClickService {
	return &OneClickService{
		tmplManager: tm,
		dbDeployer:  dd,
		envRepo:     er,
		dbRepo:      dr,
	}
}

func (s *OneClickService) ListApps() []models.OneClickApp {
	apps := []models.OneClickApp{}
	for _, id := range s.tmplManager.ListTemplates() {
		tmpl, err := s.tmplManager.GetTemplate(id)
		if err != nil {
			continue
		}
		app := extractOneClickApp(id, &tmpl)
		if app != nil {
			apps = append(apps, *app)
		}
	}
	return apps
}

func (s *OneClickService) DeployApp(ctx context.Context, appID, projectID, name string) (*models.Database, error) {
	tmpl, err := s.tmplManager.GetTemplate(appID)
	if err != nil {
		return nil, errors.New("unknown app: " + appID)
	}

	meta := findOneClickMetadata(&tmpl)
	if meta == nil {
		return nil, errors.New("app has no one-click metadata")
	}

	if projectID == "" {
		return nil, errors.New("projectId is required")
	}

	envs, err := s.envRepo.ListByProject(ctx, projectID)
	if err != nil || len(envs) == 0 {
		return nil, errors.New("project has no environments")
	}

	appName := name
	if appName == "" {
		appName = meta.Name
	}

	db := buildDatabaseRecord(projectID, envs[0].ID, appID, appName, extractPort(&tmpl))
	db.Password = uuid.New().String()[:16]

	if err := s.dbRepo.Create(ctx, db); err != nil {
		return nil, err
	}

	if s.dbDeployer == nil {
		return db, nil
	}

	containerID, err := s.dbDeployer.SpinUp(ctx, db)
	if err != nil {
		db.Status = models.DatabaseStatusError
		_ = s.dbRepo.Update(ctx, db)
		return nil, err
	}

	db.ContainerID = containerID
	db.Status = models.DatabaseStatusRunning
	_ = s.dbRepo.Update(ctx, db)
	return db, nil
}

func extractOneClickApp(id string, tmpl *engine.ComposeTemplate) *models.OneClickApp {
	if tmpl.XCodedock != nil && tmpl.XCodedock.IsOneClick {
		return &models.OneClickApp{
			ID:          id,
			Name:        tmpl.XCodedock.Name,
			Description: tmpl.XCodedock.Description,
			Port:        extractPort(tmpl),
		}
	}
	for _, svc := range tmpl.Services {
		if svc.XCodedock != nil && svc.XCodedock.IsOneClick {
			return &models.OneClickApp{
				ID:          id,
				Name:        svc.XCodedock.Name,
				Description: svc.XCodedock.Description,
				Port:        parsePortFromString(svc.Ports),
			}
		}
	}
	return nil
}

func findOneClickMetadata(tmpl *engine.ComposeTemplate) *engine.CodedockMetadata {
	if tmpl.XCodedock != nil && tmpl.XCodedock.IsOneClick {
		return tmpl.XCodedock
	}
	for _, svc := range tmpl.Services {
		if svc.XCodedock != nil && svc.XCodedock.IsOneClick {
			return svc.XCodedock
		}
	}
	return nil
}

func extractPort(tmpl *engine.ComposeTemplate) int {
	for _, svc := range tmpl.Services {
		if svc.XCodedock != nil && svc.XCodedock.IsOneClick && len(svc.Ports) > 0 {
			return parsePortFromString(svc.Ports)
		}
		if tmpl.XCodedock != nil && tmpl.XCodedock.IsOneClick && len(svc.Ports) > 0 {
			return parsePortFromString(svc.Ports)
		}
	}
	return 3000
}

func parsePortFromString(ports []string) int {
	if len(ports) == 0 {
		return 3000
	}
	var p int
	for _, c := range ports[0] {
		if c >= '0' && c <= '9' {
			p = p*10 + int(c-'0')
		} else {
			break
		}
	}
	if p <= 0 {
		return 3000
	}
	return p
}

func buildDatabaseRecord(projectID, envID, engineID, name string, port int) *models.Database {
	return &models.Database{
		ID:            uuid.New().String(),
		ProjectID:     projectID,
		EnvironmentID: envID,
		Name:          name,
		Engine:        models.DatabaseEngine(engineID),
		Port:          port,
		Status:        models.DatabaseStatusCreated,
		Username:      "codedock",
		DatabaseName:  "codedock",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
