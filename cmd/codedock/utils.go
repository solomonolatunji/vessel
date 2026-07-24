package main

import (
	"fmt"
	"os"

	"codedock.run/codedock/pkg/config"
	"codedock.run/codedock/pkg/http"
)

func getClient() *http.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	if cfg.ServerURL == "" || cfg.Token == "" {
		fmt.Println("Error: Not authenticated. Please run 'codedock login' first.")
		os.Exit(1)
	}
	return http.NewClient(cfg.ServerURL, cfg.Token)
}
