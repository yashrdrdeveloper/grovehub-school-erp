package run

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"grovehub/libs/core/database/internal/src/transaction/helpers"
	"grovehub/libs/core/database/internal/src/transaction/operations/active"
)

var ErrMissingDatabase = errors.New("transaction database is nil")

// Run joins an existing transaction or starts a new one. Callbacks are
// executed only after successful commit; rollback discards them.
func Run(ctx context.Context, db *gorm.DB, operation func(context.Context) error) error {
	if db == nil {
		return ErrMissingDatabase
	}

	if active.Active(ctx) {
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
