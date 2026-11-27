package gateway

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strings"
)

type WorkloadKeyConfig struct {
	KeyID string
	PublicKey string
	Issuer string
	Subject string
	Scopes []string
	AllowDelegation bool
}

// NewEd25519Credentials validates configuration before serving requests.
// Public keys use unpadded base64url. Private keys are never loaded here.
func NewEd25519Credentials(config []WorkloadKeyConfig) (Ed25519Credentials, error) {
	keys := make(map[string]TrustedWorkloadKey, len(config))
	if len(config) == 0 {
		return Ed25519Credentials{}, fmt.Errorf("no trusted workload keys configured")
	}
	for _, item := range config {
		if item.KeyID == "" || len(item.KeyID) > 128 || strings.ContainsAny(item.KeyID, ". \t\r\n,") || strings.TrimSpace(item.Issuer) == "" || strings.TrimSpace(item.Subject) == "" || len(item.Scopes) == 0 {
			return Ed25519Credentials{}, fmt.Errorf("invalid workload key configuration")
		}
		if _, exists := keys[item.KeyID]; exists {
			return Ed25519Credentials{}, fmt.Errorf("duplicate workload key ID")
		}
		key, err := base64.RawURLEncoding.DecodeString(item.PublicKey)
		if err != nil || len(key) != ed25519.PublicKeySize {
			return Ed25519Credentials{}, fmt.Errorf("invalid workload public key")
		}
		seen := make(map[string]bool)
		for _, scope := range item.Scopes {
			if strings.TrimSpace(scope) == "" || seen[scope] {
				return Ed25519Credentials{}, fmt.Errorf("invalid or duplicate workload scope")
			}
			seen[scope] = true
		}
		keys[item.KeyID] = TrustedWorkloadKey{PublicKey: ed25519.PublicKey(key), Issuer: item.Issuer, Subject: item.Subject, Scopes: append([]string(nil), item.Scopes...), AllowDelegation: item.AllowDelegation}
	}
	return Ed25519Credentials{Keys: keys}, nil
}
