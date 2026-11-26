package impl

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	dto "grovehub/services/audit/internal/app/dto/activity"
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

func (repository *Repository) Search(ctx context.Context, tenantID string, selfAccountID string, query dto.Search) ([]entity.Activity, error) {
	db := repository.db.WithContext(ctx).Model(&entity.Activity{}).Where("tenant_id = ?", tenantID)
	if selfAccountID != "" {
		db = db.Where("(actor_id = ? OR subject_account_id = ? OR (asset_type = 'ACCOUNT' AND asset_id = ?))", selfAccountID, selfAccountID, selfAccountID)
	}
	for _, filter := range []struct{ column string; value string }{
		{"subject_account_id",query.SubjectAccountID},{"asset_type",query.AssetType},{"asset_id",query.AssetID},{"actor_id",query.ActorID},{"person_id",query.PersonID},{"audit_type",query.AuditType},{"operation",query.Operation},{"source",query.Source},
	} {
		if filter.value != "" { db = db.Where(filter.column+" = ?", filter.value) }
	}
	if query.From != nil { db = db.Where("occurred_at >= ?", *query.From) }
	if query.To != nil { db = db.Where("occurred_at < ?", *query.To) }
	if query.Cursor != nil { db = db.Where("(occurred_at, audit_id) < (?, CAST(? AS uuid))", query.Cursor.OccurredAt, query.Cursor.AuditID) }
	rows := make([]entity.Activity,0)
	err := db.Order("occurred_at DESC, audit_id DESC").Limit(query.Limit+1).Find(&rows).Error
	return rows, err
}
