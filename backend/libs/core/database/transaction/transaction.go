// Package transaction coordinates a shared GORM transaction across repositories.
package transaction

import (
 "context"
 "errors"

 "gorm.io/gorm"
)

type contextKey struct{}
type state struct {
 db *gorm.DB
 afterCommit []func()
}

var ErrMissingDatabase = errors.New("transaction database is nil")

// Database returns the transaction bound to ctx or the repository's regular DB.
func Database(ctx context.Context, fallback *gorm.DB) *gorm.DB {
 if s, ok := ctx.Value(contextKey{}).(*state); ok { return s.db.WithContext(ctx) }
 return fallback.WithContext(ctx)
}

// AfterCommit schedules a callback if inside a transaction, otherwise returns false.
func AfterCommit(ctx context.Context, callback func()) bool {
 s, ok := ctx.Value(contextKey{}).(*state)
 if !ok { return false }
 s.afterCommit = append(s.afterCommit, callback)
 return true
}

// Run joins an existing transaction or starts a GORM transaction.
// All registered callbacks run after commit and are discarded on rollback.
func Run(ctx context.Context, db *gorm.DB, operation func(context.Context) error) error {
 if db == nil { return ErrMissingDatabase }
 if _, ok := ctx.Value(contextKey{}).(*state); ok { return operation(ctx) }
 s := &state{}
 err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
  s.db = tx
  return operation(context.WithValue(ctx, contextKey{}, s))
 })
 if err != nil { return err }
 for _, callback := range s.afterCommit { callback() }
 return nil
}
