package validators

import (
    "strings"
    "time"

    corenormalize "grovehub/libs/core/normalize"
    "grovehub/libs/core/search/contracts"
    "grovehub/libs/core/search/mappers"
)

func cursorError(invalid error) error {
    if invalid != nil { return invalid }
    return ErrInvalidQuery
}

func ValidateUUIDv7Cursor(value string, invalid error) error {
    if _, err := corenormalize.UUIDv7(value); err != nil { return cursorError(invalid) }
    return nil
}

func ValidateTimeCursor(value string, invalid error) error {
    if _, err := time.Parse(time.RFC3339Nano, value); err != nil { return cursorError(invalid) }
    return nil
}

func ValidateNonEmptyCursor(value string, invalid error) error {
    if strings.TrimSpace(value) == "" { return cursorError(invalid) }
    return nil
}

// NormalizeCursor validates a stable cursor against the requested sort.
// The owning resource chooses its sort-value and cursor-ID validators.
func NormalizeCursor(value string, sort contracts.QuerySort, validateValue contracts.CursorValueValidator, validateID contracts.CursorValueValidator, invalid error) (*contracts.Cursor[string], contracts.CursorPayload, error) {
    payload := contracts.CursorPayload{Field: sort.Field, Direction: sort.Direction}
    if strings.TrimSpace(value) == "" { return nil, payload, nil }
    if validateValue == nil || validateID == nil { return nil, contracts.CursorPayload{}, cursorError(invalid) }

    decoded, err := mappers.DecodeCursor(value)
    if err != nil || decoded.Field != sort.Field || decoded.Direction != sort.Direction {
        return nil, contracts.CursorPayload{}, cursorError(invalid)
    }
    if validateValue(decoded.Value, invalid) != nil || validateID(decoded.CursorId, invalid) != nil {
        return nil, contracts.CursorPayload{}, cursorError(invalid)
    }
    return &contracts.Cursor[string]{Value: decoded.Value, CursorId: decoded.CursorId}, decoded, nil
}
