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
	Create(ctx context.Context, person entity.Person, contacts []entity.Contact) error
	FindByGID(ctx context.Context, gid string) (entity.Person, []entity.Contact, error)
	FindByCreateIdempotencyKey(ctx context.Context, key string) (entity.Person, []entity.Contact, error)
	FindGIDsByContacts(ctx context.Context, contacts []entity.Contact) ([]string, error)
	Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], map[string][]entity.Contact, error)
	Update(ctx context.Context, person entity.Person, contacts []entity.Contact) error
	Delete(ctx context.Context, gid string) error
}
