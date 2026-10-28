package helpers

import (
	"strings"
	"time"
	"unicode/utf8"

	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/normalize"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	"grovehub/services/foundation/entity"
)

const (
	maxPersonNameLength = 100
	maxContactsPerType  = 10
	maxSearchTextLength = 256
	maxSearchTerms      = 10
	maxSearchFilters    = 20
	maxFilterValues     = 50
)

type valueNormalizer func(coresearch.Operator, string) (string, error)
type cursorValueValidator func(string) error

type searchFieldSpec struct {
	operators  map[coresearch.Operator]struct{}
	normalizer valueNormalizer
}

type sortFieldSpec struct {
	validate cursorValueValidator
}

var personSearchFields = map[string]searchFieldSpec{
	"gid": {
		operators: operatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
		),
		normalizer: normalizeGIDSearchValue,
	},
	"firstName":  nameSearchField(),
	"middleName": nameSearchField(),
	"lastName":   nameSearchField(),
	"status": {
		operators: operatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
		),
		normalizer: normalizeStatusSearchValue,
	},
	"createdAt":              timeSearchField(),
	"updatedAt":              timeSearchField(),
	"contacts.email":         contactSearchField(normalizeEmailSearchValue),
	"contacts.primaryEmail":  contactSearchField(normalizeEmailSearchValue),
	"contacts.verifiedEmail": contactSearchField(normalizeEmailSearchValue),
	"contacts.phone":         contactSearchField(normalizePhoneSearchFieldValue),
	"contacts.primaryPhone":  contactSearchField(normalizePhoneSearchFieldValue),
	"contacts.verifiedPhone": contactSearchField(normalizePhoneSearchFieldValue),
}

var personSortFields = map[string]sortFieldSpec{
	"gid":       {validate: validateUUIDv7String},
	"firstName": {validate: validateNonEmptyCursorValue},
	"lastName":  {validate: validateNonEmptyCursorValue},
	"createdAt": {validate: validateTimeCursorValue},
	"updatedAt": {validate: validateTimeCursorValue},
}

var valueCountValidators = map[coresearch.Operator]func(int) bool{
	coresearch.OperatorExists:    func(count int) bool { return count == 0 },
	coresearch.OperatorNotExists: func(count int) bool { return count == 0 },
	coresearch.OperatorIn:        func(count int) bool { return count > 0 },
	coresearch.OperatorNotIn:     func(count int) bool { return count > 0 },
}

func NormalizePersonNames(firstName, middleName, lastName string) (string, string, string, error) {
	firstName, err := NormalizeRequiredName(firstName)
	if err != nil {
		return "", "", "", err
	}
	middleName, err = NormalizeOptionalName(middleName)
	if err != nil {
		return "", "", "", err
	}
	lastName, err = NormalizeRequiredName(lastName)
	if err != nil {
		return "", "", "", err
	}
	return firstName, middleName, lastName, nil
}

func NormalizeRequiredName(value string) (string, error) {
	normalized, err := NormalizeOptionalName(value)
	if err != nil || normalized == "" {
		return "", coreerrors.NewError(personerror.InvalidPersonRequest)
	}
	return normalized, nil
}

func NormalizeOptionalName(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	if utf8.RuneCountInString(normalized) > maxPersonNameLength {
		return "", coreerrors.NewError(personerror.InvalidPersonRequest)
	}
	return normalized, nil
}

func NormalizeContacts(request foundationrequest.Contacts) (foundationrequest.Contacts, error) {
	if len(request.Emails) > maxContactsPerType || len(request.Phones) > maxContactsPerType {
		return foundationrequest.Contacts{}, coreerrors.NewError(personerror.InvalidPersonContact)
	}

	seen := make(map[string]struct{}, len(request.Emails)+len(request.Phones))
	emails, err := normalizeContactGroup(entity.ContactTypeEmail, request.Emails, normalizeEmail, seen)
	if err != nil {
		return foundationrequest.Contacts{}, err
	}
	phones, err := normalizeContactGroup(entity.ContactTypePhone, request.Phones, normalizePhone, seen)
	if err != nil {
		return foundationrequest.Contacts{}, err
	}

	return foundationrequest.Contacts{Emails: emails, Phones: phones}, nil
}

func ValidateGID(gid entity.GID) error {
	if _, err := corenormalize.UUIDv7(string(gid)); err != nil {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}
	return nil
}

