package config

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/operations/credential"

	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strings"
)

// NewEd25519Credentials validates configuration before serving requests.
// Public keys use unpadded base64url. Private keys are never loaded here.
func NewEd25519Credentials(config []contracts.WorkloadKeyConfig) (credential.Ed25519Credentials, error) {
	keys := make(map[string]contracts.TrustedWorkloadKey, len(config))
	if len(config) == 0 {
		return credential.Ed25519Credentials{}, fmt.Errorf("no trusted workload keys configured")
	}
	for _, item := range config {
		if !helpers.ValidKeyID(item.KeyID) || strings.TrimSpace(item.Issuer) == "" || strings.TrimSpace(item.Subject) == "" || len(item.Scopes) == 0 {
			return credential.Ed25519Credentials{}, fmt.Errorf("invalid workload key configuration")
		}
		if _, exists := keys[item.KeyID]; exists {
			return credential.Ed25519Credentials{}, fmt.Errorf("duplicate workload key ID")
		}
		key, err := base64.RawURLEncoding.DecodeString(item.PublicKey)
		if err != nil || len(key) != ed25519.PublicKeySize {
			return credential.Ed25519Credentials{}, fmt.Errorf("invalid workload public key")
		}
		seen := make(map[string]bool)
		for _, scope := range item.Scopes {
			if strings.TrimSpace(scope) == "" || seen[scope] {
				return credential.Ed25519Credentials{}, fmt.Errorf("invalid or duplicate workload scope")
			}
			seen[scope] = true
		}
		keys[item.KeyID] = contracts.TrustedWorkloadKey{PublicKey: ed25519.PublicKey(key), Issuer: item.Issuer, Subject: item.Subject, Scopes: append([]string(nil), item.Scopes...), AllowDelegation: item.AllowDelegation}
	}
	return credential.Ed25519Credentials{Keys: keys}, nil
}
