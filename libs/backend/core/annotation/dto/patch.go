package dto

import (
	"bytes"
	"encoding/json"
)

// Patch describes an optional JSON field in a PATCH request.
// Set=false: omitted; Set=true and Value=nil: explicitly null;
// Set=true and Value!=nil: supplied value.
//
// The owning service validates which fields may be cleared.
type Patch[T any] struct {
	Set   bool
	Value *T
}

// UnmarshalJSON is invoked by encoding/json only when a field is present.
// A new destination object should be used for each decoded request.
func (field *Patch[T]) UnmarshalJSON(data []byte) error {
	field.Set = true
	field.Value = nil
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil
}
