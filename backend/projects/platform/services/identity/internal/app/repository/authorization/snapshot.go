package authorization

import (
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
)

type RoleGrant struct {
	RoleID      string
	TenantID    string
	Active      bool
	Permissions []permission.Code
}

type PermissionState struct {
	Active  bool
	Feature feature.Code
}

type FeatureState struct {
	Active  bool
	Enabled bool
}

// Snapshot is the authoritative account and its registered role grants from one database read.
// It must not be populated from an untrusted request.
type Snapshot struct {
	AccountID   string
	TenantID    string
	Status      entity.AccountStatus
	Roles       []RoleGrant
	Permissions map[permission.Code]PermissionState
	Features    map[feature.Code]FeatureState
}

