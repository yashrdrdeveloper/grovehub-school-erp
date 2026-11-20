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
}

type Update struct {
	Username *string `json:"username"`
}

type Setup struct {
	Password string `json:"password"`
}
