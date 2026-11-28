package workload

import "grovehub/libs/core/request"

type CreateRequest struct {
	request.Request
	RequestBody CreateInput
}

type GetRequest struct {
	request.Request
	Params Params
}

type StatusRequest struct {
	request.Request
	Params Params
	RequestBody StatusInput
}

type KeyRequest struct {
	request.Request
	Params Params
	RequestBody KeyInput
}

type RevokeKeyRequest struct {
	request.Request
	Params KeyParams
}

type GrantRequest struct {
	request.Request
	Params Params
	RequestBody GrantInput
}

type RevokeGrantRequest struct {
	request.Request
	Params GrantParams
}
