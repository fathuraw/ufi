package cmd

import (
	"fmt"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/fathuraw/ufi/internal/unifi"
	"github.com/spf13/cobra"
)

var wifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "Manage WiFi broadcasts",
}

var wifiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List WiFi broadcasts",
	RunE: func(cmd *cobra.Command, args []string) error {
		broadcasts, err := client.ListWiFiBroadcasts(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(broadcasts)
		}
		t := output.NewTable("ID", "NAME", "TYPE", "SECURITY", "BANDS", "ENABLED")
		for _, b := range broadcasts {
			security := ""
			if b.SecurityConfiguration != nil {
				security = b.SecurityConfiguration.Type
			}
			bands := fmt.Sprintf("%v", b.BroadcastingFrequenciesGHz)
			t.AddRow(b.ID, b.Name, b.Type, security, bands, b.Enabled)
		}
		t.Flush()
		return nil
	},
}

var wifiGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get WiFi broadcast details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := client.GetWiFiBroadcast(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(b)
		}
		t := output.NewTable("FIELD", "VALUE")
		t.AddRow("ID", b.ID)
		t.AddRow("Name", b.Name)
		t.AddRow("Type", b.Type)
		t.AddRow("Enabled", b.Enabled)
		if b.SecurityConfiguration != nil {
			t.AddRow("Security", b.SecurityConfiguration.Type)
		}
		if b.Network != nil {
			t.AddRow("Network Type", b.Network.Type)
		}
		t.AddRow("Bands (GHz)", fmt.Sprintf("%v", b.BroadcastingFrequenciesGHz))
		t.Flush()
		return nil
	},
}

var wifiCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a WiFi broadcast",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		ssid, _ := cmd.Flags().GetString("ssid")
		security, _ := cmd.Flags().GetString("security")
		password, _ := cmd.Flags().GetString("password")
		networkID, _ := cmd.Flags().GetString("network-id")
		band, _ := cmd.Flags().GetString("band")
		enabled, _ := cmd.Flags().GetBool("enabled")
		guest, _ := cmd.Flags().GetBool("guest")

		if ssid == "" {
			ssid = name
		}

		req := unifi.WiFiCreateRequest{
			Name:      name,
			SSID:      ssid,
			Security:  security,
			Password:  password,
			NetworkID: networkID,
			Band:      band,
			Enabled:   enabled,
			IsGuest:   guest,
		}
		b, err := client.CreateWiFiBroadcast(req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(b)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "WiFi broadcast %s created (ID: %s).\n", b.Name, b.ID)
		return nil
	},
}

var wifiUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a WiFi broadcast",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		ssid, _ := cmd.Flags().GetString("ssid")
		security, _ := cmd.Flags().GetString("security")
		password, _ := cmd.Flags().GetString("password")
		networkID, _ := cmd.Flags().GetString("network-id")
		band, _ := cmd.Flags().GetString("band")
		enabled, _ := cmd.Flags().GetBool("enabled")
		guest, _ := cmd.Flags().GetBool("guest")

		req := unifi.WiFiCreateRequest{
			Name:      name,
			SSID:      ssid,
			Security:  security,
			Password:  password,
			NetworkID: networkID,
			Band:      band,
			Enabled:   enabled,
			IsGuest:   guest,
		}
		b, err := client.UpdateWiFiBroadcast(args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(b)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "WiFi broadcast %s updated.\n", args[0])
		return nil
	},
}

var wifiDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a WiFi broadcast",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteWiFiBroadcast(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "WiFi broadcast %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	addListFlags(wifiListCmd)

	wifiCreateCmd.Flags().String("name", "", "broadcast name (required)")
	wifiCreateCmd.Flags().String("ssid", "", "SSID (defaults to name)")
	wifiCreateCmd.Flags().String("security", "", "security type")
	wifiCreateCmd.Flags().String("password", "", "WiFi password")
	wifiCreateCmd.Flags().String("network-id", "", "network ID")
	wifiCreateCmd.Flags().String("band", "", "band (2g, 5g, both)")
	wifiCreateCmd.Flags().Bool("enabled", true, "enable broadcast")
	wifiCreateCmd.Flags().Bool("guest", false, "guest network")
	wifiCreateCmd.MarkFlagRequired("name")

	wifiUpdateCmd.Flags().String("name", "", "broadcast name")
	wifiUpdateCmd.Flags().String("ssid", "", "SSID")
	wifiUpdateCmd.Flags().String("security", "", "security type")
	wifiUpdateCmd.Flags().String("password", "", "WiFi password")
	wifiUpdateCmd.Flags().String("network-id", "", "network ID")
	wifiUpdateCmd.Flags().String("band", "", "band")
	wifiUpdateCmd.Flags().Bool("enabled", true, "enable broadcast")
	wifiUpdateCmd.Flags().Bool("guest", false, "guest network")

	wifiCmd.AddCommand(wifiListCmd)
	wifiCmd.AddCommand(wifiGetCmd)
	wifiCmd.AddCommand(wifiCreateCmd)
	wifiCmd.AddCommand(wifiUpdateCmd)
	wifiCmd.AddCommand(wifiDeleteCmd)

	rootCmd.AddCommand(wifiCmd)
}
