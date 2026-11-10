package helpers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
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

/** ContactSearchField constructs a trusted, read-only contact relation filter. */
func ContactSearchField(contactType entity.ContactType, conditions ...coresearch.FixedCondition) coresearch.RelationField {
	return coresearch.RelationField{
		Table: entity.Contact{}.TableName(),
		Alias: "contact",
		ParentCondition: "contact.person_gid = p.gid",
		ValueExpression: "contact.normalized_value",
		FixedConditions: append(
			[]coresearch.FixedCondition{{Expression: "contact.type = ", Value: contactType}},
			conditions...,
		),
	}
}
