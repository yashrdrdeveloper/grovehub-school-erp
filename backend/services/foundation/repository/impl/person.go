package impl

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	searchcontract "grovehub/libs/core/contracts/search"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

const uniqueContactConstraint = "uq_person_contacts_type_normalized_value"

type PersonRepository struct {
	database *gorm.DB
}

type personScalarSearchField struct {
	expression string
	cast       string
}

type personRelationSearchField struct {
	contactType  entity.ContactType
	primaryOnly  bool
	verifiedOnly bool
}

type personSortField struct {
	expression string
	cast       string
}

var personScalarSearchFields = map[string]personScalarSearchField{
	"gid":        {expression: "p.gid::text"},
	"firstName":  {expression: "LOWER(p.first_name)"},
	"middleName": {expression: "LOWER(COALESCE(p.middle_name, ''))"},
	"lastName":   {expression: "LOWER(p.last_name)"},
	"status":     {expression: "p.status"},
	"createdAt":  {expression: "p.created_at", cast: "::timestamptz"},
	"updatedAt":  {expression: "p.updated_at", cast: "::timestamptz"},
}

var personRelationSearchFields = map[string]personRelationSearchField{
	"contacts.email":         {contactType: entity.ContactTypeEmail},
	"contacts.phone":         {contactType: entity.ContactTypePhone},
	"contacts.primaryEmail":  {contactType: entity.ContactTypeEmail, primaryOnly: true},
	"contacts.primaryPhone":  {contactType: entity.ContactTypePhone, primaryOnly: true},
	"contacts.verifiedEmail": {contactType: entity.ContactTypeEmail, verifiedOnly: true},
	"contacts.verifiedPhone": {contactType: entity.ContactTypePhone, verifiedOnly: true},
}

var personSortFields = map[string]personSortField{
	"gid":       {expression: "p.gid", cast: "::uuid"},
	"firstName": {expression: "LOWER(p.first_name)", cast: "::text"},
	"lastName":  {expression: "LOWER(p.last_name)", cast: "::text"},
	"createdAt": {expression: "p.created_at", cast: "::timestamptz"},
	"updatedAt": {expression: "p.updated_at", cast: "::timestamptz"},
}

var _ foundationrepository.PersonRepository = (*PersonRepository)(nil)

func NewPersonRepository(database *gorm.DB) foundationrepository.PersonRepository {
	return &PersonRepository{database: database}
}

func (repository *PersonRepository) Create(ctx context.Context, person entity.Person) error {
	return repository.database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := toPersonModel(person)
		if err := tx.Create(&model).Error; err != nil {
			return mapPersonWriteError(err)
		}

		if err := upsertPersonContacts(tx, person, false); err != nil {
			return err
		}
		return refreshPersonSearchDetails(tx, person.GetGID())
	})
}

func (repository *PersonRepository) FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	var model personModel
	if err := repository.database.WithContext(ctx).
		Where("gid = ?", string(gid)).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, foundationrepository.ErrPersonNotFound
		}
		return entity.Person{}, err
	}

	contactsByGID, err := repository.findContactModelsByGIDs(ctx, []entity.GID{gid})
	if err != nil {
		return entity.Person{}, err
	}
	return toPersonEntity(model, contactsByGID[gid]), nil
}

func (repository *PersonRepository) FindGIDsByContacts(
	ctx context.Context,
	contacts []entity.Contact,
) ([]entity.GID, error) {
	if len(contacts) == 0 {
		return []entity.GID{}, nil
	}

	pairs := make([][]any, 0, len(contacts))
	for _, contact := range contacts {
		pairs = append(pairs, []any{string(contact.GetType()), contact.GetNormalizedValue()})
	}

	var gidValues []string
	if err := repository.database.WithContext(ctx).
		Model(&personContactModel{}).
		Distinct("person_gid").
		Where("(type, normalized_value) IN ?", pairs).
		Order("person_gid ASC").
		Pluck("person_gid", &gidValues).Error; err != nil {
		return nil, err
	}

	gids := make([]entity.GID, 0, len(gidValues))
	for _, gid := range gidValues {
		gids = append(gids, entity.GID(gid))
	}
	return gids, nil
}

