package role

import (
	"grovehub/libs/core/request"
	"grovehub/libs/core/search"
)

type ByID struct {
	RoleID string `path:"roleId" json:"roleId"`
}

type GetRequest struct {
	request.Request
	Params ByID `validate:"valid"`
}

type SearchRequest struct {
	request.Request
	RequestBody search.Request `validate:"valid"`
}

type PreviewRequest struct {
	request.Request
	Params      ByID               `validate:"valid"`
	RequestBody PreviewPermissions `validate:"valid"`
}

type CreateRequest struct {
	request.Request
	RequestBody CreateDraft `validate:"valid"`
}

type UpdateRequest struct {
	request.Request
	Params      ByID           `validate:"valid"`
	RequestBody UpdateMetadata `validate:"valid"`
}
