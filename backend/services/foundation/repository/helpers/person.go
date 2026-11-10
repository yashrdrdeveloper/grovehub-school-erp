package helpers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
)

const uniquePersonIdempotencyKeyConstraint = "uq_persons_create_idempotency_key"

/** MapPersonWriteError translates known PostgreSQL constraint failures. */
func MapPersonWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" &&
		postgresError.ConstraintName == uniquePersonIdempotencyKeyConstraint {
		return repositoryerrors.ErrPersonIdempotencyKeyAlreadyExists
	}
	return err
}
