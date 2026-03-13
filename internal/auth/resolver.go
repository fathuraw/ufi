package auth

import (
	"fmt"
	"os"
)

// ResolveAPIKey resolves the API key using the priority chain:
// 1. UFI_API_KEY env var
// 2. flagValue (--api-key flag)
// 3. System keyring
// 4. Encrypted file (prompts for passphrase via promptFn)
func ResolveAPIKey(flagValue string, promptFn func(prompt string) (string, error)) (string, error) {
	// 1. Environment variable
	if key := os.Getenv("UFI_API_KEY"); key != "" {
		return key, nil
	}

	// 2. CLI flag
	if flagValue != "" {
		return flagValue, nil
	}

	// 3. System keyring
	if KeyringAvailable() {
		key, err := GetKeyring()
		if err == nil && key != "" {
			return key, nil
		}
	}

	// 4. Encrypted file
	if EncFileExists() {
		if promptFn == nil {
			return "", fmt.Errorf("encrypted credentials found but no passphrase prompt available")
		}
		passphrase, err := promptFn("Enter passphrase to decrypt credentials: ")
		if err != nil {
			return "", fmt.Errorf("read passphrase: %w", err)
		}
		return DecryptFromFile(passphrase)
	}

	return "", fmt.Errorf("no API key found. Run 'ufi login' to authenticate, or set UFI_API_KEY")
}
