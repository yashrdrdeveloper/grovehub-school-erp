package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/id"
	coreidempotency "grovehub/libs/core/idempotency"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
	foundationmappers "grovehub/services/foundation/mappers"
	foundationrepository "grovehub/services/foundation/repository"
	foundationservices "grovehub/services/foundation/services"
	servicehelpers "grovehub/services/foundation/services/helpers"
	foundationvalidators "grovehub/services/foundation/validators"
)

type PersonService struct {
	repository foundationrepository.PersonRepository
	mapper     foundationmappers.PersonMapper
	validator  foundationvalidators.PersonValidator
}

func NewPersonService(
	repository foundationrepository.PersonRepository,
	mapper foundationmappers.PersonMapper,
	validator foundationvalidators.PersonValidator,
) foundationservices.PersonService {
	return &PersonService{repository: repository, mapper: mapper, validator: validator}
}

func (service *PersonService) Create(
	ctx context.Context,
	request foundationrequest.CreatePerson,
) (foundationresponse.Person, bool, error) {
	validatedRequest, err := service.validator.ValidateCreate(request)
	if err != nil {
		return foundationresponse.Person{}, false, err
	}
	contacts := servicehelpers.BuildContacts(validatedRequest.Contacts, nil)

	idempotencyKey := ""
	requestHash := ""
	if request.IdempotencyKey != "" {
		idempotencyKey, err = coreidempotency.ValidateKey(request.IdempotencyKey)
		if err != nil {
			return foundationresponse.Person{}, false, coreerrors.NewError(personerror.InvalidPersonIdempotencyKey)
		}
		requestHash, err = coreidempotency.Fingerprint(validatedRequest)
		if err != nil {
			return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonCreateFailed, err)
		}
		existing, existingContacts, readErr := service.repository.FindByCreateIdempotencyKey(ctx, idempotencyKey)
		if readErr == nil {
			if servicehelpers.NullableString(existing.RequestHash) != requestHash {
				return foundationresponse.Person{}, false, coreerrors.NewError(personerror.PersonIdempotencyConflict)
			}
			return service.mapper.ToPersonResponse(existing, existingContacts), false, nil
		}
		if !errors.Is(readErr, foundationrepository.ErrPersonNotFound) {
			return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonReadFailed, readErr)
		}
	}
	if len(contacts) == 0 && idempotencyKey == "" {
		return foundationresponse.Person{}, false, coreerrors.NewError(personerror.InvalidPersonIdempotencyKey)
	}

	existingPerson, existingContacts, found, err := service.findExistingPersonByContacts(ctx, contacts)
	if err != nil || found {
		if found {
			return service.mapper.ToPersonResponse(existingPerson, existingContacts), false, nil
		}
		return foundationresponse.Person{}, false, err
	}

	gidValue, err := coreid.NewUUIDv7()
	if err != nil {
		return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonCreateFailed, err)
	}
	gid := gidValue
	now := time.Now().UTC()
	person := servicehelpers.NewPerson(gid, validatedRequest, now)
	if idempotencyKey != "" {
		person.IdempotencyKey = &idempotencyKey
		person.RequestHash = &requestHash
	}

	contacts = servicehelpers.AssociateContacts(person, contacts)
	if err := service.repository.Create(ctx, person, contacts); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonIdempotencyKeyAlreadyExists) {
			existing, existingContacts, readErr := service.repository.FindByCreateIdempotencyKey(ctx, idempotencyKey)
			if readErr != nil {
				return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonReadFailed, readErr)
			}
			if servicehelpers.NullableString(existing.RequestHash) != requestHash {
				return foundationresponse.Person{}, false, coreerrors.NewError(personerror.PersonIdempotencyConflict)
			}
			return service.mapper.ToPersonResponse(existing, existingContacts), false, nil
		}
		if errors.Is(err, foundationrepository.ErrPersonContactAlreadyExists) {
			return service.resolveConcurrentCreate(ctx, contacts)
		}
		return foundationresponse.Person{}, false, coreerrors.WrapWithMetadata(
			personerror.PersonCreateFailed,
			err,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	return service.mapper.ToPersonResponse(person, contacts), true, nil
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


func (service *PersonService) Update(
	ctx context.Context,
	personRequest foundationrequest.PersonByGID,
	request foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	gid := service.mapper.ToPersonGID(personRequest)
	person, contacts, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}
	validatedRequest, err := service.validator.ValidateUpdate(request)
	if err != nil {
		return foundationresponse.Person{}, err
	}
	updatedPerson, updatedContacts, changed, err := service.applyPersonUpdate(ctx, person, contacts, validatedRequest)
	if err != nil {
		return foundationresponse.Person{}, err
	}
	if !changed {
		return service.mapper.ToPersonResponse(person, contacts), nil
	}
	updatedPerson.Version = person.Version + 1
	updatedPerson.UpdatedAt = time.Now().UTC()
	updatedContacts = servicehelpers.AssociateContacts(updatedPerson, updatedContacts)
	if err := service.repository.Update(ctx, updatedPerson, updatedContacts); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonVersionConflict) {
			latest, latestContacts, readErr := service.getPersonByGID(ctx, gid)
			if readErr != nil {
				return foundationresponse.Person{}, readErr
			}
			_, _, stillChanged, reconcileErr := service.applyPersonUpdate(ctx, latest, latestContacts, validatedRequest)
			if reconcileErr != nil {
				return foundationresponse.Person{}, reconcileErr
			}
			if !stillChanged {
				return service.mapper.ToPersonResponse(latest, latestContacts), nil
			}
		}
		return foundationresponse.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson, updatedContacts), nil
}


