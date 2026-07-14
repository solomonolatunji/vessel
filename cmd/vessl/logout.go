package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vessl.dev/vessl/pkg/config"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your Vessl account and clear credentials",
	Run: func(cmd *cobra.Command, args []string) {
		// Attempt to hit the logout endpoint if we have a valid client
		cfg, err := config.Load()
		if err == nil && cfg.Token != "" && cfg.ServerURL != "" {
			client := getClient() // getClient calls os.Exit if not logged in, but we checked token
			_ = client.Logout() // ignore error, we just want to clear local config
		}

		// Clear local config
		emptyCfg := &config.Config{
			ServerURL: "",
			Token:     "",
		}
		if err := config.Save(emptyCfg); err != nil {
			fmt.Printf("❌ Failed to clear config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("👋 Successfully logged out.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
