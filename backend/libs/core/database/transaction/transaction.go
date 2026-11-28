package transaction

import (
	"context"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

var ErrMissingDatabase = coredatabase.ErrMissingDatabase

func Database(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	return coredatabase.TransactionDB(ctx, fallback)
}

func Run(ctx context.Context, db *gorm.DB, operation func(context.Context) error) error {
	return coredatabase.RunTransaction(ctx, db, operation)
}

func Active(ctx context.Context) bool {
	return coredatabase.TransactionActive(ctx)
}

func AfterCommit(ctx context.Context, callback func()) bool {
	return coredatabase.AfterCommit(ctx, callback)
}

func DetachedContext(ctx context.Context) context.Context {
	return coredatabase.DetachedContext(ctx)
}
