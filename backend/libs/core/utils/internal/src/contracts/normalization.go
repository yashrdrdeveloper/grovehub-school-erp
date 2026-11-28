package contracts

import "errors"

var (
	ErrInvalidEmail  = errors.New("invalid email")
	ErrInvalidPhone  = errors.New("invalid phone")
	ErrInvalidUUIDv7 = errors.New("invalid UUIDv7")
	ErrInvalidTime   = errors.New("invalid RFC3339 time")
)
