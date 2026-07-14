package main

import (
	"fmt"
	"os"

	"vessl.dev/vessl/pkg/config"
	"vessl.dev/vessl/pkg/http"
)

func getClient() *http.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	if cfg.ServerURL == "" || cfg.Token == "" {
		fmt.Println("Error: Not authenticated. Please run 'vessl login' first.")
		os.Exit(1)
	}
	return http.NewClient(cfg.ServerURL, cfg.Token)
}
