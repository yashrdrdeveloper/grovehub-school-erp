package errors

import "errors"

var (
	ErrPersonNotFound = errors.New("person not found")
	ErrPersonVersionConflict = errors.New("person version conflict")
	ErrPersonIdempotencyKeyAlreadyExists = errors.New("person idempotency key already exists")
)
