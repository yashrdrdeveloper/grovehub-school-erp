// Package rbacmigration owns the staged legacy-to-normalized migration contracts.
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

var ErrStalePlan = errors.New("reviewed migration plan is stale")
var ErrParity = errors.New("migration parity check failed")

type BackfillResult struct {
	RolesCreated int
	AssignmentsCreated int
	UnassignedAccounts int
}

// Backfill is an offline operator primitive, never a school API. It creates
// only inactive role records and inert normalized assignments.
type Backfill interface {
	Apply(ctx context.Context, reviewed Plan, operatorID string) (BackfillResult, error)
}
