package activity

import (
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
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
