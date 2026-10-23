package entity

import (
	"encoding/hex"
	"errors"
	"strings"
)

var ErrInvalidGID = errors.New("gid must be a non-zero UUID in canonical 8-4-4-4-12 format")

type GID struct {
	value string
}

func ParseGID(value string) (GID, error) {
	normalized := strings.ToLower(value)

	if len(normalized) != 36 ||
		normalized[8] != '-' ||
		normalized[13] != '-' ||
		normalized[18] != '-' ||
		normalized[23] != '-' {
		return GID{}, ErrInvalidGID
	}

	decoded, err := hex.DecodeString(strings.ReplaceAll(normalized, "-", ""))
	if err != nil || len(decoded) != 16 || isZeroUUID(decoded) {
		return GID{}, ErrInvalidGID
	}

	return GID{value: normalized}, nil
}

func (gid GID) String() string {
	return gid.value
}

func (gid GID) IsZero() bool {
	return gid.value == ""
}

func isZeroUUID(value []byte) bool {
	for _, item := range value {
		if item != 0 {
			return false
		}
	}

	return true
}
