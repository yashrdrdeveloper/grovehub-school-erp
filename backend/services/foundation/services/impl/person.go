package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/id"
	coreidempotency "grovehub/libs/core/idempotency"
	headercontext "grovehub/libs/core/headercontext"
	coresearch "grovehub/libs/core/search"
	coretransaction "grovehub/libs/core/database/transaction"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request/person"
	foundationresponse "grovehub/services/foundation/dto/response/person"
	"grovehub/services/foundation/entity"
	foundationmappers "grovehub/services/foundation/mappers"
	foundationrepository "grovehub/services/foundation/repository"
	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
	foundationservices "grovehub/services/foundation/services"
	servicehelpers "grovehub/services/foundation/services/helpers"
	foundationvalidators "grovehub/services/foundation/validators"
)

type PersonService struct {
	personRepository foundationrepository.PersonRepository
	contactRepository foundationrepository.ContactRepository
	unitOfWork coretransaction.UnitOfWork
	mapper foundationmappers.PersonMapper
	validator foundationvalidators.PersonValidator
}

func NewPersonService(
	personRepository foundationrepository.PersonRepository,
	contactRepository foundationrepository.ContactRepository,
	unitOfWork coretransaction.UnitOfWork,
	mapper foundationmappers.PersonMapper,
	validator foundationvalidators.PersonValidator,
) foundationservices.PersonService {
	return &PersonService{
		personRepository: personRepository,
		contactRepository: contactRepository,
		unitOfWork: unitOfWork,
		mapper: mapper,
		validator: validator,
	}
}


func (service *PersonService) Create(ctx context.Context, request foundationrequest.CreatePerson) (foundationresponse.Person, bool, error) {
	validatedRequest, err := service.validator.ValidateCreate(request)
	if err != nil { return foundationresponse.Person{}, false, err }
	contacts := servicehelpers.BuildContacts(validatedRequest.Contacts, nil)
	idempotencyKey, requestHash := "", ""
	if headercontext.IdempotencyKey(ctx) != "" {
		idempotencyKey, err = coreidempotency.ValidateKey(headercontext.IdempotencyKey(ctx))
		if err != nil { return foundationresponse.Person{}, false, coreerrors.NewError(personerror.InvalidPersonIdempotencyKey) }
		requestHash, err = coreidempotency.Fingerprint(validatedRequest)
		if err != nil { return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonCreateFailed, err) }
		existing, storedContacts, readErr := service.findPersonByIdempotencyKey(ctx, idempotencyKey)
		if readErr == nil {
			if servicehelpers.NullableString(existing.RequestHash) != requestHash {
				return foundationresponse.Person{}, false, coreerrors.NewError(personerror.PersonIdempotencyConflict)
			}
			return service.mapper.ToPersonResponse(existing, storedContacts), false, nil
		}
		if !errors.Is(readErr, repositoryerrors.ErrPersonNotFound) {
			return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonReadFailed, readErr)
		}
	}
	if len(contacts) == 0 && idempotencyKey == "" {
		return foundationresponse.Person{}, false, coreerrors.NewError(personerror.InvalidPersonIdempotencyKey)
	}
	existingPerson, existingContacts, found, err := service.findExistingPersonByContacts(ctx, contacts)
	if err != nil || found {
		if found { return service.mapper.ToPersonResponse(existingPerson, existingContacts), false, nil }
		return foundationresponse.Person{}, false, err
	}
	gid, err := coreid.NewUUIDv7()
	if err != nil { return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonCreateFailed, err) }
	person := servicehelpers.NewPerson(gid, validatedRequest, time.Now().UTC())
	if idempotencyKey != "" { person.IdempotencyKey, person.RequestHash = &idempotencyKey, &requestHash }
	contacts = servicehelpers.AssociateContacts(person, contacts)
	person.SearchDetails = servicehelpers.BuildPersonSearchDetails(person, contacts)
	err = service.unitOfWork.Run(ctx, func(txCtx context.Context) error {
		if err := service.personRepository.Create(txCtx, person); err != nil { return err }
		return service.contactRepository.ReplaceForPerson(txCtx, person.GID, contacts)
	})
	if err != nil {
		if errors.Is(err, repositoryerrors.ErrPersonIdempotencyKeyAlreadyExists) {
			existing, storedContacts, readErr := service.findPersonByIdempotencyKey(ctx, idempotencyKey)
			if readErr != nil { return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonReadFailed, readErr) }
			if servicehelpers.NullableString(existing.RequestHash) != requestHash {
				return foundationresponse.Person{}, false, coreerrors.NewError(personerror.PersonIdempotencyConflict)
			}
			return service.mapper.ToPersonResponse(existing, storedContacts), false, nil
		}
		if errors.Is(err, repositoryerrors.ErrContactAlreadyExists) { return service.resolveConcurrentCreate(ctx, contacts) }
		return foundationresponse.Person{}, false, coreerrors.WrapWithMetadata(personerror.PersonCreateFailed, err, service.mapper.ToPersonErrorMetadata(gid))
	}
	return service.mapper.ToPersonResponse(person, contacts), true, nil
}

