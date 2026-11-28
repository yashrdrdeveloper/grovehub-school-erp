package query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/helpers"
	"grovehub/libs/core/search/internal/src/operations/page"
)

// applyMultiSort builds a lexicographic keyset predicate using only allowlisted expressions.
func applyMultiSort[TId any](db *gorm.DB, query contracts.Query[TId], config contracts.Config) (*gorm.DB, error) {
	rules := append([]contracts.QuerySort{query.Sort}, query.SecondarySorts...)
	seen := map[string]bool{}
	fields := make([]contracts.SortField, 0, len(rules)+1)
	directions := make([]contracts.Direction, 0, len(rules)+1)
	hasCursorField := false
	for _, rule := range rules {
		field, ok := config.SortFields[rule.Field]
		if !ok || seen[rule.Field] || (rule.Direction != contracts.DirectionAscending && rule.Direction != contracts.DirectionDescending) {
			return nil, contracts.ErrInvalidQuery
		}
		seen[rule.Field] = true
		if field.Expression == config.CursorExpression {
			hasCursorField = true
		}
		fields = append(fields, field)
		directions = append(directions, rule.Direction)
	}
	if !hasCursorField {
		fields = append(fields, contracts.SortField{Expression: config.CursorExpression})
		directions = append(directions, query.Sort.Direction)
	}
	if query.Cursor != nil {
		if len(query.Cursor.SecondaryValues) != len(rules)-1 {
			return nil, contracts.ErrInvalidQuery
		}
		values := append([]string{query.Cursor.Value}, query.Cursor.SecondaryValues...)
		if !hasCursorField {
			values = append(values, fmt.Sprint(query.Cursor.CursorId))
		}
		builder := &helpers.QueryBuilder{}
		parts := make([]string, 0, len(fields))
		for index, field := range fields {
			terms := make([]string, 0, index+1)
			for previous := 0; previous < index; previous++ {
				terms = append(terms, fields[previous].Expression+" = "+builder.Add(values[previous])+fields[previous].Cast)
			}
			op := ">"
			if directions[index] == contracts.DirectionDescending {
				op = "<"
			}
			terms = append(terms, field.Expression+" "+op+" "+builder.Add(values[index])+field.Cast)
			parts = append(parts, "("+strings.Join(terms, " AND ")+")")
		}
		db = db.Where("("+strings.Join(parts, " OR ")+")", builder.Args()...)
	}
	for index, field := range fields {
		db = db.Order(field.Expression + " " + string(directions[index]))
	}
	return page.ApplyPage(db, query), nil
}
