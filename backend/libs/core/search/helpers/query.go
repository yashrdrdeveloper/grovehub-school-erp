package helpers

import (
	"errors"
	"fmt"
	"strings"

	"grovehub/libs/core/search/contracts"
)

var (
	ErrInvalidField        = errors.New("invalid search field")
	ErrUnsupportedOperator = errors.New("unsupported search operator")
)

type QueryBuilder struct {
	args []any
}

func (builder *QueryBuilder) Add(value any) string {
	builder.args = append(builder.args, value)
	return "?"
}

func (builder *QueryBuilder) Args() []any {
	return builder.args
}

func BuildFilter(
	builder *QueryBuilder,
	filter contracts.QueryFilter,
	config contracts.Config,
) (string, error) {
	if field, ok := config.ScalarFields[filter.Field]; ok {
		return BuildScalarCondition(builder, field, filter.Operator, filter.Values)
	}
	if field, ok := config.RelationFields[filter.Field]; ok {
		return BuildRelationCondition(builder, field, filter.Operator, filter.Values)
	}
	return "", ErrInvalidField
}

func BuildScalarCondition(
	builder *QueryBuilder,
	field contracts.ScalarField,
	operator contracts.Operator,
	values []string,
) (string, error) {
	factory, ok := scalarOperatorFactories[operator]
	if !ok {
		return "", ErrUnsupportedOperator
	}
	return factory(builder, field, values)
}

func BuildRelationCondition(
	builder *QueryBuilder,
	field contracts.RelationField,
	operator contracts.Operator,
	values []string,
) (string, error) {
	factory, ok := relationOperatorFactories[operator]
	if !ok {
		return "", ErrUnsupportedOperator
	}
	return factory(builder, field, values)
}

func BuildCursorCondition(
	builder *QueryBuilder,
	sort contracts.SortField,
	direction contracts.Direction,
	value string,
	cursorId string,
	cursorExpression string,
) string {
	comparison := map[contracts.Direction]string{
		contracts.DirectionAscending:  ">",
		contracts.DirectionDescending: "<",
	}[direction]
	if sort.Expression == cursorExpression {
		valuePlaceholder := builder.Add(value)
		return fmt.Sprintf("%s %s %s%s", cursorExpression, comparison, valuePlaceholder, sort.Cast)
	}

	comparisonPlaceholder := builder.Add(value)
	equalityPlaceholder := builder.Add(value)
	cursorPlaceholder := builder.Add(cursorId)
	return fmt.Sprintf(
		"(%s %s %s%s OR (%s = %s%s AND %s %s %s))",
		sort.Expression,
		comparison,
		comparisonPlaceholder,
		sort.Cast,
		sort.Expression,
		equalityPlaceholder,
		sort.Cast,
		cursorExpression,
		comparison,
		cursorPlaceholder,
	)
}

func EscapeLike(value string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	)
	return replacer.Replace(value)
}

type operatorFactory func(*QueryBuilder, contracts.ScalarField, []string) (string, error)
type relationOperatorFactory func(*QueryBuilder, contracts.RelationField, []string) (string, error)

var scalarOperatorFactories = map[contracts.Operator]operatorFactory{
	contracts.OperatorEquals:      scalarBinary("="),
	contracts.OperatorNotEquals:   scalarBinary("<>"),
	contracts.OperatorGreaterThan: scalarBinary(">"),
	contracts.OperatorGTE:         scalarBinary(">="),
	contracts.OperatorLessThan:    scalarBinary("<"),
	contracts.OperatorLTE:         scalarBinary("<="),
	contracts.OperatorIn:          scalarSet("IN"),
	contracts.OperatorNotIn:       scalarSet("NOT IN"),
	contracts.OperatorExists:      scalarExists(false),
	contracts.OperatorNotExists:   scalarExists(true),
	contracts.OperatorContains:    scalarLike("contains"),
	contracts.OperatorStartsWith:  scalarLike("startsWith"),
}

var relationOperatorFactories = map[contracts.Operator]relationOperatorFactory{
	contracts.OperatorEquals:     relationBinary(false),
	contracts.OperatorNotEquals:  relationBinary(true),
	contracts.OperatorIn:         relationSet(false),
	contracts.OperatorNotIn:      relationSet(true),
	contracts.OperatorExists:     relationExists(false),
	contracts.OperatorNotExists:  relationExists(true),
	contracts.OperatorContains:   relationLike("contains"),
	contracts.OperatorStartsWith: relationLike("startsWith"),
}

