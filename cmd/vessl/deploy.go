package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [service-id]",
	Short: "Trigger a deployment for an existing service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]
		client := getClient()

		fmt.Printf("🚀 Triggering deployment for service %s...\n", serviceID)
		deployment, err := client.TriggerDeployment(serviceID)
		if err != nil {
			fmt.Printf("❌ Failed to trigger deployment: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Deployment started!\n")
		fmt.Printf("ID: %s\n", deployment.ID)
		fmt.Printf("Status: %s\n", deployment.Status)
		fmt.Println("To check status, use the dashboard or future CLI commands.")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
