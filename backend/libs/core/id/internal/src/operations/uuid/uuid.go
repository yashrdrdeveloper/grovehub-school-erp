package uuid

import googleuuid "github.com/google/uuid"

// NewUUIDv7 returns a UUIDv7 string or the original generation error.
func NewUUIDv7() (string, error) {
	value, err := googleuuid.NewV7()
	if err != nil {
		return "", err
	}
	return value.String(), nil
}
