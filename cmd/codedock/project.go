package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"codedock.dev/codedock/internal/models"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project"},
	Short:   "Manage projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projects, err := client.ListProjects()
		if err != nil {
			fmt.Printf("Error listing projects: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")
		for _, p := range projects {
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, p.Description)
		}
		w.Flush()
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")

		if name == "" {
			fmt.Println("Error: --name flag is required")
			os.Exit(1)
		}

		client := getClient()
		req := &models.CreateProjectRequest{
			Name:        name,
			Description: desc,
		}

		created, err := client.CreateProject(req)
		if err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Project %s created successfully with ID: %s\n", created.Name, created.ID)
	},
}

var projectDestroyCmd = &cobra.Command{
	Use:   "destroy [id]",
	Short: "Destroy a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		if err := client.DeleteProject(args[0]); err != nil {
			fmt.Printf("Error destroying project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Project %s destroyed successfully\n", args[0])
	},
}

func init() {
	projectCreateCmd.Flags().StringP("name", "n", "", "Project name (required)")
	projectCreateCmd.Flags().StringP("description", "d", "", "Project description")

	projectCmd.AddCommand(projectListCmd, projectCreateCmd, projectDestroyCmd)
	rootCmd.AddCommand(projectCmd)
}
