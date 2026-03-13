package cmd

import (
	"github.com/fathuraw/ufi/internal/output"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Manage UniFi sites",
}

var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sites",
	RunE: func(cmd *cobra.Command, args []string) error {
		sites, err := client.ListSites()
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(sites)
		}
		t := output.NewTable("ID", "NAME", "REF")
		for _, s := range sites {
			t.AddRow(s.ID, s.Name, s.InternalReference)
		}
		t.Flush()
		return nil
	},
}

func init() {
	siteCmd.AddCommand(siteListCmd)
	rootCmd.AddCommand(siteCmd)
}
