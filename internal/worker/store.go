package worker

import (
	"codedock.run/codedock/internal/models"
)

type WorkerLocalStore struct {
	payload models.WorkerDeployAppPayload
}

func NewWorkerLocalStore(p models.WorkerDeployAppPayload) *WorkerLocalStore {
	return &WorkerLocalStore{payload: p}
}

func (s *WorkerLocalStore) GetServerSettings() (*models.ServerSettings, error) {
	return &models.ServerSettings{}, nil
}

func (s *WorkerLocalStore) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return nil, nil
}

func (s *WorkerLocalStore) GetEnvVars(projectID string) (map[string]string, error) {
	return s.payload.Env, nil
}

func (s *WorkerLocalStore) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	return nil, nil // Return empty, assuming Env contains all compiled vars
}

func (s *WorkerLocalStore) GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error) {
	return nil, nil
}

func (s *WorkerLocalStore) UpdateAppService(app *models.AppService) error {
	return nil // No-op
}

func (s *WorkerLocalStore) ListLogDrainsByService(serviceID string) ([]*models.LogDrain, error) {
	return nil, nil // No log drains supported in standalone worker yet
}
