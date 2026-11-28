package accountrole

import (
	"grovehub/libs/core/request"
)

type ByID struct {
	AccountID string `path:"accountId" json:"accountId"`
}

type GetRequest struct {
	request.Request
	Params ByID `validate:"valid"`
}

type PreviewRequest struct {
	request.Request
	Params      ByID    `validate:"valid"`
	RequestBody Replace `validate:"valid"`
}
