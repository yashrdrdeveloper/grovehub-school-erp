package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coretransaction "grovehub/libs/core/database/transaction"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
	repositoryerrorshelpers "grovehub/services/foundation/repository/helpers/errors"
)

type PersonRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) foundationrepository.PersonRepository {
	return &PersonRepository{db: db}
}

func (repository *PersonRepository) Create(ctx context.Context, person entity.Person) error {
	return repositoryerrorshelpers.MapPersonWriteError(coretransaction.Database(ctx, repository.db).Create(&person).Error)
}

func (repository *PersonRepository) FindByGID(ctx context.Context, gid string) (entity.Person, error) {
	var person entity.Person
	if err := coretransaction.Database(ctx, repository.db).Where("gid = ?", gid).Take(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, repositoryerrors.ErrPersonNotFound
		}
		return entity.Person{}, err
	}
	return person, nil
}

func (repository *PersonRepository) FindByIdempotencyKey(ctx context.Context, key string) (entity.Person, error) {
	var person entity.Person
	if err := coretransaction.Database(ctx, repository.db).Where(&entity.Person{IdempotencyKey: &key}).Take(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, repositoryerrors.ErrPersonNotFound
		}
		return entity.Person{}, err
	}
	return person, nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	expectedVersion := person.Version - 1
	result := coretransaction.Database(ctx, repository.db).
		Model(&entity.Person{}).
		Where("gid = ? AND version = ?", person.GID, expectedVersion).
		Updates(map[string]any{
			"first_name": person.FirstName,
			"middle_name": person.MiddleName,
			"last_name": person.LastName,
			"status": person.Status,
			"version": person.Version,
			"search_details": person.SearchDetails,
			"updated_at": person.UpdatedAt,
		})
	if result.Error != nil {
		return repositoryerrorshelpers.MapPersonWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return repositoryerrors.ErrPersonVersionConflict
	}
	return nil
}

// @Searchable config=search.PersonConfig alias=p
// Search delegates the approved, cursor-safe search query to generated plumbing.
func (repository *PersonRepository) Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error) {
	return repository.searchGenerated(ctx, query)
}
