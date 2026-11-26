package impl

import (
	"context"
	"time"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/repository/approval/constants"
)

type Repository struct {
	db *gorm.DB
}

var _ approval.Repository = (*Repository)(nil)

func New(db *gorm.DB) approval.Repository {
	return &Repository{db: db}
}

func (repository *Repository) dbFor(ctx context.Context) *gorm.DB {
	return transaction.Database(ctx, repository.db)
}

func (repository *Repository) Now(ctx context.Context) (time.Time, error) {
	var now time.Time
	err := repository.dbFor(ctx).Raw(constants.Now).Scan(&now).Error
	return now.UTC(), err
}

func (repository *Repository) ActivePolicy(ctx context.Context, tenantID string, requestType kind.RequestType) (approval.Policy, error) {
	var policy approval.Policy
	db := repository.dbFor(ctx)
	result := db.Raw(constants.ActivePolicy, tenantID, requestType).Scan(&policy.PolicyID)
	if result.Error != nil {
		return approval.Policy{}, result.Error
	}
	if result.RowsAffected != 1 {
		return approval.Policy{}, approval.ErrPolicy
	}
	if err := db.Raw(constants.PolicySteps, tenantID, policy.PolicyID).Scan(&policy.Steps).Error; err != nil {
		return approval.Policy{}, err
	}
	return policy, nil
}

func (repository *Repository) FindByKey(ctx context.Context, tenantID string, requesterID string, key string) (entity.ApprovalRequest, error) {
	var request entity.ApprovalRequest
	result := repository.dbFor(ctx).Raw(constants.FindByKey, tenantID, requesterID, key).Scan(&request)
	if result.Error != nil {
		return entity.ApprovalRequest{}, result.Error
	}
	if result.RowsAffected != 1 {
		return entity.ApprovalRequest{}, approval.ErrNotFound
	}
	return request, nil
}

func (repository *Repository) Lock(ctx context.Context, tenantID string, requestID string) (entity.ApprovalRequest, error) {
	var request entity.ApprovalRequest
	result := repository.dbFor(ctx).Raw(constants.LockRequest, tenantID, requestID).Scan(&request)
	if result.Error != nil {
		return entity.ApprovalRequest{}, result.Error
	}
	if result.RowsAffected != 1 {
		return entity.ApprovalRequest{}, approval.ErrNotFound
	}
	return request, nil
}

func (repository *Repository) Step(ctx context.Context, tenantID string, requestID string, number int) (approval.Step, error) {
	var step approval.Step
	result := repository.dbFor(ctx).Raw(constants.RequestStep, tenantID, requestID, number).Scan(&step)
	if result.Error != nil {
		return approval.Step{}, result.Error
	}
	if result.RowsAffected != 1 {
		return approval.Step{}, approval.ErrPolicy
	}
	return step, nil
}

func (repository *Repository) Insert(ctx context.Context, request entity.ApprovalRequest, steps []approval.Step) error {
	db := repository.dbFor(ctx)
	if err := db.Exec(constants.InsertRequest, request.RequestID, request.TenantID, request.PolicyID, request.RequesterID, request.RequestType, string(request.Payload), request.PayloadSHA256, request.Reason, request.IdempotencyKey, request.Status, request.CurrentStep, request.StepCount, request.Version, request.ExpiresAt, request.CreatedAt, request.UpdatedAt).Error; err != nil {
		return err
	}
	for _, step := range steps {
		if err := db.Exec(constants.InsertStep, request.TenantID, request.RequestID, step.Number, step.Permission, step.RoleID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repository *Repository) InsertDecision(ctx context.Context, tenantID string, requestID string, record approval.ReviewRecord) error {
	return repository.dbFor(ctx).Exec(constants.InsertDecision, tenantID, requestID, record.RequesterID, record.Step, record.ReviewerID, record.Decision, record.Comment, record.At).Error
}

func (repository *Repository) Advance(ctx context.Context, before entity.ApprovalRequest, next entity.ApprovalRequest) error {
	result := repository.dbFor(ctx).Exec(constants.Advance, next.Status, next.CurrentStep, next.Version, next.UpdatedAt, before.TenantID, before.RequestID, before.Version, before.Status, before.CurrentStep)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return approval.ErrConflict
	}
	return nil
}

func (repository *Repository) Decisions(ctx context.Context, tenantID string, requestID string) ([]approval.ReviewRecord, error) {
	var records []approval.ReviewRecord
	err := repository.dbFor(ctx).Raw(constants.Decisions, tenantID, requestID).Scan(&records).Error
	return records, err
}
