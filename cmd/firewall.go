package cmd

import (
	"fmt"
	"strings"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/fathuraw/ufi/internal/unifi"
	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Manage firewall policies and zones",
}

// --- Policies ---

var fwPolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage firewall policies",
}

var fwPolicyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewall policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		policies, err := client.ListFirewallPolicies(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(policies)
		}
		t := output.NewTable("INDEX", "NAME", "ACTION", "ENABLED")
		for _, p := range policies {
			t.AddRow(p.Index, p.Name, p.Action.Type, p.Enabled)
		}
		t.Flush()
		return nil
	},
}

var fwPolicyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a firewall policy",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		action, _ := cmd.Flags().GetString("action")
		protocol, _ := cmd.Flags().GetString("protocol")
		srcZone, _ := cmd.Flags().GetString("source-zone")
		dstZone, _ := cmd.Flags().GetString("dest-zone")
		desc, _ := cmd.Flags().GetString("description")
		enabled, _ := cmd.Flags().GetBool("enabled")

		req := unifi.FirewallPolicyCreateRequest{
			Name:              name,
			Action:            action,
			Protocol:          protocol,
			SourceZoneID:      srcZone,
			DestinationZoneID: dstZone,
			Description:       desc,
			Enabled:           enabled,
		}
		p, err := client.CreateFirewallPolicy(req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(p)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Firewall policy %s created.\n", p.Name)
		return nil
	},
}

var fwPolicyUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a firewall policy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		action, _ := cmd.Flags().GetString("action")
		protocol, _ := cmd.Flags().GetString("protocol")
		srcZone, _ := cmd.Flags().GetString("source-zone")
		dstZone, _ := cmd.Flags().GetString("dest-zone")
		desc, _ := cmd.Flags().GetString("description")
		enabled, _ := cmd.Flags().GetBool("enabled")

		req := unifi.FirewallPolicyCreateRequest{
			Name:              name,
			Action:            action,
			Protocol:          protocol,
			SourceZoneID:      srcZone,
			DestinationZoneID: dstZone,
			Description:       desc,
			Enabled:           enabled,
		}
		p, err := client.UpdateFirewallPolicy(args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(p)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Firewall policy %s updated.\n", args[0])
		return nil
	},
}

var fwPolicyDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a firewall policy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteFirewallPolicy(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Firewall policy %s deleted.\n", args[0])
		return nil
	},
}

var fwPolicyReorderCmd = &cobra.Command{
	Use:   "reorder <id1,id2,...>",
	Short: "Reorder firewall policies",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := strings.Split(args[0], ",")
		if err := client.ReorderFirewallPolicies(ids); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Firewall policies reordered.")
		return nil
	},
}

// --- Zones ---

var fwZoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage firewall zones",
}

var fwZoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewall zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		zones, err := client.ListFirewallZones(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(zones)
		}
		t := output.NewTable("ID", "NAME", "NETWORKS")
		for _, z := range zones {
			t.AddRow(z.ID, z.Name, strings.Join(z.NetworkIDs, ", "))
		}
		t.Flush()
		return nil
	},
}

var fwZoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a firewall zone",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		networkIDs, _ := cmd.Flags().GetStringSlice("network-ids")

		req := unifi.FirewallZoneCreateRequest{
			Name:       name,
			NetworkIDs: networkIDs,
		}
		z, err := client.CreateFirewallZone(req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(z)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Firewall zone %s created (ID: %s).\n", z.Name, z.ID)
		return nil
	},
}

var fwZoneUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a firewall zone",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		networkIDs, _ := cmd.Flags().GetStringSlice("network-ids")

		req := unifi.FirewallZoneCreateRequest{
			Name:       name,
			NetworkIDs: networkIDs,
		}
		z, err := client.UpdateFirewallZone(args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(z)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Firewall zone %s updated.\n", args[0])
		return nil
	},
}

var fwZoneDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a firewall zone",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteFirewallZone(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Firewall zone %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	addListFlags(fwPolicyListCmd)
	addListFlags(fwZoneListCmd)

	fwPolicyCreateCmd.Flags().String("name", "", "policy name (required)")
	fwPolicyCreateCmd.Flags().String("action", "", "action (ALLOW, DENY, REJECT)")
	fwPolicyCreateCmd.Flags().String("protocol", "", "protocol")
	fwPolicyCreateCmd.Flags().String("source-zone", "", "source zone ID")
	fwPolicyCreateCmd.Flags().String("dest-zone", "", "destination zone ID")
	fwPolicyCreateCmd.Flags().String("description", "", "description")
	fwPolicyCreateCmd.Flags().Bool("enabled", true, "enable policy")
	fwPolicyCreateCmd.MarkFlagRequired("name")
	fwPolicyCreateCmd.MarkFlagRequired("action")

	fwPolicyUpdateCmd.Flags().String("name", "", "policy name")
	fwPolicyUpdateCmd.Flags().String("action", "", "action")
	fwPolicyUpdateCmd.Flags().String("protocol", "", "protocol")
	fwPolicyUpdateCmd.Flags().String("source-zone", "", "source zone ID")
	fwPolicyUpdateCmd.Flags().String("dest-zone", "", "destination zone ID")
	fwPolicyUpdateCmd.Flags().String("description", "", "description")
	fwPolicyUpdateCmd.Flags().Bool("enabled", true, "enable policy")

	fwZoneCreateCmd.Flags().String("name", "", "zone name (required)")
	fwZoneCreateCmd.Flags().StringSlice("network-ids", nil, "network IDs")
	fwZoneCreateCmd.MarkFlagRequired("name")

	fwZoneUpdateCmd.Flags().String("name", "", "zone name")
	fwZoneUpdateCmd.Flags().StringSlice("network-ids", nil, "network IDs")

	fwPolicyCmd.AddCommand(fwPolicyListCmd)
	fwPolicyCmd.AddCommand(fwPolicyCreateCmd)
	fwPolicyCmd.AddCommand(fwPolicyUpdateCmd)
	fwPolicyCmd.AddCommand(fwPolicyDeleteCmd)
	fwPolicyCmd.AddCommand(fwPolicyReorderCmd)

	fwZoneCmd.AddCommand(fwZoneListCmd)
	fwZoneCmd.AddCommand(fwZoneCreateCmd)
	fwZoneCmd.AddCommand(fwZoneUpdateCmd)
	fwZoneCmd.AddCommand(fwZoneDeleteCmd)

	firewallCmd.AddCommand(fwPolicyCmd)
	firewallCmd.AddCommand(fwZoneCmd)

	rootCmd.AddCommand(firewallCmd)
}