func scalarBinary(operator string) operatorFactory {
	return func(builder *QueryBuilder, field contracts.ScalarField, values []string) (string, error) {
		if len(values) != 1 {
			return "", ErrUnsupportedOperator
		}
		placeholder := builder.Add(values[0])
		return field.Expression + " " + operator + " " + placeholder + field.Cast, nil
	}
}

func scalarSet(operator string) operatorFactory {
	return func(builder *QueryBuilder, field contracts.ScalarField, values []string) (string, error) {
		if len(values) == 0 {
			return "", ErrUnsupportedOperator
		}
		placeholders := make([]string, 0, len(values))
		for _, value := range values {
			placeholders = append(placeholders, builder.Add(value)+field.Cast)
		}
		return field.Expression + " " + operator + " (" + strings.Join(placeholders, ", ") + ")", nil
	}
}

func scalarExists(negated bool) operatorFactory {
	return func(_ *QueryBuilder, field contracts.ScalarField, values []string) (string, error) {
		if len(values) != 0 {
			return "", ErrUnsupportedOperator
		}
		predicate := "IS NOT NULL"
		if negated {
			predicate = "IS NULL"
		}
		return "NULLIF(" + field.Expression + ", '') " + predicate, nil
	}
}

func scalarLike(mode string) operatorFactory {
	return func(builder *QueryBuilder, field contracts.ScalarField, values []string) (string, error) {
		if len(values) != 1 {
			return "", ErrUnsupportedOperator
		}
		value := EscapeLike(values[0])
		if mode == "contains" {
			value = "%" + value + "%"
		} else {
			value += "%"
		}
		return field.Expression + " LIKE " + builder.Add(value), nil
	}
}

func relationBase(builder *QueryBuilder, field contracts.RelationField) string {
	conditions := []string{field.ParentCondition}
	for _, condition := range field.FixedConditions {
		conditions = append(conditions, condition.Expression+builder.Add(condition.Value))
	}
	return strings.Join(conditions, " AND ")
}

func relationQuery(field contracts.RelationField, base string, extra string, negated bool) string {
	prefix := "EXISTS"
	if negated {
		prefix = "NOT EXISTS"
	}
	condition := base
	if extra != "" {
		condition += " AND " + extra
	}
	return prefix + " (SELECT 1 FROM " + field.Table + " AS " + field.Alias + " WHERE " + condition + ")"
}

func relationBinary(negated bool) relationOperatorFactory {
	return func(builder *QueryBuilder, field contracts.RelationField, values []string) (string, error) {
		if len(values) != 1 {
			return "", ErrUnsupportedOperator
		}
		base := relationBase(builder, field)
		extra := field.ValueExpression + " = " + builder.Add(values[0])
		return relationQuery(field, base, extra, negated), nil
	}
}

func relationSet(negated bool) relationOperatorFactory {
	return func(builder *QueryBuilder, field contracts.RelationField, values []string) (string, error) {
		if len(values) == 0 {
			return "", ErrUnsupportedOperator
		}
		placeholders := make([]string, 0, len(values))
		for _, value := range values {
			placeholders = append(placeholders, builder.Add(value))
		}
		base := relationBase(builder, field)
		extra := field.ValueExpression + " IN (" + strings.Join(placeholders, ", ") + ")"
		return relationQuery(field, base, extra, negated), nil
	}
}

func relationExists(negated bool) relationOperatorFactory {
	return func(builder *QueryBuilder, field contracts.RelationField, values []string) (string, error) {
		if len(values) != 0 {
			return "", ErrUnsupportedOperator
		}
		return relationQuery(field, relationBase(builder, field), "", negated), nil
	}
}

func relationLike(mode string) relationOperatorFactory {
	return func(builder *QueryBuilder, field contracts.RelationField, values []string) (string, error) {
		if len(values) != 1 {
			return "", ErrUnsupportedOperator
		}
		value := EscapeLike(values[0])
		if mode == "contains" {
			value = "%" + value + "%"
		} else {
			value += "%"
		}
		base := relationBase(builder, field)
		extra := field.ValueExpression + " LIKE " + builder.Add(value)
		return relationQuery(field, base, extra, false), nil
	}
}
