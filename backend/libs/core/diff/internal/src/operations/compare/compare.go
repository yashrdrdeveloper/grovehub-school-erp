// Package compare owns the JSON comparison operation.
package compare

import (
	"fmt"

	"grovehub/libs/core/diff/internal/src/contracts"
	"grovehub/libs/core/diff/internal/src/operations/compare/helpers"
)

func Compare[T any](old, next T) ([]contracts.Change, error) {
	before, err := helpers.JSONValue(old)
	if err != nil {
		return nil, fmt.Errorf("marshal old value: %w", err)
	}

	after, err := helpers.JSONValue(next)
	if err != nil {
		return nil, fmt.Errorf("marshal new value: %w", err)
	}

	return helpers.Changes(before, after), nil
}
