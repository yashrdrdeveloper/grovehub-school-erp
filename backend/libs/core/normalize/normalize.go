package normalize

import "grovehub/libs/core/normalize/helpers"

var (
	ErrInvalidEmail  = helpers.ErrInvalidEmail
	ErrInvalidPhone  = helpers.ErrInvalidPhone
	ErrInvalidUUIDv7 = helpers.ErrInvalidUUIDv7
	ErrInvalidTime   = helpers.ErrInvalidTime
)

func Email(value string) (string, error) {
	return helpers.Email(value)
}

func Phone(value string) (string, error) {
	return helpers.Phone(value)
}

func PhoneCharacters(value string) (string, error) {
	return helpers.PhoneCharacters(value)
}

func UUIDv7(value string) (string, error) {
	return helpers.UUIDv7(value)
}

func RFC3339UTC(value string) (string, error) {
	return helpers.RFC3339UTC(value)
}
