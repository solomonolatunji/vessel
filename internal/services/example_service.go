package services

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"codedock.dev/codedock/internal/models"
)

type ExampleService struct {
	cache       []models.ExampleApp
	lastFetched time.Time
	mu          sync.RWMutex
}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (s *ExampleService) ListExamples() ([]models.ExampleApp, error) {
	s.mu.RLock()
	cacheValid := time.Since(s.lastFetched) < 1*time.Hour && len(s.cache) > 0
	cached := s.cache
	s.mu.RUnlock()

	if cacheValid {
		return cached, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/buildwithtechx/codedock-examples/contents")
	if err != nil {
		if len(cached) > 0 {
			return cached, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	var contents []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		HtmlUrl string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		if len(cached) > 0 {
			return cached, nil
		}
		return nil, err
	}

	var examples []models.ExampleApp
	for _, c := range contents {
		if c.Type == "dir" && !strings.HasPrefix(c.Name, ".") {
			examples = append(examples, models.ExampleApp{
				ID:          c.Name,
				Name:        formatExampleName(c.Name),
				Description: "Deploy " + formatExampleName(c.Name) + " example app",
				Repo:        c.HtmlUrl,
			})
		}
	}

	s.mu.Lock()
	s.cache = examples
	s.lastFetched = time.Now()
	s.mu.Unlock()

	return examples, nil
}

func formatExampleName(name string) string {
	parts := strings.Split(name, "-")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, " ")
}
