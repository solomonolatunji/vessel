package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "codedock",
	Short: "Codedock CLI - Manage your self-hosted Codedock server",
	Long:  `A command line interface to authenticate and deploy applications to your Codedock self-hosted PaaS.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the codedock CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("codedock %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	},
}

func main() {
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
