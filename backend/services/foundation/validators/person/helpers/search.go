package helpers

import (
    "strings"

    coresearch "grovehub/libs/core/search"
)

const (
    maxSearchTextLength = 256
    maxSearchTerms = 10
    maxSearchFilters = 20
    maxFilterValues = 50
)

var valueCountValidators = map[coresearch.Operator]func(int) bool{
	coresearch.OperatorExists:    func(count int) bool { return count == 0 },
	coresearch.OperatorNotExists: func(count int) bool { return count == 0 },
	coresearch.OperatorIn:        func(count int) bool { return count > 0 },
	coresearch.OperatorNotIn:     func(count int) bool { return count > 0 },
}

func NormalizePersonSearchRequest(
	request coresearch.Request,
) (coresearch.Query[string], coresearch.CursorPayload, error) {
	searchText := strings.TrimSpace(request.SearchText)
	if len(searchText) > maxSearchTextLength {
		return coresearch.Query[string]{}, coresearch.CursorPayload{}, invalidPersonSearch()
	}
	searchTerms := strings.Fields(strings.ToLower(searchText))
	if len(searchTerms) > maxSearchTerms || len(request.Filters) > maxSearchFilters {
		return coresearch.Query[string]{}, coresearch.CursorPayload{}, invalidPersonSearch()
	}

	rows := request.PageInfo.Rows
	if rows == 0 {
		rows = coresearch.DefaultRows
	}
	if rows < 1 || rows > coresearch.MaxRows {
		return coresearch.Query[string]{}, coresearch.CursorPayload{}, invalidPersonSearch()
	}

	sort, err := normalizeSort(request.Sort)
	if err != nil {
		return coresearch.Query[string]{}, coresearch.CursorPayload{}, err
	}
	filters, err := normalizeFilters(request.Filters)
	if err != nil {
		return coresearch.Query[string]{}, coresearch.CursorPayload{}, err
	}
	cursor, payload, err := normalizeCursor(request.PageInfo.Cursor, sort)
	if err != nil {
		return coresearch.Query[string]{}, coresearch.CursorPayload{}, err
	}

	return coresearch.Query[string]{
		SearchTerms: searchTerms,
		Filters:     filters,
		Sort:        sort,
		Cursor:      cursor,
		Limit:       rows,
	}, payload, nil
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
) (*coresearch.Cursor[string], coresearch.CursorPayload, error) {
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
	return &coresearch.Cursor[string]{
		Value: decoded.Value,
		CursorId: decoded.CursorId,
	}, decoded, nil
}

func validSearchValueCount(operator coresearch.Operator, count int) bool {
	if validate, ok := valueCountValidators[operator]; ok {
		return validate(count)
	}
	return count == 1
}
