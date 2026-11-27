package accountrole

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
)

var (
	ErrScopeRequired = errors.New("account role lookup requires tenant scope")
	ErrNotFound = errors.New("account or role not found")
	ErrInvalidAssignment = errors.New("invalid account role assignment")
	ErrIneligible = errors.New("account or role is ineligible for assignment")
	ErrVersionConflict = errors.New("account role assignment version conflict")
)

type Assignment struct {
	AccountID string
	Status entity.AccountStatus
	Version int64
	PrimaryRoleID string
	SecondaryRoleIDs []string
}

type Candidate struct {
	RoleID string
	Status entity.RoleStatus
	Version int64
	Codes []permission.Code
}

type Repository interface {
	Find(ctx context.Context, tenantID string, accountID string) (Assignment, error)
	Candidates(ctx context.Context, tenantID string, roleIDs []string) ([]Candidate, error)
}
