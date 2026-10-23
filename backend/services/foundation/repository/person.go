package repository

import (
	"context"

	"grovehub/services/foundation/entity"
)

type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) error
	FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error)
	Update(ctx context.Context, person entity.Person) error
	Delete(ctx context.Context, gid entity.GID) error
}
