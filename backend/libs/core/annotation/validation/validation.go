// Compatibility API; use the annotation facade for new consumers.
package validation

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/validation"

type Error = validation.Error
type Rules[T any] = validation.Rules[T]

func Process[T any](input T, rules Rules[T]) (T, error) { return validation.Process(input, rules) }
func Normalize[T any](input T) (T, error) { return validation.Normalize(input) }
func NormalizeFields[T any](input T) (T, error) { return validation.NormalizeFields(input) }
func Validate[T any](input T) error { return validation.Validate(input) }
