package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type LogService struct {
	lokiURL    string
	httpClient *http.Client
}

func NewLogService() *LogService {
	return &LogService{
		lokiURL:    "http://127.0.0.1:3100",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type LokiQueryResponse struct {
	Data LokiQueryData `json:"data"`
}

type LokiQueryData struct {
	Result []LokiStream `json:"result"`
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// GetHistoricalLogs fetches logs from Loki for a given serviceID within the specified time range.
func (s *LogService) GetHistoricalLogs(ctx context.Context, serviceID string, start, end time.Time, limit int) ([]map[string]any, error) {
	reqURL := s.buildLokiURL(serviceID, start, end, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("loki returned status %d", resp.StatusCode)
	}

	var res LokiQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return s.parseLokiLogs(res), nil
}

func (s *LogService) buildLokiURL(serviceID string, start, end time.Time, limit int) string {
	query := fmt.Sprintf(`{service_id="%s"}`, serviceID)
	u, _ := url.Parse(s.lokiURL + "/loki/api/v1/query_range")
	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.UnixNano()))
	q.Set("end", fmt.Sprintf("%d", end.UnixNano()))
	q.Set("limit", fmt.Sprintf("%d", limit))
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *LogService) parseLokiLogs(res LokiQueryResponse) []map[string]any {
	var logs []map[string]any
	for _, stream := range res.Data.Result {
		for _, val := range stream.Values {
			if len(val) == 2 {
				logs = append(logs, map[string]any{
					"timestamp": val[0],
					"line":      val[1],
				})
			}
		}
	}
	return logs
}
