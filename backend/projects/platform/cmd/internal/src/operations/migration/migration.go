package migration

import (
	"context"

	"gorm.io/gorm"
)

// Migration is a service-owned schema migration.
type Migration interface {
	Run(context.Context, *gorm.DB) error
}
