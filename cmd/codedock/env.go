package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"codedock.run/codedock/internal/models"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:     "environments",
	Aliases: []string{"environment", "env"},
	Short:   "Manage environments",
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments for a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}

		client := getClient()
		envs, err := client.ListEnvironments(projectID)
		if err != nil {
			fmt.Printf("Error listing environments: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tPROJECT_ID\tNAME\tDEFAULT")
		for _, e := range envs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", e.ID, e.ProjectID, e.Name, e.IsDefault)
		}
		w.Flush()
	},
}

var envCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an environment",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")
		isDefault, _ := cmd.Flags().GetBool("default")

		if projectID == "" || name == "" {
			fmt.Println("Error: --project and --name flags are required")
			os.Exit(1)
		}

		client := getClient()
		req := &models.EnvironmentConfig{
			Name:      name,
			IsDefault: isDefault,
		}

		created, err := client.CreateEnvironment(projectID, req)
		if err != nil {
			fmt.Printf("Error creating environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Environment %s created successfully with ID: %s\n", created.Name, created.ID)
	},
}

var envDestroyCmd = &cobra.Command{
	Use:   "destroy [id]",
	Short: "Destroy an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		if err := client.DeleteEnvironment(args[0]); err != nil {
			fmt.Printf("Error destroying environment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Environment %s destroyed successfully\n", args[0])
	},
}

func init() {
	envListCmd.Flags().StringP("project", "p", "", "Project ID (required)")

	envCreateCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	envCreateCmd.Flags().StringP("name", "n", "", "Environment name (required)")
	envCreateCmd.Flags().Bool("default", false, "Set as default environment")

	envCmd.AddCommand(envListCmd, envCreateCmd, envDestroyCmd)
	rootCmd.AddCommand(envCmd)
}
