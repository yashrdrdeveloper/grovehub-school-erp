// Package rbacmigration plans the legacy-to-normalized role migration. It does
// not write to the database, enable permissions or establish administrator access.
package rbacmigration

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/persona"
)

var (
	ErrScopeRequired = errors.New("migration requires tenant scope")
	ErrInvalidInventory = errors.New("invalid legacy role inventory")
	ErrConflict = errors.New("migration conflicts with existing normalized RBAC state")
)

type LegacyAccount struct {
	AccountID string
	Version int64
	Primary *string
	Secondary []string
}

type ExistingRole struct {
	RoleID string
	NormalizedName string
}

type Inventory struct {
	TenantID string
	Accounts []LegacyAccount
	Roles []ExistingRole
	NormalizedAssignments int64
}

type RolePlan struct {
	RoleID string `json:"roleId"`
	Name string `json:"name"`
	Persona persona.Code `json:"persona"`
}

type AssignmentPlan struct {
	AccountID string `json:"accountId"`
	Version int64 `json:"version"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}

type Plan struct {
	TenantID string `json:"tenantId"`
	Roles []RolePlan `json:"roles"`
	Assignments []AssignmentPlan `json:"assignments"`
	UnassignedAccountIDs []string `json:"unassignedAccountIds"`
}

type Reader interface {
	Inventory(ctx context.Context, tenantID string) (Inventory, error)
}
