package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"vessl.dev/vessl/internal/models"
)

type mockAppRepo struct {
	createErr    error
	updateErr    error
	createCalled bool
	updateCalled bool
}

func (m *mockAppRepo) Create(ctx context.Context, svc *models.AppService) error {
	m.createCalled = true
	return m.createErr
}
func (m *mockAppRepo) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) ListAll(ctx context.Context) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) Update(ctx context.Context, svc *models.AppService) error {
	m.updateCalled = true
	return m.updateErr
}
func (m *mockAppRepo) Delete(ctx context.Context, id string) error                { return nil }
func (m *mockAppRepo) CreateWebhook(ctx context.Context, w *models.Webhook) error { return m.createErr }
func (m *mockAppRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	return nil, nil
}
func (m *mockAppRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error { return nil }

func TestCreateAppService(t *testing.T) {
	repo := &mockAppRepo{}
	svc := NewAppService(repo, nil) // varRepo is not used in CreateAppService

	ctx := context.Background()

	t.Run("invalid service", func(t *testing.T) {
		_, err := svc.CreateAppService(ctx, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing project id", func(t *testing.T) {
		_, err := svc.CreateAppService(ctx, &models.AppService{Name: "test"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		_, err := svc.CreateAppService(ctx, &models.AppService{ProjectID: "proj-1"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		app := &models.AppService{ProjectID: "proj-1", Name: "app-1"}
		res, err := svc.CreateAppService(ctx, app)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.ID == "" {
			t.Error("expected ID to be set")
		}
		if string(res.Status) != string(models.AppServiceStatusCreated) {
			t.Errorf("expected status %s, got %s", models.AppServiceStatusCreated, res.Status)
		}
		if !repo.createCalled {
			t.Error("expected repo Create to be called")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo.createErr = errors.New("db error")
		app := &models.AppService{ProjectID: "proj-2", Name: "app-2"}
		_, err := svc.CreateAppService(ctx, app)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		repo.createErr = nil // reset
	})
}

func TestUpdateAppService(t *testing.T) {
	repo := &mockAppRepo{}
	svc := NewAppService(repo, nil)
	ctx := context.Background()

	t.Run("invalid service", func(t *testing.T) {
		err := svc.UpdateAppService(ctx, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		err := svc.UpdateAppService(ctx, &models.AppService{Name: "app-1"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		repo.updateCalled = false
		app := &models.AppService{ID: "app-1"}
		before := time.Now()
		time.Sleep(1 * time.Millisecond)
		err := svc.UpdateAppService(ctx, app)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !repo.updateCalled {
			t.Error("expected repo Update to be called")
		}
		if !app.UpdatedAt.After(before) {
			t.Error("expected UpdatedAt to be updated")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo.updateErr = errors.New("db error")
		app := &models.AppService{ID: "app-2"}
		err := svc.UpdateAppService(ctx, app)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
