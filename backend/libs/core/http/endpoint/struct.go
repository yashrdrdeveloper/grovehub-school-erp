package endpoint

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

// CallStruct invokes the registered controller route from an endpoint-specific
// DTO. Only present exported Params, Query and RequestBody fields are used;
// no empty placeholder types or transport wrapper are required.
func CallStruct[Response, Request any](ctx context.Context, handler http.Handler, method, route string, input Request) (Response, int, error) {
	var zero Response
	value := reflect.ValueOf(input)
	if value.Kind() != reflect.Struct { return zero, 0, errors.New("endpoint request must be a struct") }
	var params any
	var queryValues any
	var body any
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if !field.IsExported() { return zero, 0, errors.New("unexported endpoint request field") }
		section := value.Field(i).Interface()
		switch field.Name {
		case "Params":
			params = section
		case "Query":
			queryValues = section
		case "RequestBody":
			body = section
		default:
			return zero, 0, fmt.Errorf("unsupported endpoint request field %q", field.Name)
		}
	}
	encodedQuery, err := encodeQuery(queryValues)
	if err != nil { return zero, 0, err }
	return call[Response](ctx, handler, method, route, params, encodedQuery, body)
}
