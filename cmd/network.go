package cmd

import (
	"fmt"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/fathuraw/ufi/internal/unifi"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all networks",
	RunE: func(cmd *cobra.Command, args []string) error {
		networks, err := client.ListNetworks(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(networks)
		}
		t := output.NewTable("ID", "NAME", "VLAN", "MGMT", "ENABLED", "DEFAULT")
		for _, n := range networks {
			t.AddRow(n.ID, n.Name, n.VlanID, n.Management, n.Enabled, n.Default)
		}
		t.Flush()
		return nil
	},
}

var networkGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get network details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := client.GetNetwork(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(n)
		}
		t := output.NewTable("FIELD", "VALUE")
		t.AddRow("ID", n.ID)
		t.AddRow("Name", n.Name)
		t.AddRow("VLAN ID", n.VlanID)
		t.AddRow("Management", n.Management)
		t.AddRow("Enabled", n.Enabled)
		t.AddRow("Default", n.Default)
		t.AddRow("Zone ID", n.ZoneID)
		if n.Metadata != nil {
			t.AddRow("Origin", n.Metadata.Origin)
		}
		t.Flush()
		return nil
	},
}

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a network",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		purpose, _ := cmd.Flags().GetString("purpose")
		vlanID, _ := cmd.Flags().GetInt("vlan-id")
		subnet, _ := cmd.Flags().GetString("subnet")
		dhcpMode, _ := cmd.Flags().GetString("dhcp-mode")

		req := unifi.NetworkCreateRequest{
			Name:     name,
			Purpose:  purpose,
			VlanID:   vlanID,
			Subnet:   subnet,
			DHCPMode: dhcpMode,
		}
		n, err := client.CreateNetwork(req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(n)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Network %s created (ID: %s).\n", n.Name, n.ID)
		return nil
	},
}

var networkUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a network",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		purpose, _ := cmd.Flags().GetString("purpose")
		vlanID, _ := cmd.Flags().GetInt("vlan-id")
		subnet, _ := cmd.Flags().GetString("subnet")
		dhcpMode, _ := cmd.Flags().GetString("dhcp-mode")

		req := unifi.NetworkCreateRequest{
			Name:     name,
			Purpose:  purpose,
			VlanID:   vlanID,
			Subnet:   subnet,
			DHCPMode: dhcpMode,
		}
		n, err := client.UpdateNetwork(args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(n)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Network %s updated.\n", args[0])
		return nil
	},
}

var networkDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a network",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteNetwork(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Network %s deleted.\n", args[0])
		return nil
	},
}

var networkRefsCmd = &cobra.Command{
	Use:   "refs <id>",
	Short: "Show network dependency references",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refs, err := client.GetNetworkReferences(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(refs)
		}
		if len(refs) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No references found.")
			return nil
		}
		t := output.NewTable("TYPE", "ID", "NAME")
		for _, r := range refs {
			t.AddRow(r.Type, r.ID, r.Name)
		}
		t.Flush()
		return nil
	},
}

func init() {
	addListFlags(networkListCmd)

	networkCreateCmd.Flags().String("name", "", "network name (required)")
	networkCreateCmd.Flags().String("purpose", "", "network purpose (corporate, guest, etc.)")
	networkCreateCmd.Flags().Int("vlan-id", 0, "VLAN ID")
	networkCreateCmd.Flags().String("subnet", "", "subnet CIDR")
	networkCreateCmd.Flags().String("dhcp-mode", "", "DHCP mode")
	networkCreateCmd.MarkFlagRequired("name")

	networkUpdateCmd.Flags().String("name", "", "network name")
	networkUpdateCmd.Flags().String("purpose", "", "network purpose")
	networkUpdateCmd.Flags().Int("vlan-id", 0, "VLAN ID")
	networkUpdateCmd.Flags().String("subnet", "", "subnet CIDR")
	networkUpdateCmd.Flags().String("dhcp-mode", "", "DHCP mode")

	networkCmd.AddCommand(networkListCmd)
	networkCmd.AddCommand(networkGetCmd)
	networkCmd.AddCommand(networkCreateCmd)
	networkCmd.AddCommand(networkUpdateCmd)
	networkCmd.AddCommand(networkDeleteCmd)
	networkCmd.AddCommand(networkRefsCmd)

	rootCmd.AddCommand(networkCmd)
}
