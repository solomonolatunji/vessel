package main

import (
	"fmt"
	"os"

	codedockhttp "codedock.run/codedock/pkg/http"
	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Manage docker-compose deployments",
}

var composeAnalyzeCmd = &cobra.Command{
	Use:   "analyze <file>",
	Short: "Analyze a docker-compose.yml file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("❌ Failed to read file: %v\n", err)
			os.Exit(1)
		}

		client := getClient()
		projectID, _ := cmd.Flags().GetString("project")

		req := &codedockhttp.ComposeAnalyzeRequest{
			ProjectID:      projectID,
			ComposeContent: string(content),
		}

		fmt.Println("🔍 Analyzing compose file...")
		res, err := client.AnalyzeCompose(req)
		if err != nil {
			fmt.Printf("❌ Failed to analyze: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Analysis complete!\n")
		fmt.Printf("Found %d app services and %d databases.\n", len(res.AppServices), len(res.Databases))
	},
}

var composeDeployCmd = &cobra.Command{
	Use:   "deploy <file>",
	Short: "Deploy a docker-compose.yml file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("❌ Failed to read file: %v\n", err)
			os.Exit(1)
		}

		client := getClient()
		projectID, _ := cmd.Flags().GetString("project")

		fmt.Printf("🚀 Deploying compose file %s...\n", filename)
		count, err := client.DeployCompose(projectID, content, filename)
		if err != nil {
			fmt.Printf("❌ Failed to deploy: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Deployment successful! Created %d resources.\n", count)
	},
}

func init() {
	composeAnalyzeCmd.Flags().StringP("project", "p", "", "Project ID (optional)")
	composeDeployCmd.Flags().StringP("project", "p", "", "Project ID (optional)")

	composeCmd.AddCommand(composeAnalyzeCmd)
	composeCmd.AddCommand(composeDeployCmd)
	rootCmd.AddCommand(composeCmd)
}
