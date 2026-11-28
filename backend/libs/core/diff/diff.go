// Package diff computes JSON-field-oriented changes without depending on domain DTOs.
package diff

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/operations/compare"
)

// Change describes a single flattened JSON field change.
type Change = contracts.Change

// Compare returns deterministic flat changes using JSON field names and array
// indices. Array elements are compared by position, not by domain identity.
func Compare[T any](old, next T) ([]Change, error) {
	return compare.Compare(old, next)
}
