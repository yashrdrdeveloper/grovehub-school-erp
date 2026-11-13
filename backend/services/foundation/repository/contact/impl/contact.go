package impl

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	coretransaction "grovehub/libs/core/database/transaction"
	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
	repositoryerrorshelpers "grovehub/services/foundation/repository/helpers/errors"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) foundationrepository.ContactRepository {
	return &ContactRepository{db: db}
}

func (repository *ContactRepository) FindByPersonGID(ctx context.Context, gid string) ([]entity.Contact, error) {
	contactsByGID, err := repository.FindByPersonGIDs(ctx, []string{gid})
	return contactsByGID[gid], err
}

func (repository *ContactRepository) FindByPersonGIDs(ctx context.Context, gids []string) (map[string][]entity.Contact, error) {
	contactsByGID := make(map[string][]entity.Contact, len(gids))
	if len(gids) == 0 {
		return contactsByGID, nil
	}
	contacts := make([]entity.Contact, 0)
	if err := coretransaction.Database(ctx, repository.db).
		Where("person_gid IN ?", gids).
		Order("person_gid, type, is_primary DESC, normalized_value").
		Find(&contacts).Error; err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		contactsByGID[contact.PersonGID] = append(contactsByGID[contact.PersonGID], contact)
	}
	return contactsByGID, nil
}

func (repository *ContactRepository) FindPersonGIDsByContacts(ctx context.Context, contacts []entity.Contact) ([]string, error) {
	if len(contacts) == 0 {
		return []string{}, nil
	}
	conditions := make([]string, 0, len(contacts))
	args := make([]any, 0, len(contacts)*2)
	for _, contact := range contacts {
		conditions = append(conditions, "(type = ? AND normalized_value = ?)")
		args = append(args, contact.Type, contact.NormalizedValue)
	}
	gids := make([]string, 0)
	err := coretransaction.Database(ctx, repository.db).
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

// ReplaceForPerson replaces the Person-owned child collection atomically.
// It joins the outer PersonService transaction or creates its own transaction
// when called without one, so primary-reset, upsert, and deletion never partially commit.
func (repository *ContactRepository) ReplaceForPerson(ctx context.Context, gid string, contacts []entity.Contact) error {
	return coretransaction.Run(ctx, repository.db, func(txCtx context.Context) error {
		db := coretransaction.Database(txCtx, repository.db)
		if err := db.Model(&entity.Contact{}).
			Where("person_gid = ?", gid).
			Update("is_primary", false).Error; err != nil {
			return err
		}
		if len(contacts) > 0 {
			if err := db.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "person_gid"}, {Name: "type"}, {Name: "normalized_value"},
				},
				DoUpdates: clause.AssignmentColumns([]string{"value", "is_primary", "is_verified", "updated_at"}),
			}).Create(&contacts).Error; err != nil {
				return repositoryerrorshelpers.MapContactWriteError(err)
			}
		}
		return deleteStalePersonContacts(db, gid, contacts)
	})
}

func deleteStalePersonContacts(tx *gorm.DB, gid string, contacts []entity.Contact) error {
	query := tx.Where("person_gid = ?", gid)
	if len(contacts) == 0 {
		return query.Delete(&entity.Contact{}).Error
	}
	conditions := make([]string, 0, len(contacts))
	args := make([]any, 0, len(contacts)*2)
	for _, contact := range contacts {
		conditions = append(conditions, "(type = ? AND normalized_value = ?)")
		args = append(args, contact.Type, contact.NormalizedValue)
	}
	return query.Where("NOT ("+strings.Join(conditions, " OR ")+")", args...).Delete(&entity.Contact{}).Error
}
