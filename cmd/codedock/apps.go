package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"codedock.run/codedock/internal/models"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage applications",
}

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications",
	Run: func(cmd *cobra.Command, args []string) {
		envID, _ := cmd.Flags().GetString("environment")
		if envID == "" {
			fmt.Println("Error: --environment flag is required")
			os.Exit(1)
		}

		client := getClient()
		apps, err := client.ListServices(envID)
		if err != nil {
			fmt.Printf("Error listing applications: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tDOMAIN")
		for _, app := range apps {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", app.ID, app.Name, app.Status, app.Domain)
		}
		w.Flush()
	},
}

var appsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an application",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		envID, _ := cmd.Flags().GetString("environment")
		name, _ := cmd.Flags().GetString("name")
		repoURL, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")

		if projectID == "" || envID == "" || name == "" {
			fmt.Println("Error: --project, --environment, and --name flags are required")
			os.Exit(1)
		}

		client := getClient()
		app := &models.AppService{
			ProjectID:     projectID,
			EnvironmentID: envID,
			Name:          name,
			RepositoryURL: repoURL,
			Branch:        branch,
		}

		created, err := client.CreateService(app)
		if err != nil {
			fmt.Printf("Error creating application: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Application %s created successfully with ID: %s\n", created.Name, created.ID)
	},
}

var appsDestroyCmd = &cobra.Command{
	Use:   "destroy [id]",
	Short: "Destroy an application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		if err := client.DeleteService(args[0]); err != nil {
			fmt.Printf("Error destroying application: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Application %s destroyed successfully\n", args[0])
	},
}

func init() {
	appsListCmd.Flags().StringP("environment", "e", "", "Environment ID (required)")

	appsCreateCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	appsCreateCmd.Flags().StringP("environment", "e", "", "Environment ID (required)")
	appsCreateCmd.Flags().StringP("name", "n", "", "Application name (required)")
	appsCreateCmd.Flags().String("repo", "", "Repository URL")
	appsCreateCmd.Flags().String("branch", "main", "Branch name")

	appsCmd.AddCommand(appsListCmd, appsCreateCmd, appsDestroyCmd)
	rootCmd.AddCommand(appsCmd)
}
