package validation

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
)

func ValidateQuery[TId any](
	query contracts.Query[TId],
	config contracts.Config,
) (contracts.SortField, error) {
	if query.Limit <= 0 || query.Offset < 0 || (query.Cursor != nil && query.Offset != 0) || config.CursorExpression == "" {
		return contracts.SortField{}, contracts.ErrInvalidQuery
	}
	if len(query.SearchTerms) > 0 && config.SearchExpression == "" {
		return contracts.SortField{}, contracts.ErrInvalidQuery
	}

	sortField, ok := config.SortFields[query.Sort.Field]
	if !ok {
		return contracts.SortField{}, contracts.ErrInvalidQuery
	}
	if query.Sort.Direction != contracts.DirectionAscending &&
		query.Sort.Direction != contracts.DirectionDescending {
		return contracts.SortField{}, contracts.ErrInvalidQuery
	}
	return sortField, nil
}
