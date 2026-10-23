package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	personerror "grovehub/services/foundation/constants/error"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

type PersonService struct {
	repository foundationrepository.PersonRepository
}

func NewPersonService(repository foundationrepository.PersonRepository) *PersonService {
	return &PersonService{repository: repository}
}

func (service *PersonService) Create(ctx context.Context) (entity.Person, error) {
	gid, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Person{}, coreerrors.Wrap(personerror.PersonCreateFailed, err)
	}

	now := time.Now().UTC()
	person := entity.Person{
		GID:       entity.GID(gid),
		Status:    entity.PersonStatusActive,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := service.repository.Create(ctx, person); err != nil {
		return entity.Person{}, coreerrors.Wrap(personerror.PersonCreateFailed, err)
	}

	return person, nil
}

func (service *PersonService) GetByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	if err := validateGID(gid); err != nil {
		return entity.Person{}, err
	}

	person, err := service.repository.FindByGID(ctx, gid)
	if errors.Is(err, foundationrepository.ErrPersonNotFound) {
		return entity.Person{}, coreerrors.NewError(personerror.PersonNotFound)
	}
	if err != nil {
		return entity.Person{}, coreerrors.Wrap(personerror.PersonReadFailed, err)
	}

	return person, nil
}

func (service *PersonService) Deactivate(ctx context.Context, gid entity.GID) (entity.Person, error) {
	person, err := service.GetByGID(ctx, gid)
	if err != nil {
		return entity.Person{}, err
	}

	if person.Status == entity.PersonStatusDeactivated {
		return person, nil
	}

	person.Status = entity.PersonStatusDeactivated
	person.Version++
	person.UpdatedAt = time.Now().UTC()

	if err := service.repository.Update(ctx, person); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonNotFound) {
			return entity.Person{}, coreerrors.NewError(personerror.PersonNotFound)
		}

		return entity.Person{}, coreerrors.Wrap(personerror.PersonUpdateFailed, err)
	}

	return person, nil
}

func validateGID(gid entity.GID) error {
	value, err := uuid.Parse(string(gid))
	if err != nil || value == uuid.Nil || value.Version() != uuid.Version(7) {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}

	return nil
}
