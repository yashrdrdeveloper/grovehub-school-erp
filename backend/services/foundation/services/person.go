package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
	"grovehub/services/foundation/mappers"
	foundationrepository "grovehub/services/foundation/repository"
)

type PersonService struct {
	repository foundationrepository.PersonRepository
}

func NewPersonService(repository foundationrepository.PersonRepository) *PersonService {
	return &PersonService{repository: repository}
}

func (service *PersonService) Create(ctx context.Context) (foundationresponse.Person, error) {
	gidValue, err := coreid.NewUUIDv7()
	if err != nil {
		return foundationresponse.Person{}, coreerrors.Wrap(personerror.PersonCreateFailed, err)
	}

	gid := entity.GID(gidValue)
	now := time.Now().UTC()
	person := entity.NewPerson(
		gid,
		entity.PersonStatusActive,
		1,
		now,
		now,
	)

	if err := service.repository.Create(ctx, person); err != nil {
		return foundationresponse.Person{}, coreerrors.WrapWithMetadata(
			personerror.PersonCreateFailed,
			err,
			mappers.ToPersonErrorMetadata(gid),
		)
	}

	return mappers.ToPersonResponse(person), nil
}

func (service *PersonService) GetByGID(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	gid := mappers.ToPersonGID(request)
	person, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}

	return mappers.ToPersonResponse(person), nil
}

func (service *PersonService) Deactivate(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	gid := mappers.ToPersonGID(request)
	person, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}

	if person.GetStatus() == entity.PersonStatusDeactivated {
		return mappers.ToPersonResponse(person), nil
	}

	updatedPerson := person.
		SetStatus(entity.PersonStatusDeactivated).
		SetVersion(person.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())

	if err := service.repository.Update(ctx, updatedPerson); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonNotFound) {
			return foundationresponse.Person{}, coreerrors.NewErrorWithMetadata(
				personerror.PersonNotFound,
				mappers.ToPersonErrorMetadata(gid),
			)
		}

		return foundationresponse.Person{}, coreerrors.WrapWithMetadata(
			personerror.PersonUpdateFailed,
			err,
			mappers.ToPersonErrorMetadata(gid),
		)
	}

	return mappers.ToPersonResponse(updatedPerson), nil
}

func (service *PersonService) getPersonByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	if err := validateGID(gid); err != nil {
		return entity.Person{}, err
	}

	person, err := service.repository.FindByGID(ctx, gid)
	if errors.Is(err, foundationrepository.ErrPersonNotFound) {
		return entity.Person{}, coreerrors.NewErrorWithMetadata(
			personerror.PersonNotFound,
			mappers.ToPersonErrorMetadata(gid),
		)
	}
	if err != nil {
		return entity.Person{}, coreerrors.WrapWithMetadata(
			personerror.PersonReadFailed,
			err,
			mappers.ToPersonErrorMetadata(gid),
		)
	}

	return person, nil
}

func validateGID(gid entity.GID) error {
	value, err := uuid.Parse(string(gid))
	if err != nil || value == uuid.Nil || value.Version() != uuid.Version(7) {
		return coreerrors.NewErrorWithMetadata(
			personerror.InvalidPersonGID,
			mappers.ToPersonErrorMetadata(gid),
		)
	}

	return nil
}
