// Package validation executes typed DTO normalization and tag validation.
// NormalizeFields and Validate are separate passes so validation always sees
// canonicalized values; Process also runs resource-owned custom rules.
package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	corenormalize "grovehub/libs/core/normalize"
)

// Error describes the invalid field path/rule without disclosing its value.
type Error struct {
	Field string
	Rule  string
}

func (err Error) Error() string { return fmt.Sprintf("%s violates %s", err.Field, err.Rule) }

// Rules describes optional resource-owned canonicalization, validation and
// error mapping. DTO normalization/validation always run around these hooks.
type Rules[T any] struct {
	Normalize   func(T) (T, error)
	Validate    func(T) error
	MapDTOError func(error) error
}

// Process runs tagged normalization -> custom normalization -> DTO validation
// -> resource validation. A failed stage prevents all subsequent stages.
func Process[T any](input T, rules Rules[T]) (T, error) {
	value, err := NormalizeFields(input)
	if err != nil {
		return zeroOnError[T](err, rules.MapDTOError)
	}
	if rules.Normalize != nil {
		value, err = rules.Normalize(value)
		if err != nil {
			var zero T
			return zero, err
		}
	}
	if err := Validate(value); err != nil {
		return zeroOnError[T](err, rules.MapDTOError)
	}
	if rules.Validate != nil {
		if err := rules.Validate(value); err != nil {
			var zero T
			return zero, err
		}
	}
	return value, nil
}

func zeroOnError[T any](err error, mapper func(error) error) (T, error) {
	if mapper != nil {
		err = mapper(err)
	}
	var zero T
	return zero, err
}

// Normalize preserves the original two-in-one DTO helper for existing callers.
// New pipelines should use Process, NormalizeFields and Validate.
func Normalize[T any](input T) (T, error) {
	return Process(input, Rules[T]{})
}

// NormalizeFields performs only normalization; it never checks validate tags.
// It copies pointers/slices, including present PATCH values, before mutation.
func NormalizeFields[T any](input T) (T, error) {
	value := input
	target := reflect.ValueOf(&value).Elem()
	if target.Kind() != reflect.Struct {
		var zero T
		return zero, fmt.Errorf("normalization target must be a struct")
	}
	if err := normalizeStruct(target, ""); err != nil {
		var zero T
		return zero, err
	}
	return value, nil
}

// Validate only checks rules. It never mutates the supplied DTO.
func Validate[T any](input T) error {
	target := reflect.ValueOf(input)
	if target.Kind() != reflect.Struct {
		return fmt.Errorf("validation target must be a struct")
	}
	return validateStruct(target, "")
}

func normalizeStruct(value reflect.Value, prefix string) error {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field, target := typ.Field(i), value.Field(i)
		if !target.CanSet() {
			continue
		}
		normalizer, rules := field.Tag.Get("normalize"), field.Tag.Get("validate")
		if normalizer == "" && rules == "" {
			if !field.Anonymous {
				continue
			}
			rules = "valid"
		}
		path := fieldPath(prefix, field.Name)
		if err := normalizeField(target, path, normalizer, rules); err != nil {
			return err
		}
	}
	return nil
}

func normalizeField(value reflect.Value, path, normalizer, rules string) error {
	if isPatch(value) {
		if !value.FieldByName("Set").Bool() {
			return nil
		}
		return normalizeField(value.FieldByName("Value"), path, normalizer, rules)
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		copied := reflect.New(value.Type().Elem())
		copied.Elem().Set(value.Elem())
		value.Set(copied)
		return normalizeField(copied.Elem(), path, normalizer, rules)
	}
	switch value.Kind() {
	case reflect.String:
		text := value.String()
		for _, operation := range strings.Split(normalizer, ",") {
			switch operation {
			case "":
			case "trim":
				text = strings.TrimSpace(text)
			case "uuidv7":
				normalized, err := corenormalize.UUIDv7(text)
				if err != nil { return Error{Field: path, Rule: "uuidv7"} }
				text = normalized
			default:
				return fmt.Errorf("%s: unsupported normalizer %q", path, operation)
			}
		}
		value.SetString(text)
	case reflect.Struct:
		if normalizer != "" {
			return fmt.Errorf("%s: unsupported normalize target", path)
		}
		if hasRule(rules, "valid") {
			return normalizeStruct(value, path)
		}
	case reflect.Slice:
		if normalizer != "" {
			return fmt.Errorf("%s: unsupported normalize target", path)
		}
		if hasRule(rules, "valid") {
			copyValue := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
			reflect.Copy(copyValue, value)
			value.Set(copyValue)
			for i := 0; i < copyValue.Len(); i++ {
				if err := normalizeField(copyValue.Index(i), fmt.Sprintf("%s[%d]", path, i), "", "valid"); err != nil {
					return err
				}
			}
		}
	default:
		if normalizer != "" {
			return fmt.Errorf("%s: unsupported normalize target", path)
		}
	}
	return nil
}

