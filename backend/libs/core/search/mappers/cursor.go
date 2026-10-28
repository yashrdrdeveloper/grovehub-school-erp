package mappers

import (
	"encoding/base64"
	"encoding/json"

	"grovehub/libs/core/search/contracts"
)

func EncodeCursor(cursor contracts.CursorPayload) (string, error) {
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
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
