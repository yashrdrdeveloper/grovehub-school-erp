package validators

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
)

type PersonValidator interface {
	CreateRules() annotationvalidation.Rules[foundationrequest.CreatePerson]
	UpdateRules() annotationvalidation.Rules[foundationrequest.UpdatePerson]
	ValidateCreate(request foundationrequest.CreatePerson) (foundationrequest.CreatePerson, error)
	ValidateUpdate(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error)
	ValidateGID(gid string) error
	ValidateSearch(request coresearch.Request) (coresearch.Query[string], coresearch.CursorPayload, error)
}
