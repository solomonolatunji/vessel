package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"vessel.dev/vessel/internal/repositories"
)

type TelemetryPayload struct {
	InstanceID    string `json:"instance_id"`
	Version       string `json:"version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	ActiveServers int    `json:"active_servers"`
	ActiveApps    int    `json:"active_apps"`
}

func StartTelemetryReporter(db *sql.DB, version string) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		// Do an initial ping shortly after startup
		time.Sleep(5 * time.Minute)
		pingTelemetry(db, version)

		for range ticker.C {
			pingTelemetry(db, version)
		}
	}()
}

func pingTelemetry(db *sql.DB, version string) {
	settingsRepo := repositories.NewSettingsSQLiteRepository(db)
	settings, err := settingsRepo.GetServerSettings(context.Background())
	if err != nil {
		return
	}

	// By default, we can assume telemetry is enabled unless explicitly disabled, 
	// or we can strictly check if it's enabled. Let's strictly check:
	if !settings.TelemetryEnabled {
		return
	}

	appRepo := repositories.NewAppServiceSQLiteRepository(db)
	apps, err := appRepo.ListAll(context.Background())
	activeApps := 0
	if err == nil {
		activeApps = len(apps) // simplified, usually you'd filter by active status
	}

	payload := TelemetryPayload{
		InstanceID:    settings.ID,
		Version:       version,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		ActiveServers: 1, // OSS is 1 server
		ActiveApps:    activeApps,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Ping the cloud endpoint
	req, err := http.NewRequest("POST", "https://cloud.vessel.dev/api/cloud/telemetry/ping", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return // Silently fail to not spam logs for offline instances
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		log.Println("Telemetry ping sent successfully.")
	}
}
