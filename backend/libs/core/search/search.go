package search

import (
	"fmt"

	"gorm.io/gorm"

	"grovehub/libs/core/search/contracts"
	"grovehub/libs/core/search/helpers"
	"grovehub/libs/core/search/mappers"
	"grovehub/libs/core/search/validators"
)

const (
	OperatorEquals      = contracts.OperatorEquals
	OperatorNotEquals   = contracts.OperatorNotEquals
	OperatorIn          = contracts.OperatorIn
	OperatorNotIn       = contracts.OperatorNotIn
	OperatorExists      = contracts.OperatorExists
	OperatorNotExists   = contracts.OperatorNotExists
	OperatorGreaterThan = contracts.OperatorGreaterThan
	OperatorGTE         = contracts.OperatorGTE
	OperatorLessThan    = contracts.OperatorLessThan
	OperatorLTE         = contracts.OperatorLTE
	OperatorContains    = contracts.OperatorContains
	OperatorStartsWith  = contracts.OperatorStartsWith

	DirectionAscending  = contracts.DirectionAscending
	DirectionDescending = contracts.DirectionDescending

	DefaultRows = contracts.DefaultRows
	MaxRows     = contracts.MaxRows
)

var (
	ErrInvalidQuery        = validators.ErrInvalidQuery
	ErrUnsupportedOperator = helpers.ErrUnsupportedOperator
)

type Operator = contracts.Operator
type Direction = contracts.Direction
type Filter = contracts.Filter
type Sort = contracts.Sort
type PageRequest = contracts.PageRequest
type PageResponse = contracts.PageResponse
type Request = contracts.Request
type Result[T any] = contracts.Result[T]
type QueryFilter = contracts.QueryFilter
type QuerySort = contracts.QuerySort
type Cursor[TId any] = contracts.Cursor[TId]
type Query[TId any] = contracts.Query[TId]
type QueryResult[T any] = contracts.QueryResult[T]
type CursorPayload = contracts.CursorPayload
type ScalarField = contracts.ScalarField
type FixedCondition = contracts.FixedCondition
type RelationField = contracts.RelationField
type SortField = contracts.SortField
type Config = contracts.Config
type QueryBuilder = helpers.QueryBuilder

func Apply[TId any](db *gorm.DB, query Query[TId], config Config) (*gorm.DB, error) {
	if db == nil {
		return nil, ErrInvalidQuery
	}

	sortField, err := validators.ValidateQuery(query, config)
	if err != nil {
		return nil, ErrInvalidQuery
	}

	for _, term := range query.SearchTerms {
		value := "%" + helpers.EscapeLike(term) + "%"
		db = db.Where(config.SearchExpression+" ILIKE ?", value)
	}

	for _, filter := range query.Filters {
		builder := &helpers.QueryBuilder{}
		condition, err := helpers.BuildFilter(builder, filter, config)
		if err != nil {
			return nil, ErrInvalidQuery
		}
		db = db.Where(condition, builder.Args()...)
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

	return db.Limit(query.Limit + 1), nil
}

func EncodeCursor(cursor CursorPayload) (string, error) {
	return mappers.EncodeCursor(cursor)
}

func DecodeCursor(value string) (CursorPayload, error) {
	return mappers.DecodeCursor(value)
}

func BuildScalarCondition(
	builder *QueryBuilder,
	field ScalarField,
	operator Operator,
	values []string,
) (string, error) {
	return helpers.BuildScalarCondition(builder, field, operator, values)
}

func BuildRelationCondition(
	builder *QueryBuilder,
	field RelationField,
	operator Operator,
	values []string,
) (string, error) {
	return helpers.BuildRelationCondition(builder, field, operator, values)
}

func BuildCursorCondition(
	builder *QueryBuilder,
	sort SortField,
	direction Direction,
	value string,
	cursorId string,
	cursorExpression string,
) string {
	return helpers.BuildCursorCondition(builder, sort, direction, value, cursorId, cursorExpression)
}

func EscapeLike(value string) string {
	return helpers.EscapeLike(value)
}
