package services

import (
	"context"
	"errors"
	"testing"

	"vessl.dev/vessl/internal/models"
)

func TestCreateWebhook(t *testing.T) {
	repo := &mockAppRepo{}
	svc := NewAppService(repo, nil)

	ctx := context.Background()

	t.Run("invalid webhook", func(t *testing.T) {
		_, err := svc.CreateWebhook(ctx, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing service ID", func(t *testing.T) {
		_, err := svc.CreateWebhook(ctx, &models.Webhook{URL: "http://example.com"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing URL", func(t *testing.T) {
		_, err := svc.CreateWebhook(ctx, &models.Webhook{ServiceID: "svc-1"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		webhook := &models.Webhook{ServiceID: "svc-1", URL: "http://example.com"}
		repo.createErr = nil

		res, err := svc.CreateWebhook(ctx, webhook)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.ID == "" {
			t.Error("expected ID to be set")
		}
		if res.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo.createErr = errors.New("db error")
		webhook := &models.Webhook{ServiceID: "svc-2", URL: "http://example.com/2"}
		_, err := svc.CreateWebhook(ctx, webhook)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		repo.createErr = nil
	})
}
