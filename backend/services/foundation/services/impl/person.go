package impl

import (
	"context"
	"errors"
	"time"

	searchcontract "grovehub/libs/core/contracts/search"
	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
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
	return &PersonService{
		repository: repository,
		mapper:     mapper,
		validator:  validator,
	}
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

	existingPerson, found, err := service.findExistingPersonByContacts(ctx, contacts)
	if err != nil || found {
		if found {
			return service.mapper.ToPersonResponse(existingPerson), false, nil
		}
		return foundationresponse.Person{}, false, err
	}

	gidValue, err := coreid.NewUUIDv7()
	if err != nil {
		return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonCreateFailed, err)
	}

	gid := entity.GID(gidValue)
	now := time.Now().UTC()
	person := entity.NewPerson(
		gid,
		validatedRequest.FirstName,
		validatedRequest.MiddleName,
		validatedRequest.LastName,
		contacts,
		entity.PersonStatusActive,
		1,
		now,
		now,
	)

	if err := service.repository.Create(ctx, person); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonContactAlreadyExists) {
			return service.resolveConcurrentCreate(ctx, contacts)
		}
		return foundationresponse.Person{}, false, coreerrors.WrapWithMetadata(
			personerror.PersonCreateFailed,
			err,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}

	return service.mapper.ToPersonResponse(person), true, nil
}

func (service *PersonService) GetByGID(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	person, err := service.getPersonByGID(ctx, service.mapper.ToPersonGID(request))
	if err != nil {
		return foundationresponse.Person{}, err
	}
	return service.mapper.ToPersonResponse(person), nil
}

func (service *PersonService) Update(
	ctx context.Context,
	personRequest foundationrequest.PersonByGID,
	request foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	gid := service.mapper.ToPersonGID(personRequest)
	person, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}

	validatedRequest, err := service.validator.ValidateUpdate(request)
	if err != nil {
		return foundationresponse.Person{}, err
	}

	updatedPerson, changed, err := service.applyPersonUpdate(ctx, person, validatedRequest)
	if err != nil {
		return foundationresponse.Person{}, err
	}
	if !changed {
		return service.mapper.ToPersonResponse(person), nil
	}

	updatedPerson = updatedPerson.
		SetVersion(person.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())

	if err := service.repository.Update(ctx, updatedPerson); err != nil {
		return foundationresponse.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson), nil
}

func (service *PersonService) Search(
	ctx context.Context,
	request searchcontract.Request,
) (searchcontract.Result[foundationresponse.Person], error) {
	query, cursorPayload, err := service.validator.ValidateSearch(request)
	if err != nil {
		return searchcontract.Result[foundationresponse.Person]{}, err
	}

	result, err := service.repository.Search(ctx, query)
	if err != nil {
		if errors.Is(err, coresearch.ErrInvalidQuery) {
			return searchcontract.Result[foundationresponse.Person]{}, coreerrors.NewError(personerror.InvalidPersonSearch)
		}
		return searchcontract.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err)
	}

	items := make([]foundationresponse.Person, 0, len(result.Items))
	for _, person := range result.Items {
		items = append(items, service.mapper.ToPersonResponse(person))
	}

	nextCursor, err := buildNextCursor(result, query, cursorPayload)
	if err != nil {
		return searchcontract.Result[foundationresponse.Person]{}, coreerrors.Wrap(personerror.PersonSearchFailed, err)
	}

	return searchcontract.Result[foundationresponse.Person]{
		Items: items,
		PageInfo: searchcontract.PageResponse{
			Rows:       query.Limit,
			NextCursor: nextCursor,
			HasNext:    result.HasNext,
		},
	}, nil
}

func (service *PersonService) Deactivate(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	gid := service.mapper.ToPersonGID(request)
	person, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}
	if person.GetStatus() == entity.PersonStatusDeactivated {
		return service.mapper.ToPersonResponse(person), nil
	}

	updatedPerson := person.
		SetStatus(entity.PersonStatusDeactivated).
		SetVersion(person.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())

	if err := service.repository.Update(ctx, updatedPerson); err != nil {
		return foundationresponse.Person{}, service.mapUpdateError(gid, err)
	}
	return service.mapper.ToPersonResponse(updatedPerson), nil
}

