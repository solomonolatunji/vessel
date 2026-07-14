package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vessl",
	Short: "Vessl CLI - Manage your self-hosted Vessl server",
	Long:  `A command line interface to authenticate and deploy applications to your Vessl self-hosted PaaS.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
