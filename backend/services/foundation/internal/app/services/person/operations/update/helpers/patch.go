// Package helpers owns PATCH application and contact ownership for Person Update.
package helpers

import (
    "context"

    coreerrors "grovehub/libs/core/errors"
    corehelpers "grovehub/libs/core/helpers"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    servicehelpers "grovehub/services/foundation/internal/app/services/person/helpers"
)

func ApplyPersonUpdate(ctx context.Context, person entity.Person, contacts []entity.Contact, request persondto.UpdatePerson, contactRepository contactrepository.ContactRepository, mapper foundationmappers.PersonMapper) (entity.Person, []entity.Contact, bool, error) {
	updatedPerson := person
	updatedContacts := servicehelpers.CopyContacts(contacts)
	changed := false
	if request.FirstName.Set && request.FirstName.Value != nil && *request.FirstName.Value != person.FirstName {
		updatedPerson.FirstName = *request.FirstName.Value
		changed = true
	}
	if request.MiddleName.Set {
		if request.MiddleName.Value == nil {
			if person.MiddleName != nil {
				updatedPerson.MiddleName = nil
				changed = true
			}
		} else if *request.MiddleName.Value != corehelpers.NullableString(person.MiddleName) {
			updatedPerson.MiddleName = corehelpers.OptionalString(*request.MiddleName.Value)
			changed = true
		}
	}
	if request.LastName.Set && request.LastName.Value != nil && *request.LastName.Value != person.LastName {
		updatedPerson.LastName = *request.LastName.Value
		changed = true
	}
	if !request.Contacts.Set {
		return updatedPerson, updatedContacts, changed, nil
	}
	var newContacts []entity.Contact
	if request.Contacts.Value != nil {
		newContacts = servicehelpers.BuildContacts(*request.Contacts.Value, contacts)
		if err := ensureContactsBelongToPerson(ctx, person.GID, newContacts, contactRepository, mapper); err != nil {
			return entity.Person{}, nil, false, err
		}
	}
	if !ContactsEqual(newContacts, contacts) {
		updatedContacts = servicehelpers.CopyContacts(newContacts)
		changed = true
	}
	return updatedPerson, updatedContacts, changed, nil
}

func ensureContactsBelongToPerson(ctx context.Context, gid string, contacts []entity.Contact, contactRepository contactrepository.ContactRepository, mapper foundationmappers.PersonMapper) error {
	if len(contacts) == 0 {
		return nil
	}
	gids, err := contactRepository.FindPersonGIDsByContacts(ctx, contacts)
	if err != nil {
		return coreerrors.WrapWithMetadata(
			personerror.PersonReadFailed,
			err,
			mapper.ToPersonErrorMetadata(gid),
		)
	}
	for _, matchedGID := range gids {
		if matchedGID != gid {
			return coreerrors.NewError(personerror.PersonContactConflict)
		}
	}
	return nil
}
