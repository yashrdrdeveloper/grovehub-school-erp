package accountrole

type Replace struct {
	Version int64 `json:"version"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}

type RoleVersion struct {
	RoleID string `json:"roleId"`
	Version int64 `json:"version"`
}

type Preview struct {
	AccountID string `json:"accountId"`
	Version int64 `json:"version"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
	RoleVersions []RoleVersion `json:"roleVersions"`
	ApprovalRequired bool `json:"approvalRequired"`
}
