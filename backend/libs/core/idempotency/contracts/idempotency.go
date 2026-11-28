package contracts

import "errors"

const MaxKeyLength = 128

var (
	ErrInvalidKey  = errors.New("invalid idempotency key")
	ErrFingerprint = errors.New("unable to fingerprint idempotent request")
)
