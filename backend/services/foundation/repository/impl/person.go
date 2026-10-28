package impl

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

const (
	uniqueContactConstraint              = "uq_person_contacts_type_normalized_value"
	uniquePersonIdempotencyKeyConstraint = "uq_persons_create_idempotency_key"
)

type PersonRepository struct {
	db *gorm.DB
}

var personSearchConfig = coresearch.Config{
	SearchExpression: "p.search_details",
	CursorExpression: "p.gid",
	ScalarFields: map[string]coresearch.ScalarField{
		"gid":        {Expression: "p.gid::text"},
		"firstName":  {Expression: "LOWER(p.first_name)"},
		"middleName": {Expression: "LOWER(COALESCE(p.middle_name, ''))"},
		"lastName":   {Expression: "LOWER(p.last_name)"},
		"status":     {Expression: "p.status"},
		"createdAt":  {Expression: "p.created_at", Cast: "::timestamptz"},
		"updatedAt":  {Expression: "p.updated_at", Cast: "::timestamptz"},
	},
	RelationFields: map[string]coresearch.RelationField{
		"contacts.email":         contactSearchField(entity.ContactTypeEmail),
		"contacts.phone":         contactSearchField(entity.ContactTypePhone),
		"contacts.primaryEmail":  contactSearchField(entity.ContactTypeEmail, coresearch.FixedCondition{Expression: "contact.is_primary = ", Value: true}),
		"contacts.primaryPhone":  contactSearchField(entity.ContactTypePhone, coresearch.FixedCondition{Expression: "contact.is_primary = ", Value: true}),
		"contacts.verifiedEmail": contactSearchField(entity.ContactTypeEmail, coresearch.FixedCondition{Expression: "contact.is_verified = ", Value: true}),
		"contacts.verifiedPhone": contactSearchField(entity.ContactTypePhone, coresearch.FixedCondition{Expression: "contact.is_verified = ", Value: true}),
	},
	SortFields: map[string]coresearch.SortField{
		"gid":       {Expression: "p.gid", Cast: "::uuid"},
		"firstName": {Expression: "LOWER(p.first_name)", Cast: "::text"},
		"lastName":  {Expression: "LOWER(p.last_name)", Cast: "::text"},
		"createdAt": {Expression: "p.created_at", Cast: "::timestamptz"},
		"updatedAt": {Expression: "p.updated_at", Cast: "::timestamptz"},
	},
}

func NewPersonRepository(db *gorm.DB) foundationrepository.PersonRepository {
	return &PersonRepository{db: db}
}

func (repository *PersonRepository) Create(ctx context.Context, person entity.Person) error {
	person.SearchDetails = buildPersonSearchDetails(person)
	return repository.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&person).Error; err != nil {
			return mapPersonWriteError(err)
		}
		return upsertPersonContacts(tx, person, false)
	})
}

func (repository *PersonRepository) FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	var person entity.Person
	if err := repository.db.WithContext(ctx).Where("gid = ?", gid).Take(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, foundationrepository.ErrPersonNotFound
		}
		return entity.Person{}, err
	}
	contacts, err := repository.findContactsByGID(ctx, gid)
	if err != nil {
		return entity.Person{}, err
	}
	return person.SetContacts(contacts), nil
}

func (repository *PersonRepository) FindByCreateIdempotencyKey(ctx context.Context, key string) (entity.Person, error) {
	var person entity.Person
	if err := repository.db.WithContext(ctx).Where(&entity.Person{CreateIdempotencyKey: &key}).Take(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, foundationrepository.ErrPersonNotFound
		}
		return entity.Person{}, err
	}
	contacts, err := repository.findContactsByGID(ctx, person.GetGID())
	if err != nil {
		return entity.Person{}, err
	}
	return person.SetContacts(contacts), nil
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
		conditions = append(conditions, "(type = ? AND normalized_value = ?)")
		args = append(args, contact.GetType(), contact.GetNormalizedValue())
	}
	gids := make([]entity.GID, 0)
	err := repository.db.WithContext(ctx).
		Model(&entity.Contact{}).
		Distinct("person_gid").
		Where(strings.Join(conditions, " OR "), args...).
		Order("person_gid").
		Pluck("person_gid", &gids).Error
	if err != nil {
		return nil, err
	}
	return gids, nil
}

