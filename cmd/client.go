package cmd

import (
	"fmt"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Manage network clients",
}

var clientListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all clients",
	RunE: func(cmd *cobra.Command, args []string) error {
		clients, err := client.ListClients(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(clients)
		}
		t := output.NewTable("ID", "NAME", "TYPE", "MAC", "IP", "CONNECTED AT")
		for _, c := range clients {
			t.AddRow(c.ID, c.Name, c.Type, c.MACAddress, c.IPAddress, c.ConnectedAt)
		}
		t.Flush()
		return nil
	},
}

var clientGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get client details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.GetClient(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(c)
		}
		t := output.NewTable("FIELD", "VALUE")
		t.AddRow("ID", c.ID)
		t.AddRow("Name", c.Name)
		t.AddRow("Type", c.Type)
		t.AddRow("MAC", c.MACAddress)
		t.AddRow("IP", c.IPAddress)
		t.AddRow("Connected At", c.ConnectedAt)
		t.AddRow("Uplink Device", c.UplinkDeviceID)
		if c.Access != nil {
			t.AddRow("Access", c.Access.Type)
		}
		t.Flush()
		return nil
	},
}

var clientBlockCmd = &cobra.Command{
	Use:   "block <id>",
	Short: "Block a client",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.BlockClient(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Client %s blocked.\n", args[0])
		return nil
	},
}

var clientUnblockCmd = &cobra.Command{
	Use:   "unblock <id>",
	Short: "Unblock a client",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.UnblockClient(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Client %s unblocked.\n", args[0])
		return nil
	},
}

func init() {
	addListFlags(clientListCmd)

	clientCmd.AddCommand(clientListCmd)
	clientCmd.AddCommand(clientGetCmd)
	clientCmd.AddCommand(clientBlockCmd)
	clientCmd.AddCommand(clientUnblockCmd)

	rootCmd.AddCommand(clientCmd)
}
