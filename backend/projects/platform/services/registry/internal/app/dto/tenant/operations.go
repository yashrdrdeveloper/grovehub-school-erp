package tenant

import "grovehub/libs/core/request"

type GetRequest struct {
	request.Request
	Params Params
}

type CreateRequest struct {
	request.Request
	Params ProjectParams
	RequestBody CreateInput
}

type StatusRequest struct {
	request.Request
	Params Params
	RequestBody StatusInput
}
