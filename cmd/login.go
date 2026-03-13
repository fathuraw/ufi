package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fathuraw/ufi/internal/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Store API key for a UniFi controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Prompt for host
		host := viper.GetString("host")
		if host == "" {
			fmt.Fprint(os.Stderr, "Controller URL (e.g. https://192.168.1.1): ")
			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			host = strings.TrimSpace(line)
		}
		if host == "" {
			return fmt.Errorf("host is required")
		}

		// Prompt for API key
		fmt.Fprint(os.Stderr, "API Key: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		key := strings.TrimSpace(line)
		if key == "" {
			return fmt.Errorf("API key is required")
		}

		// Store host in config
		viper.Set("host", host)
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath := home + "/.ufi.yaml"
		if err := viper.WriteConfigAs(configPath); err != nil {
			// If file exists, safe-write
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("write config: %w", err)
			}
		}
		fmt.Fprintf(os.Stderr, "Host saved to %s\n", configPath)

		// Try keyring first
		if auth.KeyringAvailable() {
			if err := auth.SetKeyring(key); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: keyring store failed: %v\n", err)
			} else {
				fmt.Fprintln(os.Stderr, "API key stored in system keyring.")
				return nil
			}
		}

		// Fall back to encrypted file
		fmt.Fprintln(os.Stderr, "System keyring not available, using encrypted file.")
		fmt.Fprint(os.Stderr, "Choose a passphrase for encrypting credentials: ")
		passLine, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		passphrase := strings.TrimSpace(passLine)
		if passphrase == "" {
			return fmt.Errorf("passphrase is required for encrypted storage")
		}

		if err := auth.EncryptAndStore(key, passphrase); err != nil {
			return fmt.Errorf("store credentials: %w", err)
		}
		fmt.Fprintln(os.Stderr, "API key stored in encrypted file (~/.ufi/credentials.enc).")
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		removed := false

		if auth.KeyringAvailable() {
			if err := auth.DeleteKeyring(); err == nil {
				fmt.Fprintln(os.Stderr, "API key removed from system keyring.")
				removed = true
			}
		}

		if auth.EncFileExists() {
			if err := auth.DeleteEncFile(); err == nil {
				fmt.Fprintln(os.Stderr, "Encrypted credentials file removed.")
				removed = true
			}
		}

		if !removed {
			fmt.Fprintln(os.Stderr, "No stored credentials found.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}
