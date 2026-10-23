package repository

import (
	"context"
	"errors"

	"grovehub/services/foundation/entity"
)

var ErrPersonNotFound = errors.New("person not found")

type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) error
	FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error)
	Update(ctx context.Context, person entity.Person) error
	Delete(ctx context.Context, gid entity.GID) error
}
