package services

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

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

const (
	maxPersonNameLength = 100
	maxContactsPerType  = 10
)

var phonePattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

type PersonService struct {
	repository foundationrepository.PersonRepository
}

func NewPersonService(repository foundationrepository.PersonRepository) *PersonService {
	return &PersonService{repository: repository}
}

func (service *PersonService) Create(
	ctx context.Context,
	request foundationrequest.CreatePerson,
) (foundationresponse.Person, bool, error) {
	firstName, middleName, lastName, err := normalizePersonNames(
		request.FirstName,
		request.MiddleName,
		request.LastName,
	)
	if err != nil {
		return foundationresponse.Person{}, false, err
	}

	contacts, err := buildContacts(request.Contacts, nil)
	if err != nil {
		return foundationresponse.Person{}, false, err
	}

	existingPerson, found, err := service.findExistingPersonByContacts(ctx, contacts)
	if err != nil {
		return foundationresponse.Person{}, false, err
	}
	if found {
		return mappers.ToPersonResponse(existingPerson), false, nil
	}

	gidValue, err := coreid.NewUUIDv7()
	if err != nil {
		return foundationresponse.Person{}, false, coreerrors.Wrap(personerror.PersonCreateFailed, err)
	}

	gid := entity.GID(gidValue)
	now := time.Now().UTC()
	person := entity.NewPerson(
		gid,
		firstName,
		middleName,
		lastName,
		contacts,
		entity.PersonStatusActive,
		1,
		now,
		now,
	)

	if err := service.repository.Create(ctx, person); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonContactAlreadyExists) {
			existingPerson, found, lookupErr := service.findExistingPersonByContacts(ctx, contacts)
			if lookupErr != nil {
				return foundationresponse.Person{}, false, lookupErr
			}
			if found {
				return mappers.ToPersonResponse(existingPerson), false, nil
			}
		}

		return foundationresponse.Person{}, false, coreerrors.WrapWithMetadata(
			personerror.PersonCreateFailed,
			err,
			mappers.ToPersonErrorMetadata(gid),
		)
	}

	return mappers.ToPersonResponse(person), true, nil
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

