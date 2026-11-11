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

// JSONField is a named PATCH field with its presence and value preserved.
type JSONField struct {
	Name string
	Set bool
	Value any
}

// NamedField adapts a typed PATCH field for JSON serialization.
func NamedField[T any](name string, patch Patch[T]) JSONField {
	return JSONField{Name: name, Set: patch.Set, Value: patch.Value}
}

// MarshalPatchJSON serializes only explicitly supplied PATCH fields.
// Explicit null remains JSON null; omitted fields are absent.
func MarshalPatchJSON(fields ...JSONField) ([]byte, error) {
	values := make(map[string]any, len(fields))
	for _, field := range fields {
		if field.Set { values[field.Name] = field.Value }
	}
	return json.Marshal(values)
}
