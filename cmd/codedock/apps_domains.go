package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"codedock.run/codedock/internal/models"
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Manage application domains",
}

var domainsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List domains for a service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceID, _ := cmd.Flags().GetString("service")
		if serviceID == "" {
			fmt.Println("Error: --service flag is required")
			os.Exit(1)
		}

		client := getClient()
		domains, err := client.ListDomains(serviceID)
		if err != nil {
			fmt.Printf("Error listing domains: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tDOMAIN_NAME\tPATH_PREFIX\tSSL_STATUS\tREDIRECT_TO")
		for _, domain := range domains {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", domain.ID, domain.DomainName, domain.PathPrefix, domain.SSLCertStatus, domain.RedirectTo)
		}
		w.Flush()
	},
}

var domainsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a domain to a service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceID, _ := cmd.Flags().GetString("service")
		domainName, _ := cmd.Flags().GetString("domain")
		if serviceID == "" || domainName == "" {
			fmt.Println("Error: --service and --domain flags are required")
			os.Exit(1)
		}

		redirectTo, _ := cmd.Flags().GetString("redirect")
		pathPrefix, _ := cmd.Flags().GetString("prefix")

		client := getClient()
		req := &models.DomainConfig{
			DomainName: domainName,
			ServiceID:  serviceID,
			RedirectTo: redirectTo,
			PathPrefix: pathPrefix,
		}

		created, err := client.AddDomain(serviceID, req)
		if err != nil {
			fmt.Printf("Error adding domain: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Domain %s added successfully with ID: %s\n", created.DomainName, created.ID)
	},
}

var domainsRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		if err := client.RemoveDomain(args[0]); err != nil {
			fmt.Printf("Error removing domain: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Domain %s removed successfully\n", args[0])
	},
}

func init() {
	domainsListCmd.Flags().StringP("service", "s", "", "Service ID (required)")

	domainsAddCmd.Flags().StringP("service", "s", "", "Service ID (required)")
	domainsAddCmd.Flags().StringP("domain", "d", "", "Domain name (required)")
	domainsAddCmd.Flags().String("redirect", "", "Redirect To URL")
	domainsAddCmd.Flags().String("prefix", "", "Path Prefix")

	domainsCmd.AddCommand(domainsListCmd, domainsAddCmd, domainsRemoveCmd)
	appsCmd.AddCommand(domainsCmd)
}
