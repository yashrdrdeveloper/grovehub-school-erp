package events

// A narrow schema, not a recursive blacklist. Unknown fields (including
// approval Payload hashes, emails, credentials and free-form provider data)
// are never persisted. Existing role descriptions are deliberately omitted.
type Snapshot struct {
	ID          string           `json:"id,omitempty"`
	AccountID   string           `json:"accountId,omitempty"`
	PersonGID   string           `json:"personGid,omitempty"`
	TenantID    string           `json:"tenantId,omitempty"`
	Version     int64            `json:"version"`
	Status      string           `json:"status,omitempty"`
	Username    *string          `json:"username,omitempty"`
	CurrentStep int              `json:"currentStep,omitempty"`
	StepCount   int              `json:"stepCount,omitempty"`
	Role        *RoleSnapshot    `json:"role,omitempty"`
	Account     *AccountSnapshot `json:"account,omitempty"`
}

type RoleSnapshot struct {
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Personas    []string `json:"personas"`
	Permissions []string `json:"permissionCodes"`
}

type AccountSnapshot struct {
	Status           string   `json:"status"`
	PrimaryRoleID    string   `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}

type Payload struct {
	TenantID    string    `json:"tenantId"`
	ActorID     string    `json:"actorId"`
	Operation   string    `json:"operation"`
	RequestID   string    `json:"requestId"`
	RequestType string    `json:"requestType"`
	Old         *Snapshot `json:"old"`
	New         *Snapshot `json:"new"`
}
