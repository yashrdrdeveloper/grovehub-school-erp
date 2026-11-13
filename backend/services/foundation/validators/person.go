package validators

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
)

type PersonValidator interface {
	CreateRules() annotationvalidation.Rules[foundationrequest.CreatePerson]
	UpdateRules() annotationvalidation.Rules[foundationrequest.UpdatePerson]
	SearchRules() annotationvalidation.Rules[coresearch.Request]
}
