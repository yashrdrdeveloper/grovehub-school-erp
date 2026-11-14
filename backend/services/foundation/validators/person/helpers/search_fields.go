package helpers

import (
    "strings"

    corenormalize "grovehub/libs/core/normalize"
    coresearch "grovehub/libs/core/search"
    "grovehub/services/foundation/entity"
)

// Person owns the allowed fields and their normalizers, not generic
// operator mechanics or filter iteration.
var personSearchFields = map[string]coresearch.FieldSpec{
    "gid": equalitySearchField(normalizeGIDSearchValue),
    "firstName": nameSearchField(),
    "middleName": nameSearchField(),
    "lastName": nameSearchField(),
    "status": equalitySearchField(normalizeStatusSearchValue),
    "createdAt": timeSearchField(),
    "updatedAt": timeSearchField(),
    "contacts.email": contactSearchField(normalizeEmailSearchValue),
    "contacts.primaryEmail": contactSearchField(normalizeEmailSearchValue),
    "contacts.verifiedEmail": contactSearchField(normalizeEmailSearchValue),
    "contacts.phone": contactSearchField(normalizePhoneSearchFieldValue),
    "contacts.primaryPhone": contactSearchField(normalizePhoneSearchFieldValue),
    "contacts.verifiedPhone": contactSearchField(normalizePhoneSearchFieldValue),
}

func nameSearchField() coresearch.FieldSpec {
    return coresearch.FieldSpec{
        Operators: coresearch.TextOperators(),
        Normalizer: func(_ coresearch.Operator, value string) (string, error) {
            normalized := strings.ToLower(strings.TrimSpace(value))
            if normalized == "" { return "", invalidPersonSearch() }
            return normalized, nil
        },
    }
}

func timeSearchField() coresearch.FieldSpec {
    return coresearch.FieldSpec{Operators: coresearch.TimeOperators(), Normalizer: normalizeTimeSearchValue}
}

func contactSearchField(normalizer coresearch.ValueNormalizer) coresearch.FieldSpec {
    return coresearch.FieldSpec{Operators: coresearch.TextOperators(), Normalizer: normalizer}
}

func equalitySearchField(normalizer coresearch.ValueNormalizer) coresearch.FieldSpec {
    return coresearch.FieldSpec{Operators: coresearch.EqualityOperators(), Normalizer: normalizer}
}

func normalizeGIDSearchValue(_ coresearch.Operator, value string) (string, error) {
    normalized, err := corenormalize.UUIDv7(value)
    if err != nil { return "", invalidPersonSearch() }
    return normalized, nil
}

func normalizeStatusSearchValue(_ coresearch.Operator, value string) (string, error) {
    status := strings.ToUpper(strings.TrimSpace(value))
    allowed := map[string]struct{}{
        string(entity.PersonStatusActive): {},
        string(entity.PersonStatusDeactivated): {},
    }
    if _, ok := allowed[status]; !ok { return "", invalidPersonSearch() }
    return status, nil
}

func normalizeTimeSearchValue(_ coresearch.Operator, value string) (string, error) {
    normalized, err := corenormalize.RFC3339UTC(value)
    if err != nil { return "", invalidPersonSearch() }
    return normalized, nil
}

func normalizeEmailSearchValue(operator coresearch.Operator, value string) (string, error) {
    value = strings.TrimSpace(value)
    if coresearch.IsPartialOperator(operator) {
        if value == "" { return "", invalidPersonSearch() }
        return strings.ToLower(value), nil
    }
    return normalizeEmail(value)
}

func normalizePhoneSearchFieldValue(operator coresearch.Operator, value string) (string, error) {
    if coresearch.IsPartialOperator(operator) {
        normalized, err := corenormalize.PhoneCharacters(value)
        if err != nil || normalized == "" { return "", invalidPersonSearch() }
        return normalized, nil
    }
    return normalizePhone(value)
}
