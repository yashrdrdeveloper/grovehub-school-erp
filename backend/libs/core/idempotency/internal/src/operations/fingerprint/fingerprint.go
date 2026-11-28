package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"grovehub/libs/core/idempotency/internal/src/contracts"
)

// Fingerprint hashes the JSON representation without changing its semantics.
func Fingerprint(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("%w: %v", contracts.ErrFingerprint, err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}
