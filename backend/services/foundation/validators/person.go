package validators

import (
	searchcontract "grovehub/libs/core/contracts/search"
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	"grovehub/services/foundation/entity"
)

type ValidatedCreatePerson struct {
	FirstName  string
	MiddleName string
	LastName   string
	Contacts   []entity.Contact
}

type ValidatedUpdatePerson struct {
	FirstName  *string
	MiddleName *string
	LastName   *string
	Contacts   *[]entity.Contact
}

type ValidatedPersonSearch struct {
	Query         coresearch.Query[entity.GID]
	CursorPayload coresearch.CursorPayload
}

type PersonValidator interface {
	ValidateCreate(request foundationrequest.CreatePerson) (ValidatedCreatePerson, error)
	ValidateUpdate(
		request foundationrequest.UpdatePerson,
		existingContacts []entity.Contact,
	) (ValidatedUpdatePerson, error)
	ValidateGID(gid entity.GID) error
	ValidateSearch(request searchcontract.Request) (ValidatedPersonSearch, error)
}
