package errors

import (
	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"

	repositoryconstants "grovehub/services/foundation/repository/constants"
	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
)

/** MapPersonWriteError translates known PostgreSQL constraint failures. */
func MapPersonWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if stderrors.As(err, &postgresError) && postgresError.Code == "23505" &&
		postgresError.ConstraintName == repositoryconstants.PersonIdempotencyKeyUniqueConstraint {
		return repositoryerrors.ErrPersonIdempotencyKeyAlreadyExists
	}
	return err
}
