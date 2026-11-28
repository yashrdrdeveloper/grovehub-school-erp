package helpers

import (
	"bytes"
	"encoding/json"
)

// JSONValue respects JSON tags, custom marshaling and exact numeric tokens.
func JSONValue(value any) (any, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()

	var result any
	err = decoder.Decode(&result)

	return result, err
}