func (service *PersonService) findPersonByIdempotencyKey(ctx context.Context, key string) (entity.Person, []entity.Contact, error) {
	person, err := service.personRepository.FindByIdempotencyKey(ctx, key)
	if err != nil { return entity.Person{}, nil, err }
	contacts, err := service.contactRepository.FindByPersonGID(ctx, person.GID)
	if err != nil { return entity.Person{}, nil, err }
	return person, contacts, nil
}

func (service *PersonService) GetByGID(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	person, contacts, err := service.getPersonByGID(ctx, service.mapper.ToPersonGID(request))
	if err != nil {
		return foundationresponse.Person{}, err
	}
	return service.mapper.ToPersonResponse(person, contacts), nil
}



func (service *PersonService) Update(ctx context.Context, personRequest foundationrequest.PersonByGID, request foundationrequest.UpdatePerson) (foundationresponse.Person, error) {
	gid := service.mapper.ToPersonGID(personRequest)
	person, contacts, err := service.getPersonByGID(ctx, gid)
	if err != nil { return foundationresponse.Person{}, err }
	validatedRequest, err := service.validator.ValidateUpdate(request)
	if err != nil { return foundationresponse.Person{}, err }
	updatedPerson, updatedContacts, changed, err := service.applyPersonUpdate(ctx, person, contacts, validatedRequest)
	if err != nil { return foundationresponse.Person{}, err }
	if !changed { return service.mapper.ToPersonResponse(person, contacts), nil }
	updatedPerson.Version = person.Version + 1
	updatedPerson.UpdatedAt = time.Now().UTC()
	contactsChanged := !servicehelpers.ContactsEqual(contacts, updatedContacts)
	if contactsChanged {
		updatedContacts = servicehelpers.AssociateContacts(updatedPerson, updatedContacts)
	}
	updatedPerson.SearchDetails = servicehelpers.BuildPersonSearchDetails(updatedPerson, updatedContacts)
	err = service.unitOfWork.Run(ctx, func(txCtx context.Context) error {
		if err := service.personRepository.Update(txCtx, updatedPerson); err != nil {
			return err
		}
		if contactsChanged {
			return service.contactRepository.ReplaceForPerson(txCtx, updatedPerson.GID, updatedContacts)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, repositoryerrors.ErrPersonVersionConflict) {
			latest, latestContacts, readErr := service.getPersonByGID(ctx, gid)
			if readErr != nil { return foundationresponse.Person{}, readErr }
			_, _, stillChanged, reconcileErr := service.applyPersonUpdate(ctx, latest, latestContacts, validatedRequest)
			if reconcileErr != nil { return foundationresponse.Person{}, reconcileErr }
			if !stillChanged { return service.mapper.ToPersonResponse(latest, latestContacts), nil }
		}
		return foundationresponse.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson, updatedContacts), nil
}


func (service *PersonService) Search(ctx context.Context, request coresearch.Request) (coresearch.Result[foundationresponse.Person], error) {
	query, cursorPayload, err := service.validator.ValidateSearch(request)
	if err != nil { return coresearch.Result[foundationresponse.Person]{}, err }
	result, err := service.personRepository.Search(ctx, query)
	if err != nil {
		if errors.Is(err, coresearch.ErrInvalidQuery) {
			return coresearch.Result[foundationresponse.Person]{}, coreerrors.NewError(personerror.InvalidPersonSearch)
		}
		return coresearch.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err)
	}
	gids := make([]string, 0, len(result.Items))
	for _, person := range result.Items { gids = append(gids, person.GID) }
	contactsByGID, err := service.contactRepository.FindByPersonGIDs(ctx, gids)
	if err != nil { return coresearch.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err) }
	items := make([]foundationresponse.Person, 0, len(result.Items))
	for _, person := range result.Items {
		items = append(items, service.mapper.ToPersonResponse(person, contactsByGID[person.GID]))
	}
	nextCursor, err := buildNextCursor(result, query, cursorPayload)
	if err != nil { return coresearch.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err) }
	return coresearch.Result[foundationresponse.Person]{
		Items: items,
		PageInfo: coresearch.PageResponse{Rows: query.Limit, NextCursor: nextCursor, HasNext: result.HasNext},
	}, nil
}

