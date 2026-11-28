package account

import (
	"grovehub/libs/core/request"
	"grovehub/libs/core/search"
)

type CreateRequest struct {
	request.Request
	RequestBody Create `validate:"valid"`
}

type SetupRequest struct {
	request.Request
	RequestBody Setup `validate:"valid"`
}

type UpdateRequest struct {
	request.Request
	Params      ByID   `validate:"valid"`
	RequestBody Update `validate:"valid"`
}

type GetRequest struct {
	request.Request
	Params ByID `validate:"valid"`
}

type SearchRequest struct {
	request.Request
	RequestBody search.Request `validate:"valid"`
}

type RolesRequest struct {
	request.Request
}
