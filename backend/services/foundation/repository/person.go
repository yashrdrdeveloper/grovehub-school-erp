package repository

import (
	"context"
	"errors"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
)

var (
	ErrPersonNotFound                    = errors.New("person not found")
	ErrPersonVersionConflict             = errors.New("person version conflict")
	ErrPersonIdempotencyKeyAlreadyExists = errors.New("person idempotency key already exists")
)

// @Repository entity=Person
// PersonRepository owns writes to foundation.persons.
// Contact predicates in Search are read-only aggregate search filters.
type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) error
	FindByGID(ctx context.Context, gid string) (entity.Person, error)
	FindByIdempotencyKey(ctx context.Context, key string) (entity.Person, error)
	Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error)
	Update(ctx context.Context, person entity.Person) error
	Delete(ctx context.Context, gid string) error
}