func (service *PersonService) Deactivate(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	gid := service.mapper.ToPersonGID(request)
	person, contacts, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}
	if person.Status == entity.PersonStatusDeactivated {
		return service.mapper.ToPersonResponse(person, contacts), nil
	}
	updatedPerson := person
	updatedPerson.Status = entity.PersonStatusDeactivated
	updatedPerson.Version = person.Version + 1
	updatedPerson.UpdatedAt = time.Now().UTC()
	if err := service.personRepository.Update(ctx, updatedPerson); err != nil {
		if errors.Is(err, repositoryerrors.ErrPersonVersionConflict) {
			latest, latestContacts, readErr := service.getPersonByGID(ctx, gid)
			if readErr != nil {
				return foundationresponse.Person{}, readErr
			}
			if latest.Status == entity.PersonStatusDeactivated {
				return service.mapper.ToPersonResponse(latest, latestContacts), nil
			}
		}
		return foundationresponse.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson, contacts), nil
}


func (service *PersonService) applyPersonUpdate(
	ctx context.Context,
	person entity.Person,
	contacts []entity.Contact,
	request foundationrequest.UpdatePerson,
) (entity.Person, []entity.Contact, bool, error) {
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
		} else if *request.MiddleName.Value != servicehelpers.NullableString(person.MiddleName) {
			updatedPerson.MiddleName = servicehelpers.OptionalString(*request.MiddleName.Value)
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
		if err := service.ensureContactsBelongToPerson(ctx, person.GID, newContacts); err != nil {
			return entity.Person{}, nil, false, err
		}
	}
	if !servicehelpers.ContactsEqual(newContacts, contacts) {
		updatedContacts = servicehelpers.CopyContacts(newContacts)
		changed = true
	}
	return updatedPerson, updatedContacts, changed, nil
}



func (service *PersonService) getPersonByGID(ctx context.Context, gid string) (entity.Person, []entity.Contact, error) {
	if err := service.validator.ValidateGID(gid); err != nil {
		return entity.Person{}, nil, coreerrors.NewErrorWithMetadata(personerror.InvalidPersonGID, service.mapper.ToPersonErrorMetadata(gid))
	}
	person, err := service.personRepository.FindByGID(ctx, gid)
	if errors.Is(err, repositoryerrors.ErrPersonNotFound) {
		return entity.Person{}, nil, coreerrors.NewErrorWithMetadata(personerror.PersonNotFound, service.mapper.ToPersonErrorMetadata(gid))
	}
	if err != nil {
		return entity.Person{}, nil, coreerrors.WrapWithMetadata(personerror.PersonReadFailed, err, service.mapper.ToPersonErrorMetadata(gid))
	}
	contacts, err := service.contactRepository.FindByPersonGID(ctx, gid)
	if err != nil {
		return entity.Person{}, nil, coreerrors.WrapWithMetadata(personerror.PersonReadFailed, err, service.mapper.ToPersonErrorMetadata(gid))
	}
	return person, contacts, nil
}


