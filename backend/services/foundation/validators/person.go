package validators

import (
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request/person"
)

type PersonValidator interface {
	ValidateCreate(request foundationrequest.CreatePerson) (foundationrequest.CreatePerson, error)
	ValidateUpdate(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error)
	ValidateGID(gid string) error
	ValidateSearch(
		request coresearch.Request,
	) (coresearch.Query[string], coresearch.CursorPayload, error)
}
