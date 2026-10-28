package validators

import (
	"errors"

	"grovehub/libs/core/search/contracts"
)

var ErrInvalidQuery = errors.New("invalid search query")

func ValidateQuery[TId any](
	query contracts.Query[TId],
	config contracts.Config,
) (contracts.SortField, error) {
	if query.Limit <= 0 || config.CursorExpression == "" {
		return contracts.SortField{}, ErrInvalidQuery
	}
	if len(query.SearchTerms) > 0 && config.SearchExpression == "" {
		return contracts.SortField{}, ErrInvalidQuery
	}

	sortField, ok := config.SortFields[query.Sort.Field]
	if !ok {
		return contracts.SortField{}, ErrInvalidQuery
	}
	if query.Sort.Direction != contracts.DirectionAscending &&
		query.Sort.Direction != contracts.DirectionDescending {
		return contracts.SortField{}, ErrInvalidQuery
	}
	return sortField, nil
}
