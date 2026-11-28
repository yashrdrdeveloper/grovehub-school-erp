package impl

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/repository/approval/constants"
	searchconfig "grovehub/services/identity/internal/app/repository/approval/helpers/search"
)

type Reader struct {
	db *gorm.DB
}

var _ approval.Reader = (*Reader)(nil)

func NewReader(db *gorm.DB) approval.Reader {
	return &Reader{db: db}
}

func (reader *Reader) scoped(ctx context.Context, visibility approval.Visibility) (*gorm.DB, error) {
	if reader.db == nil || visibility.TenantID == "" || visibility.AccountID == "" || visibility.Now.IsZero() {
		return nil, approval.ErrInvalid
	}
	db := coredatabase.TransactionDB(ctx, reader.db).Table("identity.approval_requests AS r").Where("r.tenant_id = ?", visibility.TenantID)
	if !visibility.Inbox {
		return db.Where("r.requester_id = ?", visibility.AccountID), nil
	}
	grants, err := json.Marshal(append([]approval.EligibleGrant{}, visibility.Grants...))
	if err != nil {
		return nil, err
	}
	return db.Where(constants.Inbox, visibility.Now, visibility.AccountID, visibility.AccountID, string(grants)), nil
}

func (reader *Reader) Search(ctx context.Context, visibility approval.Visibility, query coresearch.Query[string]) (coresearch.QueryResult[dto.Summary], error) {
	db, err := reader.scoped(ctx, visibility)
	if err != nil {
		return coresearch.QueryResult[dto.Summary]{}, err
	}
	var total *int64
	if query.PageNumber == 0 {
		countQuery := query
		countQuery.Cursor = nil
		countQuery.Offset = 0
		countDB, err := coresearch.Apply(db.Session(&gorm.Session{}), countQuery, searchconfig.Config)
		if err != nil {
			return coresearch.QueryResult[dto.Summary]{}, err
		}
		var count int64
		if err := countDB.Select("COUNT(*)").Limit(-1).Offset(-1).Count(&count).Error; err != nil {
			return coresearch.QueryResult[dto.Summary]{}, err
		}
		total = &count
	}
	rows, err := coresearch.Execute[string, dto.Summary](db.Select(constants.SummaryProjection), query, searchconfig.Config)
	if err == nil && total != nil {
		rows.TotalRecords = total
	}
	return rows, err
}

func (reader *Reader) CountInbox(ctx context.Context, visibility approval.Visibility) (int64, error) {
	visibility.Inbox = true
	db, err := reader.scoped(ctx, visibility)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Count(&count).Error
	return count, err
}

func (reader *Reader) Get(ctx context.Context, visibility approval.Visibility, requestID string) (entity.ApprovalRequest, error) {
	db, err := reader.scoped(ctx, visibility)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	var request entity.ApprovalRequest
	result := db.Select("r.*").Where("r.request_id = ?", requestID).Scan(&request)
	if result.Error != nil {
		return entity.ApprovalRequest{}, result.Error
	}
	if result.RowsAffected != 1 {
		return entity.ApprovalRequest{}, approval.ErrNotFound
	}
	return request, nil
}
