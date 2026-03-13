package cmd

import (
	"fmt"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/fathuraw/ufi/internal/unifi"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS records",
}

var dnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS records",
	RunE: func(cmd *cobra.Command, args []string) error {
		records, err := client.ListDNSRecords(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(records)
		}
		t := output.NewTable("ID", "TYPE", "DOMAIN", "IP", "ENABLED", "TTL")
		for _, r := range records {
			t.AddRow(r.ID, r.Type, r.Domain, r.IPv4Address, r.Enabled, r.TTLSeconds)
		}
		t.Flush()
		return nil
	},
}

var dnsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS record",
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		ip, _ := cmd.Flags().GetString("ip")
		recordType, _ := cmd.Flags().GetString("type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		enabled, _ := cmd.Flags().GetBool("enabled")

		if recordType == "" {
			recordType = "A_RECORD"
		}

		req := unifi.DNSRecordCreateRequest{
			Type:        recordType,
			Domain:      domain,
			IPv4Address: ip,
			TTLSeconds:  ttl,
			Enabled:     enabled,
		}
		r, err := client.CreateDNSRecord(req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(r)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "DNS record %s -> %s created (ID: %s).\n", r.Domain, r.IPv4Address, r.ID)
		return nil
	},
}

var dnsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a DNS record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		ip, _ := cmd.Flags().GetString("ip")
		recordType, _ := cmd.Flags().GetString("type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		enabled, _ := cmd.Flags().GetBool("enabled")

		req := unifi.DNSRecordCreateRequest{
			Type:        recordType,
			Domain:      domain,
			IPv4Address: ip,
			TTLSeconds:  ttl,
			Enabled:     enabled,
		}
		r, err := client.UpdateDNSRecord(args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(r)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "DNS record %s updated.\n", args[0])
		return nil
	},
}

var dnsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a DNS record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteDNSRecord(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "DNS record %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	addListFlags(dnsListCmd)

	dnsCreateCmd.Flags().String("domain", "", "domain name (required)")
	dnsCreateCmd.Flags().String("ip", "", "IPv4 address (required)")
	dnsCreateCmd.Flags().String("type", "A_RECORD", "record type")
	dnsCreateCmd.Flags().Int("ttl", 0, "TTL in seconds")
	dnsCreateCmd.Flags().Bool("enabled", true, "enable record")
	dnsCreateCmd.MarkFlagRequired("domain")
	dnsCreateCmd.MarkFlagRequired("ip")

	dnsUpdateCmd.Flags().String("domain", "", "domain name")
	dnsUpdateCmd.Flags().String("ip", "", "IPv4 address")
	dnsUpdateCmd.Flags().String("type", "", "record type")
	dnsUpdateCmd.Flags().Int("ttl", 0, "TTL in seconds")
	dnsUpdateCmd.Flags().Bool("enabled", true, "enable record")

	dnsCmd.AddCommand(dnsListCmd)
	dnsCmd.AddCommand(dnsCreateCmd)
	dnsCmd.AddCommand(dnsUpdateCmd)
	dnsCmd.AddCommand(dnsDeleteCmd)

	rootCmd.AddCommand(dnsCmd)
}
