// Compatibility API; use the annotation facade for new consumers.
package configuration

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/configuration"

type Lookup = configuration.Lookup

func FromEnv[T any]() (T, error) { return configuration.FromEnv[T]() }
func Decode[T any](lookup Lookup) (T, error) { return configuration.Decode[T](lookup) }
