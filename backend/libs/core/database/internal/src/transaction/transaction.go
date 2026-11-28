// Package transaction exposes private transaction operations.
package transaction

import (
	"context"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/active"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/aftercommit"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/detachedcontext"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/resolve"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/run"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/unitofwork"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/unitofwork/impl"
)

// UnitOfWork exposes the operation contract to the database facade.
type UnitOfWork interface {
	unitofwork.UnitOfWork
}

var ErrMissingDatabase = run.ErrMissingDatabase

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return impl.NewUnitOfWork(db)
}

func Database(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	return resolve.Database(ctx, fallback)
}

func Active(ctx context.Context) bool {
	return active.Active(ctx)
}

func AfterCommit(ctx context.Context, callback func()) bool {
	return aftercommit.AfterCommit(ctx, callback)
}

func DetachedContext(ctx context.Context) context.Context {
	return detachedcontext.DetachedContext(ctx)
}

func Run(ctx context.Context, db *gorm.DB, operation func(context.Context) error) error {
	return run.Run(ctx, db, operation)
}
