package impl

import (
	"context"
	"errors"
	"time"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	coresearch "grovehub/libs/core/search"
	coretransaction "grovehub/libs/core/database/transaction"
    servicepipeline "grovehub/libs/core/pipeline/service"
	personerror "grovehub/services/foundation/constants/error"
	persondto "grovehub/services/foundation/dto/person"
	"grovehub/services/foundation/entity"
	foundationmappers "grovehub/services/foundation/mappers/person"
	personrepository "grovehub/services/foundation/repository/person"
	contactrepository "grovehub/services/foundation/repository/contact"
	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
	foundationservices "grovehub/services/foundation/services/person"
    personflow "grovehub/services/foundation/services/person/flows"
    createstrategies "grovehub/services/foundation/services/person/strategies"
    createimpl "grovehub/services/foundation/services/person/strategies/impl"
	servicehelpers "grovehub/services/foundation/services/person/helpers"
    corehelpers "grovehub/libs/core/helpers"
	foundationvalidators "grovehub/services/foundation/validators/person"
	validatorhelpers "grovehub/services/foundation/validators/person/helpers"
)

// @Service
type PersonService struct {
	personRepository personrepository.PersonRepository
	contactRepository contactrepository.ContactRepository
	unitOfWork coretransaction.UnitOfWork
	mapper foundationmappers.PersonMapper
	validator foundationvalidators.PersonValidator
	creator createstrategies.CreateStrategy
}

func NewPersonService(
	personRepository personrepository.PersonRepository,
	contactRepository contactrepository.ContactRepository,
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
        creator: createimpl.NewCreateStrategy(personRepository, contactRepository, unitOfWork, mapper),
	}
}


// @ChainOfResponsibility
func (service *PersonService) Create(ctx context.Context, request persondto.CreatePersonRequest) (persondto.Person, bool, error) {
    flow := personflow.Create{Request: request}
    return servicepipeline.Run(ctx, &flow,
        service.creator.Prepare,
        service.creator.Idempotency,
        service.creator.Contacts,
        service.creator.Persist,
    )
}

func (service *PersonService) GetByGID(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error) {
	person, contacts, err := service.getPersonByGID(ctx, service.mapper.ToPersonGID(request.Params))
	if err != nil {
		return persondto.Person{}, err
	}
	return service.mapper.ToPersonResponse(person, contacts), nil
}



func (service *PersonService) Update(ctx context.Context, request persondto.UpdatePersonRequest) (persondto.Person, error) {
	gid := service.mapper.ToPersonGID(request.Params)
	person, contacts, err := service.getPersonByGID(ctx, gid)
	if err != nil { return persondto.Person{}, err }
	validatedRequest, err := annotationvalidation.Process(request, service.validator.UpdateRules())
	if err != nil { return persondto.Person{}, err }
	updatedPerson, updatedContacts, changed, err := service.applyPersonUpdate(ctx, person, contacts, validatedRequest.RequestBody)
	if err != nil { return persondto.Person{}, err }
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
			if readErr != nil { return persondto.Person{}, readErr }
			_, _, stillChanged, reconcileErr := service.applyPersonUpdate(ctx, latest, latestContacts, validatedRequest.RequestBody)
			if reconcileErr != nil { return persondto.Person{}, reconcileErr }
			if !stillChanged { return service.mapper.ToPersonResponse(latest, latestContacts), nil }
		}
		return persondto.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson, updatedContacts), nil
}


func (service *PersonService) Search(ctx context.Context, request persondto.SearchPersonsRequest) (coresearch.Result[persondto.Person], error) {
	query, cursorPayload, err := validatorhelpers.NormalizePersonSearchRequest(request.RequestBody)
	if err != nil { return coresearch.Result[persondto.Person]{}, err }
	result, err := service.personRepository.Search(ctx, query)
	if err != nil {
		if errors.Is(err, coresearch.ErrInvalidQuery) {
			return coresearch.Result[persondto.Person]{}, coreerrors.NewError(personerror.InvalidPersonSearch)
		}
		return coresearch.Result[persondto.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err)
	}
	gids := make([]string, 0, len(result.Items))
	for _, person := range result.Items { gids = append(gids, person.GID) }
	contactsByGID, err := service.contactRepository.FindByPersonGIDs(ctx, gids)
	if err != nil { return coresearch.Result[persondto.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err) }
	items := make([]persondto.Person, 0, len(result.Items))
	for _, person := range result.Items {
		items = append(items, service.mapper.ToPersonResponse(person, contactsByGID[person.GID]))
	}
	nextCursor, err := buildNextCursor(result, query, cursorPayload)
	if err != nil { return coresearch.Result[persondto.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err) }
	return coresearch.Result[persondto.Person]{
		Items: items,
		PageInfo: coresearch.PageResponse{Rows: query.Limit, NextCursor: nextCursor, HasNext: result.HasNext},
	}, nil
}

func (service *PersonService) Deactivate(ctx context.Context, request persondto.DeactivatePersonRequest) (persondto.Person, error) {
	gid := service.mapper.ToPersonGID(request.Params)
	person, contacts, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return persondto.Person{}, err
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
				return persondto.Person{}, readErr
			}
			if latest.Status == entity.PersonStatusDeactivated {
				return service.mapper.ToPersonResponse(latest, latestContacts), nil
			}
		}
		return persondto.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson, contacts), nil
}


func (service *PersonService) applyPersonUpdate(
	ctx context.Context,
	person entity.Person,
	contacts []entity.Contact,
	request persondto.UpdatePerson,
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
	if err := validatorhelpers.ValidateGID(gid); err != nil {
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
	return coresearch.EncodeCursor(cursor)
}
