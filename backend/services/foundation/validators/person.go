package validators

import (
	searchcontract "grovehub/libs/core/contracts/search"
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	"grovehub/services/foundation/entity"
)

type PersonValidator interface {
	ValidateCreate(request foundationrequest.CreatePerson) (foundationrequest.CreatePerson, error)
	ValidateUpdate(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error)
	ValidateGID(gid entity.GID) error
	ValidateSearch(
		request searchcontract.Request,
	) (coresearch.Query[entity.GID], coresearch.CursorPayload, error)
}
