package impl

import (
	"context"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	dto "grovehub/services/identity/internal/app/dto/audit"
	"grovehub/services/identity/internal/app/repository/auditread"
	"grovehub/services/identity/internal/app/repository/auditread/constants"
	"grovehub/services/identity/internal/app/repository/authority"
)

type Reader struct { db *gorm.DB }

func New(db *gorm.DB) auditread.Reader {
	return &Reader{db: db}
}

func (reader *Reader) Search(ctx context.Context, tenantID string, query dto.Search) ([]dto.Event, error) {
	if tenantID == "" { return nil, authority.ErrScopeRequired }
	statement := constants.Timeline
	args := []any{tenantID, tenantID, tenantID, tenantID, tenantID}
	for _, filter := range []struct{ clause string; value string }{
		{" AND source = ?", query.Source},
		{" AND resource_type = ?", query.ResourceType},
		{" AND resource_id = CAST(? AS uuid)", query.ResourceID},
		{" AND actor_id = CAST(? AS uuid)", query.ActorID},
		{" AND operation = ?", query.Operation},
	} {
		if filter.value != "" { statement += filter.clause; args = append(args, filter.value) }
	}
	if query.From != nil { statement += " AND occurred_at >= ?"; args = append(args, *query.From) }
	if query.To != nil { statement += " AND occurred_at < ?"; args = append(args, *query.To) }
	if query.Cursor != nil {
		statement += " AND (occurred_at, source, event_id) < (?, ?, CAST(? AS uuid))"
		args = append(args, query.Cursor.OccurredAt, query.Cursor.Source, query.Cursor.EventID)
	}
	statement += " ORDER BY occurred_at DESC, source DESC, event_id DESC LIMIT ?"
	args = append(args, query.Limit+1)
	rows := make([]dto.Event, 0)
	err := transaction.Database(ctx, reader.db).Raw(statement, args...).Scan(&rows).Error
	return rows, err
}
