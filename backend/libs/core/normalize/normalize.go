package normalize

import "grovehub/libs/core/normalize/helpers"

var (
	ErrInvalidEmail  = helpers.ErrInvalidEmail
	ErrInvalidPhone  = helpers.ErrInvalidPhone
	ErrInvalidUUIDv7 = helpers.ErrInvalidUUIDv7
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
