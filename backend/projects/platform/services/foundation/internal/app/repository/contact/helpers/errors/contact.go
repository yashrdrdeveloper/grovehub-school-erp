package errors

import (
	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"

	repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
	repositorykeys "grovehub/services/foundation/internal/app/repository/constants/keys"
)

/** MapContactWriteError translates known PostgreSQL contact uniqueness failures. */
func MapContactWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if stderrors.As(err, &postgresError) && postgresError.Code == "23505" &&
		postgresError.ConstraintName == repositorykeys.ContactTypeNormalizedValueUniqueConstraint {
		return repositoryerrors.ErrContactAlreadyExists
	}
	return err
}
