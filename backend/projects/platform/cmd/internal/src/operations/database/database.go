package database

import (
	"context"

	"gorm.io/gorm"
)

// Database owns operator connection acquisition, identity verification and cleanup.
type Database interface {
	Business(context.Context) (*gorm.DB, func(), error)
	Audit(context.Context, *gorm.DB) (*gorm.DB, func(), error)
	CurrentName(context.Context, *gorm.DB) (string, error)
}
