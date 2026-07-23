package services_test

import (
	"context"
	"errors"
	"testing"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
)

type mockAppRepo struct {
	apps map[string]*models.AppService
}

func (m *mockAppRepo) Create(ctx context.Context, svc *models.AppService) error {
	m.apps[svc.ID] = svc
	return nil
}
func (m *mockAppRepo) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	if app, ok := m.apps[id]; ok {
		return app, nil
	}
	return nil, errors.New("not found")
}
func (m *mockAppRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	var list []*models.AppService
	for _, a := range m.apps {
		if a.EnvironmentID == environmentID {
			list = append(list, a)
		}
	}
	return list, nil
}
func (m *mockAppRepo) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	var list []*models.AppService
	for _, a := range m.apps {
		if a.ProjectID == projectID {
			list = append(list, a)
		}
	}
	return list, nil
}
func (m *mockAppRepo) ListAll(ctx context.Context) ([]*models.AppService, error) {
	var list []*models.AppService
	for _, a := range m.apps {
		list = append(list, a)
	}
	return list, nil
}
func (m *mockAppRepo) Update(ctx context.Context, svc *models.AppService) error {
	if _, ok := m.apps[svc.ID]; !ok {
		return errors.New("not found")
	}
	m.apps[svc.ID] = svc
	return nil
}
func (m *mockAppRepo) Delete(ctx context.Context, id string) error {
	delete(m.apps, id)
	return nil
}
func (m *mockAppRepo) CreateWebhook(ctx context.Context, w *models.Webhook) error { return nil }
func (m *mockAppRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	return nil, nil
}
func (m *mockAppRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error { return nil }
func (m *mockAppRepo) CreateLogDrain(ctx context.Context, d *models.LogDrain) error  { return nil }
func (m *mockAppRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error) {
	return nil, nil
}
func (m *mockAppRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error { return nil }

type mockVarRepo struct{}

func (m *mockVarRepo) Create(ctx context.Context, v *models.Variable) error { return nil }
func (m *mockVarRepo) Update(ctx context.Context, v *models.Variable) error { return nil }
func (m *mockVarRepo) GetByID(ctx context.Context, id string) (*models.Variable, error) {
	return nil, nil
}
func (m *mockVarRepo) ListByService(ctx context.Context, serviceID string) ([]*models.Variable, error) {
	return nil, nil
}
func (m *mockVarRepo) Delete(ctx context.Context, id string) error { return nil }

type mockVolRepo struct{}

func (m *mockVolRepo) Create(ctx context.Context, volume *models.ServiceVolume) error { return nil }
func (m *mockVolRepo) GetByID(ctx context.Context, id string) (*models.ServiceVolume, error) {
	return nil, nil
}
func (m *mockVolRepo) ListByService(ctx context.Context, serviceID string) ([]models.ServiceVolume, error) {
	return nil, nil
}
func (m *mockVolRepo) Delete(ctx context.Context, id string) error { return nil }

func TestCreateAppService(t *testing.T) {
	repo := &mockAppRepo{apps: make(map[string]*models.AppService)}
	svc := services.NewAppService(repo, &mockVarRepo{}, &mockVolRepo{})

	ctx := context.Background()

	// Invalid input
	_, err := svc.CreateAppService(ctx, nil)
	if err == nil {
		t.Fatal("expected error on nil app service")
	}

	app := &models.AppService{ProjectID: "proj-1", Name: "my-app"}
	created, err := svc.CreateAppService(ctx, app)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected ID to be generated")
	}

	// Verify it was stored
	stored, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error fetching stored app: %v", err)
	}
	if stored.Name != "my-app" {
		t.Fatalf("expected name my-app, got %s", stored.Name)
	}
}

func TestGetAppService(t *testing.T) {
	repo := &mockAppRepo{apps: make(map[string]*models.AppService)}
	svc := services.NewAppService(repo, &mockVarRepo{}, &mockVolRepo{})

	ctx := context.Background()
	_, err := svc.GetAppService(ctx, "")
	if err == nil {
		t.Fatal("expected error on empty id")
	}

	app := &models.AppService{ID: "app-1", ProjectID: "proj-1", Name: "app-name"}
	repo.apps["app-1"] = app

	retrieved, err := svc.GetAppService(ctx, "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retrieved.ID != "app-1" {
		t.Fatalf("expected ID app-1, got %s", retrieved.ID)
	}
}
