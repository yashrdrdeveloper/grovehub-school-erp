package helpers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"grovehub/services/registry/internal/app/repository/constants"
)

const (
	uniqueViolation = "23505"
	foreignKeyViolation = "23503"
)

func WriteError(err error) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) || errors.Is(err, gorm.ErrForeignKeyViolated) {
		return constants.ErrConflict
	}
	var postgres *pgconn.PgError
	if errors.As(err, &postgres) && (postgres.Code == uniqueViolation || postgres.Code == foreignKeyViolation) {
		return constants.ErrConflict
	}
	return err
}

func UpdateError(result *gorm.DB) error {
	if result.Error != nil {
		return WriteError(result.Error)
	}
	if result.RowsAffected != 1 {
		return constants.ErrConflict
	}
	return nil
}
