package request

type ContactValue struct {
	Value     string `json:"value"`
	IsPrimary bool   `json:"isPrimary"`
}

type Contacts struct {
	Emails []ContactValue `json:"emails,omitempty"`
	Phones []ContactValue `json:"phones,omitempty"`
}

type CreatePerson struct {
	FirstName  string   `json:"firstName"`
	MiddleName string   `json:"middleName,omitempty"`
	LastName   string   `json:"lastName"`
	Contacts   Contacts `json:"contacts,omitempty"`
}

type UpdatePerson struct {
	FirstName  *string   `json:"firstName,omitempty"`
	MiddleName *string   `json:"middleName,omitempty"`
	LastName   *string   `json:"lastName,omitempty"`
	Contacts   *Contacts `json:"contacts,omitempty"`
}

type PersonByGID struct {
	GID string `json:"gid"`
}
