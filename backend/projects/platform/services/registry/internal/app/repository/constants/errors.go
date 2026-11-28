package constants

import "errors"

var (
	ErrNotFound = errors.New("registry resource not found")
	ErrConflict = errors.New("registry resource or version conflict")
)
