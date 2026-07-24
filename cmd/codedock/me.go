package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current logged-in user",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()

		user, err := client.Me()
		if err != nil {
			fmt.Printf("❌ Failed to fetch user profile: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("👤 Logged in as:\n")
		fmt.Printf("  Name:  %s\n", user.Name)
		fmt.Printf("  Email: %s\n", user.Email)
		fmt.Printf("  Role:  %s\n", user.Role)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