func (repository *PersonRepository) Search(
	ctx context.Context,
	query coresearch.Query[entity.GID],
) (coresearch.Result[entity.Person], error) {
	database := repository.database.WithContext(ctx).
		Table(personTable + " AS p").
		Select(`p.gid, p.first_name, p.middle_name, p.last_name,
			p.status, p.version, p.created_at, p.updated_at`)

	for _, term := range query.SearchTerms {
		database = database.Where("p.search_details ILIKE ?", "%"+coresearch.EscapeLike(term)+"%")
	}

	var err error
	for _, filter := range query.Filters {
		database, err = repository.applyPersonFilter(database, filter)
		if err != nil {
			return coresearch.Result[entity.Person]{}, err
		}
	}

	sortField, ok := personSortFields[query.Sort.Field]
	if !ok || !validSortDirection(query.Sort.Direction) {
		return coresearch.Result[entity.Person]{}, coresearch.ErrInvalidQuery
	}

	if query.Cursor != nil {
		database = applyPersonCursor(database, query, sortField)
	}

	direction := string(query.Sort.Direction)
	database = database.Order(sortField.expression + " " + direction)
	if query.Sort.Field != "gid" {
		database = database.Order("p.gid " + direction)
	}

	var models []personModel
	if err := database.Limit(query.Limit + 1).Scan(&models).Error; err != nil {
		return coresearch.Result[entity.Person]{}, err
	}

	hasNext := len(models) > query.Limit
	if hasNext {
		models = models[:query.Limit]
	}

	gids := make([]entity.GID, 0, len(models))
	for _, model := range models {
		gids = append(gids, entity.GID(model.GID))
	}
	contactsByGID, err := repository.findContactModelsByGIDs(ctx, gids)
	if err != nil {
		return coresearch.Result[entity.Person]{}, err
	}

	persons := make([]entity.Person, 0, len(models))
	for _, model := range models {
		gid := entity.GID(model.GID)
		persons = append(persons, toPersonEntity(model, contactsByGID[gid]))
	}
	return coresearch.Result[entity.Person]{Items: persons, HasNext: hasNext}, nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	return repository.database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := toPersonModel(person)
		result := tx.Model(&personModel{}).
			Where("gid = ?", model.GID).
			Updates(map[string]any{
				"first_name":  model.FirstName,
				"middle_name": model.MiddleName,
				"last_name":   model.LastName,
				"status":      model.Status,
				"version":     model.Version,
				"updated_at":  model.UpdatedAt,
			})
		if result.Error != nil {
			return mapPersonWriteError(result.Error)
		}
		if result.RowsAffected == 0 {
			return foundationrepository.ErrPersonNotFound
		}

		if err := upsertPersonContacts(tx, person, true); err != nil {
			return err
		}
		return refreshPersonSearchDetails(tx, person.GetGID())
	})
}

func (repository *PersonRepository) Delete(ctx context.Context, gid entity.GID) error {
	result := repository.database.WithContext(ctx).
		Where("gid = ?", string(gid)).
		Delete(&personModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return foundationrepository.ErrPersonNotFound
	}
	return nil
}

func (repository *PersonRepository) findContactModelsByGIDs(
	ctx context.Context,
	gids []entity.GID,
) (map[entity.GID][]personContactModel, error) {
	contactsByGID := make(map[entity.GID][]personContactModel, len(gids))
	if len(gids) == 0 {
		return contactsByGID, nil
	}

	gidValues := make([]string, 0, len(gids))
	for _, gid := range gids {
		gidValues = append(gidValues, string(gid))
	}

	var contacts []personContactModel
	if err := repository.database.WithContext(ctx).
		Where("person_gid IN ?", gidValues).
		Order("person_gid ASC").
		Order("type ASC").
		Order("is_primary DESC").
		Order("normalized_value ASC").
		Find(&contacts).Error; err != nil {
		return nil, err
	}

	for _, contact := range contacts {
		gid := entity.GID(contact.PersonGID)
		contactsByGID[gid] = append(contactsByGID[gid], contact)
	}
	return contactsByGID, nil
}

