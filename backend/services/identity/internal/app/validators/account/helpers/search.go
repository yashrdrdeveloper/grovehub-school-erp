package helpers

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	corenormalize "grovehub/libs/core/normalize"
	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
)

var errInvalidSearch = errors.New("invalid account search")

func normalizeUUID(_ coresearch.Operator, value string) (string, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return "", errInvalidSearch
	}
	return parsed.String(), nil
}

func normalizeText(_ coresearch.Operator, value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "", errInvalidSearch
	}
	return value, nil
}

func normalizeStatus(_ coresearch.Operator, value string) (string, error) {
	status := entity.AccountStatus(strings.ToUpper(strings.TrimSpace(value)))
	switch status {
	case entity.Active, entity.PendingSetup, entity.Suspended, entity.Deactivated:
		return string(status), nil
	default:
		return "", errInvalidSearch
	}
}

func normalizeTime(_ coresearch.Operator, value string) (string, error) {
	normalized, err := corenormalize.RFC3339UTC(value)
	if err != nil {
		return "", errInvalidSearch
	}
	return normalized, nil
}

func sortFields() map[string]coresearch.CursorValueValidator {
	return map[string]coresearch.CursorValueValidator{
		accountdto.FieldCreatedAt: coresearch.ValidateTimeCursor,
		accountdto.FieldUpdatedAt: coresearch.ValidateTimeCursor,
		accountdto.FieldAccountID: coresearch.ValidateUUIDv7Cursor,
		accountdto.FieldStatus: coresearch.ValidateNonEmptyCursor,
		accountdto.FieldUsername: coresearch.ValidateNonEmptyCursor,
	}
}

func filterFields() map[string]coresearch.FieldSpec {
	return map[string]coresearch.FieldSpec{
		accountdto.FieldAccountID: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeUUID},
		accountdto.FieldPersonGID: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeUUID},
		accountdto.FieldUsername: {Operators: coresearch.TextOperators(), Normalizer: normalizeText},
		accountdto.FieldEmail: {Operators: coresearch.TextOperators(), Normalizer: normalizeText},
		accountdto.FieldStatus: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeStatus},
		accountdto.FieldCreatedAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
		accountdto.FieldUpdatedAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
	}
}

// NormalizeSearch validates all public search keys, values, sorts and cursors.
func NormalizeSearch(request coresearch.Request) (coresearch.Query[string], error) {
	terms, err := coresearch.NormalizeSearchTerms(request.SearchText, 256, 10, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	rows, err := coresearch.ResolveRows(request.PageInfo.Rows, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	sorts := sortFields()
	sort, err := coresearch.NormalizeSort(request.Sort, coresearch.QuerySort{
		Field: accountdto.FieldCreatedAt,
		Direction: coresearch.DirectionDescending,
	}, sorts, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	filters, err := coresearch.NormalizeFilters(request.Filters, filterFields(), 20, 50, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	cursor, _, err := coresearch.NormalizeCursor(request.PageInfo.Cursor, sort, sorts[sort.Field], coresearch.ValidateUUIDv7Cursor, errInvalidSearch)
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