func (service *PersonService) applyPersonUpdate(
	ctx context.Context,
	person entity.Person,
	request foundationrequest.UpdatePerson,
) (entity.Person, bool, error) {
	updatedPerson := person
	changed := false

	nameUpdates := []struct {
		value   *string
		current string
		apply   func(entity.Person, string) entity.Person
	}{
		{request.FirstName, person.GetFirstName(), func(value entity.Person, name string) entity.Person { return value.SetFirstName(name) }},
		{request.MiddleName, person.GetMiddleName(), func(value entity.Person, name string) entity.Person { return value.SetMiddleName(name) }},
		{request.LastName, person.GetLastName(), func(value entity.Person, name string) entity.Person { return value.SetLastName(name) }},
	}

	for _, update := range nameUpdates {
		if update.value == nil || *update.value == update.current {
			continue
		}
		updatedPerson = update.apply(updatedPerson, *update.value)
		changed = true
	}

	if request.Contacts == nil {
		return updatedPerson, changed, nil
	}

	contacts := servicehelpers.BuildContacts(*request.Contacts, person.GetContacts())
	if err := service.ensureContactsBelongToPerson(ctx, person.GetGID(), contacts); err != nil {
		return entity.Person{}, false, err
	}
	if !servicehelpers.ContactsEqual(contacts, person.GetContacts()) {
		updatedPerson = updatedPerson.SetContacts(contacts)
		changed = true
	}
	return updatedPerson, changed, nil
}

func (service *PersonService) getPersonByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	if err := service.validator.ValidateGID(gid); err != nil {
		return entity.Person{}, coreerrors.NewErrorWithMetadata(
			personerror.InvalidPersonGID,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}

	person, err := service.repository.FindByGID(ctx, gid)
	if errors.Is(err, foundationrepository.ErrPersonNotFound) {
		return entity.Person{}, coreerrors.NewErrorWithMetadata(
			personerror.PersonNotFound,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	if err != nil {
		return entity.Person{}, coreerrors.WrapWithMetadata(
			personerror.PersonReadFailed,
			err,
			service.mapper.ToPersonErrorMetadata(gid),
		)
	}
	return person, nil
}

func (service *PersonService) findExistingPersonByContacts(
	ctx context.Context,
	contacts []entity.Contact,
) (entity.Person, bool, error) {
	if len(contacts) == 0 {
		return entity.Person{}, false, nil
	}

	gids, err := service.repository.FindGIDsByContacts(ctx, contacts)
	if err != nil {
		return entity.Person{}, false, coreerrors.Wrap(personerror.PersonReadFailed, err)
	}
	if len(gids) > 1 {
		return entity.Person{}, false, coreerrors.NewError(personerror.PersonContactConflict)
	}
	if len(gids) == 0 {
		return entity.Person{}, false, nil
	}

	person, err := service.repository.FindByGID(ctx, gids[0])
	if err != nil {
		return entity.Person{}, false, coreerrors.Wrap(personerror.PersonReadFailed, err)
	}
	return person, true, nil
}

func (service *PersonService) ensureContactsBelongToPerson(
	ctx context.Context,
	gid entity.GID,
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
	existingPerson, found, err := service.findExistingPersonByContacts(ctx, contacts)
	if err != nil {
		return foundationresponse.Person{}, false, err
	}
	if found {
		return service.mapper.ToPersonResponse(existingPerson), false, nil
	}
	return foundationresponse.Person{}, false, coreerrors.NewError(personerror.PersonCreateFailed)
}

func (service *PersonService) mapUpdateError(gid entity.GID, err error) error {
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
	result coresearch.Result[entity.Person],
	query coresearch.Query[entity.GID],
	cursor coresearch.CursorPayload,
) (string, error) {
	if !result.HasNext || len(result.Items) == 0 {
		return "", nil
	}
	lastPerson := result.Items[len(result.Items)-1]
	cursor.Value = servicehelpers.PersonCursorValue(lastPerson, query.Sort.Field)
	cursor.CursorId = string(lastPerson.GetGID())
	return servicehelpers.EncodePersonSearchCursor(cursor)
}
