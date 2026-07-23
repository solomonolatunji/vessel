package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MetricsService struct {
	tsdbURL    string
	httpClient *http.Client
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		tsdbURL:    "http://127.0.0.1:8428",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type ServiceMetricsOpts struct {
	ServiceID string
	Start     time.Time
	End       time.Time
	Step      string
}

func (s *MetricsService) GetServiceMetrics(ctx context.Context, opts ServiceMetricsOpts) (map[string]any, error) {
	cpuQuery := fmt.Sprintf(`rate(container_cpu_usage_seconds_total{container_label_codedock_service="%s"}[5m])`, opts.ServiceID)
	memQuery := fmt.Sprintf(`container_memory_usage_bytes{container_label_codedock_service="%s"}`, opts.ServiceID)
	netRxQuery := fmt.Sprintf(`rate(container_network_receive_bytes_total{container_label_codedock_service="%s"}[5m])`, opts.ServiceID)
	netTxQuery := fmt.Sprintf(`rate(container_network_transmit_bytes_total{container_label_codedock_service="%s"}[5m])`, opts.ServiceID)

	cpuData, err := s.queryRange(ctx, cpuQuery, opts.Start, opts.End, opts.Step)
	if err != nil {
		return nil, err
	}

	memData, err := s.queryRange(ctx, memQuery, opts.Start, opts.End, opts.Step)
	if err != nil {
		return nil, err
	}

	netRxData, err := s.queryRange(ctx, netRxQuery, opts.Start, opts.End, opts.Step)
	if err != nil {
		return nil, err
	}

	netTxData, err := s.queryRange(ctx, netTxQuery, opts.Start, opts.End, opts.Step)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"cpu":        cpuData,
		"memory":     memData,
		"network_rx": netRxData,
		"network_tx": netTxData,
	}, nil
}

func (s *MetricsService) queryRange(ctx context.Context, query string, start, end time.Time, step string) (any, error) {
	u, _ := url.Parse(s.tsdbURL + "/api/v1/query_range")
	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.Unix()))
	q.Set("end", fmt.Sprintf("%d", end.Unix()))
	q.Set("step", step)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("tsdb returned status %d", resp.StatusCode)
	}

	var res struct {
		Data any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
