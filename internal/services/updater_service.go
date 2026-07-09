package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"vessel.dev/vessel/internal/repositories"
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
	repo       repositories.SettingsRepository
	httpClient *http.Client
	mu         sync.Mutex
	stopCh     chan struct{}
}

func NewUpdaterService(repo repositories.SettingsRepository) *UpdaterService {
	return &UpdaterService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopCh: make(chan struct{}),
	}
}

func (u *UpdaterService) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		go func() {
			time.Sleep(30 * time.Second)
			settingsCfg, err := u.repo.GetServerSettings(ctx)
			if err == nil && strings.TrimSpace(settingsCfg.LastUpdateCheck) == "" {
				_, _ = u.CheckForUpdates(ctx)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-u.stopCh:
				return
			case <-ticker.C:
				info, err := u.CheckForUpdates(ctx)
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

func (u *UpdaterService) GetStatus() *UpdateInfo {
	u.mu.Lock()
	defer u.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settingsCfg, err := u.repo.GetServerSettings(ctx)
	if err != nil {
		return &UpdateInfo{}
	}

	return &UpdateInfo{
		CurrentVersion:  settingsCfg.CurrentVersion,
		LatestVersion:   settingsCfg.LatestVersion,
		HasUpdate:       isNewerVersion(settingsCfg.CurrentVersion, settingsCfg.LatestVersion),
		LastChecked:     settingsCfg.LastUpdateCheck,
		AutoUpdate:      settingsCfg.AutoUpdateEnabled,
		UpdateCheckCron: settingsCfg.UpdateCheckCron,
	}
}

func (u *UpdaterService) CheckForUpdates(ctx context.Context) (*UpdateInfo, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	settingsCfg, err := u.repo.GetServerSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed fetching server settings: %w", err)
	}

	currentVer := settingsCfg.CurrentVersion
	if strings.TrimSpace(currentVer) == "" {
		currentVer = defaultVersion
		settingsCfg.CurrentVersion = currentVer
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

	if latestVer == currentVer && strings.HasSuffix(currentVer, "-dev") {
		latestVer = strings.TrimSuffix(currentVer, "-dev")
	}

	hasUpdate := isNewerVersion(currentVer, latestVer)
	settingsCfg.LatestVersion = latestVer
	settingsCfg.LastUpdateCheck = time.Now().Format(time.RFC3339)

	if err := u.repo.UpdateServerSettings(ctx, settingsCfg); err != nil {
		return nil, fmt.Errorf("failed saving updated version info: %w", err)
	}

	return &UpdateInfo{
		CurrentVersion:  currentVer,
		LatestVersion:   latestVer,
		HasUpdate:       hasUpdate,
		ReleaseNotes:    releaseNotes,
		DownloadURL:     downloadURL,
		LastChecked:     settingsCfg.LastUpdateCheck,
		AutoUpdate:      settingsCfg.AutoUpdateEnabled,
		UpdateCheckCron: settingsCfg.UpdateCheckCron,
	}, nil
}

func (u *UpdaterService) DeployUpdate(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	settingsCfg, err := u.repo.GetServerSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed loading server settings: %w", err)
	}

	if settingsCfg.LatestVersion == "" || settingsCfg.LatestVersion == settingsCfg.CurrentVersion {
		return nil
	}

	settingsCfg.CurrentVersion = settingsCfg.LatestVersion
	settingsCfg.LastUpdateCheck = time.Now().Format(time.RFC3339)

	if err := u.repo.UpdateServerSettings(ctx, settingsCfg); err != nil {
		return fmt.Errorf("failed finalizing update deployment: %w", err)
	}

	return nil
}

func isNewerVersion(current, latest string) bool {
	cClean := strings.TrimPrefix(strings.TrimSpace(current), "v")
	lClean := strings.TrimPrefix(strings.TrimSpace(latest), "v")
	return lClean != "" && lClean != cClean
}
