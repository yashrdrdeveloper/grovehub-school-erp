package account

// Role describes an assignable operational role, not an access permission.
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleAssignment struct {
	PrimaryRoleID    string   `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}
