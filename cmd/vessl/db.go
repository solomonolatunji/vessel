package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"vessl.dev/vessl/internal/models"
)

var dbCmd = &cobra.Command{
	Use:     "db",
	Aliases: []string{"database"},
	Short:   "Manage databases",
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List databases",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")

		client := getClient()
		dbs, err := client.ListDatabases(projectID)
		if err != nil {
			fmt.Printf("Error listing databases: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tENGINE\tSTATUS")
		for _, db := range dbs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", db.ID, db.Name, db.Engine, db.Status)
		}
		w.Flush()
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a database",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		envID, _ := cmd.Flags().GetString("environment")
		name, _ := cmd.Flags().GetString("name")
		engine, _ := cmd.Flags().GetString("engine")

		if projectID == "" || envID == "" || name == "" || engine == "" {
			fmt.Println("Error: --project, --environment, --name, and --engine flags are required")
			os.Exit(1)
		}

		client := getClient()
		req := &models.CreateDatabaseRequest{
			ProjectID:     projectID,
			EnvironmentID: envID,
			Name:          name,
			Engine:        models.DatabaseEngine(engine),
		}

		created, err := client.CreateDatabase(req)
		if err != nil {
			fmt.Printf("Error creating database: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Database %s created successfully with ID: %s\n", created.Name, created.ID)
	},
}

var dbDestroyCmd = &cobra.Command{
	Use:   "destroy [id]",
	Short: "Destroy a database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		if err := client.DeleteDatabase(args[0]); err != nil {
			fmt.Printf("Error destroying database: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Database %s destroyed successfully\n", args[0])
	},
}

var dbImportCmd = &cobra.Command{
	Use:   "import [id]",
	Short: "Import database data from a remote URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourceURL, _ := cmd.Flags().GetString("source")
		if sourceURL == "" {
			fmt.Println("Error: --source flag is required")
			os.Exit(1)
		}

		client := getClient()
		req := &models.ImportDatabaseRequest{SourceURL: sourceURL}
		if err := client.ImportDatabase(args[0], req); err != nil {
			fmt.Printf("Error importing database data: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Data import initiated for database %s\n", args[0])
	},
}

func init() {
	dbListCmd.Flags().StringP("project", "p", "", "Project ID (optional)")

	dbCreateCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	dbCreateCmd.Flags().StringP("environment", "e", "", "Environment ID (required)")
	dbCreateCmd.Flags().StringP("name", "n", "", "Database name (required)")
	dbCreateCmd.Flags().String("engine", "postgres", "Database engine (e.g. postgres, mysql)")

	dbImportCmd.Flags().StringP("source", "s", "", "Source public connection URL (required)")

	dbCmd.AddCommand(dbListCmd, dbCreateCmd, dbDestroyCmd, dbImportCmd)
	rootCmd.AddCommand(dbCmd)
}
