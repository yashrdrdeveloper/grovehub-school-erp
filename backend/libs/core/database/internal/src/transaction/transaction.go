// Package transaction coordinates shared GORM transactions.
package transaction

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/helpers"
)

var ErrMissingDatabase = errors.New("transaction database is nil")

func Database(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if state, ok := helpers.FromContext(ctx); ok {
		return state.Database(ctx)
	}

	return fallback.WithContext(ctx)
}

func Active(ctx context.Context) bool {
	_, ok := helpers.FromContext(ctx)

	return ok
}

func AfterCommit(ctx context.Context, callback func()) bool {
	state, ok := helpers.FromContext(ctx)
	if !ok {
		return false
	}

	state.AddCallback(callback)

	return true
}

func DetachedContext(ctx context.Context) context.Context {
	return helpers.DetachedContext(ctx)
}

// Run joins an existing transaction or starts a new one. Callbacks are executed
// only after a successful commit; rollback discards them.
func Run(ctx context.Context, db *gorm.DB, operation func(context.Context) error) error {
	if db == nil {
		return ErrMissingDatabase
	}

	if Active(ctx) {
		return operation(ctx)
	}

	var state *helpers.State

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transactionContext, boundState := helpers.Bind(ctx, tx)
		state = boundState

		return operation(transactionContext)
	})
	if err != nil {
		return err
	}

	state.PublishCallbacks()

	return nil
}
