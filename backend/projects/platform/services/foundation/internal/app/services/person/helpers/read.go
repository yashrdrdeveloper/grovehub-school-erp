package helpers

import (
    "context"
    "errors"

    coreerrors "grovehub/libs/core/errors"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    "grovehub/services/foundation/internal/app/entity"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    validatorhelpers "grovehub/services/foundation/internal/app/validators/person/helpers"
)

// ReadPerson loads the Person aggregate for GetByGID, Update and Deactivate.
func ReadPerson(ctx context.Context, gid string, personRepository personrepository.PersonRepository, contactRepository contactrepository.ContactRepository, mapper foundationmappers.PersonMapper) (entity.Person, []entity.Contact, error) {
	if err := validatorhelpers.ValidateGID(gid); err != nil {
		return entity.Person{}, nil, coreerrors.NewErrorWithMetadata(personerror.InvalidPersonGID, mapper.ToPersonErrorMetadata(gid))
	}
	person, err := personRepository.FindByGID(ctx, gid)
	if errors.Is(err, repositoryerrors.ErrPersonNotFound) {
		return entity.Person{}, nil, coreerrors.NewErrorWithMetadata(personerror.PersonNotFound, mapper.ToPersonErrorMetadata(gid))
	}
	if err != nil {
		return entity.Person{}, nil, coreerrors.WrapWithMetadata(personerror.PersonReadFailed, err, mapper.ToPersonErrorMetadata(gid))
	}
	contacts, err := contactRepository.FindByPersonGID(ctx, gid)
	if err != nil {
		return entity.Person{}, nil, coreerrors.WrapWithMetadata(personerror.PersonReadFailed, err, mapper.ToPersonErrorMetadata(gid))
	}
	return person, contacts, nil
}
