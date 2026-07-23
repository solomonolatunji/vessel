package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [deployment_id]",
	Short: "View logs for a deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		logs, err := client.GetDeploymentLogs(args[0])
		if err != nil {
			fmt.Printf("Error fetching logs: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(logs)
	},
}

func init() {
	appsCmd.AddCommand(logsCmd)
}
