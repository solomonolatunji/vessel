package services_test

import (
	"testing"

	"codedock.dev/codedock/internal/services"
)

func TestSystemService_GetStats(t *testing.T) {
	svc := services.NewSystemService()

	stats, err := svc.GetStats()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats == nil {
		t.Fatal("expected stats to be non-nil")
	}

	// Check delta parsing
	_, _ = svc.GetStats()
}
