package repository

import (
	"context"
	"errors"

	searchcontract "grovehub/libs/core/contracts/search"
	"grovehub/services/foundation/entity"
)

var (
	ErrPersonNotFound             = errors.New("person not found")
	ErrPersonContactAlreadyExists = errors.New("person contact already exists")
	ErrInvalidPersonSearchQuery   = errors.New("invalid person search query")
)

type PersonSearchFilter struct { Field string; Operator searchcontract.Operator; Values []string }
type PersonSearchSort struct { Field string; Direction searchcontract.Direction }
type PersonSearchCursor struct { Value string; GID entity.GID }
type PersonSearchQuery struct { SearchTerms []string; Filters []PersonSearchFilter; Sort PersonSearchSort; Cursor *PersonSearchCursor; Limit int }
type PersonSearchResult struct { Persons []entity.Person; HasNext bool }

type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) error
	FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error)
	FindGIDsByContacts(ctx context.Context, contacts []entity.Contact) ([]entity.GID, error)
	Search(ctx context.Context, query PersonSearchQuery) (PersonSearchResult, error)
	Update(ctx context.Context, person entity.Person) error
	Delete(ctx context.Context, gid entity.GID) error
}