func NormalizePersonSearchRequest(
	request coresearch.Request,
) (coresearch.Query[entity.GID], coresearch.CursorPayload, error) {
	searchText := strings.TrimSpace(request.SearchText)
	if len(searchText) > maxSearchTextLength {
		return coresearch.Query[entity.GID]{}, coresearch.CursorPayload{}, invalidPersonSearch()
	}
	searchTerms := strings.Fields(strings.ToLower(searchText))
	if len(searchTerms) > maxSearchTerms || len(request.Filters) > maxSearchFilters {
		return coresearch.Query[entity.GID]{}, coresearch.CursorPayload{}, invalidPersonSearch()
	}

	rows := request.PageInfo.Rows
	if rows == 0 {
		rows = coresearch.DefaultRows
	}
	if rows < 1 || rows > coresearch.MaxRows {
		return coresearch.Query[entity.GID]{}, coresearch.CursorPayload{}, invalidPersonSearch()
	}

	sort, err := normalizeSort(request.Sort)
	if err != nil {
		return coresearch.Query[entity.GID]{}, coresearch.CursorPayload{}, err
	}
	filters, err := normalizeFilters(request.Filters)
	if err != nil {
		return coresearch.Query[entity.GID]{}, coresearch.CursorPayload{}, err
	}
	cursor, payload, err := normalizeCursor(request.PageInfo.Cursor, sort)
	if err != nil {
		return coresearch.Query[entity.GID]{}, coresearch.CursorPayload{}, err
	}

	return coresearch.Query[entity.GID]{
		SearchTerms: searchTerms,
		Filters:     filters,
		Sort:        sort,
		Cursor:      cursor,
		Limit:       rows,
	}, payload, nil
}

func normalizeContactGroup(
	contactType entity.ContactType,
	values []foundationrequest.ContactValue,
	normalize func(string) (string, error),
	seen map[string]struct{},
) ([]foundationrequest.ContactValue, error) {
	normalized := make([]foundationrequest.ContactValue, 0, len(values))
	primaryCount := 0

	for _, contact := range values {
		value, err := normalize(contact.Value)
		if err != nil {
			return nil, err
		}
		key := contactKey(contactType, value)
		if _, exists := seen[key]; exists {
			return nil, coreerrors.NewError(personerror.InvalidPersonContact)
		}
		seen[key] = struct{}{}

		if contact.IsPrimary {
			primaryCount++
			if primaryCount > 1 {
				return nil, coreerrors.NewError(personerror.InvalidPersonContact)
			}
		}

		normalized = append(normalized, foundationrequest.ContactValue{
			Value:     value,
			IsPrimary: contact.IsPrimary,
		})
	}
	return normalized, nil
}

func normalizeSort(sortRequest []coresearch.Sort) (coresearch.QuerySort, error) {
	sort := coresearch.QuerySort{Field: "createdAt", Direction: coresearch.DirectionDescending}
	if len(sortRequest) > 1 {
		return coresearch.QuerySort{}, invalidPersonSearch()
	}
	if len(sortRequest) == 1 {
		sort.Field = strings.TrimSpace(sortRequest[0].Field)
		sort.Direction = sortRequest[0].Direction
		if sort.Direction == "" {
			sort.Direction = coresearch.DirectionDescending
		}
	}
	if _, ok := personSortFields[sort.Field]; !ok ||
		(sort.Direction != coresearch.DirectionAscending && sort.Direction != coresearch.DirectionDescending) {
		return coresearch.QuerySort{}, invalidPersonSearch()
	}
	return sort, nil
}

func normalizeFilters(filters []coresearch.Filter) ([]coresearch.QueryFilter, error) {
	normalized := make([]coresearch.QueryFilter, 0, len(filters))
	for _, filter := range filters {
		field := strings.TrimSpace(filter.Field)
		spec, ok := personSearchFields[field]
		if !ok {
			return nil, invalidPersonSearch()
		}
		if _, ok := spec.operators[filter.Operator]; !ok ||
			!validSearchValueCount(filter.Operator, len(filter.Values)) ||
			len(filter.Values) > maxFilterValues {
			return nil, invalidPersonSearch()
		}
		values := make([]string, 0, len(filter.Values))
		for _, value := range filter.Values {
			normalizedValue, err := spec.normalizer(filter.Operator, value)
			if err != nil {
				return nil, err
			}
			values = append(values, normalizedValue)
		}
		normalized = append(normalized, coresearch.QueryFilter{
			Field: field,
			Operator: filter.Operator,
			Values: values,
		})
	}
	return normalized, nil
}

func normalizeCursor(
	value string,
	sort coresearch.QuerySort,
) (*coresearch.Cursor[entity.GID], coresearch.CursorPayload, error) {
	payload := coresearch.CursorPayload{Field: sort.Field, Direction: sort.Direction}
	if strings.TrimSpace(value) == "" {
		return nil, payload, nil
	}
	decoded, err := coresearch.DecodeCursor(value)
	if err != nil || decoded.Field != sort.Field || decoded.Direction != sort.Direction {
		return nil, coresearch.CursorPayload{}, invalidPersonSearchCursor()
	}
	spec := personSortFields[sort.Field]
	if err := spec.validate(decoded.Value); err != nil || validateUUIDv7String(decoded.CursorId) != nil {
		return nil, coresearch.CursorPayload{}, invalidPersonSearchCursor()
	}
	return &coresearch.Cursor[entity.GID]{
		Value: decoded.Value,
		CursorId: entity.GID(decoded.CursorId),
	}, decoded, nil
}

