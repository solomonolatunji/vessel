package main

import (
	"context"
	"fmt"
	"os"

	"vessl.dev/vessl/internal/repositories"
)

func runConfig() {
	_, db, _ := initDataDir()
	repo := repositories.NewSettingsSQLiteRepository(db)
	settings, err := repo.GetServerSettings(context.Background())
	if err != nil {
		exitError("Failed to load settings: %v", err)
	}

	if len(os.Args) < 3 {
		fmt.Println("Current configuration:")
		fmt.Printf("  site-name:         %s\n", settings.SiteName)
		fmt.Printf("  registration:      %v\n", settings.RegistrationEnabled)
		fmt.Printf("  telemetry:         %v\n", settings.TelemetryEnabled)
		fmt.Printf("  domain:            %s\n", os.Getenv("VESSL_DOMAIN"))
		fmt.Printf("  smtp-enabled:      %v\n", settings.SMTPEnabled)
		fmt.Printf("  resend-enabled:    %v\n", settings.ResendEnabled)
		fmt.Println("\nUsage: vessld config <key>=<value>")
		fmt.Println("  e.g.  vessld config site-name=MyVessl")
		fmt.Println("        vessld config registration=true")
		fmt.Println("  Note: panel-domain and wildcard-domain are set via .env, not here.")
		return
	}

	key, value, ok := stringsCut(os.Args[2], "=")
	if !ok {
		exitError("Usage: vessld config <key>=<value>")
	}

	switch key {
	case "site-name":
		settings.SiteName = value
	case "registration":
		settings.RegistrationEnabled = value == "true"
	case "telemetry":
		settings.TelemetryEnabled = value == "true"
	default:
		exitError("Unknown config key: %s (try: site-name, registration, telemetry)", key)
	}

	if err := repo.UpdateServerSettings(context.Background(), settings); err != nil {
		exitError("Failed to update settings: %v", err)
	}
	fmt.Printf("✅ %s set to %s\n", key, value)
}

func stringsCut(s, sep string) (string, string, bool) {
	for i := 0; i < len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return s[:i], s[i+len(sep):], true
		}
	}
	return s, "", false
}
