package impl

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	corenormalize "grovehub/libs/core/normalize"
	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
)

// @Validator
type Validator struct{}

var _ accountvalidator.Validator = Validator{}

func New() accountvalidator.Validator {
	return Validator{}
}

func (Validator) NormalizeCreate(request accountdto.Create) (accountdto.Create, error) {
	email, err := corenormalize.Email(request.Email)
	if err != nil {
		return request, errors.New("invalid email")
	}
	request.Email = email

	if strings.TrimSpace(request.PhoneNumber) != "" {
		phone, phoneErr := corenormalize.Phone(request.PhoneNumber)
		if phoneErr != nil {
			return request, errors.New("invalid phone")
		}
		request.PhoneNumber = phone
	}

	request.Person.FirstName = strings.TrimSpace(request.Person.FirstName)
	request.Person.LastName = strings.TrimSpace(request.Person.LastName)
	if request.Person.MiddleName != nil {
		middleName := strings.TrimSpace(*request.Person.MiddleName)
		request.Person.MiddleName = &middleName
	}
	if request.Person.FirstName == "" || request.Person.LastName == "" ||
		len(request.Person.FirstName) > 100 || len(request.Person.LastName) > 100 {
		return request, errors.New("invalid registration")
	}
	return request, nil
}

func (Validator) NormalizeUsername(username string) (string, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	if len(username) < 3 || len(username) > 40 {
		return "", errors.New("invalid username")
	}
	for _, char := range username {
		letter := char >= 'a' && char <= 'z'
		digit := char >= '0' && char <= '9'
		punctuation := char == '.' || char == '_' || char == '-'
		if !letter && !digit && !punctuation {
			return "", errors.New("invalid username")
		}
	}
	return username, nil
}

func (Validator) Search(request coresearch.Request) (coresearch.Query[string], error) {
	invalid := errors.New("invalid account search")

	terms, err := coresearch.NormalizeSearchTerms(request.SearchText, 256, 10, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	rows, err := coresearch.ResolveRows(request.PageInfo.Rows, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	uuidCursor := func(value string) error {
		return coresearch.ValidateUUIDv7Cursor(value, invalid)
	}
	timeCursor := func(value string) error {
		return coresearch.ValidateTimeCursor(value, invalid)
	}
	textCursor := func(value string) error {
		return coresearch.ValidateNonEmptyCursor(value, invalid)
	}

	sortFields := map[string]coresearch.CursorValueValidator{
		"createdAt": timeCursor,
		"updatedAt": timeCursor,
		"accountId": uuidCursor,
		"status": textCursor,
		"username": textCursor,
	}

	sort, err := coresearch.NormalizeSort(request.Sort, coresearch.QuerySort{
		Field: "createdAt",
		Direction: coresearch.DirectionDescending,
	}, sortFields, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	normalizeUUID := func(_ coresearch.Operator, value string) (string, error) {
		if _, err := uuid.Parse(value); err != nil {
			return "", invalid
		}
		return value, nil
	}
	normalizeText := func(_ coresearch.Operator, value string) (string, error) {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			return "", invalid
		}
		return value, nil
	}
	normalizeStatus := func(_ coresearch.Operator, value string) (string, error) {
		value = strings.ToUpper(strings.TrimSpace(value))
		switch value {
		case "ACTIVE", "PENDING_SETUP", "SUSPENDED", "DEACTIVATED":
			return value, nil
		default:
			return "", invalid
		}
	}
	normalizeTime := func(_ coresearch.Operator, value string) (string, error) {
		normalized, err := corenormalize.RFC3339UTC(value)
		if err != nil {
			return "", invalid
		}
		return normalized, nil
	}

	fields := map[string]coresearch.FieldSpec{
		"accountId": {Operators: coresearch.EqualityOperators(), Normalizer: normalizeUUID},
		"personGid": {Operators: coresearch.EqualityOperators(), Normalizer: normalizeUUID},
		"username": {Operators: coresearch.TextOperators(), Normalizer: normalizeText},
		"email": {Operators: coresearch.TextOperators(), Normalizer: normalizeText},
		"status": {Operators: coresearch.EqualityOperators(), Normalizer: normalizeStatus},
		"createdAt": {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
		"updatedAt": {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
	}
	filters, err := coresearch.NormalizeFilters(request.Filters, fields, 20, 50, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	cursor, _, err := coresearch.NormalizeCursor(
		request.PageInfo.Cursor,
		sort,
		sortFields[sort.Field],
		coresearch.ValidateUUIDv7Cursor,
		invalid,
	)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	return coresearch.Query[string]{
		SearchTerms: terms,
		Filters: filters,
		Sort: sort,
		Cursor: cursor,
		Limit: rows,
	}, nil
}
