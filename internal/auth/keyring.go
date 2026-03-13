package auth

import "github.com/zalando/go-keyring"

const (
	serviceName = "ufi-cli"
	accountName = "api-key"
)

// KeyringAvailable checks if the system keyring is available.
func KeyringAvailable() bool {
	// Try a no-op get; if the error is "not found" the keyring works.
	// Any other error means no keyring.
	_, err := keyring.Get(serviceName, "__probe__")
	if err == keyring.ErrNotFound {
		return true
	}
	// If we got a value back (unlikely) or nil error, keyring is available.
	if err == nil {
		return true
	}
	return false
}

// SetKeyring stores the API key in the system keyring.
func SetKeyring(apiKey string) error {
	return keyring.Set(serviceName, accountName, apiKey)
}

// GetKeyring retrieves the API key from the system keyring.
func GetKeyring() (string, error) {
	return keyring.Get(serviceName, accountName)
}

// DeleteKeyring removes the API key from the system keyring.
func DeleteKeyring() error {
	return keyring.Delete(serviceName, accountName)
}
