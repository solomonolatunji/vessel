package services

import (
	"strings"
	"testing"
	"time"
)

func TestBuildLokiURL(t *testing.T) {
	svc := NewLogService()
	svc.lokiURL = "http://localhost:3100"

	start := time.Unix(1600000000, 0)
	end := time.Unix(1600003600, 0)

	opts := HistoricalLogsOpts{
		ServiceID: "test-service-123",
		Start:     start,
		End:       end,
		Limit:     100,
	}

	u := svc.buildLokiURL(opts)

	if !strings.HasPrefix(u, "http://localhost:3100/loki/api/v1/query_range?") {
		t.Errorf("unexpected URL prefix: %s", u)
	}

	if !strings.Contains(u, "query=%7Bcontainer_label_vessl_service%3D%22test-service-123%22%7D") {
		t.Errorf("URL missing or incorrect query param: %s", u)
	}

	if !strings.Contains(u, "start=1600000000000000000") {
		t.Errorf("URL missing or incorrect start param: %s", u)
	}

	if !strings.Contains(u, "end=1600003600000000000") {
		t.Errorf("URL missing or incorrect end param: %s", u)
	}

	if !strings.Contains(u, "limit=100") {
		t.Errorf("URL missing or incorrect limit param: %s", u)
	}
}
