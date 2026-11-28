package role

import (
	"strings"

	coresearch "grovehub/libs/core/search"
	roledto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/validators/role/helpers"
)

var ErrInvalidSearch = helpers.ErrInvalidSearch

func NormalizeSearch(request coresearch.Request) (coresearch.Query[string], error) {
	if strings.ContainsRune(request.SearchText, 0) {
		return coresearch.Query[string]{}, ErrInvalidSearch
	}
	terms, err := coresearch.NormalizeSearchTerms(request.SearchText, 256, 10, ErrInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	page, rows, offset, err := coresearch.ResolvePage(request.PageInfo, ErrInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	sorts := helpers.SortFields()
	ordered, err := coresearch.NormalizeSorts(request.Sort, coresearch.QuerySort{
		Field: roledto.FieldCreatedAt, Direction: coresearch.DirectionDescending,
	}, sorts, 5, ErrInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	filters, err := coresearch.NormalizeFilters(request.Filters, helpers.FilterFields(), 20, 50, ErrInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	cursor, err := coresearch.NormalizeMultiCursor(request.PageInfo.Cursor, ordered, sorts, helpers.ValidateRoleCursor, ErrInvalidSearch)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	if err := coresearch.ValidatePageCursor(request.PageInfo, page, rows, ErrInvalidSearch); err != nil {
		return coresearch.Query[string]{}, err
	}
	if request.PageInfo.PageNumber == 0 {
		page = 0
	}
	return coresearch.Query[string]{
		SearchTerms: terms, Filters: filters, Sort: ordered[0], SecondarySorts: ordered[1:],
		Cursor: cursor, Limit: rows, PageNumber: page, Offset: offset,
	}, nil
}
