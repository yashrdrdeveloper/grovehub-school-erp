package contracts

import "errors"

var (
	ErrInvalidKey  = errors.New("invalid idempotency key")
	ErrFingerprint = errors.New("unable to fingerprint idempotent request")
)
