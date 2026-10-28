package request

type Register struct {
	IdempotencyKey string  `json:"-"`
	Gid            string  `json:"gid,omitempty"`
	FirstName      string  `json:"firstName"`
	MiddleName     *string `json:"middleName"`
	LastName       string  `json:"lastName"`
	Email          string  `json:"email"`
	PhoneNumber    string  `json:"phoneNumber,omitempty"`
	Password       string  `json:"password"`
}
