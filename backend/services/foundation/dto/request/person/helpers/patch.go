// Package helpers provides a reusable presence-aware JSON PATCH value.
// The parent Person DTO owns its field and contact contracts.
package helpers

import (
	"bytes"
	"encoding/json"
)

// Patch distinguishes omitted (Set=false), explicit null (Set=true, Value=nil)
// and a supplied value (Set=true, Value!=nil).
type Patch[T any] struct {
	Set bool
	Value *T
}

// UnmarshalJSON is called only for a present JSON field.
func (field *Patch[T]) UnmarshalJSON(data []byte) error {
	field.Set = true
	field.Value = nil
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) { return nil }
	var value T
	if err := json.Unmarshal(data, &value); err != nil { return err }
	field.Value = &value
	return nil
}
