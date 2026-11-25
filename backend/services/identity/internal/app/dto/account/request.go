package account

type PersonNames struct {
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
}

type Create struct {
	Person      PersonNames `json:"person"`
	Email       string      `json:"email"`
	PhoneNumber string      `json:"phoneNumber,omitempty"`
	RoleAssignment *RoleAssignment `json:"roleAssignment,omitempty"`
}

type Update struct {
	Username *string `json:"username,omitempty"`
	RoleAssignment *RoleAssignment `json:"roleAssignment,omitempty"`
	Version *int64 `json:"version,omitempty"`
}

type Setup struct {
	Password string `json:"password"`
}
