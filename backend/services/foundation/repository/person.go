package repository

import (
	"context"
	"errors"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
)

var (
	ErrPersonNotFound                    = errors.New("person not found")
	ErrPersonContactAlreadyExists        = errors.New("person contact already exists")
	ErrPersonVersionConflict             = errors.New("person version conflict")
	ErrPersonIdempotencyKeyAlreadyExists = errors.New("person idempotency key already exists")
)

type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) error
	FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error)
	FindByCreateIdempotencyKey(ctx context.Context, key string) (entity.Person, error)
	FindGIDsByContacts(ctx context.Context, contacts []entity.Contact) ([]entity.GID, error)
	Search(ctx context.Context, query coresearch.Query[entity.GID]) (coresearch.QueryResult[entity.Person], error)
	Update(ctx context.Context, person entity.Person) error
	Delete(ctx context.Context, gid entity.GID) error
}