func (repository *PersonRepository) Search(
	ctx context.Context,
	query coresearch.Query[entity.GID],
) (coresearch.QueryResult[entity.Person], error) {
	searchDB, err := coresearch.Apply(
		repository.db.WithContext(ctx).
			Model(&entity.Person{}).
			Table(entity.Person{}.TableName()+" AS p").
			Select("p.*"),
		query,
		personSearchConfig,
	)
	if err != nil {
		return coresearch.QueryResult[entity.Person]{}, err
	}
	persons := make([]entity.Person, 0, query.Limit+1)
	if err := searchDB.Find(&persons).Error; err != nil {
		return coresearch.QueryResult[entity.Person]{}, err
	}
	hasNext := len(persons) > query.Limit
	if hasNext {
		persons = persons[:query.Limit]
	}
	gids := make([]entity.GID, 0, len(persons))
	for _, person := range persons {
		gids = append(gids, person.GetGID())
	}
	contactsByGID, err := repository.findContactsByGIDs(ctx, gids)
	if err != nil {
		return coresearch.QueryResult[entity.Person]{}, err
	}
	for index := range persons {
		persons[index] = persons[index].SetContacts(contactsByGID[persons[index].GetGID()])
	}
	return coresearch.QueryResult[entity.Person]{Items: persons, HasNext: hasNext}, nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	searchDetails := buildPersonSearchDetails(person)
	expectedVersion := person.GetVersion() - 1
	return repository.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entity.Person{}).
			Where("gid = ? AND version = ?", person.GetGID(), expectedVersion).
			Updates(map[string]any{
				"first_name":     person.FirstName,
				"middle_name":    person.MiddleName,
				"last_name":      person.LastName,
				"status":         person.Status,
				"version":        person.Version,
				"search_details": searchDetails,
				"updated_at":     person.UpdatedAt,
			})
		if result.Error != nil {
			return mapPersonWriteError(result.Error)
		}
		if result.RowsAffected == 0 {
			return foundationrepository.ErrPersonVersionConflict
		}
		return upsertPersonContacts(tx, person, true)
	})
}

func (repository *PersonRepository) Delete(ctx context.Context, gid entity.GID) error {
	result := repository.db.WithContext(ctx).Where("gid = ?", gid).Delete(&entity.Person{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return foundationrepository.ErrPersonNotFound
	}
	return nil
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
	contacts := make([]entity.Contact, 0)
	if err := repository.db.WithContext(ctx).
		Where("person_gid IN ?", gids).
		Order("person_gid, type, is_primary DESC, normalized_value").
		Find(&contacts).Error; err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		gid := contact.GetPersonGID()
		contactsByGID[gid] = append(contactsByGID[gid], contact)
	}
	return contactsByGID, nil
}

func contactSearchField(contactType entity.ContactType, conditions ...coresearch.FixedCondition) coresearch.RelationField {
	return coresearch.RelationField{
		Table:           entity.Contact{}.TableName(),
		Alias:           "contact",
		ParentCondition: "contact.person_gid = p.gid",
		ValueExpression: "contact.normalized_value",
		FixedConditions: append(
			[]coresearch.FixedCondition{{Expression: "contact.type = ", Value: contactType}},
			conditions...,
		),
	}
}

func upsertPersonContacts(tx *gorm.DB, person entity.Person, clearPrimary bool) error {
	gid := person.GetGID()
	if clearPrimary {
		if err := tx.Model(&entity.Contact{}).Where("person_gid = ?", gid).Update("is_primary", false).Error; err != nil {
			return err
		}
	}
	contacts := person.BuildContacts()
	if len(contacts) > 0 {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "person_gid"}, {Name: "type"}, {Name: "normalized_value"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "is_primary", "is_verified", "updated_at"}),
		}).Create(&contacts).Error; err != nil {
			return mapPersonWriteError(err)
		}
	}
	return deleteStalePersonContacts(tx, gid, contacts)
}

func deleteStalePersonContacts(tx *gorm.DB, gid entity.GID, contacts []entity.Contact) error {
	query := tx.Where("person_gid = ?", gid)
	if len(contacts) == 0 {
		return query.Delete(&entity.Contact{}).Error
	}
	conditions := make([]string, 0, len(contacts))
	args := make([]any, 0, len(contacts)*2)
	for _, contact := range contacts {
		conditions = append(conditions, "(type = ? AND normalized_value = ?)")
		args = append(args, contact.GetType(), contact.GetNormalizedValue())
	}
	return query.Where("NOT ("+strings.Join(conditions, " OR ")+")", args...).Delete(&entity.Contact{}).Error
}

func buildPersonSearchDetails(person entity.Person) string {
	contacts := person.GetContacts()
	sort.Slice(contacts, func(leftIndex int, rightIndex int) bool {
		left := contacts[leftIndex]
		right := contacts[rightIndex]
		if left.GetType() == right.GetType() {
			return left.GetNormalizedValue() < right.GetNormalizedValue()
		}
		return left.GetType() < right.GetType()
	})
	parts := make([]string, 0, 4+(len(contacts)*2))
	parts = append(parts, string(person.GetGID()), person.GetFirstName())
	if middleName := person.GetMiddleName(); middleName != "" {
		parts = append(parts, middleName)
	}
	parts = append(parts, person.GetLastName())
	for _, contact := range contacts {
		parts = append(parts, contact.GetValue(), contact.GetNormalizedValue())
	}
	return strings.ToLower(strings.TrimSpace(strings.Join(parts, " ")))
}

func mapPersonWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		switch postgresError.ConstraintName {
		case uniqueContactConstraint:
			return foundationrepository.ErrPersonContactAlreadyExists
		case uniquePersonIdempotencyKeyConstraint:
			return foundationrepository.ErrPersonIdempotencyKeyAlreadyExists
		}
	}
	return err
}
