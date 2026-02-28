package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/nexus-id/backend/internal/config"
	"go.uber.org/zap"
)

var (
	ErrNoPrivateKey = errors.New("no private key found")
	ErrNoPublicKey  = errors.New("no public key found")
	ErrInvalidKey   = errors.New("invalid key format")
)

// KeyManager handles RSA key pair generation, loading, and storage
type KeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	config     *config.JWTConfig
	logger     *zap.Logger
	keyID      string // Key ID for JWKS
}

// NewKeyManager creates a new KeyManager
func NewKeyManager(cfg *config.JWTConfig, logger *zap.Logger) (*KeyManager, error) {
	km := &KeyManager{
		config: cfg,
		logger: logger,
		keyID:  "key1", // Simple key ID, can be enhanced with rotation
	}

	// Ensure keys directory exists
	if err := os.MkdirAll(filepath.Dir(cfg.PrivateKeyPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	// Try to load existing keys, or generate new ones
	if err := km.loadOrGenerateKeys(); err != nil {
		return nil, err
	}

	logger.Info("JWT keys loaded successfully")
	return km, nil
}

// loadOrGenerateKeys loads existing keys or generates new ones
func (km *KeyManager) loadOrGenerateKeys() error {
	// Try to load existing keys
	privateKey, err := km.loadPrivateKey()
	if err == nil {
		publicKey, err := km.loadPublicKey()
		if err == nil {
			km.privateKey = privateKey
			km.publicKey = publicKey
			km.logger.Info("Loaded existing JWT keys")
			return nil
		}
	}

	// Keys don't exist, generate new ones
	km.logger.Info("Generating new RSA key pair")
	if err := km.generateKeys(); err != nil {
		return err
	}

	return nil
}

// generateKeys generates a new RSA key pair
func (km *KeyManager) generateKeys() error {
	// Generate RSA 2048-bit key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	km.privateKey = privateKey
	km.publicKey = &privateKey.PublicKey

	// Save keys to disk
	if err := km.savePrivateKey(privateKey); err != nil {
		return err
	}

	if err := km.savePublicKey(&privateKey.PublicKey); err != nil {
		return err
	}

	km.logger.Info("Generated and saved new RSA key pair")
	return nil
}

// savePrivateKey saves the private key to disk
func (km *KeyManager) savePrivateKey(key *rsa.PrivateKey) error {
	// Convert to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Write to file with restricted permissions
	if err := os.WriteFile(km.config.PrivateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	return nil
}

// savePublicKey saves the public key to disk
func (km *KeyManager) savePublicKey(key *rsa.PublicKey) error {
	// Convert to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Write to file
	if err := os.WriteFile(km.config.PublicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to save public key: %w", err)
	}

	return nil
}

// loadPrivateKey loads the private key from disk
func (km *KeyManager) loadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(km.config.PrivateKeyPath)
	if err != nil {
		return nil, ErrNoPrivateKey
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, ErrInvalidKey
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// loadPublicKey loads the public key from disk
func (km *KeyManager) loadPublicKey() (*rsa.PublicKey, error) {
	data, err := os.ReadFile(km.config.PublicKeyPath)
	if err != nil {
		return nil, ErrNoPublicKey
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, ErrInvalidKey
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKey
	}

	return rsaPublicKey, nil
}

// GetPrivateKey returns the private key for signing
func (km *KeyManager) GetPrivateKey() *rsa.PrivateKey {
	return km.privateKey
}

// GetPublicKey returns the public key for verification
func (km *KeyManager) GetPublicKey() *rsa.PublicKey {
	return km.publicKey
}

// GetKeyID returns the key ID for JWKS
func (km *KeyManager) GetKeyID() string {
	return km.keyID
}

// GetPublicJWK returns the public key in JWK format
func (km *KeyManager) GetPublicJWK() map[string]interface{} {
	// Get the public key parameters
	n := km.publicKey.N
	e := km.publicKey.E

	return map[string]interface{}{
		"kty": "RSA",
		"kid": km.keyID,
		"n":   n.Text(62), // Base64 URL encoded without padding
		"e":   big.NewInt(int64(e)).Text(62),
		"alg": "RS256",
		"use": "sig",
	}
}

// RotateKeys generates a new key pair (placeholder for key rotation)
func (km *KeyManager) RotateKeys() error {
	km.logger.Info("Key rotation requested")
	// TODO: Implement proper key rotation with multiple active keys
	return km.generateKeys()
}
