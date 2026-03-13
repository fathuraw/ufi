package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

const (
	credDir  = ".config/ufi"
	credFile = "credentials.enc"
	saltSize = 16
	keyLen   = 32 // AES-256
)

func credPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, credDir, credFile), nil
}

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, keyLen)
}

// EncryptAndStore encrypts the API key with a passphrase and writes to ~/.ufi/credentials.enc.
func EncryptAndStore(apiKey, passphrase string) error {
	path, err := credPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create credentials directory: %w", err)
	}

	// Generate salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)

	// File format: salt (16 bytes) + nonce+ciphertext
	data := append(salt, ciphertext...)

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write credentials file: %w", err)
	}
	return nil
}

// DecryptFromFile reads and decrypts the API key from ~/.ufi/credentials.enc.
func DecryptFromFile(passphrase string) (string, error) {
	path, err := credPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read credentials file: %w", err)
	}

	if len(data) < saltSize {
		return "", fmt.Errorf("credentials file is corrupted")
	}

	salt := data[:saltSize]
	ciphertext := data[saltSize:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("credentials file is corrupted")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed (wrong passphrase?): %w", err)
	}

	return string(plaintext), nil
}

// DeleteEncFile removes the encrypted credentials file.
func DeleteEncFile() error {
	path, err := credPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// EncFileExists checks if the encrypted credentials file exists.
func EncFileExists() bool {
	path, err := credPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
