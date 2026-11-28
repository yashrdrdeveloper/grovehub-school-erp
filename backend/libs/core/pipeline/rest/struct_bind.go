package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// bindParamsInto and bindQueryInto reuse the scalar/slice/pointer binder used
// by the earlier typed REST pipeline. Missing or malformed input is rejected.
func bindParamsInto(target reflect.Value, request *http.Request) error {
	return bindStructInto(target, "path", func(name string) ([]string, bool) {
		value := request.PathValue(name)
		return []string{value}, value != ""
	})
}

func bindQueryInto(target reflect.Value, request *http.Request, strict bool) error {
	values, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil { return BindError{Source:"query",Field:"query",Reason:"malformed"} }
	if strict {
		typ := target.Type()
		if typ.Kind() != reflect.Struct { return fmt.Errorf("query DTO must be a struct") }
		known := make(map[string]bool)
		for i := 0; i < typ.NumField(); i++ {
			key := strings.Split(typ.Field(i).Tag.Get("query"), ",")[0]
			if key != "" && key != "-" { known[key] = true }
		}
		for key := range values {
			if !known[key] { return BindError{Source:"query",Field:key,Reason:"unknown"} }
		}
	}
	return bindStructInto(target, "query", func(name string) ([]string, bool) {
		raw, ok := values[name]
		return raw, ok
	})
}

func bindStructInto(target reflect.Value, source string, get func(string) ([]string, bool)) error {
	if target.Kind() != reflect.Struct { return fmt.Errorf("%s request field must be a struct",source) }
	typ := target.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := strings.Split(field.Tag.Get(source), ",")[0]
		if name == "" || name == "-" { continue }
		if !field.IsExported() || !target.Field(i).CanSet() { return fmt.Errorf("unexported %s DTO field",source) }
		raw, exists := get(name)
		if !exists {
			if source=="path" { return BindError{Source:source, Field:name, Reason:"missing"} }
			continue
		}
		if err := bindValue(target.Field(i), raw); err != nil {
			return BindError{Source:source, Field:name, Reason:"invalid value or multiplicity"}
		}
	}
	return nil
}
