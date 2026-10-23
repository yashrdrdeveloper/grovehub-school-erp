package impl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

const uniqueContactConstraint = "uq_person_contacts_type_normalized_value"

type PersonRepository struct {
	pool *pgxpool.Pool
}

type personSearchRow struct {
	gid        entity.GID
	firstName  string
	middleName string
	lastName   string
	status     entity.PersonStatus
	version    int64
	createdAt  time.Time
	updatedAt  time.Time
}

var personScalarSearchFields = map[string]coresearch.ScalarField{
	"gid":        {Expression: "p.gid::text"},
	"firstName":  {Expression: "LOWER(p.first_name)"},
	"middleName": {Expression: "LOWER(COALESCE(p.middle_name, ''))"},
	"lastName":   {Expression: "LOWER(p.last_name)"},
	"status":     {Expression: "p.status"},
	"createdAt":  {Expression: "p.created_at", Cast: "::timestamptz"},
	"updatedAt":  {Expression: "p.updated_at", Cast: "::timestamptz"},
}

var personRelationSearchFields = map[string]coresearch.RelationField{
	"contacts.email":         contactSearchField(entity.ContactTypeEmail),
	"contacts.phone":         contactSearchField(entity.ContactTypePhone),
	"contacts.primaryEmail":  contactSearchField(entity.ContactTypeEmail, coresearch.FixedCondition{Expression: "contact.is_primary = ", Value: true}),
	"contacts.primaryPhone":  contactSearchField(entity.ContactTypePhone, coresearch.FixedCondition{Expression: "contact.is_primary = ", Value: true}),
	"contacts.verifiedEmail": contactSearchField(entity.ContactTypeEmail, coresearch.FixedCondition{Expression: "contact.is_verified = ", Value: true}),
	"contacts.verifiedPhone": contactSearchField(entity.ContactTypePhone, coresearch.FixedCondition{Expression: "contact.is_verified = ", Value: true}),
}

var personSortFields = map[string]coresearch.SortField{
	"gid":       {Expression: "p.gid", Cast: "::uuid"},
	"firstName": {Expression: "LOWER(p.first_name)", Cast: "::text"},
	"lastName":  {Expression: "LOWER(p.last_name)", Cast: "::text"},
	"createdAt": {Expression: "p.created_at", Cast: "::timestamptz"},
	"updatedAt": {Expression: "p.updated_at", Cast: "::timestamptz"},
}

func NewPersonRepository(pool *pgxpool.Pool) foundationrepository.PersonRepository {
	return &PersonRepository{pool: pool}
}

