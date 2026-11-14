package helpers

import coresearch "grovehub/libs/core/search"

const (
    maxSearchTextLength = 256
    maxSearchTerms = 10
    maxSearchFilters = 20
    maxFilterValues = 50
)

// NormalizePersonSearchRequest assembles generic search rules with Person policy.
func NormalizePersonSearchRequest(request coresearch.Request) (coresearch.Query[string], coresearch.CursorPayload, error) {
    terms, err := coresearch.NormalizeSearchTerms(request.SearchText, maxSearchTextLength, maxSearchTerms, invalidPersonSearch())
    if err != nil { return coresearch.Query[string]{}, coresearch.CursorPayload{}, err }

    rows, err := coresearch.ResolveRows(request.PageInfo.Rows, invalidPersonSearch())
    if err != nil { return coresearch.Query[string]{}, coresearch.CursorPayload{}, err }

    fallback := coresearch.QuerySort{Field: "createdAt", Direction: coresearch.DirectionDescending}
    sort, err := coresearch.NormalizeSort(request.Sort, fallback, personSortFields, invalidPersonSearch())
    if err != nil { return coresearch.Query[string]{}, coresearch.CursorPayload{}, err }

    filters, err := coresearch.NormalizeFilters(request.Filters, personSearchFields, maxSearchFilters, maxFilterValues, invalidPersonSearch())
    if err != nil { return coresearch.Query[string]{}, coresearch.CursorPayload{}, err }

    cursor, payload, err := coresearch.NormalizeCursor(request.PageInfo.Cursor, sort, personSortFields[sort.Field], coresearch.ValidateUUIDv7Cursor, invalidPersonSearchCursor())
    if err != nil { return coresearch.Query[string]{}, coresearch.CursorPayload{}, err }

    return coresearch.Query[string]{
        SearchTerms: terms,
        Filters: filters,
        Sort: sort,
        Cursor: cursor,
        Limit: rows,
    }, payload, nil
}
