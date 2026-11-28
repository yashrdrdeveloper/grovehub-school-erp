package contracts

import "reflect"

// Request is the common, zero-sized marker embedded in endpoint request DTOs.
type Request struct{}

// Contract ensures every endpoint request shares the same foundation.
type Contract interface {
	IsRequest()
}

func (Request) IsRequest() {}

// Sections exposes the fields actually declared by a concrete request.
type Sections struct {
	Params         reflect.Value
	Query          reflect.Value
	RequestBody    reflect.Value
	HasParams      bool
	HasQuery       bool
	HasRequestBody bool
}

// Empty and Parts support legacy typed requests with explicit absent sections.
type Empty struct{}

type Parts[Path, Query, Body any] struct {
	Path  Path
	Query Query
	Body  Body
}
