package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

type Lookup func(string) (string, bool)

func DecodeField(target reflect.Value, field reflect.StructField, lookup Lookup) error {
	key, exists := field.Tag.Lookup("env")
	if !exists {
		return nil
	}
	if field.PkgPath != "" {
		return fmt.Errorf("configuration field %s must be exported", field.Name)
	}
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("configuration field %s has an empty environment key", field.Name)
	}

	text := resolveValue(key, field, lookup)
	if text == "" && field.Tag.Get("required") == "true" {
		return fmt.Errorf("required configuration %s is missing", key)
	}
	if text == "" {
		return nil
	}

	return assignValue(target, field.Name, key, text)
}

func resolveValue(key string, field reflect.StructField, lookup Lookup) string {
	text, exists := lookup(key)
	text = strings.TrimSpace(text)

	if !exists || text == "" {
		if fallback, found := field.Tag.Lookup("default"); found {
			return fallback
		}
	}

	return text
}
