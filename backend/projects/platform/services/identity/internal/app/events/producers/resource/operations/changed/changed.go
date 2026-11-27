package changed

import (
	"grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/entity"
)

type Role struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Status entity.RoleStatus `json:"status"`
	Personas []persona.Code `json:"personas"`
	Permissions []permission.Code `json:"permissionCodes"`
}

type Account struct {
	Status entity.AccountStatus `json:"status"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}

// Snapshot contains normalized authority fields only; no credentials,
// contact data or legacy account role columns are included.
type Snapshot struct {
	ID string `json:"id"`
	Version int64 `json:"version"`
	Role *Role `json:"role,omitempty"`
	Account *Account `json:"account,omitempty"`
}

type Changed struct {
	TenantID string `json:"tenantId"`
	RequestID string `json:"requestId"`
	ActorID string `json:"actorId"`
	ResourceType string `json:"resourceType"`
	Operation approval.RequestType `json:"operation"`
	Old Snapshot `json:"old"`
	New Snapshot `json:"new"`
}
