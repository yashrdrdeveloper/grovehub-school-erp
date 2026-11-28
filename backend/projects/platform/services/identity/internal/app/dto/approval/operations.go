package approval

import (
	"grovehub/libs/core/request"
	"grovehub/libs/core/search"
)

type ByID struct {
	RequestID string `path:"requestId" json:"requestId"`
}

type GetRequest struct {
	request.Request
	Params ByID `validate:"valid"`
}

type SearchRequest struct {
	request.Request
	RequestBody search.Request `validate:"valid"`
}

type CountRequest struct {
	request.Request
}

type Count struct {
	Pending int64 `json:"pending"`
}

type LifecycleRequest struct {
	request.Request
	RequestBody CreateLifecycle `validate:"valid"`
}

type PermissionRequest struct {
	request.Request
	RequestBody CreatePermissions `validate:"valid"`
}

type AssignmentRequest struct {
	request.Request
	RequestBody CreateAssignment `validate:"valid"`
}

type DecisionBody struct {
	Version int64 `json:"version"`
	Comment string `json:"comment"`
}

type DecisionRequest struct {
	request.Request
	Params ByID `validate:"valid"`
	RequestBody DecisionBody `validate:"valid"`
}

type VersionBody struct {
	Version int64 `json:"version"`
}

type VersionRequest struct {
	request.Request
	Params ByID `validate:"valid"`
	RequestBody VersionBody `validate:"valid"`
}
