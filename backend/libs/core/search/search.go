package search

import (
	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/cursor"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/validation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/query"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/page"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/constants"
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

	DefaultRows = constants.DefaultRows
	MaxRows     = constants.MaxRows
)

var (
	ErrInvalidQuery        = contracts.ErrInvalidQuery
	ErrUnsupportedOperator = contracts.ErrUnsupportedOperator
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
type CursorValueValidator = contracts.CursorValueValidator
type ValueNormalizer = contracts.ValueNormalizer
type FieldSpec = contracts.FieldSpec
type QueryBuilder = helpers.QueryBuilder

func Apply[TId any](db *gorm.DB, input Query[TId], config Config) (*gorm.DB, error) {
	return query.Apply(db, input, config)
}

func Execute[TId any, T any](db *gorm.DB, input Query[TId], config Config) (QueryResult[T], error) {
	return query.Execute[TId, T](db, input, config)
}

func EncodeCursor(payload CursorPayload) (string, error) {
	return cursor.EncodeCursor(payload)
}

func DecodeCursor(value string) (CursorPayload, error) {
	return cursor.DecodeCursor(value)
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

 
// Shared request-validation helpers use resource-supplied errors and policies.
func OperatorSet(operators ...Operator) map[Operator]struct{} { return validation.OperatorSet(operators...) }
func EqualityOperators() map[Operator]struct{} { return validation.EqualityOperators() }
func TextOperators() map[Operator]struct{} { return validation.TextOperators() }
func TimeOperators() map[Operator]struct{} { return validation.TimeOperators() }
func IsPartialOperator(operator Operator) bool { return validation.IsPartialOperator(operator) }
func ValidValueCount(operator Operator, count int) bool { return validation.ValidValueCount(operator, count) }

func ValidateUUIDv7Cursor(value string, invalid error) error { return validation.ValidateUUIDv7Cursor(value, invalid) }
func ValidateTimeCursor(value string, invalid error) error { return validation.ValidateTimeCursor(value, invalid) }
func ValidateNonEmptyCursor(value string, invalid error) error { return validation.ValidateNonEmptyCursor(value, invalid) }

func NormalizeSearchTerms(value string, maxLength, maxTerms int, invalid error) ([]string, error) {
    return validation.NormalizeSearchTerms(value, maxLength, maxTerms, invalid)
}

func ResolveRows(rows int, invalid error) (int, error) { return validation.ResolveRows(rows, invalid) }

func NormalizeSort(request []Sort, fallback QuerySort, allowed map[string]CursorValueValidator, invalid error) (QuerySort, error) {
    return validation.NormalizeSort(request, fallback, allowed, invalid)
}

func NormalizeFilters(filters []Filter, fields map[string]FieldSpec, maxFilters, maxValues int, invalid error) ([]QueryFilter, error) {
    return validation.NormalizeFilters(filters, fields, maxFilters, maxValues, invalid)
}

func NormalizeCursor(value string, sort QuerySort, validateValue CursorValueValidator, validateID CursorValueValidator, invalid error) (*Cursor[string], CursorPayload, error) {
    return validation.NormalizeCursor(value, sort, validateValue, validateID, invalid)
}

func NormalizeSorts(request []Sort, fallback QuerySort, allowed map[string]CursorValueValidator, max int, invalid error) ([]QuerySort, error) {
	return validation.NormalizeSorts(request, fallback, allowed, max, invalid)
}

func NormalizeMultiCursor(value string, rules []QuerySort, allowed map[string]CursorValueValidator, validateID CursorValueValidator, invalid error) (*Cursor[string], error) {
	return validation.NormalizeMultiCursor(value, rules, allowed, validateID, invalid)
}

func ResolvePage(request PageRequest, invalid error) (int, int, int, error) {
	return page.ResolvePage(request, invalid)
}

func ValidatePageCursor(request PageRequest, number, rows int, invalid error) error {
	return page.ValidatePageCursor(request, number, rows, invalid)
}

func NextPageNumber(request PageRequest, number int) int {
	return page.NextPageNumber(request, number)
}

func NewPageResponse[TId any,T any](input Query[TId], result QueryResult[T]) PageResponse {
	return page.NewPageResponse(input, result)
}
