package approval

import (
	"grovehub/services/identity/internal/app/constants/permission"
	accountrole "grovehub/services/identity/internal/app/dto/accountrole"
)

type PermissionIntent struct {
	RoleID string `json:"roleId"`
	Version int64 `json:"version"`
	Codes []permission.Code `json:"permissionCodes"`
}

type AssignmentIntent struct {
	AccountID string `json:"accountId"`
	Version int64 `json:"version"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
	RoleVersions []accountrole.RoleVersion `json:"roleVersions"`
}

type CreatePermissions struct {
	Intent PermissionIntent `json:"intent"`
	Reason string `json:"reason"`
	IdempotencyKey string `json:"-"`
}

type CreateAssignment struct {
	Intent AssignmentIntent `json:"intent"`
	Reason string `json:"reason"`
	IdempotencyKey string `json:"-"`
}
