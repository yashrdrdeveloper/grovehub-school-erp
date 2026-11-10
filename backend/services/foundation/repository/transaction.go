package repository

import "context"

// UnitOfWork coordinates Person and Contact repositories in one database
// transaction without exposing GORM to the domain service contract.
type UnitOfWork interface {
	Run(ctx context.Context, operation func(context.Context) error) error
}
