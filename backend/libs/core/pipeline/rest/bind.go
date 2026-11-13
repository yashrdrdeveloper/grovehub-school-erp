package rest

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// BindError does not contain the rejected parameter value.
type BindError struct { Source string; Field string; Reason string }
func (err BindError) Error() string { return fmt.Sprintf("invalid %s parameter %s: %s", err.Source, err.Field, err.Reason) }

// bindPath and bindQuery support scalar string/bool/int/uint/float fields,
// pointers for presence-aware optional values, and repeated query slices.
// Tagged fields must be exported. Unknown query parameters are ignored.
func bindPath[T any](request *http.Request) (T, error) {
	return bindFields[T]("path", func(name string) ([]string, bool) {
		value := request.PathValue(name)
		return []string{value}, value != ""
	})
}

func bindQuery[T any](request *http.Request) (T, error) {
	values := request.URL.Query()
	return bindFields[T]("query", func(name string) ([]string, bool) {
		raw, exists := values[name]
		return raw, exists
	})
}

func bindFields[T any](source string, get func(string) ([]string, bool)) (T, error) {
	var output T
	target := reflect.ValueOf(&output).Elem()
	if target.Kind() != reflect.Struct {
		return output, fmt.Errorf("%s DTO must be a struct", source)
	}
	typ := target.Type()
	for index := 0; index < typ.NumField(); index++ {
		field := typ.Field(index)
		tag := strings.Split(field.Tag.Get(source), ",")[0]
		if tag == "" || tag == "-" { continue }
		if !field.IsExported() || !target.Field(index).CanSet() {
			return output, fmt.Errorf("%s DTO field %s must be exported", source, field.Name)
		}
		values, exists := get(tag)
		if !exists {
			if source == "path" { return output, BindError{Source: source, Field: tag, Reason: "missing"} }
			continue
		}
		if err := bindValue(target.Field(index), values); err != nil {
			return output, BindError{Source: source, Field: tag, Reason: "invalid value or multiplicity"}
		}
	}
	return output, nil
}

func bindValue(target reflect.Value, values []string) error {
	if target.Kind() == reflect.Pointer {
		element := reflect.New(target.Type().Elem())
		if err := bindValue(element.Elem(), values); err != nil { return err }
		target.Set(element)
		return nil
	}
	if target.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(target.Type(), 0, len(values))
		for _, value := range values {
			element := reflect.New(target.Type().Elem()).Elem()
			if err := bindScalar(element, value); err != nil { return err }
			slice = reflect.Append(slice, element)
		}
		target.Set(slice)
		return nil
	}
	if len(values) != 1 { return fmt.Errorf("expected one value") }
	return bindScalar(target, values[0])
}

func bindScalar(target reflect.Value, value string) error {
	switch target.Kind() {
	case reflect.String:
		target.SetString(value)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil { return err }
		target.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, target.Type().Bits())
		if err != nil { return err }
		target.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(value, 10, target.Type().Bits())
		if err != nil { return err }
		target.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(value, target.Type().Bits())
		if err != nil { return err }
		target.SetFloat(parsed)
	default:
		return fmt.Errorf("unsupported parameter type")
	}
	return nil
}

func isEmpty[T any]() bool { return reflect.TypeFor[T]() == reflect.TypeFor[Empty]() }
