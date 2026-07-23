package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var deploymentsCmd = &cobra.Command{
	Use:     "deployments",
	Aliases: []string{"deployment"},
	Short:   "Manage app deployments",
}

var deploymentsListCmd = &cobra.Command{
	Use:   "list [service_id]",
	Short: "List deployments for a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		deployments, err := client.ListDeployments(args[0])
		if err != nil {
			fmt.Printf("Error listing deployments: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tBRANCH\tCOMMIT\tCREATED")
		for _, d := range deployments {
			createdAt := d.CreatedAt.Format("2006-01-02 15:04:05")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", d.ID, d.Status, d.Branch, d.CommitHash, createdAt)
		}
		w.Flush()
	},
}

func init() {
	deploymentsCmd.AddCommand(deploymentsListCmd)
	appsCmd.AddCommand(deploymentsCmd)
}
