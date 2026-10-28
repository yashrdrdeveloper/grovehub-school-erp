package idempotency

import (
	"grovehub/libs/core/idempotency/contracts"
	"grovehub/libs/core/idempotency/helpers"
	"grovehub/libs/core/idempotency/validators"
)

const MaxKeyLength = contracts.MaxKeyLength

var (
	ErrInvalidKey  = contracts.ErrInvalidKey
	ErrFingerprint = contracts.ErrFingerprint
)

func ValidateKey(value string) (string, error) {
	return validators.ValidateKey(value)
}

func Fingerprint(value any) (string, error) {
	return helpers.Fingerprint(value)
}
