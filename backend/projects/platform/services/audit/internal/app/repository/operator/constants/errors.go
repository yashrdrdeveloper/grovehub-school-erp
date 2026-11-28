package constants

import "errors"

const ConflictMessage = "operator occurrence conflict"

var ErrConflict = errors.New(ConflictMessage)
