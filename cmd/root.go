package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fathuraw/ufi/internal/auth"
	"github.com/fathuraw/ufi/internal/unifi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	jsonOut  bool
	insecure bool
	apiKey   string

	client *unifi.Client
)

var rootCmd = &cobra.Command{
	Use:   "ufi",
	Short: "CLI for UniFi OS Network controllers",
	Long:  "ufi is a command-line interface for managing UniFi OS network controllers via the Integration API.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip auth for login, logout, version, help, completion
		skip := map[string]bool{
			"login": true, "logout": true, "version": true,
			"help": true, "completion": true,
		}
		if skip[cmd.Name()] {
			return nil
		}
		// Also skip if it's just the root command (shows help)
		if cmd == cmd.Root() {
			return nil
		}

		host := viper.GetString("host")
		if host == "" {
			return fmt.Errorf("no host configured. Run 'ufi login' or set --host")
		}

		siteID := viper.GetString("site")
		if siteID == "" {
			siteID = "default"
		}

		key, err := auth.ResolveAPIKey(apiKey, promptPassphrase)
		if err != nil {
			return err
		}

		client = unifi.NewClient(host, key, siteID, insecure)
		return nil
	},
}

func promptPassphrase(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.config/ufi/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides keyring)")
	rootCmd.PersistentFlags().String("host", "", "controller URL (e.g. https://192.168.1.1)")
	rootCmd.PersistentFlags().String("site", "", "site ID (default: from config)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "skip TLS certificate verification")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("site", rootCmd.PersistentFlags().Lookup("site"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		configDir := filepath.Join(home, ".config", "ufi")
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("UFI")
	viper.AutomaticEnv()
	viper.ReadInConfig() // ignore error if no config file
}

// listParams returns ListParams from common flags on the command.
func listParams(cmd *cobra.Command) unifi.ListParams {
	limit, _ := cmd.Flags().GetInt("limit")
	offset, _ := cmd.Flags().GetInt("offset")
	filter, _ := cmd.Flags().GetString("filter")
	return unifi.ListParams{Limit: limit, Offset: offset, Filter: filter}
}

// addListFlags adds --limit, --offset, --filter to a command.
func addListFlags(cmd *cobra.Command) {
	cmd.Flags().Int("limit", 0, "maximum number of results")
	cmd.Flags().Int("offset", 0, "result offset for pagination")
	cmd.Flags().String("filter", "", "filter expression (UniFi filter DSL)")
}
