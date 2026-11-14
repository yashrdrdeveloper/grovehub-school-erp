package helpers

import (
    "strings"

    corenormalize "grovehub/libs/core/normalize"
    coresearch "grovehub/libs/core/search"
    "grovehub/services/foundation/entity"
)

type valueNormalizer func(coresearch.Operator, string) (string, error)

type searchFieldSpec struct {
	operators  map[coresearch.Operator]struct{}
	normalizer valueNormalizer
}


var personSearchFields = map[string]searchFieldSpec{
	"gid": equalitySearchField(normalizeGIDSearchValue),
	"firstName":  nameSearchField(),
	"middleName": nameSearchField(),
	"lastName":   nameSearchField(),
	"status": equalitySearchField(normalizeStatusSearchValue),
	"createdAt":              timeSearchField(),
	"updatedAt":              timeSearchField(),
	"contacts.email":         contactSearchField(normalizeEmailSearchValue),
	"contacts.primaryEmail":  contactSearchField(normalizeEmailSearchValue),
	"contacts.verifiedEmail": contactSearchField(normalizeEmailSearchValue),
	"contacts.phone":         contactSearchField(normalizePhoneSearchFieldValue),
	"contacts.primaryPhone":  contactSearchField(normalizePhoneSearchFieldValue),
	"contacts.verifiedPhone": contactSearchField(normalizePhoneSearchFieldValue),
}

func nameSearchField() searchFieldSpec {
	return searchFieldSpec{
		operators: textSearchOperators(),
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
		operators: textSearchOperators(),
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
	normalized, err := corenormalize.RFC3339UTC(value)
	if err != nil {
		return "", invalidPersonSearch()
	}
	return normalized, nil
}

func normalizeEmailSearchValue(operator coresearch.Operator, value string) (string, error) {
	value = strings.TrimSpace(value)
	if isPartialSearchOperator(operator) {
		if value == "" {
			return "", invalidPersonSearch()
		}
		return strings.ToLower(value), nil
	}
	return normalizeEmail(value)
}

func normalizePhoneSearchFieldValue(operator coresearch.Operator, value string) (string, error) {
	if isPartialSearchOperator(operator) {
		normalized, err := corenormalize.PhoneCharacters(value)
		if err != nil || normalized == "" {
			return "", invalidPersonSearch()
		}
		return normalized, nil
	}
	return normalizePhone(value)
}

func equalitySearchField(normalizer valueNormalizer) searchFieldSpec {
    return searchFieldSpec{
        operators: operatorSet(coresearch.OperatorEquals, coresearch.OperatorNotEquals, coresearch.OperatorIn, coresearch.OperatorNotIn),
        normalizer: normalizer,
    }
}

func textSearchOperators() map[coresearch.Operator]struct{} {
    return operatorSet(coresearch.OperatorEquals, coresearch.OperatorNotEquals, coresearch.OperatorIn, coresearch.OperatorNotIn,
        coresearch.OperatorExists, coresearch.OperatorNotExists, coresearch.OperatorContains, coresearch.OperatorStartsWith)
}

func isPartialSearchOperator(operator coresearch.Operator) bool {
    return operator == coresearch.OperatorContains || operator == coresearch.OperatorStartsWith
}
