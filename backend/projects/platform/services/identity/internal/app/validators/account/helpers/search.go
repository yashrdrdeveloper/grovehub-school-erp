package helpers

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	coreutils "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	roles "grovehub/services/identity/internal/app/constants/role"
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

func normalizeRole(_ coresearch.Operator, value string) (string, error) {
	id := strings.TrimSpace(value)
	if !roles.IsAssignable(id) {
		return "", errInvalidSearch
	}
	return id, nil
}

func normalizeTime(_ coresearch.Operator, value string) (string, error) {
	normalized, err := coreutils.NormalizeRFC3339UTC(value)
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
		accountdto.FieldUsername: validateUsernameCursor,
	}
}

func validateUsernameCursor(value string, invalid error) error {
	if value != strings.ToLower(value) { return invalid }
	return nil
}

func filterFields(normalized bool) map[string]coresearch.FieldSpec {
	fields := map[string]coresearch.FieldSpec{
		accountdto.FieldAccountID: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeUUID},
		accountdto.FieldPersonGID: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeUUID},
		accountdto.FieldUsername: {Operators: coresearch.TextOperators(), Normalizer: normalizeText},
		accountdto.FieldEmail: {Operators: coresearch.TextOperators(), Normalizer: normalizeText},
		accountdto.FieldStatus: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeStatus},
		accountdto.FieldPrimaryRoleID: {Operators: coresearch.EqualityOperators(), Normalizer: normalizeRole},
		accountdto.FieldCreatedAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
		accountdto.FieldUpdatedAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
	}
	if normalized {
		fields[accountdto.FieldPrimaryRoleID] = coresearch.FieldSpec{Operators: coresearch.EqualityOperators(), Normalizer: normalizeNormalizedRole}
	}
	return fields
}

func normalizeNormalizedRole(_ coresearch.Operator, value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return "", errInvalidSearch
	}
	return id.String(), nil
}

// NormalizeSearch validates all public search keys, values, sorts and cursors.
func NormalizeSearch(request coresearch.Request) (coresearch.Query[string], error) {
	return normalizeSearch(request, false)
}

func NormalizeNormalizedSearch(request coresearch.Request) (coresearch.Query[string], error) {
	return normalizeSearch(request, true)
}

func normalizeSearch(request coresearch.Request, normalized bool) (coresearch.Query[string], error) {
	terms, err := coresearch.NormalizeSearchTerms(request.SearchText, 256, 10, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	page, rows, offset, err := coresearch.ResolvePage(request.PageInfo, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	sorts := sortFields()
	ordered, err := coresearch.NormalizeSorts(request.Sort, coresearch.QuerySort{
		Field: accountdto.FieldCreatedAt,
		Direction: coresearch.DirectionDescending,
	}, sorts, 5, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	filters, err := coresearch.NormalizeFilters(request.Filters, filterFields(normalized), 20, 50, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	cursor, err := coresearch.NormalizeMultiCursor(request.PageInfo.Cursor, ordered, sorts, coresearch.ValidateUUIDv7Cursor, errInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	if err := coresearch.ValidatePageCursor(request.PageInfo, page, rows, errInvalidSearch); err != nil { return coresearch.Query[string]{}, err }
	if request.PageInfo.PageNumber == 0 { page = 0 }
	return coresearch.Query[string]{
		SearchTerms: terms,
		Filters: filters,
		Sort: ordered[0],
		SecondarySorts: ordered[1:],
		Cursor: cursor,
		Limit: rows,
		PageNumber: page,
		Offset: offset,
	}, nil
}
