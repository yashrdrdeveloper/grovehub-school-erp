package activity

import "grovehub/libs/core/request"

type SearchRequest struct {
	request.Request
	RequestBody Search `validate:"valid"`
}
