package impl

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	searchconfig "grovehub/services/audit/internal/app/repository/activity/helpers/search"
	"grovehub/services/audit/internal/app/entity"
	"grovehub/services/audit/internal/app/repository/activity"
)

type Repository struct { db *gorm.DB }

func New(db *gorm.DB) activity.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Insert(ctx context.Context, entry entity.Activity) error {
	entry.IngestedAt = time.Now().UTC()
	db := repository.db.WithContext(ctx)
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entry)
	if result.Error != nil { return result.Error }
	if result.RowsAffected == 1 { return nil }
	var matches []entity.Activity
	if err := db.Where("audit_id = ? OR source_event_id = ?", entry.AuditID, entry.SourceEventID).Limit(2).Find(&matches).Error; err != nil { return err }
	if len(matches) != 1 || matches[0].AuditID != entry.AuditID || matches[0].SemanticHash != entry.SemanticHash { return activity.ErrConflict }
	return nil
}

// Tenant and self-visibility predicates must be attached before shared search.
// Do not expose an unscoped generated repository method.
func (repository *Repository) Search(ctx context.Context, tenantID string, selfAccountID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Activity], error) {
	if tenantID == "" { return coresearch.QueryResult[entity.Activity]{}, coresearch.ErrInvalidQuery }
	db := repository.db.WithContext(ctx).Table("audit.activities AS a").Where("a.tenant_id = ?", tenantID)
	if selfAccountID != "" {
		db = db.Where("((a.asset_type = 'ACCOUNT' AND a.asset_id = ?) OR (a.asset_type = 'INVITATION' AND a.subject_account_id = ?))", selfAccountID, selfAccountID)
	}
	return coresearch.Execute[string,entity.Activity](db,query,searchconfig.Config)
}
