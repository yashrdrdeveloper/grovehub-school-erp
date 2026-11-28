package validation

import (
	"strings"

	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/cursor"
)

// NormalizeSorts preserves request priority while rejecting duplicates.
func NormalizeSorts(request []contracts.Sort, fallback contracts.QuerySort, allowed map[string]contracts.CursorValueValidator, max int, invalid error) ([]contracts.QuerySort, error) {
	if invalid == nil {
		invalid = contracts.ErrInvalidQuery
	}
	if len(request) > max {
		return nil, invalid
	}
	if len(request) == 0 {
		rule, err := NormalizeSort(nil, fallback, allowed, invalid)
		if err != nil {
			return nil, err
		}
		return []contracts.QuerySort{rule}, nil
	}
	rules := make([]contracts.QuerySort, 0, len(request))
	seen := map[string]bool{}
	for _, item := range request {
		rule, err := NormalizeSort([]contracts.Sort{item}, fallback, allowed, invalid)
		if err != nil {
			return nil, err
		}
		if seen[rule.Field] {
			return nil, invalid
		}
		seen[rule.Field] = true
		rules = append(rules, rule)
	}
	return rules, nil
}

// NormalizeMultiCursor binds every cursor value to its requested field and direction.
func NormalizeMultiCursor(value string, rules []contracts.QuerySort, allowed map[string]contracts.CursorValueValidator, validateID contracts.CursorValueValidator, invalid error) (*contracts.Cursor[string], error) {
	if invalid == nil {
		invalid = contracts.ErrInvalidQuery
	}
	if len(rules) == 0 || validateID == nil {
		return nil, invalid
	}
	for _, rule := range rules {
		if allowed[rule.Field] == nil {
			return nil, invalid
		}
	}
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	payload, err := cursor.DecodeCursor(value)
	if err != nil || payload.Field != rules[0].Field || payload.Direction != rules[0].Direction || len(payload.SecondarySorts) != len(rules)-1 || len(payload.SecondaryValues) != len(rules)-1 {
		return nil, invalid
	}
	if validateID(payload.CursorId, invalid) != nil || allowed[rules[0].Field](payload.Value, invalid) != nil {
		return nil, invalid
	}
	for index, rule := range rules[1:] {
		if payload.SecondarySorts[index].Field != rule.Field || payload.SecondarySorts[index].Direction != rule.Direction || allowed[rule.Field](payload.SecondaryValues[index], invalid) != nil {
			return nil, invalid
		}
	}
	return &contracts.Cursor[string]{Value: payload.Value, SecondaryValues: payload.SecondaryValues, CursorId: payload.CursorId}, nil
}