func (repository *PersonRepository) Create(ctx context.Context, person entity.Person) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO foundation.persons (
			gid, first_name, middle_name, last_name, status, version, created_at, updated_at
		) VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7, $8)`,
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
	row, err := repository.findPersonRow(ctx, gid)
	if err != nil {
		return entity.Person{}, err
	}
	contacts, err := repository.findContactsByGID(ctx, gid)
	if err != nil {
		return entity.Person{}, err
	}
	return row.toPerson(contacts), nil
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
		conditions = append(
			conditions,
			fmt.Sprintf("(type = $%d AND normalized_value = $%d)", len(args)+1, len(args)+2),
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
	return gids, rows.Err()
}

func (repository *PersonRepository) Search(
	ctx context.Context,
	query foundationrepository.PersonSearchQuery,
) (foundationrepository.PersonSearchResult, error) {
	builder := &coresearch.QueryBuilder{}
	statement := strings.Builder{}
	statement.WriteString(`SELECT
		p.gid, p.first_name, p.middle_name, p.last_name,
		p.status, p.version, p.created_at, p.updated_at
	 FROM foundation.persons AS p
	 WHERE TRUE`)

	for _, term := range query.SearchTerms {
		statement.WriteString(" AND p.search_details ILIKE " + builder.Add("%"+coresearch.EscapeLike(term)+"%"))
	}

	for _, filter := range query.Filters {
		condition, err := buildPersonFilter(builder, filter)
		if err != nil {
			return foundationrepository.PersonSearchResult{}, err
		}
		statement.WriteString(" AND " + condition)
	}

	sortField, ok := personSortFields[query.Sort.Field]
	if !ok {
		return foundationrepository.PersonSearchResult{}, foundationrepository.ErrInvalidPersonSearchQuery
	}

	if query.Cursor != nil {
		statement.WriteString(" AND " + coresearch.BuildCursorCondition(
			builder,
			sortField,
			query.Sort.Direction,
			query.Cursor.Value,
			string(query.Cursor.GID),
			"p.gid",
		))
	}

	direction := string(query.Sort.Direction)
	statement.WriteString(" ORDER BY " + sortField.Expression + " " + direction)
	if query.Sort.Field != "gid" {
		statement.WriteString(", p.gid " + direction)
	}
	statement.WriteString(" LIMIT " + builder.Add(query.Limit+1))

	rows, err := repository.pool.Query(ctx, statement.String(), builder.Args()...)
	if err != nil {
		return foundationrepository.PersonSearchResult{}, err
	}
	defer rows.Close()

	personRows, err := scanPersonSearchRows(rows, query.Limit)
	if err != nil {
		return foundationrepository.PersonSearchResult{}, err
	}
	hasNext := len(personRows) > query.Limit
	if hasNext {
		personRows = personRows[:query.Limit]
	}

	gids := make([]entity.GID, 0, len(personRows))
	for _, row := range personRows {
		gids = append(gids, row.gid)
	}
	contactsByGID, err := repository.findContactsByGIDs(ctx, gids)
	if err != nil {
		return foundationrepository.PersonSearchResult{}, err
	}

	persons := make([]entity.Person, 0, len(personRows))
	for _, row := range personRows {
		persons = append(persons, row.toPerson(contactsByGID[row.gid]))
	}
	return foundationrepository.PersonSearchResult{Persons: persons, HasNext: hasNext}, nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

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
	result, err := repository.pool.Exec(ctx, `DELETE FROM foundation.persons WHERE gid = $1`, string(gid))
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return foundationrepository.ErrPersonNotFound
	}
	return nil
}

func (repository *PersonRepository) findPersonRow(ctx context.Context, gid entity.GID) (personSearchRow, error) {
	var row personSearchRow
	var gidValue string
	var middleName *string
	var status string
	err := repository.pool.QueryRow(
		ctx,
		`SELECT gid, first_name, middle_name, last_name, status, version, created_at, updated_at
		 FROM foundation.persons WHERE gid = $1`,
		string(gid),
	).Scan(
		&gidValue,
		&row.firstName,
		&middleName,
		&row.lastName,
		&status,
		&row.version,
		&row.createdAt,
		&row.updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return personSearchRow{}, foundationrepository.ErrPersonNotFound
	}
	if err != nil {
		return personSearchRow{}, err
	}
	row.gid = entity.GID(gidValue)
	row.status = entity.PersonStatus(status)
	if middleName != nil {
		row.middleName = *middleName
	}
	return row, nil
}

func (repository *PersonRepository) findContactsByGID(ctx context.Context, gid entity.GID) ([]entity.Contact, error) {
	contactsByGID, err := repository.findContactsByGIDs(ctx, []entity.GID{gid})
	return contactsByGID[gid], err
}

func (repository *PersonRepository) findContactsByGIDs(
	ctx context.Context,
	gids []entity.GID,
) (map[entity.GID][]entity.Contact, error) {
	contactsByGID := make(map[entity.GID][]entity.Contact, len(gids))
	if len(gids) == 0 {
		return contactsByGID, nil
	}

	placeholders := make([]string, 0, len(gids))
	args := make([]any, 0, len(gids))
	for _, gid := range gids {
		args = append(args, string(gid))
		placeholders = append(placeholders, fmt.Sprintf("$%d::uuid", len(args)))
	}

	rows, err := repository.pool.Query(
		ctx,
		`SELECT person_gid, type, value, normalized_value, is_primary, is_verified
		 FROM foundation.person_contacts
		 WHERE person_gid IN (`+strings.Join(placeholders, ", ")+`)
		 ORDER BY person_gid, type, is_primary DESC, normalized_value`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gid, contactType, value, normalizedValue string
		var isPrimary, isVerified bool
		if err := rows.Scan(&gid, &contactType, &value, &normalizedValue, &isPrimary, &isVerified); err != nil {
			return nil, err
		}
		personGID := entity.GID(gid)
		contactsByGID[personGID] = append(
			contactsByGID[personGID],
			entity.NewContact(entity.ContactType(contactType), value, normalizedValue, isPrimary, isVerified),
		)
	}
	return contactsByGID, rows.Err()
}

func buildPersonFilter(
	builder *coresearch.QueryBuilder,
	filter foundationrepository.PersonSearchFilter,
) (string, error) {
	if field, ok := personScalarSearchFields[filter.Field]; ok {
		condition, err := coresearch.BuildScalarCondition(builder, field, filter.Operator, filter.Values)
		return mapSearchBuildError(condition, err)
	}
	if field, ok := personRelationSearchFields[filter.Field]; ok {
		condition, err := coresearch.BuildRelationCondition(builder, field, filter.Operator, filter.Values)
		return mapSearchBuildError(condition, err)
	}
	return "", foundationrepository.ErrInvalidPersonSearchQuery
}

func mapSearchBuildError(condition string, err error) (string, error) {
	if errors.Is(err, coresearch.ErrUnsupportedOperator) {
		return "", foundationrepository.ErrInvalidPersonSearchQuery
	}
	return condition, err
}

func contactSearchField(contactType entity.ContactType, conditions ...coresearch.FixedCondition) coresearch.RelationField {
	return coresearch.RelationField{
		Table:           "foundation.person_contacts",
		Alias:           "contact",
		ParentCondition: "contact.person_gid = p.gid",
		ValueExpression: "contact.normalized_value",
		FixedConditions: append(
			[]coresearch.FixedCondition{{Expression: "contact.type = ", Value: string(contactType)}},
			conditions...,
		),
	}
}

func scanPersonSearchRows(rows pgx.Rows, limit int) ([]personSearchRow, error) {
	personRows := make([]personSearchRow, 0, limit+1)
	for rows.Next() {
		var row personSearchRow
		var gid, status string
		var middleName *string
		if err := rows.Scan(
			&gid,
			&row.firstName,
			&middleName,
			&row.lastName,
			&status,
			&row.version,
			&row.createdAt,
			&row.updatedAt,
		); err != nil {
			return nil, err
		}
		row.gid = entity.GID(gid)
		row.status = entity.PersonStatus(status)
		if middleName != nil {
			row.middleName = *middleName
		}
		personRows = append(personRows, row)
	}
	return personRows, rows.Err()
}

func (row personSearchRow) toPerson(contacts []entity.Contact) entity.Person {
	return entity.NewPerson(
		row.gid,
		row.firstName,
		row.middleName,
		row.lastName,
		contacts,
		row.status,
		row.version,
		row.createdAt,
		row.updatedAt,
	)
}

func upsertPersonContacts(ctx context.Context, tx pgx.Tx, person entity.Person, clearPrimary bool) error {
	gid := string(person.GetGID())
	if clearPrimary {
		if _, err := tx.Exec(ctx, `UPDATE foundation.person_contacts SET is_primary = FALSE WHERE person_gid = $1`, gid); err != nil {
			return err
		}
	}

	contacts := person.GetContacts()
	for _, contact := range contacts {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO foundation.person_contacts (
				person_gid, type, value, normalized_value, is_primary, is_verified, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
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
	return deleteStalePersonContacts(ctx, tx, gid, contacts)
}

func deleteStalePersonContacts(ctx context.Context, tx pgx.Tx, gid string, contacts []entity.Contact) error {
	if len(contacts) == 0 {
		_, err := tx.Exec(ctx, `DELETE FROM foundation.person_contacts WHERE person_gid = $1`, gid)
		return err
	}

	conditions := make([]string, 0, len(contacts))
	args := []any{gid}
	for _, contact := range contacts {
		conditions = append(
			conditions,
			fmt.Sprintf("(type = $%d AND normalized_value = $%d)", len(args)+1, len(args)+2),
		)
		args = append(args, string(contact.GetType()), contact.GetNormalizedValue())
	}
	_, err := tx.Exec(
		ctx,
		`DELETE FROM foundation.person_contacts
		 WHERE person_gid = $1 AND NOT (`+strings.Join(conditions, " OR ")+`)`,
		args...,
	)
	return err
}

func refreshPersonSearchDetails(ctx context.Context, tx pgx.Tx, gid entity.GID) error {
	_, err := tx.Exec(
		ctx,
		`UPDATE foundation.persons AS person
		 SET search_details = LOWER(TRIM(CONCAT_WS(
			' ', person.gid::text, person.first_name, person.middle_name, person.last_name,
			COALESCE((
				SELECT STRING_AGG(
					CONCAT_WS(' ', contact.value, contact.normalized_value),
					' ' ORDER BY contact.type, contact.normalized_value
				)
				FROM foundation.person_contacts AS contact
				WHERE contact.person_gid = person.gid
			), '')
		 )))
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