func (service *PersonService) findExistingPersonByContacts(ctx context.Context, contacts []entity.Contact) (entity.Person, []entity.Contact, bool, error) {
	if len(contacts) == 0 { return entity.Person{}, nil, false, nil }
	gids, err := service.contactRepository.FindPersonGIDsByContacts(ctx, contacts)
	if err != nil { return entity.Person{}, nil, false, coreerrors.Wrap(personerror.PersonReadFailed, err) }
	if len(gids) > 1 { return entity.Person{}, nil, false, coreerrors.NewError(personerror.PersonContactConflict) }
	if len(gids) == 0 { return entity.Person{}, nil, false, nil }
	person, storedContacts, err := service.getPersonByGID(ctx, gids[0])
	if err != nil { return entity.Person{}, nil, false, coreerrors.Wrap(personerror.PersonReadFailed, err) }
	return person, storedContacts, true, nil
}

func (service *PersonService) ensureContactsBelongToPerson(
	ctx context.Context,
	gid string,
	contacts []entity.Contact,
) error {
	if len(contacts) == 0 {
		return nil
	}
	gids, err := service.contactRepository.FindPersonGIDsByContacts(ctx, contacts)
	if err != nil {
		return coreerrors.WrapWithMetadata(
			personerror.PersonReadFailed,
			err,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	for _, matchedGID := range gids {
		if matchedGID != gid {
			return coreerrors.NewError(personerror.PersonContactConflict)
		}
	}
	return nil
}


func (service *PersonService) resolveConcurrentCreate(
	ctx context.Context,
	contacts []entity.Contact,
) (foundationresponse.Person, bool, error) {
	existingPerson, existingContacts, found, err := service.findExistingPersonByContacts(ctx, contacts)
	if err != nil {
		return foundationresponse.Person{}, false, err
	}
	if found {
		return service.mapper.ToPersonResponse(existingPerson, existingContacts), false, nil
	}
	return foundationresponse.Person{}, false, coreerrors.NewError(personerror.PersonCreateFailed)
}

func (service *PersonService) mapUpdateError(gid string, err error) error {
	mapped := map[error]func() error{
		repositoryerrors.ErrPersonNotFound: func() error {
			return coreerrors.NewErrorWithMetadata(
				personerror.PersonNotFound,
				service.mapper.ToPersonErrorMetadata(gid),
			)
		},
		repositoryerrors.ErrContactAlreadyExists: func() error {
			return coreerrors.NewError(personerror.PersonContactConflict)
		},
		repositoryerrors.ErrPersonVersionConflict: func() error {
			return coreerrors.NewErrorWithMetadata(
				personerror.PersonConcurrentModification,
				service.mapper.ToPersonErrorMetadata(gid),
			)
		},
	}
	for repositoryError, factory := range mapped {
		if errors.Is(err, repositoryError) {
			return factory()
		}
	}
	return coreerrors.WrapWithMetadata(
		personerror.PersonUpdateFailed,
		err,
		service.mapper.ToPersonErrorMetadata(gid),
	)
}

func buildNextCursor(
	result coresearch.QueryResult[entity.Person],
	query coresearch.Query[string],
	cursor coresearch.CursorPayload,
) (string, error) {
	if !result.HasNext || len(result.Items) == 0 {
		return "", nil
	}
	lastPerson := result.Items[len(result.Items)-1]
	cursor.Value = servicehelpers.PersonCursorValue(lastPerson, query.Sort.Field)
	cursor.CursorId = lastPerson.GID
	return servicehelpers.EncodePersonSearchCursor(cursor)
}
