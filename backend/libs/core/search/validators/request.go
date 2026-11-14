package validators

import (
    "strings"

    "grovehub/libs/core/search/contracts"
)

func requestError(invalid error) error {
    if invalid != nil { return invalid }
    return ErrInvalidQuery
}

// NormalizeSearchTerms applies resource-owned length and term limits.
func NormalizeSearchTerms(value string, maxLength, maxTerms int, invalid error) ([]string, error) {
    searchText := strings.TrimSpace(value)
    if len(searchText) > maxLength || maxLength < 0 || maxTerms < 0 {
        return nil, requestError(invalid)
    }
    terms := strings.Fields(strings.ToLower(searchText))
    if len(terms) > maxTerms { return nil, requestError(invalid) }
    return terms, nil
}

// ResolveRows uses the common pagination defaults and maximum.
func ResolveRows(rows int, invalid error) (int, error) {
    if rows == 0 { rows = contracts.DefaultRows }
    if rows < 1 || rows > contracts.MaxRows { return 0, requestError(invalid) }
    return rows, nil
}

// NormalizeSort permits at most one sort and checks the owning resource's
// allowlisted sort fields; the default field and direction are domain policy.
func NormalizeSort(sortRequest []contracts.Sort, fallback contracts.QuerySort, allowed map[string]contracts.CursorValueValidator, invalid error) (contracts.QuerySort, error) {
    if len(sortRequest) > 1 { return contracts.QuerySort{}, requestError(invalid) }
    sort := fallback
    if len(sortRequest) == 1 {
        sort.Field = strings.TrimSpace(sortRequest[0].Field)
        sort.Direction = sortRequest[0].Direction
        if sort.Direction == "" { sort.Direction = fallback.Direction }
    }
    if allowed[sort.Field] == nil || (sort.Direction != contracts.DirectionAscending && sort.Direction != contracts.DirectionDescending) {
        return contracts.QuerySort{}, requestError(invalid)
    }
    return sort, nil
}

// NormalizeFilters owns generic field/operator/value-count mechanics while
// each resource supplies approved fields and normalization functions.
func NormalizeFilters(filters []contracts.Filter, fields map[string]contracts.FieldSpec, maxFilters, maxValues int, invalid error) ([]contracts.QueryFilter, error) {
    if maxFilters < 0 || maxValues < 0 || len(filters) > maxFilters { return nil, requestError(invalid) }
    normalized := make([]contracts.QueryFilter, 0, len(filters))
    for _, filter := range filters {
        field := strings.TrimSpace(filter.Field)
        spec, ok := fields[field]
        if !ok || spec.Normalizer == nil { return nil, requestError(invalid) }
        if _, allowed := spec.Operators[filter.Operator]; !allowed || !ValidValueCount(filter.Operator, len(filter.Values)) || len(filter.Values) > maxValues {
            return nil, requestError(invalid)
        }
        values := make([]string, 0, len(filter.Values))
        for _, value := range filter.Values {
            parsed, err := spec.Normalizer(filter.Operator, value)
            if err != nil { return nil, err }
            values = append(values, parsed)
        }
        normalized = append(normalized, contracts.QueryFilter{Field: field, Operator: filter.Operator, Values: values})
    }
    return normalized, nil
}
