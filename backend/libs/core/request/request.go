// Package request defines the shared, transport-independent contract for
// an endpoint request. Concrete requests embed Request, declaring only the
// applicable Params, Query and RequestBody fields. Missing fields are absent,
// rather than requiring explicit Empty placeholder type arguments.
package request

import (
	"fmt"
	"reflect"
)

// Request is the common, zero-sized request contract. Embed it in each
// concrete request DTO alongside only the sections the endpoint accepts.
//
//   type UpdateRequest struct {
//       request.Request
//       Params ByID
//       RequestBody UpdateDTO
//   }
//
// Query is automatically absent in that request; no placeholder is needed.
type Request struct{}

// Contract ensures every endpoint request shares the same foundation.
type Contract interface{ IsRequest() }
func (Request) IsRequest() {}

// Sections exposes only the fields actually declared by a concrete request.
// Reflect values are used internally by binders; callers use typed DTOs.
type Sections struct {
	Params reflect.Value
	Query reflect.Value
	RequestBody reflect.Value
	HasParams bool
	HasQuery bool
	HasRequestBody bool
}

// Inspect validates the common request shape. It accepts a struct or pointer
// to struct with an embedded Request and at most one of each named section.
func Inspect(input any) (Sections, error) {
	var sections Sections
	if input == nil { return sections, fmt.Errorf("nil request") }
	if _,ok := input.(Contract); !ok { return sections, fmt.Errorf("request must implement common Request contract") }
	value := reflect.ValueOf(input)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() { return sections, fmt.Errorf("nil request pointer") }
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct { return sections, fmt.Errorf("request must be a struct") }
	if value.Type() == reflect.TypeOf(Request{}) { return sections, fmt.Errorf("request cannot consist only of the common marker") }
	markerFound := false
	for i:=0; i<value.NumField(); i++ {
		field := value.Type().Field(i)
		if field.Anonymous && field.Type == reflect.TypeOf(Request{}) {
			markerFound = true
			continue
		}
		if !field.IsExported() { return Sections{}, fmt.Errorf("request includes an unexported field") }
		switch field.Name {
		case "Params":
			sections.Params,sections.HasParams=value.Field(i),true
		case "Query":
			sections.Query,sections.HasQuery=value.Field(i),true
		case "RequestBody":
			sections.RequestBody,sections.HasRequestBody=value.Field(i),true
		default:
			return Sections{}, fmt.Errorf("unsupported request section %q",field.Name)
		}
	}
	if !markerFound { return Sections{}, fmt.Errorf("request must embed the common request.Request") }
	if !sections.HasParams && !sections.HasQuery && !sections.HasRequestBody {
		return Sections{},fmt.Errorf("request must declare at least one section")
	}
	return sections,nil
}

// Empty and Parts exist solely for callers of the original legacy typed
// pipeline; new endpoint DTOs use embedded Request and omit missing fields.
type Empty struct{}
type Parts[Path, Query, Body any] struct {
	Path Path
	Query Query
	Body Body
}
