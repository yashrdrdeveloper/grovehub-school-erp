package request

type Register struct {
	Gid         string `json:"gid,omitempty"`
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName,omitempty"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Password    string `json:"password"`
}
