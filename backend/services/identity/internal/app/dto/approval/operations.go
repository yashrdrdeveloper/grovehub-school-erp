package approval

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
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
