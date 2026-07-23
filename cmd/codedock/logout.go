package main

import (
	"fmt"
	"os"

	"codedock.dev/codedock/pkg/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your Codedock account and clear credentials",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err == nil && cfg.Token != "" && cfg.ServerURL != "" {
			client := getClient()
			_ = client.Logout()
		}

		if err := clearConfig(); err != nil {
			fmt.Printf("❌ Failed to clear config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("👋 Successfully logged out.")
	},
}

func clearConfig() error {
	return config.Save(&config.Config{})
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
