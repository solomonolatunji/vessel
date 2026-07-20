package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"vessl.dev/vessl/internal/models"
)

var backupCmd = &cobra.Command{
	Use:     "backups",
	Aliases: []string{"backup"},
	Short:   "Manage database backups",
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List database backups",
	Run: func(cmd *cobra.Command, args []string) {
		databaseID, _ := cmd.Flags().GetString("database")

		client := getClient()
		backups, err := client.ListBackups(databaseID)
		if err != nil {
			fmt.Printf("Error listing backups: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSCHEDULE\tSTATUS")
		for _, b := range backups {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", b.ID, b.Name, b.Schedule, b.Status)
		}
		w.Flush()
	},
}

var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString("project")
		databaseID, _ := cmd.Flags().GetString("database")
		name, _ := cmd.Flags().GetString("name")
		schedule, _ := cmd.Flags().GetString("schedule")

		if projectID == "" || databaseID == "" || name == "" || schedule == "" {
			fmt.Println("Error: --project, --database, --name, and --schedule flags are required")
			os.Exit(1)
		}

		client := getClient()
		req := &models.BackupConfig{
			DatabaseID: databaseID,
			Name:       name,
			Schedule:   schedule,
		}

		created, err := client.CreateBackup(req)
		if err != nil {
			fmt.Printf("Error creating backup: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Backup %s created successfully with ID: %s\n", created.Name, created.ID)
	},
}

var backupTriggerCmd = &cobra.Command{
	Use:   "trigger [id]",
	Short: "Trigger a database backup manually",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		record, err := client.TriggerBackup(args[0])
		if err != nil {
			fmt.Printf("Error triggering backup: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Backup triggered successfully with record ID: %s\n", record.ID)
	},
}

var backupHistoryCmd = &cobra.Command{
	Use:   "history [id]",
	Short: "List backup records for a backup config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		records, err := client.ListBackupRecords(args[0])
		if err != nil {
			fmt.Printf("Error fetching backup history: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tSTARTED\tCOMPLETED\tSIZE")
		for _, r := range records {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n", r.ID, r.Status, r.StartedAt, r.CompletedAt, r.FileSizeBytes)
		}
		w.Flush()
	},
}

func init() {
	backupListCmd.Flags().StringP("database", "d", "", "Database ID (optional)")

	backupCreateCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	backupCreateCmd.Flags().StringP("database", "d", "", "Database ID (required)")
	backupCreateCmd.Flags().StringP("name", "n", "", "Backup name (required)")
	backupCreateCmd.Flags().StringP("schedule", "s", "", "Backup schedule cron (required)")

	backupCmd.AddCommand(backupListCmd, backupCreateCmd, backupTriggerCmd, backupHistoryCmd)
	dbCmd.AddCommand(backupCmd)
}
