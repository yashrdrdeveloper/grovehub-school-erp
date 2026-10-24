package helpers

import (
	"strings"
	"time"
	"unicode/utf8"

	searchcontract "grovehub/libs/core/contracts/search"
	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/helpers/normalize"
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

type valueNormalizer func(searchcontract.Operator, string) (string, error)
type cursorValueValidator func(string) error

type searchFieldSpec struct {
	operators  map[searchcontract.Operator]struct{}
	normalizer valueNormalizer
}

type sortFieldSpec struct {
	validate cursorValueValidator
}

var personSearchFields = map[string]searchFieldSpec{
	"gid": {
		operators:  operatorSet(searchcontract.OperatorEquals, searchcontract.OperatorNotEquals, searchcontract.OperatorIn, searchcontract.OperatorNotIn),
		normalizer: normalizeGIDSearchValue,
	},
	"firstName":  nameSearchField(),
	"middleName": nameSearchField(),
	"lastName":   nameSearchField(),
	"status": {
		operators:  operatorSet(searchcontract.OperatorEquals, searchcontract.OperatorNotEquals, searchcontract.OperatorIn, searchcontract.OperatorNotIn),
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

var valueCountValidators = map[searchcontract.Operator]func(int) bool{
	searchcontract.OperatorExists:    func(count int) bool { return count == 0 },
	searchcontract.OperatorNotExists: func(count int) bool { return count == 0 },
	searchcontract.OperatorIn:        func(count int) bool { return count > 0 },
	searchcontract.OperatorNotIn:     func(count int) bool { return count > 0 },
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

func BuildContacts(request foundationrequest.Contacts, existing []entity.Contact) ([]entity.Contact, error) {
	if len(request.Emails) > maxContactsPerType || len(request.Phones) > maxContactsPerType {
		return nil, coreerrors.NewError(personerror.InvalidPersonContact)
	}

	existingVerification := make(map[string]bool, len(existing))
	for _, contact := range existing {
		existingVerification[contactKey(contact.GetType(), contact.GetNormalizedValue())] = contact.GetIsVerified()
	}

	contacts := make([]entity.Contact, 0, len(request.Emails)+len(request.Phones))
	seen := make(map[string]struct{}, cap(contacts))
	groups := []struct {
		contactType entity.ContactType
		values      []foundationrequest.ContactValue
		normalize   func(string) (string, error)
	}{
		{entity.ContactTypeEmail, request.Emails, normalizeEmail},
		{entity.ContactTypePhone, request.Phones, normalizePhone},
	}

	for _, group := range groups {
		built, err := buildContactGroup(group.contactType, group.values, group.normalize, existingVerification, seen)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, built...)
	}
	return contacts, nil
}

func ValidateGID(gid entity.GID) error {
	if _, err := corenormalize.UUIDv7(string(gid)); err != nil {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}
	return nil
}

func NormalizePersonSearchRequest(
	request searchcontract.Request,
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
		rows = searchcontract.DefaultRows
	}
	if rows < 1 || rows > searchcontract.MaxRows {
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

func buildContactGroup(
	contactType entity.ContactType,
	values []foundationrequest.ContactValue,
	normalize func(string) (string, error),
	existingVerification map[string]bool,
	seen map[string]struct{},
) ([]entity.Contact, error) {
	contacts := make([]entity.Contact, 0, len(values))
	primaryCount := 0
	for _, contactRequest := range values {
		normalizedValue, err := normalize(contactRequest.Value)
		if err != nil {
			return nil, err
		}
		key := contactKey(contactType, normalizedValue)
		if _, exists := seen[key]; exists {
			return nil, coreerrors.NewError(personerror.InvalidPersonContact)
		}
		seen[key] = struct{}{}
		if contactRequest.IsPrimary {
			primaryCount++
			if primaryCount > 1 {
				return nil, coreerrors.NewError(personerror.InvalidPersonContact)
			}
		}
		contacts = append(contacts, entity.NewContact(
			contactType,
			normalizedValue,
			normalizedValue,
			contactRequest.IsPrimary,
			existingVerification[key],
		))
	}
	return contacts, nil
}

func normalizeSort(sortRequest []searchcontract.Sort) (coresearch.Sort, error) {
	sort := coresearch.Sort{Field: "createdAt", Direction: searchcontract.DirectionDescending}
	if len(sortRequest) > 1 {
		return coresearch.Sort{}, invalidPersonSearch()
	}
	if len(sortRequest) == 1 {
		sort.Field = strings.TrimSpace(sortRequest[0].Field)
		sort.Direction = sortRequest[0].Direction
		if sort.Direction == "" {
			sort.Direction = searchcontract.DirectionDescending
		}
	}
	if _, ok := personSortFields[sort.Field]; !ok ||
		(sort.Direction != searchcontract.DirectionAscending && sort.Direction != searchcontract.DirectionDescending) {
		return coresearch.Sort{}, invalidPersonSearch()
	}
	return sort, nil
}

func normalizeFilters(filters []searchcontract.Filter) ([]coresearch.Filter, error) {
	normalized := make([]coresearch.Filter, 0, len(filters))
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
		normalized = append(normalized, coresearch.Filter{Field: field, Operator: filter.Operator, Values: values})
	}
	return normalized, nil
}

func normalizeCursor(
	value string,
	sort coresearch.Sort,
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
	return &coresearch.Cursor[entity.GID]{Value: decoded.Value, CursorId: entity.GID(decoded.CursorId)}, decoded, nil
}

func validSearchValueCount(operator searchcontract.Operator, count int) bool {
	if validate, ok := valueCountValidators[operator]; ok {
		return validate(count)
	}
	return count == 1
}

func nameSearchField() searchFieldSpec {
	return searchFieldSpec{
		operators: operatorSet(
			searchcontract.OperatorEquals,
			searchcontract.OperatorNotEquals,
			searchcontract.OperatorIn,
			searchcontract.OperatorNotIn,
			searchcontract.OperatorExists,
			searchcontract.OperatorNotExists,
			searchcontract.OperatorContains,
			searchcontract.OperatorStartsWith,
		),
		normalizer: func(_ searchcontract.Operator, value string) (string, error) {
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
			searchcontract.OperatorEquals,
			searchcontract.OperatorGreaterThan,
			searchcontract.OperatorGTE,
			searchcontract.OperatorLessThan,
			searchcontract.OperatorLTE,
		),
		normalizer: normalizeTimeSearchValue,
	}
}

func contactSearchField(normalizer valueNormalizer) searchFieldSpec {
	return searchFieldSpec{
		operators: operatorSet(
			searchcontract.OperatorEquals,
			searchcontract.OperatorNotEquals,
			searchcontract.OperatorIn,
			searchcontract.OperatorNotIn,
			searchcontract.OperatorExists,
			searchcontract.OperatorNotExists,
			searchcontract.OperatorContains,
			searchcontract.OperatorStartsWith,
		),
		normalizer: normalizer,
	}
}

func operatorSet(operators ...searchcontract.Operator) map[searchcontract.Operator]struct{} {
	set := make(map[searchcontract.Operator]struct{}, len(operators))
	for _, operator := range operators {
		set[operator] = struct{}{}
	}
	return set
}

func normalizeGIDSearchValue(_ searchcontract.Operator, value string) (string, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", invalidPersonSearch()
	}
	return normalized, nil
}

func normalizeStatusSearchValue(_ searchcontract.Operator, value string) (string, error) {
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

func normalizeTimeSearchValue(_ searchcontract.Operator, value string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return "", invalidPersonSearch()
	}
	return parsed.UTC().Format(time.RFC3339Nano), nil
}

func normalizeEmailSearchValue(operator searchcontract.Operator, value string) (string, error) {
	value = strings.TrimSpace(value)
	partialOperators := map[searchcontract.Operator]struct{}{
		searchcontract.OperatorContains:   {},
		searchcontract.OperatorStartsWith: {},
	}
	if _, ok := partialOperators[operator]; ok {
		if value == "" {
			return "", invalidPersonSearch()
		}
		return strings.ToLower(value), nil
	}
	return normalizeEmail(value)
}

func normalizePhoneSearchFieldValue(operator searchcontract.Operator, value string) (string, error) {
	partialOperators := map[searchcontract.Operator]struct{}{
		searchcontract.OperatorContains:   {},
		searchcontract.OperatorStartsWith: {},
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
