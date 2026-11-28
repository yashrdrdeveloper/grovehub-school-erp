package inspect

import (
	"fmt"
	"reflect"

	"grovehub/libs/core/request/internal/src/contracts"
)

// Inspect validates the common request shape and discovers declared sections.
// A pointer input preserves the addressability needed by transport binders.
func Inspect(input any) (contracts.Sections, error) {
	var sections contracts.Sections
	if input == nil {
		return sections, fmt.Errorf("nil request")
	}
	if _, ok := input.(contracts.Contract); !ok {
		return sections, fmt.Errorf("request must implement common Request contract")
	}

	value := reflect.ValueOf(input)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return sections, fmt.Errorf("nil request pointer")
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return sections, fmt.Errorf("request must be a struct")
	}

	markerType := reflect.TypeFor[contracts.Request]()
	if value.Type() == markerType {
		return sections, nil
	}

	markerFound := false
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if field.Anonymous && field.Type == markerType {
			markerFound = true
			continue
		}
		if !field.IsExported() {
			return contracts.Sections{}, fmt.Errorf("request includes an unexported field")
		}
		switch field.Name {
		case "Params":
			sections.Params, sections.HasParams = value.Field(i), true
		case "Query":
			sections.Query, sections.HasQuery = value.Field(i), true
		case "RequestBody":
			sections.RequestBody, sections.HasRequestBody = value.Field(i), true
		default:
			return contracts.Sections{}, fmt.Errorf("unsupported request section %q", field.Name)
		}
	}
	if !markerFound {
		return contracts.Sections{}, fmt.Errorf("request must embed the common request.Request")
	}
	return sections, nil
}
