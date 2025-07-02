package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
)

const (
	keyLen    = 32 // AES-256
	saltLen   = 16
	nonceLen  = 12 // GCM standard nonce length
	VaultFile = "vault.enc"
)

type Vault struct {
	Passwords map[string]string `json:"passwords"`
}

// DeriveKey generates a strong encryption key from a password using scrypt.
func DeriveKey(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 32768, 8, 1, keyLen)
}

// Encrypt encrypts data using AES-256-GCM.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-256-GCM.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// GetVaultPath returns the absolute path to the vault file.
func GetVaultPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".griphook", VaultFile), nil
}

// LoadVault loads and decrypts the vault from the file system.
func LoadVault(masterPassword []byte) (*Vault, error) {
	vaultPath, err := GetVaultPath()
	if err != nil {
		return nil, err
	}

	encryptedData, err := os.ReadFile(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("could not read vault file: %w", err)
	}

	if len(encryptedData) < saltLen+nonceLen {
		return nil, fmt.Errorf("encrypted data too short")
	}

	salt := encryptedData[:saltLen]
	encryptedData = encryptedData[saltLen:]

	key, err := DeriveKey(masterPassword, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	decrypted, err := Decrypt(key, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault: %w", err)
	}

	v := &Vault{}
	if err := json.Unmarshal(decrypted, v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault data: %w", err)
	}

	return v, nil
}

// SaveVault encrypts and saves the vault to the file system.
func SaveVault(v *Vault, masterPassword []byte) error {
	vaultPath, err := GetVaultPath()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	vaultDir := filepath.Dir(vaultPath)
	if err := os.MkdirAll(vaultDir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	plaintext, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal vault data: %w", err)
	}

	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	key, err := DeriveKey(masterPassword, salt)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Prepend salt to the encrypted data
	finalData := append(salt, encrypted...)

	if err := os.WriteFile(vaultPath, finalData, 0600); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	return nil
}
