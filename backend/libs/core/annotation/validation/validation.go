// Package validation implements the executable struct-tag subset behind
// checked declarative DTO annotations; comments alone do not validate data.
package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Error contains only a field path and rule; never a request value or secret.
type Error struct {
	Field string
	Rule string
}

func (e Error) Error() string { return fmt.Sprintf("%s violates %s", e.Field, e.Rule) }

// Normalize returns a new DTO with nested validation and safe copies of
// pointer and slice values. Nil PATCH pointers remain nil.
func Normalize[T any](input T) (T, error) {
	output := input
	value := reflect.ValueOf(&output).Elem()
	if value.Kind() != reflect.Struct { return output, fmt.Errorf("validation target must be a struct") }
	if err := normalizeStruct(value, ""); err != nil {
		var zero T
		return zero, err
	}
	return output, nil
}

func normalizeStruct(value reflect.Value, prefix string) error {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		target := value.Field(i)
		if !target.CanSet() { continue }
		path := field.Name
		if prefix != "" { path = prefix + "." + path }
		normalizer, rules := field.Tag.Get("normalize"), field.Tag.Get("validate")
		if rules == "" && normalizer == "" {
			if !field.Anonymous { continue }
			rules = "valid"
		}
		if err := normalizeField(target, path, normalizer, rules); err != nil { return err }
	}
	return nil
}

func normalizeField(value reflect.Value, path, normalizer, rules string) error {
	original := value
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			if hasRule(rules, "nonnull") || hasRule(rules, "required") { return Error{Field:path,Rule:"nonnull"} }
			return nil
		}
		copied := reflect.New(value.Type().Elem())
		copied.Elem().Set(value.Elem())
		value.Set(copied)
		value = copied.Elem()
	}
	if value.Kind() == reflect.String {
		text := value.String()
		switch normalizer {
		case "":
		case "trim": text = strings.TrimSpace(text)
		default: return fmt.Errorf("%s: unsupported normalizer %s", path, normalizer)
		}
		value.SetString(text)
		for _, rule := range strings.Split(rules, ",") {
			switch {
			case rule == "", rule == "valid", rule == "nonnull":
			case rule == "required", rule == "notblank":
				if strings.TrimSpace(text) == "" { return Error{Field:path,Rule:rule} }
			case rule == "notempty":
				if text == "" { return Error{Field:path,Rule:rule} }
			case strings.HasPrefix(rule, "maxrunes="):
				maximum, err := strconv.Atoi(strings.TrimPrefix(rule, "maxrunes="))
				if err != nil || maximum < 0 { return fmt.Errorf("%s: invalid maxrunes rule", path) }
				if utf8.RuneCountInString(text) > maximum { return Error{Field:path,Rule:"maxrunes"} }
			default:
				return fmt.Errorf("%s: unsupported validation rule %s", path, rule)
			}
		}
		return nil
	}
	if normalizer != "" { return fmt.Errorf("%s: unsupported normalize target", path) }
	for _, rule := range strings.Split(rules, ",") {
		switch rule {
		case "", "valid", "nonnull", "required":
		case "notempty":
			if (value.Kind() == reflect.Slice || value.Kind() == reflect.Map) && value.Len() == 0 {
				return Error{Field:path,Rule:rule}
			}
		default:
			return fmt.Errorf("%s: unsupported validation rule %s", path, rule)
		}
	}
	if !hasRule(rules, "valid") { return nil }
	switch value.Kind() {
	case reflect.Struct:
		return normalizeStruct(value, path)
	case reflect.Slice:
		copied := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		reflect.Copy(copied, value)
		original.Set(copied)
		for i := 0; i < copied.Len(); i++ {
			if err := normalizeField(copied.Index(i), fmt.Sprintf("%s[%d]", path, i), "", "valid"); err != nil { return err }
		}
		return nil
	default:
		return fmt.Errorf("%s: @Valid requires a struct or slice", path)
	}
}

func hasRule(rules, want string) bool {
	for _, rule := range strings.Split(rules, ",") {
		if rule == want { return true }
	}
	return false
}