func (service *PersonService) Update(
	ctx context.Context,
	personRequest foundationrequest.PersonByGID,
	request foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	gid := mappers.ToPersonGID(personRequest)
	person, err := service.getPersonByGID(ctx, gid)
	if err != nil {
		return foundationresponse.Person{}, err
	}

	updatedPerson := person
	changed := false

	if request.FirstName != nil {
		firstName, err := normalizeRequiredName(*request.FirstName)
		if err != nil {
			return foundationresponse.Person{}, err
		}
		if firstName != person.GetFirstName() {
			updatedPerson = updatedPerson.SetFirstName(firstName)
			changed = true
		}
	}

	if request.MiddleName != nil {
		middleName, err := normalizeOptionalName(*request.MiddleName)
		if err != nil {
			return foundationresponse.Person{}, err
		}
		if middleName != person.GetMiddleName() {
			updatedPerson = updatedPerson.SetMiddleName(middleName)
			changed = true
		}
	}

	if request.LastName != nil {
		lastName, err := normalizeRequiredName(*request.LastName)
		if err != nil {
			return foundationresponse.Person{}, err
		}
		if lastName != person.GetLastName() {
			updatedPerson = updatedPerson.SetLastName(lastName)
			changed = true
		}
	}

	if request.Contacts != nil {
		contacts, err := buildContacts(*request.Contacts, person.GetContacts())
		if err != nil {
			return foundationresponse.Person{}, err
		}
		if err := service.ensureContactsBelongToPerson(ctx, gid, contacts); err != nil {
			return foundationresponse.Person{}, err
		}
		if !contactsEqual(contacts, person.GetContacts()) {
			updatedPerson = updatedPerson.SetContacts(contacts)
			changed = true
		}
	}

	if !changed {
		return mappers.ToPersonResponse(person), nil
	}

	updatedPerson = updatedPerson.
		SetVersion(person.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())

	if err := service.repository.Update(ctx, updatedPerson); err != nil {
		if errors.Is(err, foundationrepository.ErrPersonNotFound) {
			return foundationresponse.Person{}, coreerrors.NewErrorWithMetadata(
				personerror.PersonNotFound,
				mappers.ToPersonErrorMetadata(gid),
			)
		}
		if errors.Is(err, foundationrepository.ErrPersonContactAlreadyExists) {
			return foundationresponse.Person{}, coreerrors.NewError(
				personerror.PersonContactConflict,
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
			mappers.ToPersonErrorMetadata(gid),
		)
	}

	for _, matchedGID := range gids {
		if matchedGID != gid {
			return coreerrors.NewError(personerror.PersonContactConflict)
		}
	}
	return nil
}

func normalizePersonNames(firstName string, middleName string, lastName string) (string, string, string, error) {
	normalizedFirstName, err := normalizeRequiredName(firstName)
	if err != nil {
		return "", "", "", err
	}
	normalizedMiddleName, err := normalizeOptionalName(middleName)
	if err != nil {
		return "", "", "", err
	}
	normalizedLastName, err := normalizeRequiredName(lastName)
	if err != nil {
		return "", "", "", err
	}
	return normalizedFirstName, normalizedMiddleName, normalizedLastName, nil
}

func normalizeRequiredName(value string) (string, error) {
	normalized, err := normalizeOptionalName(value)
	if err != nil || normalized == "" {
		return "", coreerrors.NewError(personerror.InvalidPersonRequest)
	}
	return normalized, nil
}

func normalizeOptionalName(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	if utf8.RuneCountInString(normalized) > maxPersonNameLength {
		return "", coreerrors.NewError(personerror.InvalidPersonRequest)
	}
	return normalized, nil
}

func buildContacts(
	request foundationrequest.Contacts,
	existing []entity.Contact,
) ([]entity.Contact, error) {
	if len(request.Emails) > maxContactsPerType || len(request.Phones) > maxContactsPerType {
		return nil, coreerrors.NewError(personerror.InvalidPersonContact)
	}

	existingVerification := make(map[string]bool, len(existing))
	for _, contact := range existing {
		existingVerification[contactKey(contact.GetType(), contact.GetNormalizedValue())] = contact.GetIsVerified()
	}

	contacts := make([]entity.Contact, 0, len(request.Emails)+len(request.Phones))
	seen := make(map[string]struct{}, len(request.Emails)+len(request.Phones))

	appendContacts := func(
		contactType entity.ContactType,
		values []foundationrequest.ContactValue,
		normalize func(string) (string, error),
	) error {
		primaryCount := 0
		for _, contactRequest := range values {
			normalizedValue, err := normalize(contactRequest.Value)
			if err != nil {
				return err
			}

			key := contactKey(contactType, normalizedValue)
			if _, exists := seen[key]; exists {
				return coreerrors.NewError(personerror.InvalidPersonContact)
			}
			seen[key] = struct{}{}

			if contactRequest.IsPrimary {
				primaryCount++
				if primaryCount > 1 {
					return coreerrors.NewError(personerror.InvalidPersonContact)
				}
			}

			contacts = append(
				contacts,
				entity.NewContact(
					contactType,
					normalizedValue,
					normalizedValue,
					contactRequest.IsPrimary,
					existingVerification[key],
				),
			)
		}
		return nil
	}

	if err := appendContacts(entity.ContactTypeEmail, request.Emails, normalizeEmail); err != nil {
		return nil, err
	}
	if err := appendContacts(entity.ContactTypePhone, request.Phones, normalizePhone); err != nil {
		return nil, err
	}

	return contacts, nil
}

func normalizeEmail(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	address, err := mail.ParseAddress(normalized)
	if err != nil || strings.ToLower(address.Address) != normalized {
		return "", coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return normalized, nil
}

func normalizePhone(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	builder := strings.Builder{}
	for index, character := range trimmed {
		switch {
		case character == '+' && index == 0:
			builder.WriteRune(character)
		case character >= '0' && character <= '9':
			builder.WriteRune(character)
		case character == ' ' || character == '-' || character == '(' || character == ')':
			continue
		default:
			return "", coreerrors.NewError(personerror.InvalidPersonContact)
		}
	}

	normalized := builder.String()
	if !phonePattern.MatchString(normalized) {
		return "", coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return normalized, nil
}

func contactsEqual(left []entity.Contact, right []entity.Contact) bool {
	if len(left) != len(right) {
		return false
	}

	rightByKey := make(map[string]entity.Contact, len(right))
	for _, contact := range right {
		rightByKey[contactKey(contact.GetType(), contact.GetNormalizedValue())] = contact
	}

	for _, contact := range left {
		other, exists := rightByKey[contactKey(contact.GetType(), contact.GetNormalizedValue())]
		if !exists ||
			other.GetValue() != contact.GetValue() ||
			other.GetIsPrimary() != contact.GetIsPrimary() ||
			other.GetIsVerified() != contact.GetIsVerified() {
			return false
		}
	}
	return true
}

func contactKey(contactType entity.ContactType, normalizedValue string) string {
	return string(contactType) + "\x00" + normalizedValue
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
