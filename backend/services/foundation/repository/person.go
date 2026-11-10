package repository

import (
	"context"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
)

/**
 * PersonRepository persists foundation.persons records, including the
 * denormalized SearchDetails projection. Contact persistence belongs to
 * ContactRepository; read-only contact filters may be used during search.
 */
// @Repository entity=Person
type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) error
	FindByGID(ctx context.Context, gid string) (entity.Person, error)
	FindByIdempotencyKey(ctx context.Context, key string) (entity.Person, error)

	/**
	 * Search returns validated, cursor-paginated Person records.
	 * Allowed filter and sort fields are specified in searchdetails.PersonConfig.
	 */
	// @Searchable config=searchdetails.PersonConfig alias=p
	Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error)

	Update(ctx context.Context, person entity.Person) error
}