func (repository *PersonRepository) applyPersonFilter(
	database *gorm.DB,
	filter coresearch.Filter,
) (*gorm.DB, error) {
	if field, ok := personScalarSearchFields[filter.Field]; ok {
		return applyScalarFilter(database, field, filter)
	}
	if field, ok := personRelationSearchFields[filter.Field]; ok {
		return repository.applyRelationFilter(database, field, filter)
	}
	return nil, coresearch.ErrInvalidQuery
}

func applyScalarFilter(
	database *gorm.DB,
	field personScalarSearchField,
	filter coresearch.Filter,
) (*gorm.DB, error) {
	switch filter.Operator {
	case searchcontract.OperatorEquals:
		return scalarBinary(database, field, "=", filter.Values)
	case searchcontract.OperatorNotEquals:
		return scalarBinary(database, field, "<>", filter.Values)
	case searchcontract.OperatorGreaterThan:
		return scalarBinary(database, field, ">", filter.Values)
	case searchcontract.OperatorGTE:
		return scalarBinary(database, field, ">=", filter.Values)
	case searchcontract.OperatorLessThan:
		return scalarBinary(database, field, "<", filter.Values)
	case searchcontract.OperatorLTE:
		return scalarBinary(database, field, "<=", filter.Values)
	case searchcontract.OperatorIn, searchcontract.OperatorNotIn:
		if len(filter.Values) == 0 {
			return nil, coresearch.ErrInvalidQuery
		}
		operator := "IN"
		if filter.Operator == searchcontract.OperatorNotIn {
			operator = "NOT IN"
		}
		placeholders := make([]string, len(filter.Values))
		args := make([]any, len(filter.Values))
		for index, value := range filter.Values {
			placeholders[index] = "?" + field.cast
			args[index] = value
		}
		condition := field.expression + " " + operator + " (" + strings.Join(placeholders, ", ") + ")"
		return database.Where(condition, args...), nil
	case searchcontract.OperatorExists, searchcontract.OperatorNotExists:
		if len(filter.Values) != 0 {
			return nil, coresearch.ErrInvalidQuery
		}
		predicate := "IS NOT NULL"
		if filter.Operator == searchcontract.OperatorNotExists {
			predicate = "IS NULL"
		}
		return database.Where("NULLIF(" + field.expression + ", '') " + predicate), nil
	case searchcontract.OperatorContains, searchcontract.OperatorStartsWith:
		if len(filter.Values) != 1 {
			return nil, coresearch.ErrInvalidQuery
		}
		value := coresearch.EscapeLike(filter.Values[0])
		if filter.Operator == searchcontract.OperatorContains {
			value = "%" + value + "%"
		} else {
			value += "%"
		}
		return database.Where(field.expression+" LIKE ?", value), nil
	default:
		return nil, coresearch.ErrInvalidQuery
	}
}

func scalarBinary(
	database *gorm.DB,
	field personScalarSearchField,
	operator string,
	values []string,
) (*gorm.DB, error) {
	if len(values) != 1 {
		return nil, coresearch.ErrInvalidQuery
	}
	return database.Where(field.expression+" "+operator+" ?"+field.cast, values[0]), nil
}

