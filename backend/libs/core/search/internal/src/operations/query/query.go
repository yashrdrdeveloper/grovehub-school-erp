package query

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/page"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/validation"
)

func Apply[TId any](db *gorm.DB, query contracts.Query[TId], config contracts.Config) (*gorm.DB, error) {
	if db == nil {
		return nil, contracts.ErrInvalidQuery
	}

	sortField, err := validation.ValidateQuery(query, config)
	if err != nil {
		return nil, contracts.ErrInvalidQuery
	}

	for _, term := range query.SearchTerms {
		value := "%" + helpers.EscapeLike(term) + "%"
		db = db.Where(config.SearchExpression+" ILIKE ?", value)
	}

	for _, filter := range query.Filters {
		builder := &helpers.QueryBuilder{}
		condition, err := helpers.BuildFilter(builder, filter, config)
		if err != nil {
			return nil, contracts.ErrInvalidQuery
		}
		db = db.Where(condition, builder.Args()...)
	}

	if len(query.SecondarySorts) > 0 {
		return applyMultiSort(db, query, config)
	}
	if query.Cursor != nil {
		builder := &helpers.QueryBuilder{}
		condition := helpers.BuildCursorCondition(
			builder,
			sortField,
			query.Sort.Direction,
			query.Cursor.Value,
			fmt.Sprint(query.Cursor.CursorId),
			config.CursorExpression,
		)
		db = db.Where(condition, builder.Args()...)
	}

	direction := string(query.Sort.Direction)
	db = db.Order(sortField.Expression + " " + direction)
	if sortField.Expression != config.CursorExpression {
		db = db.Order(config.CursorExpression + " " + direction)
	}

	return page.ApplyPage(db, query), nil
}

// Execute shares the safe, parameterized search and keyset-pagination mechanics.
// The repository supplies its own explicit model/table and allowlisted contracts.Config.
func Execute[TId any, T any](db *gorm.DB, query contracts.Query[TId], config contracts.Config) (contracts.QueryResult[T], error) {
	var total *int64
	if query.PageNumber > 0 {
		countQuery := query
		countQuery.Cursor = nil
		countQuery.Offset = 0
		countDB, err := Apply(db.Session(&gorm.Session{}), countQuery, config)
		if err != nil { return contracts.QueryResult[T]{}, err }
		var count int64
		if err := countDB.Model(new(T)).Select("COUNT(*)").Limit(-1).Offset(-1).Count(&count).Error; err != nil { return contracts.QueryResult[T]{}, err }
		total = &count
	}
	searchDB, err := Apply(db.Session(&gorm.Session{}), query, config)
	if err != nil {
		return contracts.QueryResult[T]{}, err
	}
	items := make([]T, 0, query.Limit+1)
	if err := searchDB.Find(&items).Error; err != nil {
		return contracts.QueryResult[T]{}, err
	}
	hasNext := len(items) > query.Limit
	if hasNext {
		items = items[:query.Limit]
	}
	return contracts.QueryResult[T]{Items: items, HasNext: hasNext, TotalRecords: total}, nil
}

