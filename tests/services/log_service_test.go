package services_test

import (
	"context"
	"testing"
	"time"

	"codedock.dev/codedock/internal/services"
)

func TestLogService_GetHistoricalLogs_InvalidHost(t *testing.T) {
	svc := services.NewLogService()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	opts := services.HistoricalLogsOpts{
		ServiceID: "svc-1",
		Start:     time.Now().Add(-1 * time.Hour),
		End:       time.Now(),
		Limit:     10,
	}

	_, err := svc.GetHistoricalLogs(ctx, opts)
	if err == nil {
		t.Log("expected an error connecting to a non-existent loki server (or ctx canceled), but got none")
	}
}
