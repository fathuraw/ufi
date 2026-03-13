package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manage UniFi devices",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		devices, err := client.ListDevices(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(devices)
		}
		t := output.NewTable("ID", "NAME", "MODEL", "MAC", "IP", "STATE", "FW")
		for _, d := range devices {
			t.AddRow(d.ID, d.Name, d.Model, d.MACAddress, d.IPAddress, d.State, d.FirmwareVersion)
		}
		t.Flush()
		return nil
	},
}

var deviceGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get device details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := client.GetDevice(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(d)
		}
		t := output.NewTable("FIELD", "VALUE")
		t.AddRow("ID", d.ID)
		t.AddRow("Name", d.Name)
		t.AddRow("Model", d.Model)
		t.AddRow("MAC", d.MACAddress)
		t.AddRow("IP", d.IPAddress)
		t.AddRow("State", d.State)
		t.AddRow("Firmware", d.FirmwareVersion)
		t.AddRow("Updatable", d.FirmwareUpdatable)
		t.AddRow("Adopted At", d.AdoptedAt)
		if d.Interfaces != nil {
			for _, p := range d.Interfaces.Ports {
				poe := ""
				if p.PoE != nil {
					poe = fmt.Sprintf(" (PoE: %s)", p.PoE.State)
				}
				t.AddRow(fmt.Sprintf("Port %d", p.Idx), fmt.Sprintf("%s %dMbps %s%s", p.Connector, p.SpeedMbps, p.State, poe))
			}
			for _, r := range d.Interfaces.Radios {
				t.AddRow(fmt.Sprintf("Radio %.1fGHz", r.FrequencyGHz), fmt.Sprintf("%s ch%d %dMHz", r.WLANStandard, r.Channel, r.ChannelWidthMHz))
			}
		}
		t.Flush()
		return nil
	},
}

var deviceAdoptCmd = &cobra.Command{
	Use:   "adopt <mac>",
	Short: "Adopt a device by MAC address",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.AdoptDevice(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Device %s adoption initiated.\n", args[0])
		return nil
	},
}

var deviceRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.RemoveDevice(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Device %s removed.\n", args[0])
		return nil
	},
}

var deviceRestartCmd = &cobra.Command{
	Use:   "restart <id>",
	Short: "Restart a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.RestartDevice(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Device %s restart initiated.\n", args[0])
		return nil
	},
}

var deviceStatsCmd = &cobra.Command{
	Use:   "stats <id>",
	Short: "Show device statistics",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := client.GetDeviceStatistics(args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(stats)
		}
		t := output.NewTable("METRIC", "VALUE")
		t.AddRow("Uptime", formatDuration(stats.UptimeSec))
		t.AddRow("CPU", fmt.Sprintf("%.1f%%", stats.CPUUtilizationPct))
		t.AddRow("Memory", fmt.Sprintf("%.1f%%", stats.MemoryUtilizationPct))
		t.AddRow("Load (1/5/15m)", fmt.Sprintf("%.2f / %.2f / %.2f", stats.LoadAverage1Min, stats.LoadAverage5Min, stats.LoadAverage15Min))
		if stats.Uplink != nil {
			t.AddRow("Uplink TX", formatRate(stats.Uplink.TxRateBps))
			t.AddRow("Uplink RX", formatRate(stats.Uplink.RxRateBps))
		}
		t.AddRow("Last Heartbeat", stats.LastHeartbeatAt)
		if stats.Interfaces != nil {
			for _, r := range stats.Interfaces.Radios {
				t.AddRow(fmt.Sprintf("Radio %.1fGHz TX Retries", r.FrequencyGHz), fmt.Sprintf("%.1f%%", r.TxRetriesPct))
			}
		}
		t.Flush()
		return nil
	},
}

var devicePortCmd = &cobra.Command{
	Use:   "port",
	Short: "Device port operations",
}

var devicePortPowerCycleCmd = &cobra.Command{
	Use:   "power-cycle <deviceId> <portIdx>",
	Short: "Power cycle a device port",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		portIdx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid port index: %s", args[1])
		}
		if err := client.PowerCyclePort(args[0], portIdx); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Port %d on device %s power cycle initiated.\n", portIdx, args[0])
		return nil
	},
}

func formatDuration(seconds int64) string {
	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatRate(bps int64) string {
	units := []string{"bps", "Kbps", "Mbps", "Gbps"}
	f := float64(bps)
	i := 0
	for f >= 1000 && i < len(units)-1 {
		f /= 1000
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d %s", bps, units[0])
	}
	s := fmt.Sprintf("%.1f", f)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	return s + " " + units[i]
}

func init() {
	addListFlags(deviceListCmd)

	devicePortCmd.AddCommand(devicePortPowerCycleCmd)

	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceGetCmd)
	deviceCmd.AddCommand(deviceAdoptCmd)
	deviceCmd.AddCommand(deviceRemoveCmd)
	deviceCmd.AddCommand(deviceRestartCmd)
	deviceCmd.AddCommand(deviceStatsCmd)
	deviceCmd.AddCommand(devicePortCmd)

	rootCmd.AddCommand(deviceCmd)
}
