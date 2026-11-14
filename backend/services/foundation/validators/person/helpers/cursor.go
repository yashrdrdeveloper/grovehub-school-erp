package helpers

import (
    "strings"
    "time"

    corenormalize "grovehub/libs/core/normalize"
)

type cursorValueValidator func(string) error

type sortFieldSpec struct {
    validate cursorValueValidator
}


var personSortFields = map[string]sortFieldSpec{
	"gid":       {validate: validateUUIDv7String},
	"firstName": {validate: validateNonEmptyCursorValue},
	"lastName":  {validate: validateNonEmptyCursorValue},
	"createdAt": {validate: validateTimeCursorValue},
	"updatedAt": {validate: validateTimeCursorValue},
}

func validateUUIDv7String(value string) error {
	if _, err := corenormalize.UUIDv7(value); err != nil {
		return invalidPersonSearchCursor()
	}
	return nil
}

func validateTimeCursorValue(value string) error {
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		return invalidPersonSearchCursor()
	}
	return nil
}

func validateNonEmptyCursorValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return invalidPersonSearchCursor()
	}
	return nil
}
