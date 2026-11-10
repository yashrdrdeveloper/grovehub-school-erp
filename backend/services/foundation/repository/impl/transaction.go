package impl

import (
	"context"

	"gorm.io/gorm"

	coretransaction "grovehub/libs/core/database/transaction"
	foundationrepository "grovehub/services/foundation/repository"
)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) foundationrepository.UnitOfWork {
	return &UnitOfWork{db: db}
}

func (unit *UnitOfWork) Run(ctx context.Context, operation func(context.Context) error) error {
	return coretransaction.Run(ctx, unit.db, operation)
}
