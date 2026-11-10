package helpers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	repositoryerrors "grovehub/services/foundation/repository/constants/errors"
)

const uniqueContactConstraint = "uq_person_contacts_type_normalized_value"

/** MapContactWriteError translates known PostgreSQL contact uniqueness failures. */
func MapContactWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" &&
		postgresError.ConstraintName == uniqueContactConstraint {
		return repositoryerrors.ErrContactAlreadyExists
	}
	return err
}
