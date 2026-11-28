package cursor

import (
	"encoding/base64"
	"encoding/json"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
)

func EncodeCursor(cursor contracts.CursorPayload) (string, error) {
	payload, err := json.Marshal(cursor)
	// Marshal returns nil bytes on failure; encoding nil yields an empty token.
	// Preserve the error without duplicating the serializer's failure branch.
	return base64.RawURLEncoding.EncodeToString(payload), err
}

func DecodeCursor(value string) (contracts.CursorPayload, error) {
	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return contracts.CursorPayload{}, err
	}
	var cursor contracts.CursorPayload
	if err := json.Unmarshal(payload, &cursor); err != nil {
		return contracts.CursorPayload{}, err
	}
	return cursor, nil
}
