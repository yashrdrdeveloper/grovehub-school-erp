package helpers

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	coreutils "grovehub/libs/core/utils"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/internal/app/constants/persona"
	roledto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/entity"
)

var ErrInvalidSearch = errors.New("invalid role search")

func normalizeID(_ coresearch.Operator, value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return "", ErrInvalidSearch
	}
	return id.String(), nil
}

func normalizeName(_ coresearch.Operator, value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || !utf8.ValidString(value) || strings.ContainsRune(value, 0) {
		return "", ErrInvalidSearch
	}
	return value, nil
}

func normalizeStatus(_ coresearch.Operator, value string) (string, error) {
	status := entity.RoleStatus(strings.ToUpper(strings.TrimSpace(value)))
	if status != entity.RoleActive && status != entity.RoleInactive {
		return "", ErrInvalidSearch
	}
	return string(status), nil
}

func normalizePersona(_ coresearch.Operator, value string) (string, error) {
	code := persona.Code(strings.TrimSpace(value))
	if !persona.Known(code) {
		return "", ErrInvalidSearch
	}
	return string(code), nil
}

func normalizeTime(_ coresearch.Operator, value string) (string, error) {
	value, err := coreutils.NormalizeRFC3339UTC(value)
	if err != nil {
		return "", ErrInvalidSearch
	}
	return value, nil
}

func validateNameCursor(value string, invalid error) error {
	normalized, err := normalizeName(coresearch.OperatorEquals, value)
	if err != nil || normalized != value {
		return invalid
	}
	return nil
}

func validateStatusCursor(value string, invalid error) error {
	normalized, err := normalizeStatus(coresearch.OperatorEquals, value)
	if err != nil || normalized != value {
		return invalid
	}
	return nil
}

func SortFields() map[string]coresearch.CursorValueValidator {
	return map[string]coresearch.CursorValueValidator{
		roledto.FieldRoleID:    ValidateRoleCursor,
		roledto.FieldName:      validateNameCursor,
		roledto.FieldStatus:    validateStatusCursor,
		roledto.FieldCreatedAt: coresearch.ValidateTimeCursor,
		roledto.FieldUpdatedAt: coresearch.ValidateTimeCursor,
	}
}

func FilterFields() map[string]coresearch.FieldSpec {
	equality := coresearch.OperatorSet(coresearch.OperatorEquals, coresearch.OperatorNotEquals, coresearch.OperatorIn, coresearch.OperatorNotIn)
	return map[string]coresearch.FieldSpec{
		roledto.FieldRoleID:    {Operators: equality, Normalizer: normalizeID},
		roledto.FieldName:      {Operators: coresearch.TextOperators(), Normalizer: normalizeName},
		roledto.FieldStatus:    {Operators: equality, Normalizer: normalizeStatus},
		roledto.FieldPersonas:  {Operators: coresearch.EqualityOperators(), Normalizer: normalizePersona},
		roledto.FieldCreatedAt: {Operators: coresearch.OperatorSet(coresearch.OperatorEquals, coresearch.OperatorNotEquals, coresearch.OperatorGreaterThan, coresearch.OperatorGTE, coresearch.OperatorLessThan, coresearch.OperatorLTE), Normalizer: normalizeTime},
		roledto.FieldUpdatedAt: {Operators: coresearch.OperatorSet(coresearch.OperatorEquals, coresearch.OperatorNotEquals, coresearch.OperatorGreaterThan, coresearch.OperatorGTE, coresearch.OperatorLessThan, coresearch.OperatorLTE), Normalizer: normalizeTime},
	}
}

func ValidateRoleCursor(value string, invalid error) error {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil || id.String() != value || (id.Version() != 5 && id.Version() != 7) {
		return invalid
	}
	return nil
}
