package impl

import (
	"context"

	"gorm.io/gorm"

	"grovehub/cmd/internal/src/operations/migration"
)

type operation func(context.Context, *gorm.DB) error

// New adapts a service-owned public migration function.
func New(run func(context.Context, *gorm.DB) error) migration.Migration {
	return operation(run)
}

func (o operation) Run(ctx context.Context, db *gorm.DB) error {
	return o(ctx, db)
}