func validSearchValueCount(operator coresearch.Operator, count int) bool {
	if validate, ok := valueCountValidators[operator]; ok {
		return validate(count)
	}
	return count == 1
}

func nameSearchField() searchFieldSpec {
	return searchFieldSpec{
		operators: operatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
			coresearch.OperatorExists,
			coresearch.OperatorNotExists,
			coresearch.OperatorContains,
			coresearch.OperatorStartsWith,
		),
		normalizer: func(_ coresearch.Operator, value string) (string, error) {
			value = strings.ToLower(strings.TrimSpace(value))
			if value == "" {
				return "", invalidPersonSearch()
			}
			return value, nil
		},
	}
}

func timeSearchField() searchFieldSpec {
	return searchFieldSpec{
		operators: operatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorGreaterThan,
			coresearch.OperatorGTE,
			coresearch.OperatorLessThan,
			coresearch.OperatorLTE,
		),
		normalizer: normalizeTimeSearchValue,
	}
}

func contactSearchField(normalizer valueNormalizer) searchFieldSpec {
	return searchFieldSpec{
		operators: operatorSet(
			coresearch.OperatorEquals,
			coresearch.OperatorNotEquals,
			coresearch.OperatorIn,
			coresearch.OperatorNotIn,
			coresearch.OperatorExists,
			coresearch.OperatorNotExists,
			coresearch.OperatorContains,
			coresearch.OperatorStartsWith,
		),
		normalizer: normalizer,
	}
}

func operatorSet(operators ...coresearch.Operator) map[coresearch.Operator]struct{} {
	set := make(map[coresearch.Operator]struct{}, len(operators))
	for _, operator := range operators {
		set[operator] = struct{}{}
	}
	return set
}

func normalizeGIDSearchValue(_ coresearch.Operator, value string) (string, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", invalidPersonSearch()
	}
	return normalized, nil
}

func normalizeStatusSearchValue(_ coresearch.Operator, value string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	allowed := map[string]struct{}{
		string(entity.PersonStatusActive):      {},
		string(entity.PersonStatusDeactivated): {},
	}
	if _, ok := allowed[status]; !ok {
		return "", invalidPersonSearch()
	}
	return status, nil
}

func normalizeTimeSearchValue(_ coresearch.Operator, value string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return "", invalidPersonSearch()
	}
	return parsed.UTC().Format(time.RFC3339Nano), nil
}

func normalizeEmailSearchValue(operator coresearch.Operator, value string) (string, error) {
	value = strings.TrimSpace(value)
	partialOperators := map[coresearch.Operator]struct{}{
		coresearch.OperatorContains:   {},
		coresearch.OperatorStartsWith: {},
	}
	if _, ok := partialOperators[operator]; ok {
		if value == "" {
			return "", invalidPersonSearch()
		}
		return strings.ToLower(value), nil
	}
	return normalizeEmail(value)
}

func normalizePhoneSearchFieldValue(operator coresearch.Operator, value string) (string, error) {
	partialOperators := map[coresearch.Operator]struct{}{
		coresearch.OperatorContains:   {},
		coresearch.OperatorStartsWith: {},
	}
	if _, ok := partialOperators[operator]; ok {
		normalized, err := corenormalize.PhoneCharacters(value)
		if err != nil || normalized == "" {
			return "", invalidPersonSearch()
		}
		return normalized, nil
	}
	return normalizePhone(value)
}

func normalizeEmail(value string) (string, error) {
	normalized, err := corenormalize.Email(value)
	if err != nil {
		return "", coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return normalized, nil
}

func normalizePhone(value string) (string, error) {
	normalized, err := corenormalize.Phone(value)
	if err != nil {
		return "", coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return normalized, nil
}

func validateUUIDv7String(value string) error {
	if _, err := corenormalize.UUIDv7(value); err != nil {
		return invalidPersonSearchCursor()
	}
	return nil
}

func validateTimeCursorValue(value string) error {
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		return invalidPersonSearchCursor()
	}
	return nil
}

func validateNonEmptyCursorValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return invalidPersonSearchCursor()
	}
	return nil
}

func contactKey(contactType entity.ContactType, normalizedValue string) string {
	return string(contactType) + "\x00" + normalizedValue
}

func invalidPersonSearch() error {
	return coreerrors.NewError(personerror.InvalidPersonSearch)
}

func invalidPersonSearchCursor() error {
	return coreerrors.NewError(personerror.InvalidPersonSearchCursor)
}
