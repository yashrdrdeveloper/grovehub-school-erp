package migration

import (
	"context"

	"gorm.io/gorm"
)

// Migration is a service-owned schema migration.
type Migration interface {
	Run(context.Context, *gorm.DB) error
}

// Operation adapts the service's public migration function to the interface.
type Operation func(context.Context, *gorm.DB) error

func (operation Operation) Run(ctx context.Context, db *gorm.DB) error {
	return operation(ctx, db)
}
