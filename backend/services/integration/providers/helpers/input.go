package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	errorcode "grovehub/services/integration/constants/error"
)

// DecodeStrictActionInput rejects unknown fields and multiple JSON documents.
// Empty input means an empty object for parameter-free capabilities.
func DecodeStrictActionInput(input json.RawMessage, target any) error {
	if len(input) == 0 {
		input = json.RawMessage("{}")
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.Join(errorcode.ErrInvalidActionInput, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errorcode.ErrInvalidActionInput
	}
	return nil
}