func (service *PersonService) Search(
	ctx context.Context,
	request coresearch.Request,
) (coresearch.Result[foundationresponse.Person], error) {
	query, cursorPayload, err := service.validator.ValidateSearch(request)
	if err != nil {
		return coresearch.Result[foundationresponse.Person]{}, err
	}
	result, contactsByGID, err := service.repository.Search(ctx, query)
	if err != nil {
		if errors.Is(err, coresearch.ErrInvalidQuery) {
			return coresearch.Result[foundationresponse.Person]{}, coreerrors.NewError(personerror.InvalidPersonSearch)
		}
		return coresearch.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err)
	}
	items := make([]foundationresponse.Person, 0, len(result.Items))
	for _, person := range result.Items {
		items = append(items, service.mapper.ToPersonResponse(person, contactsByGID[person.GID]))
	}
	nextCursor, err := buildNextCursor(result, query, cursorPayload)
	if err != nil {
		return coresearch.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err)
	}
	return coresearch.Result[foundationresponse.Person]{
		Items: items,
		PageInfo: coresearch.PageResponse{
			Rows: query.Limit, NextCursor: nextCursor, HasNext: result.HasNext,
		},
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
	if err := service.repository.Update(ctx, updatedPerson, contacts); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonVersionConflict) {
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
	if request.FirstName != nil && *request.FirstName != person.FirstName {
		updatedPerson.FirstName = *request.FirstName
		changed = true
	}
	if request.MiddleName != nil && *request.MiddleName != servicehelpers.NullableString(person.MiddleName) {
		updatedPerson.MiddleName = servicehelpers.OptionalString(*request.MiddleName)
		changed = true
	}
	if request.LastName != nil && *request.LastName != person.LastName {
		updatedPerson.LastName = *request.LastName
		changed = true
	}
	if request.Contacts == nil {
		return updatedPerson, updatedContacts, changed, nil
	}
	newContacts := servicehelpers.BuildContacts(*request.Contacts, contacts)
	if err := service.ensureContactsBelongToPerson(ctx, person.GID, newContacts); err != nil {
		return entity.Person{}, nil, false, err
	}
	if !servicehelpers.ContactsEqual(newContacts, contacts) {
		updatedContacts = servicehelpers.CopyContacts(newContacts)
		changed = true
	}
	return updatedPerson, updatedContacts, changed, nil
}


func (service *PersonService) getPersonByGID(ctx context.Context, gid string) (entity.Person, []entity.Contact, error) {
	if err := service.validator.ValidateGID(gid); err != nil {
		return entity.Person{}, nil, coreerrors.NewErrorWithMetadata(
			personerror.InvalidPersonGID,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	person, contacts, err := service.repository.FindByGID(ctx, gid)
	if errors.Is(err, foundationrepository.ErrPersonNotFound) {
		return entity.Person{}, nil, coreerrors.NewErrorWithMetadata(
			personerror.PersonNotFound,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	if err != nil {
		return entity.Person{}, nil, coreerrors.WrapWithMetadata(
			personerror.PersonReadFailed,
			err,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	return person, contacts, nil
}


func (service *PersonService) findExistingPersonByContacts(
	ctx context.Context,
	contacts []entity.Contact,
) (entity.Person, []entity.Contact, bool, error) {
	if len(contacts) == 0 {
		return entity.Person{}, nil, false, nil
	}
	gids, err := service.repository.FindGIDsByContacts(ctx, contacts)
	if err != nil {
		return entity.Person{}, nil, false, coreerrors.Wrap(personerror.PersonReadFailed, err)
	}
	if len(gids) > 1 {
		return entity.Person{}, nil, false, coreerrors.NewError(personerror.PersonContactConflict)
	}
	if len(gids) == 0 {
		return entity.Person{}, nil, false, nil
	}
	person, storedContacts, err := service.repository.FindByGID(ctx, gids[0])
	if err != nil {
		return entity.Person{}, nil, false, coreerrors.Wrap(personerror.PersonReadFailed, err)
	}
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
	gids, err := service.repository.FindGIDsByContacts(ctx, contacts)
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
		foundationrepository.ErrPersonNotFound: func() error {
			return coreerrors.NewErrorWithMetadata(
				personerror.PersonNotFound,
				service.mapper.ToPersonErrorMetadata(gid),
			)
		},
		foundationrepository.ErrPersonContactAlreadyExists: func() error {
			return coreerrors.NewError(personerror.PersonContactConflict)
		},
		foundationrepository.ErrPersonVersionConflict: func() error {
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