func (repository *PersonRepository) applyRelationFilter(
	database *gorm.DB,
	field personRelationSearchField,
	filter coresearch.Filter,
) (*gorm.DB, error) {
	subquery := repository.database.Session(&gorm.Session{NewDB: true}).
		Table(personContactTable + " AS contact").
		Select("1").
		Where("contact.person_gid = p.gid").
		Where("contact.type = ?", string(field.contactType))
	if field.primaryOnly {
		subquery = subquery.Where("contact.is_primary = ?", true)
	}
	if field.verifiedOnly {
		subquery = subquery.Where("contact.is_verified = ?", true)
	}

	negated := false
	switch filter.Operator {
	case searchcontract.OperatorEquals, searchcontract.OperatorNotEquals:
		if len(filter.Values) != 1 {
			return nil, coresearch.ErrInvalidQuery
		}
		subquery = subquery.Where("contact.normalized_value = ?", filter.Values[0])
		negated = filter.Operator == searchcontract.OperatorNotEquals
	case searchcontract.OperatorIn, searchcontract.OperatorNotIn:
		if len(filter.Values) == 0 {
			return nil, coresearch.ErrInvalidQuery
		}
		subquery = subquery.Where("contact.normalized_value IN ?", filter.Values)
		negated = filter.Operator == searchcontract.OperatorNotIn
	case searchcontract.OperatorExists, searchcontract.OperatorNotExists:
		if len(filter.Values) != 0 {
			return nil, coresearch.ErrInvalidQuery
		}
		negated = filter.Operator == searchcontract.OperatorNotExists
	case searchcontract.OperatorContains, searchcontract.OperatorStartsWith:
		if len(filter.Values) != 1 {
			return nil, coresearch.ErrInvalidQuery
		}
		value := coresearch.EscapeLike(filter.Values[0])
		if filter.Operator == searchcontract.OperatorContains {
			value = "%" + value + "%"
		} else {
			value += "%"
		}
		subquery = subquery.Where("contact.normalized_value LIKE ?", value)
	default:
		return nil, coresearch.ErrInvalidQuery
	}

	predicate := "EXISTS (?)"
	if negated {
		predicate = "NOT EXISTS (?)"
	}
	return database.Where(predicate, subquery), nil
}

func applyPersonCursor(
	database *gorm.DB,
	query coresearch.Query[entity.GID],
	sortField personSortField,
) *gorm.DB {
	comparison := ">"
	if query.Sort.Direction == searchcontract.DirectionDescending {
		comparison = "<"
	}

	if query.Sort.Field == "gid" {
		return database.Where("p.gid "+comparison+" ?::uuid", query.Cursor.Value)
	}

	condition := fmt.Sprintf(
		"(%s %s ?%s OR (%s = ?%s AND p.gid %s ?::uuid))",
		sortField.expression,
		comparison,
		sortField.cast,
		sortField.expression,
		sortField.cast,
		comparison,
	)
	return database.Where(condition, query.Cursor.Value, query.Cursor.Value, string(query.Cursor.CursorId))
}

func validSortDirection(direction searchcontract.Direction) bool {
	return direction == searchcontract.DirectionAscending || direction == searchcontract.DirectionDescending
}

func upsertPersonContacts(tx *gorm.DB, person entity.Person, clearPrimary bool) error {
	gid := string(person.GetGID())
	if clearPrimary {
		if err := tx.Model(&personContactModel{}).
			Where("person_gid = ?", gid).
			Update("is_primary", false).Error; err != nil {
			return err
		}
	}

	contacts := toPersonContactModels(person)
	if len(contacts) > 0 {
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "person_gid"},
				{Name: "type"},
				{Name: "normalized_value"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"value",
				"is_primary",
				"is_verified",
				"updated_at",
			}),
		}).Create(&contacts).Error; err != nil {
			return mapPersonWriteError(err)
		}
	}

	return deleteStalePersonContacts(tx, gid, contacts)
}

func deleteStalePersonContacts(tx *gorm.DB, gid string, contacts []personContactModel) error {
	query := tx.Where("person_gid = ?", gid)
	if len(contacts) == 0 {
		return query.Delete(&personContactModel{}).Error
	}

	pairs := make([][]any, 0, len(contacts))
	for _, contact := range contacts {
		pairs = append(pairs, []any{contact.Type, contact.NormalizedValue})
	}
	return query.
		Where("(type, normalized_value) NOT IN ?", pairs).
		Delete(&personContactModel{}).Error
}

func refreshPersonSearchDetails(tx *gorm.DB, gid entity.GID) error {
	return tx.Exec(
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
		 WHERE person.gid = ?`,
		string(gid),
	).Error
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
