package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

const uniqueContactConstraint = "uq_person_contacts_type_normalized_value"

type PersonRepository struct {
	pool *pgxpool.Pool
}

func NewPersonRepository(pool *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{pool: pool}
}

func (repository *PersonRepository) Create(ctx context.Context, person entity.Person) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO foundation.persons (
			gid,
			first_name,
			middle_name,
			last_name,
			status,
			version,
			created_at,
			updated_at
		)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7, $8)`,
		string(person.GetGID()),
		person.GetFirstName(),
		person.GetMiddleName(),
		person.GetLastName(),
		string(person.GetStatus()),
		person.GetVersion(),
		person.GetCreatedAt(),
		person.GetUpdatedAt(),
	)
	if err != nil {
		return mapPersonWriteError(err)
	}

	if err := upsertPersonContacts(ctx, tx, person, false); err != nil {
		return err
	}
	if err := refreshPersonSearchDetails(ctx, tx, person.GetGID()); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PersonRepository) FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	var gidValue string
	var firstName string
	var middleName *string
	var lastName string
	var status string
	var version int64
	var createdAt time.Time
	var updatedAt time.Time

	err := repository.pool.QueryRow(
		ctx,
		`SELECT gid, first_name, middle_name, last_name, status, version, created_at, updated_at
		 FROM foundation.persons
		 WHERE gid = $1`,
		string(gid),
	).Scan(
		&gidValue,
		&firstName,
		&middleName,
		&lastName,
		&status,
		&version,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Person{}, foundationrepository.ErrPersonNotFound
	}
	if err != nil {
		return entity.Person{}, err
	}

	contacts, err := repository.findContactsByGID(ctx, gid)
	if err != nil {
		return entity.Person{}, err
	}

	middleNameValue := ""
	if middleName != nil {
		middleNameValue = *middleName
	}

	return entity.NewPerson(
		entity.GID(gidValue),
		firstName,
		middleNameValue,
		lastName,
		contacts,
		entity.PersonStatus(status),
		version,
		createdAt,
		updatedAt,
	), nil
}

func (repository *PersonRepository) FindGIDsByContacts(
	ctx context.Context,
	contacts []entity.Contact,
) ([]entity.GID, error) {
	if len(contacts) == 0 {
		return []entity.GID{}, nil
	}

	conditions := make([]string, 0, len(contacts))
	args := make([]any, 0, len(contacts)*2)
	for _, contact := range contacts {
		typePosition := len(args) + 1
		valuePosition := len(args) + 2
		conditions = append(
			conditions,
			fmt.Sprintf("(type = $%d AND normalized_value = $%d)", typePosition, valuePosition),
		)
		args = append(args, string(contact.GetType()), contact.GetNormalizedValue())
	}

	rows, err := repository.pool.Query(
		ctx,
		`SELECT DISTINCT person_gid
		 FROM foundation.person_contacts
		 WHERE `+strings.Join(conditions, " OR ")+`
		 ORDER BY person_gid`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gids := make([]entity.GID, 0)
	for rows.Next() {
		var gid string
		if err := rows.Scan(&gid); err != nil {
			return nil, err
		}
		gids = append(gids, entity.GID(gid))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gids, nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	result, err := tx.Exec(
		ctx,
		`UPDATE foundation.persons
		 SET first_name = $2,
		     middle_name = NULLIF($3, ''),
		     last_name = $4,
		     status = $5,
		     version = $6,
		     updated_at = $7
		 WHERE gid = $1`,
		string(person.GetGID()),
		person.GetFirstName(),
		person.GetMiddleName(),
		person.GetLastName(),
		string(person.GetStatus()),
		person.GetVersion(),
		person.GetUpdatedAt(),
	)
	if err != nil {
		return mapPersonWriteError(err)
	}
	if result.RowsAffected() == 0 {
		return foundationrepository.ErrPersonNotFound
	}

	if err := upsertPersonContacts(ctx, tx, person, true); err != nil {
		return err
	}
	if err := refreshPersonSearchDetails(ctx, tx, person.GetGID()); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PersonRepository) Delete(ctx context.Context, gid entity.GID) error {
	result, err := repository.pool.Exec(
		ctx,
		`DELETE FROM foundation.persons
		 WHERE gid = $1`,
		string(gid),
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return foundationrepository.ErrPersonNotFound
	}

	return nil
}

func (repository *PersonRepository) findContactsByGID(
	ctx context.Context,
	gid entity.GID,
) ([]entity.Contact, error) {
	rows, err := repository.pool.Query(
		ctx,
		`SELECT type, value, normalized_value, is_primary, is_verified
		 FROM foundation.person_contacts
		 WHERE person_gid = $1
		 ORDER BY type, is_primary DESC, normalized_value`,
		string(gid),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]entity.Contact, 0)
	for rows.Next() {
		var contactType string
		var value string
		var normalizedValue string
		var isPrimary bool
		var isVerified bool

		if err := rows.Scan(
			&contactType,
			&value,
			&normalizedValue,
			&isPrimary,
			&isVerified,
		); err != nil {
			return nil, err
		}

		contacts = append(
			contacts,
			entity.NewContact(
				entity.ContactType(contactType),
				value,
				normalizedValue,
				isPrimary,
				isVerified,
			),
		)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func upsertPersonContacts(
	ctx context.Context,
	tx pgx.Tx,
	person entity.Person,
	clearPrimary bool,
) error {
	gid := string(person.GetGID())
	if clearPrimary {
		if _, err := tx.Exec(
			ctx,
			`UPDATE foundation.person_contacts
			 SET is_primary = FALSE
			 WHERE person_gid = $1`,
			gid,
		); err != nil {
			return err
		}
	}

	contacts := person.GetContacts()
	for _, contact := range contacts {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO foundation.person_contacts (
				person_gid,
				type,
				value,
				normalized_value,
				is_primary,
				is_verified,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
			ON CONFLICT (person_gid, type, normalized_value)
			DO UPDATE SET
				value = EXCLUDED.value,
				is_primary = EXCLUDED.is_primary,
				is_verified = EXCLUDED.is_verified,
				updated_at = EXCLUDED.updated_at`,
			gid,
			string(contact.GetType()),
			contact.GetValue(),
			contact.GetNormalizedValue(),
			contact.GetIsPrimary(),
			contact.GetIsVerified(),
			person.GetUpdatedAt(),
		)
		if err != nil {
			return mapPersonWriteError(err)
		}
	}

	if len(contacts) == 0 {
		_, err := tx.Exec(
			ctx,
			`DELETE FROM foundation.person_contacts WHERE person_gid = $1`,
			gid,
		)
		return err
	}

	conditions := make([]string, 0, len(contacts))
	args := make([]any, 0, 1+len(contacts)*2)
	args = append(args, gid)
	for _, contact := range contacts {
		typePosition := len(args) + 1
		valuePosition := len(args) + 2
		conditions = append(
			conditions,
			fmt.Sprintf("(type = $%d AND normalized_value = $%d)", typePosition, valuePosition),
		)
		args = append(args, string(contact.GetType()), contact.GetNormalizedValue())
	}

	_, err := tx.Exec(
		ctx,
		`DELETE FROM foundation.person_contacts
		 WHERE person_gid = $1
		   AND NOT (`+strings.Join(conditions, " OR ")+`)`,
		args...,
	)
	return err
}

func refreshPersonSearchDetails(ctx context.Context, tx pgx.Tx, gid entity.GID) error {
	_, err := tx.Exec(
		ctx,
		`UPDATE foundation.persons AS person
		 SET search_details = LOWER(
			TRIM(
				CONCAT_WS(
					' ',
					person.gid::text,
					person.first_name,
					person.middle_name,
					person.last_name,
					COALESCE(
						(
							SELECT STRING_AGG(
								CONCAT_WS(' ', contact.value, contact.normalized_value),
								' '
								ORDER BY contact.type, contact.normalized_value
							)
							FROM foundation.person_contacts AS contact
							WHERE contact.person_gid = person.gid
						),
						''
					)
				)
			)
		 )
		 WHERE person.gid = $1`,
		string(gid),
	)
	return err
}

func mapPersonWriteError(err error) error {
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) &&
		postgresError.Code == "23505" &&
		postgresError.ConstraintName == uniqueContactConstraint {
		return foundationrepository.ErrPersonContactAlreadyExists
	}

	return err
}
