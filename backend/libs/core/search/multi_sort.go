package search

import (
	"fmt"
	"strings"
	"gorm.io/gorm"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/helpers"
)

// NormalizeSorts preserves request priority while rejecting duplicates.
func NormalizeSorts(request []Sort, fallback QuerySort, allowed map[string]CursorValueValidator, max int, invalid error) ([]QuerySort, error) {
	if invalid == nil { invalid = ErrInvalidQuery }
	if len(request) > max { return nil, invalid }
	if len(request) == 0 {
		rule, err := NormalizeSort(nil, fallback, allowed, invalid)
		if err != nil { return nil, err }
		return []QuerySort{rule}, nil
	}
	rules := make([]QuerySort, 0, len(request))
	seen := map[string]bool{}
	for _, item := range request {
		rule, err := NormalizeSort([]Sort{item}, fallback, allowed, invalid)
		if err != nil { return nil, err }
		if seen[rule.Field] { return nil, invalid }
		seen[rule.Field] = true
		rules = append(rules, rule)
	}
	return rules, nil
}

// NormalizeMultiCursor binds every cursor value to its requested field and direction.
func NormalizeMultiCursor(value string, rules []QuerySort, allowed map[string]CursorValueValidator, validateID CursorValueValidator, invalid error) (*Cursor[string], error) {
	if invalid == nil { invalid = ErrInvalidQuery }
	if len(rules) == 0 || validateID == nil { return nil, invalid }
	for _, rule := range rules { if allowed[rule.Field] == nil { return nil, invalid } }
	if strings.TrimSpace(value) == "" { return nil, nil }
	payload, err := DecodeCursor(value)
	if err != nil || payload.Field != rules[0].Field || payload.Direction != rules[0].Direction || len(payload.SecondarySorts) != len(rules)-1 || len(payload.SecondaryValues) != len(rules)-1 {
		return nil, invalid
	}
	if validateID(payload.CursorId, invalid) != nil || allowed[rules[0].Field](payload.Value, invalid) != nil { return nil, invalid }
	for index, rule := range rules[1:] {
		if payload.SecondarySorts[index].Field != rule.Field || payload.SecondarySorts[index].Direction != rule.Direction || allowed[rule.Field](payload.SecondaryValues[index], invalid) != nil {
			return nil, invalid
		}
	}
	return &Cursor[string]{Value: payload.Value, SecondaryValues: payload.SecondaryValues, CursorId: payload.CursorId}, nil
}

// applyMultiSort builds a lexicographic keyset predicate using only allowlisted expressions.
func applyMultiSort[TId any](db *gorm.DB, query Query[TId], config Config) (*gorm.DB, error) {
	rules := append([]QuerySort{query.Sort}, query.SecondarySorts...)
	seen := map[string]bool{}
	fields := make([]SortField, 0, len(rules)+1)
	directions := make([]Direction, 0, len(rules)+1)
	hasCursorField := false
	for _, rule := range rules {
		field, ok := config.SortFields[rule.Field]
		if !ok || seen[rule.Field] || (rule.Direction != DirectionAscending && rule.Direction != DirectionDescending) { return nil, ErrInvalidQuery }
		seen[rule.Field] = true
		if field.Expression == config.CursorExpression { hasCursorField = true }
		fields = append(fields, field)
		directions = append(directions, rule.Direction)
	}
	if !hasCursorField {
		fields = append(fields, SortField{Expression: config.CursorExpression})
		directions = append(directions, query.Sort.Direction)
	}
	if query.Cursor != nil {
		if len(query.Cursor.SecondaryValues) != len(rules)-1 { return nil, ErrInvalidQuery }
		values := append([]string{query.Cursor.Value}, query.Cursor.SecondaryValues...)
		if !hasCursorField { values = append(values, fmt.Sprint(query.Cursor.CursorId)) }
		builder := &helpers.QueryBuilder{}
		parts := make([]string, 0, len(fields))
		for index, field := range fields {
			terms := make([]string, 0, index+1)
			for previous := 0; previous < index; previous++ {
				terms = append(terms, fields[previous].Expression+" = "+builder.Add(values[previous])+fields[previous].Cast)
			}
			op := ">"
			if directions[index] == DirectionDescending { op = "<" }
			terms = append(terms, field.Expression+" "+op+" "+builder.Add(values[index])+field.Cast)
			parts = append(parts, "("+strings.Join(terms, " AND ")+")")
		}
		db = db.Where("("+strings.Join(parts, " OR ")+")", builder.Args()...)
	}
	for index, field := range fields { db = db.Order(field.Expression+" "+string(directions[index])) }
	return applyPage(db, query), nil
}
