package services

import (
	"context"
	"database/sql"
	"log"
	"runtime"
	"time"

	"vessl.dev/vessl/internal/repositories"
)

func StartTelemetryReporter(db *sql.DB, version string) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

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

	if !settings.TelemetryEnabled {
		return
	}

	appRepo := repositories.NewAppServiceSQLiteRepository(db)
	apps, err := appRepo.ListAll(context.Background())
	activeApps := 0
	if err == nil {
		activeApps = len(apps)
	}

	log.Printf("Telemetry: instance=%s version=%s os=%s arch=%s apps=%d", settings.ID, version, runtime.GOOS, runtime.GOARCH, activeApps)
}
