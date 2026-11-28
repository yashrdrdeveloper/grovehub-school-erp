package activity

import (
	coresearch "grovehub/libs/core/search"
	"grovehub/services/audit/internal/app/entity"
)

type Page = coresearch.Result[entity.Activity]
