package idempotency

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/validators"
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
