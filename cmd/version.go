package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags:
//
//	go build -ldflags "-X github.com/fathuraw/ufi/cmd.Version=1.0.0"
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the ufi version",
	Run: func(cmd *cobra.Command, args []string) {
		v := Version
		if v == "dev" {
			if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
				v = info.Main.Version
			}
		}
		fmt.Printf("ufi %s\n", v)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
