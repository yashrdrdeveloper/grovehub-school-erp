// Compatibility API; use annotation.NewEntity and annotation.FromEntity.
package entity

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/entity"

type Builder[T any] = entity.Builder[T]

func New[T any]() Builder[T] { return entity.New[T]() }
func From[T any](value T) Builder[T] { return entity.From(value) }
