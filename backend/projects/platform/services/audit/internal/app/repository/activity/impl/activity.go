package impl

import (
	"context"
	coredatabase "grovehub/libs/core/database"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/audit/internal/app/entity"
	"grovehub/services/audit/internal/app/repository/activity"
	searchconfig "grovehub/services/audit/internal/app/repository/activity/helpers/search"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) activity.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Insert(ctx context.Context, entry entity.Activity) error {
	return coredatabase.RunTransaction(ctx, repository.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, repository.db)
		entry.IngestedAt = time.Now().UTC()
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entry)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			var matches []entity.Activity
			if err := db.Where("audit_id = ? OR source_event_id = ?", entry.AuditID, entry.SourceEventID).Limit(2).Find(&matches).Error; err != nil {
				return err
			}
			if len(matches) != 1 || matches[0].AuditID != entry.AuditID || matches[0].SemanticHash != entry.SemanticHash {
				return activity.ErrConflict
			}
		}
		// Reserve every observed source ID, including logically equivalent
		// events with a different envelope ID, in the same transaction.
		if err := db.Exec("INSERT INTO audit.event_receipts(source_event_id,audit_id) VALUES (?,?) ON CONFLICT (source_event_id) DO NOTHING", entry.SourceEventID, entry.AuditID).Error; err != nil {
			return err
		}
		var receipt struct{ AuditID string }
		if err := db.Table("audit.event_receipts").Select("audit_id").Where("source_event_id = ?", entry.SourceEventID).Take(&receipt).Error; err != nil {
			return err
		}
		if receipt.AuditID != entry.AuditID {
			return activity.ErrConflict
		}
		return nil
	})
}

// Tenant and self-visibility predicates must be attached before shared search.
// Do not expose an unscoped generated repository method.
func (repository *Repository) Search(ctx context.Context, tenantID string, selfAccountID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Activity], error) {
	if tenantID == "" {
		return coresearch.QueryResult[entity.Activity]{}, coresearch.ErrInvalidQuery
	}
	db := repository.db.WithContext(ctx).Table("audit.activities AS a").Where("a.tenant_id = ?", tenantID)
	if selfAccountID != "" {
		db = db.Where("((a.asset_type = 'ACCOUNT' AND a.asset_id = ?) OR (a.asset_type = 'INVITATION' AND a.subject_account_id = ?))", selfAccountID, selfAccountID)
	}
	return coresearch.Execute[string, entity.Activity](db, query, searchconfig.Config)
}
