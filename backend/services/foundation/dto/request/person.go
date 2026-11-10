package request

type ContactValue struct {
	Value     string `json:"value" validate:"notblank"`
	IsPrimary bool   `json:"isPrimary"`
}

type Contacts struct {
	Emails []ContactValue `json:"emails,omitempty" validate:"valid"`
	Phones []ContactValue `json:"phones,omitempty" validate:"valid"`
}

// Person groups the common string-valued names for create requests.
type Person struct {
	FirstName  string `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=100"`
	MiddleName string `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	LastName   string `json:"lastName" normalize:"trim" validate:"notblank,maxrunes=100"`
}

type CreatePerson struct {
	Person
	Contacts Contacts `json:"contacts,omitempty" validate:"valid"`
}

// UpdatePerson keeps pointer fields so absent PATCH values remain unchanged.
type UpdatePerson struct {
	FirstName  *string   `json:"firstName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
	MiddleName *string   `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	LastName   *string   `json:"lastName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
	Contacts   *Contacts `json:"contacts,omitempty" validate:"valid"`
}

type PersonByGID struct {
	GID string `json:"gid"`
}
