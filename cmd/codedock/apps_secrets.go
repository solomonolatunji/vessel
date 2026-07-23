package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"codedock.dev/codedock/internal/models"
	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage application secrets",
}

var secretsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets for a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}

		client := getClient()
		secrets, err := client.GetSecrets(projectID)
		if err != nil {
			fmt.Printf("Error listing secrets: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "KEY\tVALUE")
		for k, v := range secrets {
			fmt.Fprintf(w, "%s\t%s\n", k, v)
		}
		w.Flush()
	},
}

var secretsSetCmd = &cobra.Command{
	Use:   "set [KEY=VALUE...]",
	Short: "Set secrets for a project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}

		req := make(models.SetEnvVarsRequest)
		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				fmt.Printf("Error: invalid format for secret: %s. Expected KEY=VALUE\n", arg)
				os.Exit(1)
			}
			req[parts[0]] = parts[1]
		}

		client := getClient()
		if err := client.SetSecrets(projectID, req); err != nil {
			fmt.Printf("Error setting secrets: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Secrets set successfully")
	},
}

func init() {
	secretsListCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	secretsSetCmd.Flags().StringP("project", "p", "", "Project ID (required)")

	secretsCmd.AddCommand(secretsListCmd, secretsSetCmd)
	appsCmd.AddCommand(secretsCmd)
}
