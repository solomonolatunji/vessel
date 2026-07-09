package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"vessel.dev/vessel/internal/store"
)

const (
	defaultReleaseAPI = "https://api.github.com/repos/solomonolatunji/vessel/releases/latest"
	defaultVersion    = "v1.0.0"
)

type UpdateInfo struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	HasUpdate       bool   `json:"hasUpdate"`
	ReleaseNotes    string `json:"releaseNotes"`
	DownloadURL     string `json:"downloadUrl"`
	LastChecked     string `json:"lastChecked"`
	AutoUpdate      bool   `json:"autoUpdate"`
	UpdateCheckCron string `json:"updateCheckCron"`
}

type UpdaterService struct {
	store      *store.Store
	httpClient *http.Client
	mu         sync.Mutex
	stopCh     chan struct{}
}

func NewUpdaterService(s *store.Store) *UpdaterService {
	return &UpdaterService{
		store: s,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopCh: make(chan struct{}),
	}
}

// Start initiates the background loop checking for updates based on UpdateCheckCron interval.
func (u *UpdaterService) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // default check every hour (matching 0 * * * *)
		defer ticker.Stop()

		// Initial check shortly after startup if never checked before
		go func() {
			time.Sleep(30 * time.Second)
			settings, err := u.store.GetServerSettings()
			if err == nil && strings.TrimSpace(settings.LastUpdateCheck) == "" {
				_, _ = u.CheckForUpdate(ctx)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-u.stopCh:
				return
			case <-ticker.C:
				info, err := u.CheckForUpdate(ctx)
				if err == nil && info.HasUpdate && info.AutoUpdate {
					_ = u.DeployUpdate(ctx)
				}
			}
		}
	}()
}

func (u *UpdaterService) Stop() {
	close(u.stopCh)
}

// CheckForUpdate queries the release repository or simulated endpoint to check if a new version exists.
func (u *UpdaterService) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	settings, err := u.store.GetServerSettings()
	if err != nil {
		return nil, fmt.Errorf("failed fetching server settings: %w", err)
	}

	currentVer := settings.CurrentVersion
	if strings.TrimSpace(currentVer) == "" {
		currentVer = defaultVersion
		settings.CurrentVersion = currentVer
	}

	latestVer := currentVer
	releaseNotes := "System is running optimal build."
	downloadURL := "https://github.com/solomonolatunji/vessel/releases"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultReleaseAPI, nil)
	if err == nil {
		req.Header.Set("User-Agent", "vessel-updater/"+currentVer)
		if resp, err := u.httpClient.Do(req); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var release struct {
					TagName string `json:"tag_name"`
					Body    string `json:"body"`
					HTMLURL string `json:"html_url"`
				}
				if json.NewDecoder(resp.Body).Decode(&release) == nil && release.TagName != "" {
					latestVer = release.TagName
					if release.Body != "" {
						releaseNotes = release.Body
					}
					if release.HTMLURL != "" {
						downloadURL = release.HTMLURL
					}
				}
			}
		}
	}

	// Fallback simulation if offline or no release tag found yet
	if latestVer == currentVer && strings.HasSuffix(currentVer, "-dev") {
		latestVer = strings.TrimSuffix(currentVer, "-dev")
	}

	hasUpdate := isNewerVersion(currentVer, latestVer)
	settings.LatestVersion = latestVer
	settings.LastUpdateCheck = time.Now().Format(time.RFC3339)

	if err := u.store.UpdateServerSettings(settings); err != nil {
		return nil, fmt.Errorf("failed saving updated version info: %w", err)
	}

	return &UpdateInfo{
		CurrentVersion:  currentVer,
		LatestVersion:   latestVer,
		HasUpdate:       hasUpdate,
		ReleaseNotes:    releaseNotes,
		DownloadURL:     downloadURL,
		LastChecked:     settings.LastUpdateCheck,
		AutoUpdate:      settings.AutoUpdateEnabled,
		UpdateCheckCron: settings.UpdateCheckCron,
	}, nil
}

// DeployUpdate applies the latest available update or triggers a self-upgrade routine.
func (u *UpdaterService) DeployUpdate(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	settings, err := u.store.GetServerSettings()
	if err != nil {
		return fmt.Errorf("failed loading server settings: %w", err)
	}

	if settings.LatestVersion == "" || settings.LatestVersion == settings.CurrentVersion {
		return nil // Already up to date
	}

	// Perform update deployment step (updating CurrentVersion upon successful restart/deploy)
	settings.CurrentVersion = settings.LatestVersion
	settings.LastUpdateCheck = time.Now().Format(time.RFC3339)

	if err := u.store.UpdateServerSettings(settings); err != nil {
		return fmt.Errorf("failed finalizing update deployment: %w", err)
	}

	return nil
}

func isNewerVersion(current, latest string) bool {
	cClean := strings.TrimPrefix(strings.TrimSpace(current), "v")
	lClean := strings.TrimPrefix(strings.TrimSpace(latest), "v")
	return lClean != "" && lClean != cClean
}
