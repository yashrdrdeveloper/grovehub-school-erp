package activity

import (
	"grovehub/libs/core/request"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/audit/internal/app/entity"
)

type Search struct {
	coresearch.Request
	Scope string `json:"scope,omitempty"`
}

type SearchRequest struct {
	request.Request
	RequestBody Search `validate:"valid"`
}

type Page = coresearch.Result[entity.Activity]
