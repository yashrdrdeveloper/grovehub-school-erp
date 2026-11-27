package approval

import (
	"context"
	"errors"
	"time"

	"grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
)

var (
	ErrNotFound = errors.New("approval request not found")
	ErrPolicy = errors.New("no valid active approval policy")
	ErrConflict = errors.New("approval request conflict")
	ErrInvalid = errors.New("invalid approval request")
)

type Step struct {
	Number int
	Permission permission.Code
	RoleID *string
}

type Policy struct {
	PolicyID string
	Steps []Step
}

type ReviewRecord struct {
	RequesterID string
	ReviewerID string
	Decision approval.Decision
	Comment string
	Step int
	At time.Time
}

// Repository is an internal persistence contract. Every mutation and its
// prerequisite reads must use the same tenant authority transaction.
type Repository interface {
	Now(ctx context.Context) (time.Time, error)
	ActivePolicy(ctx context.Context, tenantID string, kind approval.RequestType) (Policy, error)
	FindByKey(ctx context.Context, tenantID string, requesterID string, key string) (entity.ApprovalRequest, error)
	Lock(ctx context.Context, tenantID string, requestID string) (entity.ApprovalRequest, error)
	Step(ctx context.Context, tenantID string, requestID string, number int) (Step, error)
	Decisions(ctx context.Context, tenantID string, requestID string) ([]ReviewRecord, error)
	Insert(ctx context.Context, request entity.ApprovalRequest, steps []Step) error
	InsertDecision(ctx context.Context, tenantID string, requestID string, record ReviewRecord) error
	Advance(ctx context.Context, before entity.ApprovalRequest, next entity.ApprovalRequest) error
}