func validateStruct(value reflect.Value, prefix string) error {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field, target := typ.Field(i), value.Field(i)
		rules := field.Tag.Get("validate")
		if rules == "" {
			if !field.Anonymous {
				continue
			}
			rules = "valid"
		}
		if err := validateField(target, fieldPath(prefix, field.Name), rules); err != nil {
			return err
		}
	}
	return nil
}

func validateField(value reflect.Value, path, rules string) error {
	if isPatch(value) {
		if !value.FieldByName("Set").Bool() {
			return nil // omitted PATCH properties do not require validation
		}
		return validateField(value.FieldByName("Value"), path, rules)
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			if hasRule(rules, "nonnull") || hasRule(rules, "required") {
				return Error{Field: path, Rule: "nonnull"}
			}
			return nil
		}
		return validateField(value.Elem(), path, rules)
	}
	if value.Kind() == reflect.String {
		for _, rule := range strings.Split(rules, ",") {
			switch {
			case rule == "", rule == "valid", rule == "nonnull":
			case rule == "notblank", rule == "required":
				if strings.TrimSpace(value.String()) == "" {
					return Error{Field: path, Rule: rule}
				}
			case rule == "notempty":
				if value.String() == "" {
					return Error{Field: path, Rule: rule}
				}
			case rule == "uuidv7":
				if _, err := corenormalize.UUIDv7(value.String()); err != nil {
					return Error{Field: path, Rule: rule}
				}
			case strings.HasPrefix(rule, "maxrunes="):
				maximum, err := strconv.Atoi(strings.TrimPrefix(rule, "maxrunes="))
				if err != nil || maximum < 0 {
					return fmt.Errorf("%s: invalid maxrunes rule", path)
				}
				if utf8.RuneCountInString(value.String()) > maximum {
					return Error{Field: path, Rule: "maxrunes"}
				}
			default:
				return fmt.Errorf("%s: unsupported validation rule %q", path, rule)
			}
		}
		return nil
	}
	for _, rule := range strings.Split(rules, ",") {
		switch rule {
		case "", "valid", "nonnull", "required":
		case "notempty":
			if (value.Kind() == reflect.Slice || value.Kind() == reflect.Map) && value.Len() == 0 {
				return Error{Field: path, Rule: rule}
			}
		default:
			if strings.HasPrefix(rule, "min=") || strings.HasPrefix(rule, "max=") {
				if err := validateNumber(value, path, rule); err != nil { return err }
				continue
			}
			return fmt.Errorf("%s: unsupported validation rule %q", path, rule)
		}
	}
	if !hasRule(rules, "valid") {
		return nil
	}
	switch value.Kind() {
	case reflect.Struct:
		return validateStruct(value, path)
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if err := validateField(value.Index(i), fmt.Sprintf("%s[%d]", path, i), "valid"); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%s: validate:valid requires a struct or slice", path)
	}
	return nil
}

func validateNumber(value reflect.Value, path, rule string) error {
	parts := strings.SplitN(rule, "=", 2)
	limit, err := strconv.ParseFloat(parts[1], 64)
	if err != nil { return fmt.Errorf("%s: invalid numeric constraint %q", path, rule) }
	var number float64
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		number = float64(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		number = float64(value.Uint())
	case reflect.Float32, reflect.Float64:
		number = value.Float()
	default:
		return fmt.Errorf("%s: numeric constraint on unsupported type", path)
	}
	if parts[0] == "min" && number < limit { return Error{Field: path, Rule: "min"} }
	if parts[0] == "max" && number > limit { return Error{Field: path, Rule: "max"} }
	return nil
}

func isPatch(value reflect.Value) bool {
	if value.Kind() != reflect.Struct {
		return false
	}
	typ := value.Type()
	return typ.PkgPath() == "grovehub/libs/core/annotation/dto" &&
		strings.HasPrefix(typ.Name(), "Patch[")
}

func fieldPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "." + name
}

func hasRule(rules, want string) bool {
	for _, rule := range strings.Split(rules, ",") {
		if rule == want {
			return true
		}
	}
	return false
}
