package search

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	searchcontract "grovehub/libs/core/contracts/search"
)

var ErrUnsupportedOperator = errors.New("unsupported search operator")

type Cursor struct {
	Field      string                   `json:"field"`
	Direction  searchcontract.Direction `json:"direction"`
	Value      string                   `json:"value"`
	TieBreaker string                   `json:"tieBreaker"`
}

func EncodeCursor(cursor Cursor) (string, error) {
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecodeCursor(value string) (Cursor, error) {
	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, err
	}
	var cursor Cursor
	if err := json.Unmarshal(payload, &cursor); err != nil {
		return Cursor{}, err
	}
	return cursor, nil
}

type QueryBuilder struct {
	args []any
}

func (builder *QueryBuilder) Add(value any) string {
	builder.args = append(builder.args, value)
	return fmt.Sprintf("$%d", len(builder.args))
}

func (builder *QueryBuilder) Args() []any {
	return builder.args
}

type ScalarField struct {
	Expression string
	Cast       string
}

type FixedCondition struct {
	Expression string
	Value      any
}

type RelationField struct {
	Table           string
	Alias           string
	ParentCondition string
	ValueExpression string
	FixedConditions []FixedCondition
}

type SortField struct {
	Expression string
	Cast       string
}

type operatorFactory func(builder *QueryBuilder, field ScalarField, values []string) (string, error)

type relationOperatorFactory func(builder *QueryBuilder, field RelationField, values []string) (string, error)

var scalarOperatorFactories = map[searchcontract.Operator]operatorFactory{
	searchcontract.OperatorEquals:      scalarBinary("="),
	searchcontract.OperatorNotEquals:   scalarBinary("<>"),
	searchcontract.OperatorGreaterThan: scalarBinary(">"),
	searchcontract.OperatorGTE:         scalarBinary(">="),
	searchcontract.OperatorLessThan:    scalarBinary("<"),
	searchcontract.OperatorLTE:         scalarBinary("<="),
	searchcontract.OperatorIn:          scalarSet("IN"),
	searchcontract.OperatorNotIn:       scalarSet("NOT IN"),
	searchcontract.OperatorExists:      scalarExists(false),
	searchcontract.OperatorNotExists:   scalarExists(true),
	searchcontract.OperatorContains:    scalarLike("contains"),
	searchcontract.OperatorStartsWith:  scalarLike("startsWith"),
}

var relationOperatorFactories = map[searchcontract.Operator]relationOperatorFactory{
	searchcontract.OperatorEquals:     relationBinary(false),
	searchcontract.OperatorNotEquals:  relationBinary(true),
	searchcontract.OperatorIn:         relationSet(false),
	searchcontract.OperatorNotIn:      relationSet(true),
	searchcontract.OperatorExists:     relationExists(false),
	searchcontract.OperatorNotExists:  relationExists(true),
	searchcontract.OperatorContains:   relationLike("contains"),
	searchcontract.OperatorStartsWith: relationLike("startsWith"),
}

func BuildScalarCondition(
	builder *QueryBuilder,
	field ScalarField,
	operator searchcontract.Operator,
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
	field RelationField,
	operator searchcontract.Operator,
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
	sort SortField,
	direction searchcontract.Direction,
	value string,
	cursorId string,
	cursorExpression string,
) string {
	comparison := map[searchcontract.Direction]string{
		searchcontract.DirectionAscending:  ">",
		searchcontract.DirectionDescending: "<",
	}[direction]
	valuePlaceholder := builder.Add(value)
	cursorPlaceholder := builder.Add(cursorId)

	if sort.Expression == cursorExpression {
		return fmt.Sprintf("%s %s %s%s", cursorExpression, comparison, valuePlaceholder, sort.Cast)
	}

	return fmt.Sprintf(
		"(%s %s %s%s OR (%s = %s%s AND %s %s %s))",
		sort.Expression,
		comparison,
		valuePlaceholder,
		sort.Cast,
		sort.Expression,
		valuePlaceholder,
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

func scalarBinary(operator string) operatorFactory {
	return func(builder *QueryBuilder, field ScalarField, values []string) (string, error) {
		if len(values) != 1 {
			return "", ErrUnsupportedOperator
		}
		placeholder := builder.Add(values[0])
		return field.Expression + " " + operator + " " + placeholder + field.Cast, nil
	}
}

func scalarSet(operator string) operatorFactory {
	return func(builder *QueryBuilder, field ScalarField, values []string) (string, error) {
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
	return func(_ *QueryBuilder, field ScalarField, values []string) (string, error) {
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
	return func(builder *QueryBuilder, field ScalarField, values []string) (string, error) {
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

func relationBase(builder *QueryBuilder, field RelationField) string {
	conditions := []string{field.ParentCondition}
	for _, condition := range field.FixedConditions {
		conditions = append(conditions, condition.Expression+builder.Add(condition.Value))
	}
	return strings.Join(conditions, " AND ")
}

func relationQuery(field RelationField, base string, extra string, negated bool) string {
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
	return func(builder *QueryBuilder, field RelationField, values []string) (string, error) {
		if len(values) != 1 {
			return "", ErrUnsupportedOperator
		}
		base := relationBase(builder, field)
		extra := field.ValueExpression + " = " + builder.Add(values[0])
		return relationQuery(field, base, extra, negated), nil
	}
}

func relationSet(negated bool) relationOperatorFactory {
	return func(builder *QueryBuilder, field RelationField, values []string) (string, error) {
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
	return func(builder *QueryBuilder, field RelationField, values []string) (string, error) {
		if len(values) != 0 {
			return "", ErrUnsupportedOperator
		}
		return relationQuery(field, relationBase(builder, field), "", negated), nil
	}
}

func relationLike(mode string) relationOperatorFactory {
	return func(builder *QueryBuilder, field RelationField, values []string) (string, error) {
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
