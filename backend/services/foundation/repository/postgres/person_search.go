package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	searchcontract "grovehub/libs/core/contracts/search"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

type searchQueryBuilder struct {
	args []any
}

func (builder *searchQueryBuilder) add(value any) string {
	builder.args = append(builder.args, value)
	return fmt.Sprintf("$%d", len(builder.args))
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

func (repository *PersonRepository) Search(
	ctx context.Context,
	query foundationrepository.PersonSearchQuery,
) (foundationrepository.PersonSearchResult, error) {
	builder := &searchQueryBuilder{}
	statement := strings.Builder{}
	statement.WriteString(`SELECT
		p.gid,
		p.first_name,
		p.middle_name,
		p.last_name,
		p.status,
		p.version,
		p.created_at,
		p.updated_at
	 FROM foundation.persons AS p
	 WHERE TRUE`)

	for _, term := range query.SearchTerms {
		placeholder := builder.add("%" + escapeLike(term) + "%")
		statement.WriteString(" AND p.search_details ILIKE " + placeholder)
	}

	for _, filter := range query.Filters {
		condition, err := personFilterCondition(builder, filter)
		if err != nil {
			return foundationrepository.PersonSearchResult{}, err
		}
		statement.WriteString(" AND " + condition)
	}

	sortExpression, sortCast, err := personSortExpression(query.Sort.Field)
	if err != nil {
		return foundationrepository.PersonSearchResult{}, err
	}

	if query.Cursor != nil {
		valuePlaceholder := builder.add(query.Cursor.Value)
		gidPlaceholder := builder.add(string(query.Cursor.GID))
		comparison := ">"
		if query.Sort.Direction == searchcontract.DirectionDescending {
			comparison = "<"
		}

		if query.Sort.Field == "gid" {
			statement.WriteString(
				fmt.Sprintf(" AND p.gid %s %s::uuid", comparison, valuePlaceholder),
			)
		} else {
			statement.WriteString(
				fmt.Sprintf(
					" AND (%s %s %s%s OR (%s = %s%s AND p.gid %s %s::uuid))",
					sortExpression,
					comparison,
					valuePlaceholder,
					sortCast,
					sortExpression,
					valuePlaceholder,
					sortCast,
					comparison,
					gidPlaceholder,
				),
			)
		}
	}

	direction := string(query.Sort.Direction)
	statement.WriteString(" ORDER BY " + sortExpression + " " + direction)
	if query.Sort.Field != "gid" {
		statement.WriteString(", p.gid " + direction)
	}

	limitPlaceholder := builder.add(query.Limit + 1)
	statement.WriteString(" LIMIT " + limitPlaceholder)

	rows, err := repository.pool.Query(ctx, statement.String(), builder.args...)
	if err != nil {
		return foundationrepository.PersonSearchResult{}, err
	}
	defer rows.Close()

	personRows := make([]personSearchRow, 0, query.Limit+1)
	for rows.Next() {
		var gid string
		var middleName *string
		var status string
		var row personSearchRow

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
			return foundationrepository.PersonSearchResult{}, err
		}

		row.gid = entity.GID(gid)
		row.status = entity.PersonStatus(status)
		if middleName != nil {
			row.middleName = *middleName
		}
		personRows = append(personRows, row)
	}
	if err := rows.Err(); err != nil {
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
		persons = append(
			persons,
			entity.NewPerson(
				row.gid,
				row.firstName,
				row.middleName,
				row.lastName,
				contactsByGID[row.gid],
				row.status,
				row.version,
				row.createdAt,
				row.updatedAt,
			),
		)
	}

	return foundationrepository.PersonSearchResult{
		Persons: persons,
		HasNext: hasNext,
	}, nil
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
		var gid string
		var contactType string
		var value string
		var normalizedValue string
		var isPrimary bool
		var isVerified bool
		if err := rows.Scan(
			&gid,
			&contactType,
			&value,
			&normalizedValue,
			&isPrimary,
			&isVerified,
		); err != nil {
			return nil, err
		}

		personGID := entity.GID(gid)
		contactsByGID[personGID] = append(
			contactsByGID[personGID],
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

	return contactsByGID, nil
}

func personFilterCondition(
	builder *searchQueryBuilder,
	filter foundationrepository.PersonSearchFilter,
) (string, error) {
	switch filter.Field {
	case "gid":
		return scalarFilterCondition(builder, "p.gid::text", "", filter)
	case "firstName":
		return scalarFilterCondition(builder, "LOWER(p.first_name)", "", filter)
	case "middleName":
		return scalarFilterCondition(builder, "LOWER(COALESCE(p.middle_name, ''))", "", filter)
	case "lastName":
		return scalarFilterCondition(builder, "LOWER(p.last_name)", "", filter)
	case "status":
		return scalarFilterCondition(builder, "p.status", "", filter)
	case "createdAt":
		return scalarFilterCondition(builder, "p.created_at", "::timestamptz", filter)
	case "updatedAt":
		return scalarFilterCondition(builder, "p.updated_at", "::timestamptz", filter)
	case "contacts.email":
		return contactFilterCondition(builder, entity.ContactTypeEmail, "", filter)
	case "contacts.phone":
		return contactFilterCondition(builder, entity.ContactTypePhone, "", filter)
	case "contacts.primaryEmail":
		return contactFilterCondition(builder, entity.ContactTypeEmail, " AND contact.is_primary = TRUE", filter)
	case "contacts.primaryPhone":
		return contactFilterCondition(builder, entity.ContactTypePhone, " AND contact.is_primary = TRUE", filter)
	case "contacts.verifiedEmail":
		return contactFilterCondition(builder, entity.ContactTypeEmail, " AND contact.is_verified = TRUE", filter)
	case "contacts.verifiedPhone":
		return contactFilterCondition(builder, entity.ContactTypePhone, " AND contact.is_verified = TRUE", filter)
	default:
		return "", foundationrepository.ErrInvalidPersonSearchQuery
	}
}

func scalarFilterCondition(
	builder *searchQueryBuilder,
	expression string,
	cast string,
	filter foundationrepository.PersonSearchFilter,
) (string, error) {
	switch filter.Operator {
	case searchcontract.OperatorEquals:
		placeholder := builder.add(filter.Values[0])
		return expression + " = " + placeholder + cast, nil
	case searchcontract.OperatorNotEquals:
		placeholder := builder.add(filter.Values[0])
		return expression + " <> " + placeholder + cast, nil
	case searchcontract.OperatorIn, searchcontract.OperatorNotIn:
		placeholders := make([]string, 0, len(filter.Values))
		for _, value := range filter.Values {
			placeholder := builder.add(value)
			placeholders = append(placeholders, placeholder+cast)
		}
		operator := "IN"
		if filter.Operator == searchcontract.OperatorNotIn {
			operator = "NOT IN"
		}
		return expression + " " + operator + " (" + strings.Join(placeholders, ", ") + ")", nil
	case searchcontract.OperatorExists:
		return "NULLIF(" + expression + ", '') IS NOT NULL", nil
	case searchcontract.OperatorNotExists:
		return "NULLIF(" + expression + ", '') IS NULL", nil
	case searchcontract.OperatorGreaterThan,
		searchcontract.OperatorGTE,
		searchcontract.OperatorLessThan,
		searchcontract.OperatorLTE:
		placeholder := builder.add(filter.Values[0])
		operator := mapComparisonOperator(filter.Operator)
		return expression + " " + operator + " " + placeholder + cast, nil
	case searchcontract.OperatorContains:
		placeholder := builder.add("%" + escapeLike(filter.Values[0]) + "%")
		return expression + " LIKE " + placeholder, nil
	case searchcontract.OperatorStartsWith:
		placeholder := builder.add(escapeLike(filter.Values[0]) + "%")
		return expression + " LIKE " + placeholder, nil
	default:
		return "", foundationrepository.ErrInvalidPersonSearchQuery
	}
}

func contactFilterCondition(
	builder *searchQueryBuilder,
	contactType entity.ContactType,
	extraCondition string,
	filter foundationrepository.PersonSearchFilter,
) (string, error) {
	typePlaceholder := builder.add(string(contactType))
	base := "contact.person_gid = p.gid AND contact.type = " + typePlaceholder + extraCondition

	switch filter.Operator {
	case searchcontract.OperatorExists:
		return "EXISTS (SELECT 1 FROM foundation.person_contacts AS contact WHERE " + base + ")", nil
	case searchcontract.OperatorNotExists:
		return "NOT EXISTS (SELECT 1 FROM foundation.person_contacts AS contact WHERE " + base + ")", nil
	case searchcontract.OperatorEquals, searchcontract.OperatorNotEquals:
		valuePlaceholder := builder.add(filter.Values[0])
		condition := base + " AND contact.normalized_value = " + valuePlaceholder
		prefix := "EXISTS"
		if filter.Operator == searchcontract.OperatorNotEquals {
			prefix = "NOT EXISTS"
		}
		return prefix + " (SELECT 1 FROM foundation.person_contacts AS contact WHERE " + condition + ")", nil
	case searchcontract.OperatorIn, searchcontract.OperatorNotIn:
		placeholders := make([]string, 0, len(filter.Values))
		for _, value := range filter.Values {
			placeholders = append(placeholders, builder.add(value))
		}
		condition := base + " AND contact.normalized_value IN (" + strings.Join(placeholders, ", ") + ")"
		prefix := "EXISTS"
		if filter.Operator == searchcontract.OperatorNotIn {
			prefix = "NOT EXISTS"
		}
		return prefix + " (SELECT 1 FROM foundation.person_contacts AS contact WHERE " + condition + ")", nil
	case searchcontract.OperatorContains:
		valuePlaceholder := builder.add("%" + escapeLike(filter.Values[0]) + "%")
		return "EXISTS (SELECT 1 FROM foundation.person_contacts AS contact WHERE " + base +
			" AND contact.normalized_value LIKE " + valuePlaceholder + ")", nil
	case searchcontract.OperatorStartsWith:
		valuePlaceholder := builder.add(escapeLike(filter.Values[0]) + "%")
		return "EXISTS (SELECT 1 FROM foundation.person_contacts AS contact WHERE " + base +
			" AND contact.normalized_value LIKE " + valuePlaceholder + ")", nil
	default:
		return "", foundationrepository.ErrInvalidPersonSearchQuery
	}
}

func personSortExpression(field string) (string, string, error) {
	switch field {
	case "gid":
		return "p.gid", "::uuid", nil
	case "firstName":
		return "LOWER(p.first_name)", "::text", nil
	case "lastName":
		return "LOWER(p.last_name)", "::text", nil
	case "createdAt":
		return "p.created_at", "::timestamptz", nil
	case "updatedAt":
		return "p.updated_at", "::timestamptz", nil
	default:
		return "", "", foundationrepository.ErrInvalidPersonSearchQuery
	}
}

func mapComparisonOperator(operator searchcontract.Operator) string {
	switch operator {
	case searchcontract.OperatorGreaterThan:
		return ">"
	case searchcontract.OperatorGTE:
		return ">="
	case searchcontract.OperatorLessThan:
		return "<"
	case searchcontract.OperatorLTE:
		return "<="
	default:
		return "="
	}
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	)
	return replacer.Replace(value)
}
