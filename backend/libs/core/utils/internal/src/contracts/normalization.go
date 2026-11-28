package contracts

import "errors"

var (
	ErrInvalidEmail  = errors.New("invalid email")
	ErrInvalidPhone  = errors.New("invalid phone")
	ErrInvalidUUIDv7 = errors.New("invalid UUIDv7")
	ErrInvalidUUID   = errors.New("invalid canonical UUID")
	ErrInvalidText   = errors.New("invalid text")
	ErrInvalidTime   = errors.New("invalid RFC3339 time")
)
