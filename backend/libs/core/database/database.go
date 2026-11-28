// Package database exposes generic PostgreSQL and transaction mechanics.
package database

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/impl"
)

// UnitOfWork coordinates repository operations in a context-bound transaction.
type UnitOfWork interface {
	transaction.UnitOfWork
}

var ErrMissingDatabase = transaction.ErrMissingDatabase

func NewDB(databaseURL string) (*gorm.DB, error) {
	return postgres.NewDB(databaseURL)
}

func OpenProjectPool(ctx context.Context, dsn string) (*sql.DB, error) {
	return postgres.OpenProjectPool(ctx, dsn)
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return impl.NewUnitOfWork(db)
}

func TransactionDB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	return transaction.Database(ctx, fallback)
}

func RunTransaction(ctx context.Context, db *gorm.DB, operation func(context.Context) error) error {
	return transaction.Run(ctx, db, operation)
}

func TransactionActive(ctx context.Context) bool {
	return transaction.Active(ctx)
}

func AfterCommit(ctx context.Context, callback func()) bool {
	return transaction.AfterCommit(ctx, callback)
}

func DetachedContext(ctx context.Context) context.Context {
	return transaction.DetachedContext(ctx)
}
