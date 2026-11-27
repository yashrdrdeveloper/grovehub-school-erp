package role

import "grovehub/services/identity/internal/app/constants/permission"

type PermissionSet struct {
	RoleID  string            `json:"roleId"`
	Version int64             `json:"version"`
	Codes   []permission.Code `json:"permissionCodes"`
}

type PreviewPermissions struct {
	Version int64             `json:"version"`
	Codes   []permission.Code `json:"permissionCodes"`
}

type PermissionPreview struct {
	RoleID           string            `json:"roleId"`
	Version          int64             `json:"version"`
	Added            []permission.Code `json:"added"`
	Removed          []permission.Code `json:"removed"`
	ApprovalRequired bool              `json:"approvalRequired"`
}
