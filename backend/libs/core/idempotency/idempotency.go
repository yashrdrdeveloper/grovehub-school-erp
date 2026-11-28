// Package idempotency exposes request key validation and fingerprinting.
package idempotency

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/constants"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/operations/fingerprint"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/operations/key"
)

const MaxKeyLength = constants.MaxKeyLength

var (
	ErrInvalidKey  = contracts.ErrInvalidKey
	ErrFingerprint = contracts.ErrFingerprint
)

func ValidateKey(value string) (string, error) {
	return key.ValidateKey(value)
}

func Fingerprint(value any) (string, error) {
	return fingerprint.Fingerprint(value)
}
