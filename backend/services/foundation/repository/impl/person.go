package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	coretransaction "grovehub/libs/core/database/transaction"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

const uniquePersonIdempotencyKeyConstraint = "uq_persons_create_idempotency_key"

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
	return mapPersonWriteError(coretransaction.Database(ctx, repository.db).Create(&person).Error)
}

func (repository *PersonRepository) FindByGID(ctx context.Context, gid string) (entity.Person, error) {
	var person entity.Person
	if err := coretransaction.Database(ctx, repository.db).Where("gid = ?", gid).Take(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, foundationrepository.ErrPersonNotFound
		}
		return entity.Person{}, err
	}
	return person, nil
}

func (repository *PersonRepository) FindByIdempotencyKey(ctx context.Context, key string) (entity.Person, error) {
	var person entity.Person
	if err := coretransaction.Database(ctx, repository.db).Where(&entity.Person{IdempotencyKey: &key}).Take(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Person{}, foundationrepository.ErrPersonNotFound
		}
		return entity.Person{}, err
	}
	return person, nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	expectedVersion := person.Version - 1
	result := coretransaction.Database(ctx, repository.db).
		Model(&entity.Person{}).
		Where("gid = ? AND version = ?", person.GID, expectedVersion).
		Updates(map[string]any{
			"first_name": person.FirstName,
			"middle_name": person.MiddleName,
			"last_name": person.LastName,
			"status": person.Status,
			"version": person.Version,
			"search_details": person.SearchDetails,
			"updated_at": person.UpdatedAt,
		})
	if result.Error != nil {
		return mapPersonWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return foundationrepository.ErrPersonVersionConflict
	}
	return nil
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


func mapPersonWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" &&
		postgresError.ConstraintName == uniquePersonIdempotencyKeyConstraint {
		return foundationrepository.ErrPersonIdempotencyKeyAlreadyExists
	}
	return err
}
